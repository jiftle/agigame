package engine

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"

	"agigame/emulator/agent"
	"agigame/emulator/agent/games"
	"agigame/emulator/core/gb"
)

// framePush carries a rendered frame pointer added to the encode queue with
// its tick.
type framePush struct {
	data *[gb.ScreenWidth][gb.ScreenHeight][3]uint8
	tick uint64
}

// Session owns one running GameBoy: emulator, decision agent, input
// aggregation and the encoded frame pipeline. It is transport-agnostic; hosts
// register callbacks and forward them over their own protocol.
type Session struct {
	cfg   Config
	gb    *gb.Gameboy
	agent *agent.Agent
	input *InputState

	auto          atomic.Bool
	frameSkip     uint64
	stateInterval uint64

	frames chan framePush

	mu      sync.RWMutex
	onFrame func(Frame)
	onState func(StateUpdate)
}

// New builds a Session and loads the configured ROM. It does not start
// emulation; call Start.
func New(cfg Config) (*Session, error) {
	cfg.normalize()

	plugin, ok := games.Get(cfg.Game)
	if !ok {
		return nil, fmt.Errorf("unknown game plugin %q", cfg.Game)
	}

	var opts []gb.GameboyOption
	if cfg.CGB {
		opts = append(opts, gb.WithCGBEnabled())
	}
	gameboy, err := gb.New(cfg.ROM, opts...)
	if err != nil {
		return nil, err
	}

	idx, err := PaletteIndex(cfg.Palette)
	if err != nil {
		cfg.Logf("warn", "%v (falling back to greyscale)", err)
		idx = gb.PaletteGreyscale
	}
	gb.SetDMGPalette(idx)

	agentCfg := cfg.Agent
	if agentCfg.Logf == nil {
		agentCfg.Logf = func(format string, args ...any) {
			cfg.Logf("agent", format, args...)
		}
	}

	s := &Session{
		cfg:           cfg,
		gb:            gameboy,
		input:         NewInputState(),
		frameSkip:     uint64(cfg.FrameSkip),
		stateInterval: cfg.StateInterval,
		frames:        make(chan framePush, 1),
	}
	s.agent = agent.New(agentCfg, plugin, s.input)
	return s, nil
}

// SetFrameCallback registers the callback invoked with each encoded frame on
// the encoder goroutine. Must not block.
func (s *Session) SetFrameCallback(cb func(Frame)) {
	s.mu.Lock()
	s.onFrame = cb
	s.mu.Unlock()
}

// SetStateCallback registers the callback invoked with each periodic state
// update on the emulator goroutine. Must not block.
func (s *Session) SetStateCallback(cb func(StateUpdate)) {
	s.mu.Lock()
	s.onState = cb
	s.mu.Unlock()
}

// Start wires the emulator callbacks and runs the emulation loop until ctx is
// cancelled. It blocks; hosts normally call it from a goroutine.
func (s *Session) Start(ctx context.Context) error {
	s.gb.SetInputProvider(s.input)
	s.gb.SetPreFrameCallback(s.onPreFrame)
	s.gb.SetFrameCallback(s.onRawFrame)
	if s.stateInterval > 0 {
		s.gb.SetStateCallback(s.onGBState, s.stateInterval)
	}

	go s.runEncoder(ctx)
	go s.agent.StartLLM(ctx)

	s.cfg.Logf("info", "emulation started")
	s.gb.Run(ctx)
	s.cfg.Logf("info", "emulation stopped")
	return nil
}

// onPreFrame runs on the emulator goroutine before each frame advances and
// fires the decision agent when auto mode is on.
func (s *Session) onPreFrame() {
	if s.auto.Load() {
		s.agent.Tick(s.gb)
	}
}

// onRawFrame copies the fresh frame into the encode queue (latest-wins) so the
// emulator goroutine never blocks.
func (s *Session) onRawFrame(frame *[gb.ScreenWidth][gb.ScreenHeight][3]uint8, tick uint64) {
	if tick%s.frameSkip != 0 {
		return
	}
	cp := *frame
	push := framePush{data: &cp, tick: tick}
	select {
	case s.frames <- push:
	default:
		select {
		case <-s.frames:
		default:
		}
		select {
		case s.frames <- push:
		default:
		}
	}
}

// runEncoder encodes queued frames to PNG and invokes the frame callback.
func (s *Session) runEncoder(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case push := <-s.frames:
			png := encodePNG(push.data)
			if png == nil {
				continue
			}
			s.mu.RLock()
			cb := s.onFrame
			s.mu.RUnlock()
			if cb != nil {
				cb(Frame{PNG: png, Tick: push.tick})
			}
		}
	}
}

// onGBState assembles a state update (with the agent summary in auto mode) and
// invokes the state callback on the emulator goroutine.
func (s *Session) onGBState(state gb.State) {
	u := StateUpdate{State: state, Auto: s.auto.Load()}
	if u.Auto {
		u.Agent = s.agent.Status()
	}
	s.mu.RLock()
	cb := s.onState
	s.mu.RUnlock()
	if cb != nil {
		cb(u)
	}
}

// --- lifecycle ---------------------------------------------------------------

// Reset requests an emulator reset on the next frame.
func (s *Session) Reset() { s.gb.RequestReset() }

// SetPaused pauses or resumes emulation.
func (s *Session) SetPaused(paused bool) { s.gb.SetPaused(paused) }

// IsPaused reports whether emulation is paused.
func (s *Session) IsPaused() bool { return s.gb.IsPaused() }

// --- input -------------------------------------------------------------------

// Press marks the given buttons as held.
func (s *Session) Press(buttons ...gb.Button) { s.input.Press(buttons...) }

// Release marks the given buttons as released.
func (s *Session) Release(buttons ...gb.Button) { s.input.Release(buttons...) }

// SetInput sets a single button's held state (agent InputSink compatible).
func (s *Session) SetInput(button gb.Button, down bool) { s.input.Set(button, down) }

// ClearInput releases every button.
func (s *Session) ClearInput() { s.input.Clear() }

// --- agent / config ----------------------------------------------------------

// SetAuto switches agent control on/off. Turning it off releases agent-held
// buttons and hands control back to manual input.
func (s *Session) SetAuto(v bool) {
	s.auto.Store(v)
	if !v {
		s.input.Clear()
		s.agent.ReleaseAll()
	}
}

// Auto reports whether the agent controls the game.
func (s *Session) Auto() bool { return s.auto.Load() }

// SetMode switches the agent decision layer.
func (s *Session) SetMode(m agent.Mode) { s.agent.SetMode(m) }

// Mode returns the current agent decision layer.
func (s *Session) Mode() agent.Mode { return s.agent.Mode() }

// SetPalette switches the DMG palette by name.
func (s *Session) SetPalette(name string) error {
	idx, err := PaletteIndex(name)
	if err != nil {
		return err
	}
	gb.SetDMGPalette(idx)
	return nil
}

// --- accessors ---------------------------------------------------------------

// Gameboy returns the underlying emulator.
func (s *Session) Gameboy() *gb.Gameboy { return s.gb }

// Agent returns the decision agent.
func (s *Session) Agent() *agent.Agent { return s.agent }

// CartName returns the loaded cartridge name.
func (s *Session) CartName() string { return s.gb.GetLoadedCart().GetName() }

// FrameNumber returns the number of frames emulated so far.
func (s *Session) FrameNumber() uint64 { return s.gb.FrameNumber() }

// Status returns a health/overview snapshot suitable for healthz / the UI.
func (s *Session) Status() map[string]any {
	return map[string]any{
		"cart":   s.CartName(),
		"frames": s.gb.FrameNumber(),
		"auto":   s.auto.Load(),
		"paused": s.gb.IsPaused(),
		"mode":   string(s.agent.Mode()),
		"stats":  s.agent.Status(),
	}
}

// defaultLogf is used when Config.Logf is nil.
func defaultLogf(level, format string, args ...any) {
	log.Printf("[%s] %s", level, fmt.Sprintf(format, args...))
}

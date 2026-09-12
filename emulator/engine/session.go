package engine

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"agigame/emulator/agent"
	"agigame/emulator/agent/games"
	"agigame/emulator/core/gb"
	"agigame/emulator/gba"
)

// framePush carries a rendered frame pointer added to the encode queue with
// its tick.
type framePush struct {
	data *[gb.ScreenWidth][gb.ScreenHeight][3]uint8
	tick uint64
}

// Session owns one running handheld: emulator, decision agent, input
// aggregation and the encoded frame pipeline. It is transport-agnostic; hosts
// register callbacks and forward them over their own protocol.
type Session struct {
	cfg     Config
	console string
	gb      *gb.Gameboy
	gba     *gba.Console
	agent   *agent.Agent
	input   *InputState

	auto          atomic.Bool
	frameSkip     uint64
	stateInterval uint64

	frames chan framePush

	mu      sync.RWMutex
	onFrame func(Frame)
	onState func(StateUpdate)
	onAudio func(Audio)
}

// New builds a Session and loads the configured ROM. It does not start
// emulation; call Start.
func New(cfg Config) (*Session, error) {
	cfg.normalize()

	if cfg.Console == "gba" {
		core, err := gba.New(cfg.ROM)
		if err != nil {
			return nil, err
		}
		return &Session{
			cfg:           cfg,
			console:       "gba",
			gba:           core,
			frameSkip:     uint64(cfg.FrameSkip),
			stateInterval: cfg.StateInterval,
		}, nil
	}

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
		console:       "gb",
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

// SetAudioCallback registers the callback invoked with each PCM chunk on the
// emulator goroutine. It must not block; hosts should copy or enqueue quickly.
func (s *Session) SetAudioCallback(cb func(Audio)) {
	s.mu.Lock()
	s.onAudio = cb
	s.mu.Unlock()
}

// emitAudio forwards a PCM chunk to the audio callback (if any).
func (s *Session) emitAudio(a Audio) {
	s.mu.RLock()
	cb := s.onAudio
	s.mu.RUnlock()
	if cb != nil {
		cb(a)
	}
}

// Start wires the emulator callbacks and runs the emulation loop until ctx is
// cancelled. It blocks; hosts normally call it from a goroutine.
func (s *Session) Start(ctx context.Context) error {
	if s.gba != nil {
		return s.startGBA(ctx)
	}
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

// startGBA drives the GBA core with our own frame loop. The agent layer is
// not supported for GBA yet, so frames and state are emitted without it.
func (s *Session) startGBA(ctx context.Context) error {
	fps := gba.FPS
	interval := time.Duration(float64(time.Second) / fps)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	s.cfg.Logf("info", "emulation started (gba %s)", s.CartName())
	for {
		select {
		case <-ctx.Done():
			s.cfg.Logf("info", "emulation stopped")
			return nil
		case <-ticker.C:
			if s.gba.Paused() {
				continue
			}
			s.gba.Step()
			tick := s.gba.FrameCount()

			// Always drain PCM to release the APU ring's backpressure.
			if pcm := s.gba.DrainPCM(); len(pcm) > 0 {
				s.emitAudio(Audio{PCM: pcm, SampleRate: s.gba.SampleRate()})
			}

			if tick%s.frameSkip == 0 {
				if png := encodeRGBA(s.gba.Pixels(), gba.Width, gba.Height); png != nil {
					s.mu.RLock()
					cb := s.onFrame
					s.mu.RUnlock()
					if cb != nil {
						cb(Frame{PNG: png, Tick: tick})
					}
				}
			}
			if s.stateInterval > 0 && tick%s.stateInterval == 0 {
				s.emitState(StateUpdate{
					Console: "gba",
					Width:   gba.Width,
					Height:  gba.Height,
				})
			}
		}
	}
}

// emitState invokes the state callback (if any).
func (s *Session) emitState(u StateUpdate) {
	s.mu.RLock()
	cb := s.onState
	s.mu.RUnlock()
	if cb != nil {
		cb(u)
	}
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
	u := StateUpdate{
		State:   state,
		Auto:    s.auto.Load(),
		Console: "gb",
		Width:   gb.ScreenWidth,
		Height:  gb.ScreenHeight,
	}
	if u.Auto {
		u.Agent = s.agent.Status()
	}
	s.emitState(u)
}

// --- lifecycle ---------------------------------------------------------------

// Reset requests an emulator reset on the next frame.
func (s *Session) Reset() {
	if s.gba != nil {
		s.gba.Reset()
		return
	}
	s.gb.RequestReset()
}

// SetPaused pauses or resumes emulation.
func (s *Session) SetPaused(paused bool) {
	if s.gba != nil {
		s.gba.SetPaused(paused)
		return
	}
	s.gb.SetPaused(paused)
}

// IsPaused reports whether emulation is paused.
func (s *Session) IsPaused() bool {
	if s.gba != nil {
		return s.gba.Paused()
	}
	return s.gb.IsPaused()
}

// --- input -------------------------------------------------------------------

// Press marks the given buttons as held.
func (s *Session) Press(buttons ...gb.Button) {
	if s.gba != nil {
		for _, b := range buttons {
			s.gba.SetButton(b, true)
		}
		return
	}
	s.input.Press(buttons...)
}

// Release marks the given buttons as released.
func (s *Session) Release(buttons ...gb.Button) {
	if s.gba != nil {
		for _, b := range buttons {
			s.gba.SetButton(b, false)
		}
		return
	}
	s.input.Release(buttons...)
}

// SetInput sets a single button's held state (agent InputSink compatible).
func (s *Session) SetInput(button gb.Button, down bool) {
	if s.gba != nil {
		s.gba.SetButton(button, down)
		return
	}
	s.input.Set(button, down)
}

// ClearInput releases every button.
func (s *Session) ClearInput() {
	if s.gba != nil {
		s.gba.ReleaseAll()
		return
	}
	s.input.Clear()
}

// --- agent / config ----------------------------------------------------------

// SetAuto switches agent control on/off. Turning it off releases agent-held
// buttons and hands control back to manual input. No-op for consoles without
// an agent plugin (gba).
func (s *Session) SetAuto(v bool) {
	if s.gba != nil || s.agent == nil {
		return
	}
	s.auto.Store(v)
	if !v {
		s.input.Clear()
		s.agent.ReleaseAll()
	}
}

// Auto reports whether the agent controls the game.
func (s *Session) Auto() bool {
	if s.gba != nil || s.agent == nil {
		return false
	}
	return s.auto.Load()
}

// SetMode switches the agent decision layer.
func (s *Session) SetMode(m agent.Mode) {
	if s.agent != nil {
		s.agent.SetMode(m)
	}
}

// Mode returns the current agent decision layer.
func (s *Session) Mode() agent.Mode {
	if s.agent == nil {
		return agent.ModeManual
	}
	return s.agent.Mode()
}

// SetPalette switches the DMG palette by name (no-op for gba).
func (s *Session) SetPalette(name string) error {
	if s.gba != nil {
		return nil
	}
	idx, err := PaletteIndex(name)
	if err != nil {
		return err
	}
	gb.SetDMGPalette(idx)
	return nil
}

// --- accessors ---------------------------------------------------------------

// Gameboy returns the underlying emulator (nil for gba).
func (s *Session) Gameboy() *gb.Gameboy { return s.gb }

// Agent returns the decision agent (nil for gba).
func (s *Session) Agent() *agent.Agent { return s.agent }

// Console returns the active console kind ("gb" or "gba").
func (s *Session) Console() string { return s.console }

// Width returns the active screen width.
func (s *Session) Width() int {
	if s.gba != nil {
		return gba.Width
	}
	return gb.ScreenWidth
}

// Height returns the active screen height.
func (s *Session) Height() int {
	if s.gba != nil {
		return gba.Height
	}
	return gb.ScreenHeight
}

// CartName returns the loaded cartridge name.
func (s *Session) CartName() string {
	if s.gba != nil {
		return s.gba.CartName()
	}
	return s.gb.GetLoadedCart().GetName()
}

// FrameNumber returns the number of frames emulated so far.
func (s *Session) FrameNumber() uint64 {
	if s.gba != nil {
		return s.gba.FrameCount()
	}
	return s.gb.FrameNumber()
}

// Status returns a health/overview snapshot suitable for healthz / the UI.
func (s *Session) Status() map[string]any {
	out := map[string]any{
		"console": s.console,
		"cart":    s.CartName(),
		"frames":  s.FrameNumber(),
		"width":   s.Width(),
		"height":  s.Height(),
		"auto":    s.Auto(),
		"paused":  s.IsPaused(),
		"mode":    string(s.Mode()),
	}
	if s.agent != nil {
		out["stats"] = s.agent.Status()
	}
	return out
}

// defaultLogf is used when Config.Logf is nil.
func defaultLogf(level, format string, args ...any) {
	log.Printf("[%s] %s", level, fmt.Sprintf(format, args...))
}

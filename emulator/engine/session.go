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
	gbconsole "agigame/emulator/gb"
	"agigame/emulator/gba"
)

// targetFPS is the real refresh rate shared by both handhelds (59.7275 Hz).
const targetFPS = 59.727500569606

// consoleCore is the transport-agnostic interface implemented by each
// emulated handheld (guac-backed gb and gba adapters).
type consoleCore interface {
	Name() string
	Width() int
	Height() int
	Step()
	Pixels() []byte
	DrainPCM() []byte
	SampleRate() int
	SetButton(gb.Button, bool)
	Set(gb.Button, bool) // agent InputSink
	ReleaseAll()
	Reset()
	SetPaused(bool)
	Paused() bool
	FrameCount() uint64
	CartName() string
	ReadMemory(uint16) byte // agent GameReader
	Snapshot() gb.State
}

// framePush carries a copy of an RGBA framebuffer for off-thread encoding.
type framePush struct {
	pixels []byte
	width  int
	height int
	tick   uint64
}

// Session owns one running handheld: emulator, optional decision agent, input
// and the encoded frame/audio pipeline. It is transport-agnostic; hosts
// register callbacks and forward them over their own protocol.
type Session struct {
	cfg     Config
	console string
	core    consoleCore
	agent   *agent.Agent // non-nil only for gb (guac gba has no agent plugin)

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

	s := &Session{
		cfg:           cfg,
		frameSkip:     uint64(cfg.FrameSkip),
		stateInterval: cfg.StateInterval,
		frames:        make(chan framePush, 1),
	}

	if cfg.Console == "gba" {
		core, err := gba.New(cfg.ROM)
		if err != nil {
			return nil, err
		}
		s.console = "gba"
		s.core = core
		return s, nil
	}

	plugin, ok := games.Get(cfg.Game)
	if !ok {
		return nil, fmt.Errorf("unknown game plugin %q", cfg.Game)
	}
	core, err := gbconsole.New(cfg.ROM)
	if err != nil {
		return nil, err
	}
	if err := core.SetPalette(cfg.Palette); err != nil {
		cfg.Logf("warn", "%v (keeping default palette)", err)
	}

	agentCfg := cfg.Agent
	if agentCfg.Logf == nil {
		agentCfg.Logf = func(format string, args ...any) {
			cfg.Logf("agent", format, args...)
		}
	}

	s.console = "gb"
	s.core = core
	s.agent = agent.New(agentCfg, plugin, core)
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
// emulator goroutine. Must not block.
func (s *Session) SetAudioCallback(cb func(Audio)) {
	s.mu.Lock()
	s.onAudio = cb
	s.mu.Unlock()
}

func (s *Session) emitState(u StateUpdate) {
	s.mu.RLock()
	cb := s.onState
	s.mu.RUnlock()
	if cb != nil {
		cb(u)
	}
}

func (s *Session) emitAudio(a Audio) {
	s.mu.RLock()
	cb := s.onAudio
	s.mu.RUnlock()
	if cb != nil {
		cb(a)
	}
}

// Start runs the emulation loop until ctx is cancelled. It blocks; hosts
// normally call it from a goroutine.
func (s *Session) Start(ctx context.Context) error {
	go s.runEncoder(ctx)
	if s.agent != nil {
		go s.agent.StartLLM(ctx)
	}
	return s.run(ctx)
}

// run drives the core with accumulated-time pacing: the number of frames that
// should have elapsed since the start is computed from the real frame rate and
// any shortfall is caught up (bounded), so a slow frame cannot accumulate into
// drift and the game never runs faster or slower than real hardware.
func (s *Session) run(ctx context.Context) error {
	start := time.Now()
	var frame uint64
	const maxCatchUp = 4

	s.cfg.Logf("info", "emulation started (%s %s)", s.console, s.core.CartName())
	for {
		frame++
		target := start.Add(time.Duration(float64(frame) / targetFPS * float64(time.Second)))
		if wait := time.Until(target); wait > 0 {
			timer := time.NewTimer(wait)
			select {
			case <-ctx.Done():
				timer.Stop()
				s.cfg.Logf("info", "emulation stopped")
				return nil
			case <-timer.C:
			}
		}

		s.step()
		for i := 0; i < maxCatchUp; i++ {
			if time.Since(target) < 0 {
				break
			}
			s.step()
			frame++
		}
	}
}

// step advances one frame and emits audio/frames/state.
func (s *Session) step() {
	if s.core.Paused() {
		return
	}
	if s.agent != nil && s.auto.Load() {
		s.agent.Tick(s.core)
	}

	s.core.Step()
	tick := s.core.FrameCount()

	// Always drain PCM to release the APU ring's backpressure.
	if pcm := s.core.DrainPCM(); len(pcm) > 0 {
		s.emitAudio(Audio{PCM: pcm, SampleRate: s.core.SampleRate()})
	}

	if s.frameSkip == 0 || tick%s.frameSkip == 0 {
		s.queueFrame(tick)
	}
	if s.stateInterval > 0 && tick%s.stateInterval == 0 {
		u := StateUpdate{
			State:   s.core.Snapshot(),
			Console: s.console,
			Width:   s.core.Width(),
			Height:  s.core.Height(),
		}
		if s.agent != nil && s.auto.Load() {
			u.Auto = true
			u.Agent = s.agent.Status()
		}
		s.emitState(u)
	}
}

// queueFrame copies the framebuffer into the encode queue (latest-wins) so PNG
// encoding and broadcasting stay off the emulation goroutine.
func (s *Session) queueFrame(tick uint64) {
	pixels := s.core.Pixels()
	cp := make([]byte, len(pixels))
	copy(cp, pixels)

	push := framePush{pixels: cp, width: s.core.Width(), height: s.core.Height(), tick: tick}
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

// runEncoder PNG-encodes queued frames and invokes the frame callback.
func (s *Session) runEncoder(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case push := <-s.frames:
			png := encodeRGBA(push.pixels, push.width, push.height)
			if png == nil {
				continue
			}
			s.mu.RLock()
			cb := s.onFrame
			s.mu.RUnlock()
			if cb != nil {
				cb(Frame{
					PNG:    png,
					Tick:   push.tick,
					RGBA:   push.pixels,
					Width:  push.width,
					Height: push.height,
				})
			}
		}
	}
}

// --- lifecycle ---------------------------------------------------------------

// Reset reboots the emulator.
func (s *Session) Reset() { s.core.Reset() }

// SetPaused pauses or resumes emulation.
func (s *Session) SetPaused(paused bool) { s.core.SetPaused(paused) }

// IsPaused reports whether emulation is paused.
func (s *Session) IsPaused() bool { return s.core.Paused() }

// --- input -------------------------------------------------------------------

// Press marks the given buttons as held.
func (s *Session) Press(buttons ...gb.Button) {
	for _, b := range buttons {
		s.core.SetButton(b, true)
	}
}

// Release marks the given buttons as released.
func (s *Session) Release(buttons ...gb.Button) {
	for _, b := range buttons {
		s.core.SetButton(b, false)
	}
}

// SetInput sets a single button's held state.
func (s *Session) SetInput(button gb.Button, down bool) { s.core.SetButton(button, down) }

// ClearInput releases every button.
func (s *Session) ClearInput() { s.core.ReleaseAll() }

// --- agent / config ----------------------------------------------------------

// SetAuto switches agent control on/off (gb only).
func (s *Session) SetAuto(v bool) {
	if s.agent == nil {
		return
	}
	s.auto.Store(v)
	if !v {
		s.core.ReleaseAll()
		s.agent.ReleaseAll()
	}
}

// Auto reports whether the agent controls the game.
func (s *Session) Auto() bool {
	if s.agent == nil {
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
	if p, ok := s.core.(interface{ SetPalette(string) error }); ok {
		return p.SetPalette(name)
	}
	return nil
}

// --- accessors ---------------------------------------------------------------

// Gameboy returns nil (the GoBoy core has been replaced by the guac adapter).
func (s *Session) Gameboy() *gb.Gameboy { return nil }

// Agent returns the decision agent (nil for gba).
func (s *Session) Agent() *agent.Agent { return s.agent }

// Console returns the active console kind ("gb" or "gba").
func (s *Session) Console() string { return s.console }

// Width returns the active screen width.
func (s *Session) Width() int { return s.core.Width() }

// Height returns the active screen height.
func (s *Session) Height() int { return s.core.Height() }

// CartName returns the loaded cartridge name.
func (s *Session) CartName() string { return s.core.CartName() }

// FrameNumber returns the number of frames emulated so far.
func (s *Session) FrameNumber() uint64 { return s.core.FrameCount() }

// Status returns a health/overview snapshot suitable for healthz / the UI.
func (s *Session) Status() map[string]any {
	out := map[string]any{
		"console": s.console,
		"cart":    s.core.CartName(),
		"frames":  s.core.FrameCount(),
		"width":   s.core.Width(),
		"height":  s.core.Height(),
		"auto":    s.Auto(),
		"paused":  s.core.Paused(),
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

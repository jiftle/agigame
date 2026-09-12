// Package agent implements the decision layer that drives the emulator:
//
//   - hard rules run synchronously on the emulator goroutine once per frame
//   - an optional LLM runs asynchronously (2Hz by default) and its latest
//     decision is consumed by the tick when the mode asks for it
//   - reward / stats aggregate progress for the UI and future P4 experiments
//
// Modes: "rules", "llm", "hybrid" (rules first, LLM otherwise; rules can be
// an immediate safety net), and "manual" (agent paused; keyboard only).
package agent

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"agigame/agent/games"
	"agigame/core/gb"
)

// Mode selects which decision layer drives the buttons.
type Mode string

const (
	ModeRules  Mode = "rules"
	ModeLLM    Mode = "llm"
	ModeHybrid Mode = "hybrid"
	ModeManual Mode = "manual"
)

// InputSink is anything that can set a GameBoy button state. server.InputState
// implements it (thread safe).
type InputSink interface {
	Set(button gb.Button, down bool)
}

// Config carries the agent settings (mirrors config.yaml agent section).
type Config struct {
	Mode          Mode
	LLMIntervalMs int
	// SafetyNetFrames is the stall length in frames tolerated before the
	// safety net forces a rescue (0 = default 120).
	SafetyNetFrames int
	// Provider factory hook so P3+ can swap in a real LLM.
	Provider LLMProvider
	// Logf optionally receives agent log lines; defaults to log.Printf.
	Logf func(format string, args ...any)
}

// Agent owns the per-frame decision pipeline.
type Agent struct {
	cfg     Config
	plugin  games.GamePlugin
	sink    InputSink

	mode atomic.Value // Mode

	mu        sync.Mutex
	lastState any

	lastDecision atomic.Pointer[LLMDecision]

	prevButtons map[gb.Button]bool
	prevState   any
	net         *SafetyNet
	stats       *Stats
}

// New creates an Agent. plugin and sink are required; a nil provider falls
// back to the StubLLM.
func New(cfg Config, plugin games.GamePlugin, sink InputSink) *Agent {
	if cfg.Provider == nil {
		cfg.Provider = &StubLLM{}
	}
	if cfg.Logf == nil {
		cfg.Logf = log.Printf
	}
	if cfg.Mode == "" {
		cfg.Mode = ModeManual
	}
	a := &Agent{
		cfg:         cfg,
		plugin:      plugin,
		sink:        sink,
		prevButtons: map[gb.Button]bool{},
		net:         NewSafetyNet(cfg.SafetyNetFrames),
		stats:       &Stats{},
	}
	a.mode.Store(cfg.Mode)
	return a
}

// SetMode switches the decision layer at runtime (webui Auto/Manual toggles).
func (a *Agent) SetMode(m Mode) {
	a.mode.Store(m)
	a.cfg.Logf("agent: mode -> %s", m)
}

// Mode returns the current mode.
func (a *Agent) Mode() Mode {
	if m, ok := a.mode.Load().(Mode); ok {
		return m
	}
	return ModeManual
}

// Status returns a copy of agent stats plus the game's extra fields for the
// UI state message / healthz. Safe from any goroutine.
func (a *Agent) Status() map[string]any {
	out := a.stats.Snapshot()
	out["mode"] = string(a.Mode())
	out["provider"] = a.cfg.Provider.Name()
	a.mu.Lock()
	if a.prevState != nil {
		for k, v := range a.plugin.ExtraState(a.prevState) {
			out[k] = v
		}
	}
	a.mu.Unlock()
	if d := a.lastDecision.Load(); d != nil {
		out["decision"] = d.Commands
		out["decisionRationale"] = d.Rationale
	}
	return out
}

// Tick is called once per emulated frame on the emulator goroutine. It reads
// the current state, computes reward, fires the synchronous decision layer,
// and diffs the applied buttons into the input sink. It must stay fast and
// must not block.
func (a *Agent) Tick(r games.GameReader) {
	if r == nil {
		return
	}

	state, err := a.plugin.ExtractState(r)
	if err != nil {
		a.cfg.Logf("agent: extract state: %v", err)
		return
	}

	// Keep a copy for the async LLM goroutine.
	a.mu.Lock()
	a.lastState = state
	a.mu.Unlock()

	// Reward from previous -> current.
	if a.prevState != nil {
		ri := a.plugin.Reward(a.prevState, state)
		a.stats.Update(ri)
		if ri.GameOver {
			a.cfg.Logf("agent: game over (progress=%d)", ri.Progress)
		}
	}
	a.prevState = state

	// Decision (rules layer always available; LLM consumed per mode).
	buttons := a.decide(state)

	a.applyInput(buttons)
}

// decide runs the mode-appropriate decision and applies the safety net.
func (a *Agent) decide(state any) games.Buttons {
	var b games.Buttons
	switch a.Mode() {
	case ModeManual:
		return games.Buttons{}
	case ModeLLM:
		b = decisionToButtons(a.lastDecision.Load())
	case ModeHybrid:
		b = a.plugin.Decide(state)
		if !a.plugin.NeedsLLM(state, b) {
			break
		}
		b = decisionToButtons(a.lastDecision.Load())
	default: // ModeRules
		b = a.plugin.Decide(state)
	}

	// Stuck detection uses the reward progress.
	progress := 0
	if ps, ok := state.(interface{ Progress() int }); ok {
		progress = ps.Progress()
	}
	if a.net.Observe(progress) {
		b = a.net.Rescue(b)
	}
	return b
}

// applyInput diffs the desired held-buttons against what is currently held
// and pushes the deltas to the sink.
func (a *Agent) applyInput(b games.Buttons) {
	desired := b.ToSet()
	for btn, down := range desired {
		if a.prevButtons[btn] != down {
			a.sink.Set(btn, down)
		}
	}
	a.prevButtons = desired
}

// ReleaseAll releases every button the agent currently holds and clears the
// latest LLM decision. Call it when handing control back to manual play.
func (a *Agent) ReleaseAll() {
	for btn, down := range a.prevButtons {
		if down {
			a.sink.Set(btn, false)
		}
	}
	a.prevButtons = map[gb.Button]bool{}
	a.lastDecision.Store(nil)
	a.net = NewSafetyNet(a.cfg.SafetyNetFrames)
	a.mu.Lock()
	a.lastState = nil
	a.mu.Unlock()
}

// StartLLM runs the async LLM loop at the configured interval. It copies the
// latest extracted state, renders a prompt, calls the provider, and stores the
// decision. It never touches the emulator directly.
func (a *Agent) StartLLM(ctx context.Context) {
	if a.cfg.LLMIntervalMs <= 0 {
		a.cfg.LLMIntervalMs = 500
	}
	tick := time.NewTicker(time.Duration(a.cfg.LLMIntervalMs) * time.Millisecond)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			a.llmStep(ctx)
		}
	}
}

func (a *Agent) llmStep(ctx context.Context) {
	a.mu.Lock()
	state := a.lastState
	a.mu.Unlock()
	if state == nil {
		return
	}
	prompt := BuildGamePrompt(a.plugin, state)
	dec, err := a.cfg.Provider.Decide(ctx, prompt)
	if err != nil {
		a.cfg.Logf("agent: llm %s: %v", a.cfg.Provider.Name(), err)
		return
	}
	a.lastDecision.Store(&dec)
	a.stats.SetDecision(dec.Rationale)
	a.cfg.Logf("agent: llm decision: %v", dec.Commands)
}
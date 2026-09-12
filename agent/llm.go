package agent

import (
	"context"
	"strings"
	"time"

	"agigame/agent/games"
	"agigame/core/gb"
)

// LLMDecision is the parsed output of an LLM provider.
type LLMDecision struct {
	// Commands are button combos such as "A", "Right", "A+Right".
	Commands []string `json:"commands"`
	// Rationale is the model's explanation (shown in the UI log).
	Rationale string `json:"rationale"`
}

// ButtonNames maps a canonical button name to the gb.Button it controls.
// The LLM prompt and the command tokens use these names.
var ButtonNames = map[string]gb.Button{
	"A":      gb.ButtonA,
	"B":      gb.ButtonB,
	"Start":  gb.ButtonStart,
	"Select": gb.ButtonSelect,
	"Up":     gb.ButtonUp,
	"Down":   gb.ButtonDown,
	"Left":   gb.ButtonLeft,
	"Right":  gb.ButtonRight,
}

// LLMProvider turns a rendered prompt into a button decision. Only the Stub
// is wired up for now (P3 skeleton, no real API required).
type LLMProvider interface {
	// Name identifies the provider for logging / healthz.
	Name() string
	Decide(ctx context.Context, prompt string) (LLMDecision, error)
}

// StubLLM is the placeholder provider. It never calls an external API:
// it just echoes the rules' obvious default after a simulated latency, so the
// full pipeline (async loop, prompt build, decision -> buttons) can be
// exercised end to end without credentials.
type StubLLM struct{}

// Name implements LLMProvider.
func (s *StubLLM) Name() string { return "stub" }

// Decide implements LLMProvider.
func (s *StubLLM) Decide(ctx context.Context, _ string) (LLMDecision, error) {
	select {
	case <-time.After(30 * time.Millisecond):
	case <-ctx.Done():
		return LLMDecision{}, ctx.Err()
	}
	return LLMDecision{
		Commands:  []string{"Right"},
		Rationale: "stub provider: replace with a real LLM in P3+",
	}, nil
}

// decisionToButtons flattens an LLMDecision's command combos into a Buttons set.
func decisionToButtons(d *LLMDecision) games.Buttons {
	out := games.Buttons{}
	if d == nil {
		return out
	}
	set := map[gb.Button]bool{}
	for _, cmd := range d.Commands {
		for _, name := range strings.Split(cmd, "+") {
			if b, ok := ButtonNames[strings.TrimSpace(name)]; ok {
				set[b] = true
			}
		}
	}
	out.A = set[gb.ButtonA]
	out.B = set[gb.ButtonB]
	out.Start = set[gb.ButtonStart]
	out.Select = set[gb.ButtonSelect]
	out.Up = set[gb.ButtonUp]
	out.Down = set[gb.ButtonDown]
	out.Left = set[gb.ButtonLeft]
	out.Right = set[gb.ButtonRight]
	return out
}

// AvailableButtonsText lists the buttons for the prompt schema.
func AvailableButtonsText() string {
	names := make([]string, 0, len(ButtonNames))
	for n := range ButtonNames {
		names = append(names, n)
	}
	return strings.Join(names, " ")
}
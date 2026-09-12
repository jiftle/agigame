package agent

import (
	"context"
	"sync"
	"testing"

	"agigame/emulator/agent/games"
	"agigame/emulator/core/gb"
)

// --- fakes ----------------------------------------------------------------

type fakeState struct {
	cam    int
	ground bool
}

func (s *fakeState) Progress() int { return s.cam }

type fakeReader struct{}

func (fakeReader) ReadMemory(uint16) byte { return 0 }

type fakePlugin struct {
	cam int
}

func (f *fakePlugin) ID() string   { return "fake" }
func (f *fakePlugin) Name() string { return "Fake" }
func (f *fakePlugin) ExtractState(games.GameReader) (any, error) {
	return &fakeState{cam: f.cam, ground: true}, nil
}
func (f *fakePlugin) Decide(any) games.Buttons         { return games.Buttons{Right: true} }
func (f *fakePlugin) NeedsLLM(any, games.Buttons) bool { return false }
func (f *fakePlugin) Prompt(any) string                { return "fake prompt" }
func (f *fakePlugin) Reward(prev, cur any) games.RewardInfo {
	b := prev.(*fakeState)
	a := cur.(*fakeState)
	return games.RewardInfo{Progress: a.cam, Delta: float64(a.cam - b.cam)}
}
func (f *fakePlugin) ExtraState(any) map[string]any { return map[string]any{"cam": f.cam} }

type fakeSink struct {
	mu   sync.Mutex
	held map[gb.Button]bool
}

func newFakeSink() *fakeSink { return &fakeSink{held: map[gb.Button]bool{}} }
func (f *fakeSink) Set(b gb.Button, down bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.held[b] = down
}
func (f *fakeSink) Down(b gb.Button) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.held[b]
}

// --- tests ---------------------------------------------------------------

func TestTickRulesMode(t *testing.T) {
	sink := newFakeSink()
	a := New(Config{Mode: ModeRules}, &fakePlugin{cam: 100}, sink)
	defer func() {
		a.ReleaseAll()
	}()

	a.Tick(fakeReader{})
	if !sink.Down(gb.ButtonRight) {
		t.Fatalf("rules mode should hold right")
	}
	if sink.Down(gb.ButtonA) {
		t.Fatalf("no jump expected while scouting")
	}
}

func TestApplyButtonDiff(t *testing.T) {
	sink := newFakeSink()
	a := New(Config{Mode: ModeLLM}, newFakePluginLLM(nil), sink)

	// No decision yet → nothing held.
	a.Tick(fakeReader{})
	if sink.Down(gb.ButtonRight) {
		t.Fatalf("nothing should be held without an LLM decision")
	}

	// Seed a decision and switch to llm mode; the decision drives input.
	a.lastDecision.Store(&LLMDecision{Commands: []string{"A+Right"}})
	a.Tick(fakeReader{})
	if !sink.Down(gb.ButtonRight) || !sink.Down(gb.ButtonA) {
		t.Fatalf("llm decision should map to A+Right, got right=%v a=%v",
			sink.Down(gb.ButtonRight), sink.Down(gb.ButtonA))
	}
}

func TestReleaseAll(t *testing.T) {
	sink := newFakeSink()
	a := New(Config{Mode: ModeRules}, &fakePlugin{cam: 90}, sink)
	a.lastDecision.Store(&LLMDecision{Commands: []string{"Right"}})

	a.Tick(fakeReader{})
	if !sink.Down(gb.ButtonRight) {
		t.Fatalf("precondition: right should be held")
	}
	a.ReleaseAll()
	for _, b := range []gb.Button{
		gb.ButtonA, gb.ButtonB, gb.ButtonStart, gb.ButtonSelect,
		gb.ButtonUp, gb.ButtonDown, gb.ButtonLeft, gb.ButtonRight,
	} {
		if sink.Down(b) {
			t.Fatalf("ReleaseAll should release %v", b)
		}
	}
	if a.lastDecision.Load() != nil {
		t.Fatalf("ReleaseAll should clear the last decision")
	}
}

func TestSafetyNet(t *testing.T) {
	sink := newFakeSink()
	a := New(Config{Mode: ModeRules, SafetyNetFrames: 3}, &fakePlugin{cam: 1}, sink)

	// Same progress every tick: after the initial observation, 3 further
	// stifled frames trip the threshold of 3.
	a.Tick(fakeReader{})
	a.Tick(fakeReader{})
	a.Tick(fakeReader{})
	a.Tick(fakeReader{})
	if !a.net.Fired() {
		t.Fatalf("safety net should fire after 3 stalled frames")
	}
	if !sink.Down(gb.ButtonRight) {
		t.Fatalf("rescue should hold right")
	}
}

func TestLLMStep(f *testing.T) {
	sink := newFakeSink()
	a := New(Config{Mode: ModeLLM, LLMIntervalMs: 500},
		&fakePlugin{cam: 1}, sink)

	// Seed a state first so llmStep has something to prompt with.
	a.Tick(fakeReader{})
	if a.lastState == nil {
		f.Fatalf("lastState should be cached after tick")
	}

	a.llmStep(context.Background())
	d := a.lastDecision.Load()
	if d == nil || len(d.Commands) == 0 {
		f.Fatalf("stub llm should produce a decision, got %v", d)
	}
	a.Tick(fakeReader{})
	if !sink.Down(gb.ButtonRight) {
		f.Fatalf("llm decision should ultimately hold right")
	}
}

func TestDecisionToButtons(t *testing.T) {
	cases := []struct {
		combo string
		want  games.Buttons
	}{
		{combo: "Right", want: games.Buttons{Right: true}},
		{combo: "A+Right", want: games.Buttons{A: true, Right: true}},
		{combo: "A", want: games.Buttons{A: true}},
		{combo: "", want: games.Buttons{}},
	}
	for _, c := range cases {
		got := decisionToButtons(&LLMDecision{Commands: []string{c.combo}})
		if got != c.want {
			t.Errorf("decisionToButtons(%q) = %+v, want %+v", c.combo, got, c.want)
		}
	}
}

// newFakePluginLLM is a thin wrapper to reuse fakePlugin in llm mode test.
func newFakePluginLLM(cam *int) *fakePlugin {
	return &fakePlugin{cam: 1}
}

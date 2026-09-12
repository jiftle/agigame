package agent

import (
	"sync"

	"agigame/agent/games"
)

// Stats aggregates reward signals and notable events for a run. It is safe
// for concurrent use (LLM goroutine vs emulator goroutine).
type Stats struct {
	mu sync.Mutex

	Frames        int
	Progress      int
	MaxProgress   int
	RewardTotal   float64
	Deaths        int
	GameOver      bool
	Victory       bool
	LastDecision  string // last LLM rationale, for the UI
	Decisions     int    // number of LLM decisions consumed
	LastEvent     string
	FrozenFrames  int // frames with no reward progress (>0 => stuck)
	UpdatedFrozen bool
}

// Update folds one RewardInfo transition into the stats.
func (s *Stats) Update(ri games.RewardInfo) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Frames++
	if ri.GameOver {
		s.GameOver = true
	}
	if ri.Victory {
		s.Victory = true
	}
	if ri.Progress > s.MaxProgress {
		s.MaxProgress = ri.Progress
	}
	if ri.Progress > 0 {
		// FrozenFrames only grows while progress stalls (used by the
		// safety net's UI). Any forward movement resets it.
		if ri.Progress > s.Progress {
			s.Progress = ri.Progress
			s.FrozenFrames = 0
		} else {
			s.FrozenFrames++
		}
	}
	for _, ev := range ri.Events {
		if ev == "death" {
			s.Deaths++
		}
		s.LastEvent = ev
	}
	s.RewardTotal += ri.Delta
}

// SetDecision records the most recent consumed LLM decision.
func (s *Stats) SetDecision(rationale string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Decisions++
	s.LastDecision = rationale
}

// Snapshot returns a copy of the stats for the UI / healthz.
func (s *Stats) Snapshot() map[string]any {
	s.mu.Lock()
	defer s.mu.Unlock()
	return map[string]any{
		"frames":       s.Frames,
		"progress":     s.Progress,
		"maxProgress":  s.MaxProgress,
		"rewardTotal":  s.RewardTotal,
		"deaths":       s.Deaths,
		"gameOver":     s.GameOver,
		"victory":      s.Victory,
		"decisions":    s.Decisions,
		"lastDecision": s.LastDecision,
		"lastEvent":    s.LastEvent,
		"frozenFrames": s.FrozenFrames,
	}
}

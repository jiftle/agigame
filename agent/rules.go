package agent

import "agigame/agent/games"

// SafetyNet is a fallback that kicks in when the active decision layer stops
// making progress (stuck in a wall, corner, or softlock). It forces forward
// movement and a periodic jump until reward progress resumes.
type SafetyNet struct {
	// ThresholdFrames is how long a stalled reward is tolerated before
	// the net forces a rescue (frames, ~16ms each).
	ThresholdFrames int

	stuckFrames  int
	lastProgress int
	fired        bool
	pulse        int
}

// NewSafetyNet builds a net that rescues after the given stall (in frames).
func NewSafetyNet(thresholdFrames int) *SafetyNet {
	if thresholdFrames <= 0 {
		thresholdFrames = 120 // ~2s at 60fps
	}
	return &SafetyNet{ThresholdFrames: thresholdFrames}
}

// Observe feeds the current progress metric in; returns true when the stall
// threshold has been exceeded and a rescue is needed.
func (n *SafetyNet) Observe(progress int) bool {
	if progress > n.lastProgress {
		n.lastProgress = progress
		n.stuckFrames = 0
		n.fired = false
		n.pulse = 0
		return false
	}
	if progress < n.lastProgress {
		// A level restart can reset progress; give it a fresh window.
		n.lastProgress = progress
		n.stuckFrames = 0
		n.fired = false
		n.pulse = 0
		return false
	}
	n.stuckFrames++
	return n.stuckFrames >= n.ThresholdFrames
}

// Rescue overrides a decision to break out of a stall: always move right and
// pulse A for 8 frames every 30 so that platforming gets un-stuck.
func (n *SafetyNet) Rescue(d games.Buttons) games.Buttons {
	n.fired = true
	n.stuckFrames = 0
	n.pulse++
	d.Right = true
	d.Left = false
	d.Start = false
	d.A = n.pulse%30 < 8
	return d
}

// Fired reports whether the net is currently in rescue.
func (n *SafetyNet) Fired() bool { return n.fired }

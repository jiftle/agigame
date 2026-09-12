package server

import (
	"sync"

	"agigame/emulator/core/gb"
)

// InputState aggregates button state from any number of producers (manual
// websocket clients, and from P2 onwards the agent layer) and presents the
// diff to the emulator once per frame.
type InputState struct {
	mu   sync.Mutex
	held map[gb.Button]bool // currently held buttons
	prev map[gb.Button]bool // button state reported last frame
}

// NewInputState returns an empty InputState.
func NewInputState() *InputState {
	return &InputState{
		held: make(map[gb.Button]bool),
		prev: make(map[gb.Button]bool),
	}
}

// Set marks a single gameplay button as pressed (true) or released (false).
func (in *InputState) Set(button gb.Button, down bool) {
	in.mu.Lock()
	defer in.mu.Unlock()
	if !button.IsGameBoyButton() {
		return
	}
	in.held[button] = down
}

// Press marks the given buttons as held.
func (in *InputState) Press(buttons ...gb.Button) {
	in.mu.Lock()
	defer in.mu.Unlock()
	for _, b := range buttons {
		in.held[b] = true
	}
}

// Release marks the given buttons as released.
func (in *InputState) Release(buttons ...gb.Button) {
	in.mu.Lock()
	defer in.mu.Unlock()
	for _, b := range buttons {
		in.held[b] = false
	}
}

// Clear releases every button.
func (in *InputState) Clear() {
	in.mu.Lock()
	defer in.mu.Unlock()
	for b := range in.held {
		in.held[b] = false
	}
}

// ReadButtons implements gb.InputProvider: it returns the set of buttons whose
// state changed since the last call.
func (in *InputState) ReadButtons() gb.ButtonInput {
	in.mu.Lock()
	defer in.mu.Unlock()

	var input gb.ButtonInput
	for b := gb.ButtonA; b <= gb.ButtonDown; b++ {
		if in.held[b] && !in.prev[b] {
			input.Pressed = append(input.Pressed, b)
		} else if !in.held[b] && in.prev[b] {
			input.Released = append(input.Released, b)
		}
	}
	for b, down := range in.held {
		in.prev[b] = down
	}
	return input
}

// Held returns the set of buttons currently held (for diagnostics).
func (in *InputState) Held() []string {
	in.mu.Lock()
	defer in.mu.Unlock()
	var out []string
	for name, b := range buttonNames {
		if in.held[b] {
			out = append(out, name)
		}
	}
	return out
}
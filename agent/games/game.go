// Package games defines the pluggable game interface consumed by the agent.
//
// A game plugin owns everything game-specific:
//   - the WRAM/OAM address table and state extraction
//   - the high-frequency rule decision
//   - the prompt rendering for the (low-frequency) LLM
//   - reward computation for training / statistics
//
// Adding a new game (Pokemon, Zelda, ...) means adding one more plugin whose
// ID maps to the "game" value in config.yaml.
package games

import "agigame/core/gb"

// GameReader gives a plugin everything it needs to read the current game
// state. It is implemented by *gb.Gameboy and only read on the emulator
// goroutine (synchronously), so no locking is required.
type GameReader interface {
	// ReadMemory reads one byte from the emulated address space.
	ReadMemory(address uint16) byte
	// Snapshot returns the emulator's general state snapshot.
	Snapshot() gb.State
}

// Buttons is the set of GameBoy buttons a decision wants to hold on the next
// emulated frame.
type Buttons struct {
	A      bool
	B      bool
	Start  bool
	Select bool
	Up     bool
	Down   bool
	Left   bool
	Right  bool
}

// ToSet converts the buttons into a map keyed by gb.Button held state.
func (b Buttons) ToSet() map[gb.Button]bool {
	set := map[gb.Button]bool{
		gb.ButtonA: b.A, gb.ButtonB: b.B,
		gb.ButtonStart: b.Start, gb.ButtonSelect: b.Select,
		gb.ButtonUp: b.Up, gb.ButtonDown: b.Down,
		gb.ButtonLeft: b.Left, gb.ButtonRight: b.Right,
	}
	return set
}

// RewardInfo is the reward delta computed between two consecutive states.
type RewardInfo struct {
	// Progress is the current absolute progress metric (e.g. camera X).
	Progress int
	// Delta is the reward delta for this transition.
	Delta float64
	// Events are human readable events triggered this transition.
	Events []string
	// GameOver is set when the run is finished (e.g. all lives lost).
	GameOver bool
	// Victory is set when the objective was reached.
	Victory bool
}

// GamePlugin is implemented once per supported game.
type GamePlugin interface {
	// ID is the config.yaml "game" value, e.g. "sml".
	ID() string
	// Name is the human readable game name.
	Name() string

	// ExtractState reads the game state for a single decision tick.
	// The returned value is opaque to the agent and passed back to
	// Decide/Prompt/Reward as-is.
	ExtractState(r GameReader) (any, error)

	// Decide returns the buttons the hard rules want held for a state.
	// Runs every frame (16ms) — must be fast and allocation-light.
	Decide(state any) Buttons

	// NeedsLLM reports whether the hard rules are insufficient for this
	// state and the LLM should be consulted instead (used in "hybrid" mode).
	// When false, the rule decision is used as-is.
	NeedsLLM(state any, decided Buttons) bool

	// Prompt renders the state for the LLM, describing the situation and
	// the expected output JSON schema.
	Prompt(state any) string

	// Reward compares the previous and current states and returns a reward
	// delta plus any notable events.
	Reward(prev, cur any) RewardInfo

	// ExtraState renders plugin-specific fields for the UI state message.
	ExtraState(cur any) map[string]any
}
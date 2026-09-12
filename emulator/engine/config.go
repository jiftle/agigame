package engine

import (
	"fmt"

	"agigame/emulator/agent"
	"agigame/emulator/core/gb"
)

// Defaults applied when the corresponding Config field is zero.
const (
	DefaultGame          = "sml"
	DefaultConsole       = "gb" // gb | gba
	DefaultFrameSkip     = 1
	DefaultStateInterval = 5
)

// Config configures a Session.
type Config struct {
	// ROM is the path to a legally owned GameBoy ROM dump.
	ROM string
	// Console selects the emulated handheld: "gb" (default) or "gba".
	Console string
	// Game is the game plugin id, e.g. "sml". Empty defaults to DefaultGame.
	// Ignored for consoles without an agent plugin (gba).
	Game string
	// Palette is one of greyscale | original | bgb. Empty defaults to greyscale.
	Palette string
	// CGB enables GameBoy Color mode when the ROM supports it.
	CGB bool
	// FrameSkip emits one frame every N emulated frames (0/1 = every frame).
	FrameSkip int
	// StateInterval emits a state update every N frames (default 5).
	StateInterval uint64
	// Agent holds the decision layer settings.
	Agent agent.Config

	// Logf receives structured log lines. Defaults to the standard logger.
	Logf func(level, format string, args ...any)
}

func (c *Config) normalize() {
	if c.Console == "" {
		c.Console = DefaultConsole
	}
	if c.Game == "" {
		c.Game = DefaultGame
	}
	if c.FrameSkip <= 0 {
		c.FrameSkip = DefaultFrameSkip
	}
	if c.StateInterval == 0 {
		c.StateInterval = DefaultStateInterval
	}
	if c.Logf == nil {
		c.Logf = defaultLogf
	}
}

// PaletteIndex maps a palette name to its DMG palette index.
func PaletteIndex(name string) (byte, error) {
	switch name {
	case "", "greyscale", "gray":
		return gb.PaletteGreyscale, nil
	case "original":
		return gb.PaletteOriginal, nil
	case "bgb":
		return gb.PaletteBGB, nil
	default:
		return 0, fmt.Errorf("unknown palette %q (want greyscale|original|bgb)", name)
	}
}

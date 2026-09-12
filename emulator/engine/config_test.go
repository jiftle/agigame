package engine

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"agigame/emulator/agent"
	"agigame/emulator/agent/games"
	"agigame/emulator/core/gb"
)

func TestConfigNormalize(t *testing.T) {
	c := Config{}
	c.normalize()
	if c.Game != DefaultGame || c.FrameSkip != DefaultFrameSkip ||
		c.StateInterval != DefaultStateInterval || c.Logf == nil {
		t.Fatalf("normalize = %+v", c)
	}
}

func TestPaletteIndex(t *testing.T) {
	cases := map[string]byte{
		"": gb.PaletteGreyscale, "greyscale": gb.PaletteGreyscale,
		"original": gb.PaletteOriginal, "bgb": gb.PaletteBGB,
	}
	for name, want := range cases {
		got, err := PaletteIndex(name)
		if err != nil || got != want {
			t.Fatalf("PaletteIndex(%q) = %d, %v; want %d", name, got, err, want)
		}
	}
	if _, err := PaletteIndex("rainbow"); err == nil {
		t.Fatalf("expected error for unknown palette")
	}
}

func TestGameRegistry(t *testing.T) {
	p, ok := games.Get("sml")
	if !ok || p == nil || p.ID() != "sml" {
		t.Fatalf("games.Get(sml) = %v, %v", p, ok)
	}
	if _, ok := games.Get("nope"); ok {
		t.Fatalf("expected miss for unknown game")
	}
}

// makeSmokeROM builds a minimal valid DMG ROM (NOP loop) in memory.
func makeSmokeROM() []byte {
	rom := make([]byte, 0x8000)
	copy(rom[0x134:0x13E], "SMOKE-TEST")
	rom[0x147] = 0x00
	rom[0x148] = 0x00
	rom[0x100] = 0x00
	rom[0x101] = 0x00
	rom[0x102] = 0x00
	rom[0x103] = 0x18
	rom[0x104] = 0xFD
	logo := []byte{
		0xCE, 0xED, 0x66, 0x66, 0xCC, 0x0D, 0x00, 0x0B, 0x03, 0x73, 0x00, 0x83,
		0x00, 0x0C, 0x00, 0x0D, 0x00, 0x08, 0x11, 0x1F, 0x88, 0x89, 0x00, 0x0E,
		0xDC, 0xCC, 0x6E, 0xE6, 0xDD, 0xDD, 0xD9, 0x99, 0xBB, 0xBB, 0x67, 0x63,
		0x6E, 0x0E, 0xEC, 0xCC, 0xDD, 0xDC, 0x99, 0x9F, 0xBB, 0xB9, 0x33, 0x3E,
	}
	copy(rom[0x104:0x134], logo)
	return rom
}

func TestSessionLifecycle(t *testing.T) {
	romPath := filepath.Join(t.TempDir(), "smoke.gb")
	if err := os.WriteFile(romPath, makeSmokeROM(), 0o644); err != nil {
		t.Fatal(err)
	}

	s, err := New(Config{
		ROM:           romPath,
		Game:          "sml",
		Palette:       "greyscale",
		FrameSkip:     1,
		StateInterval: 1,
		Agent:         agent.Config{Mode: agent.ModeManual},
	})
	if err != nil {
		t.Fatalf("engine.New: %v", err)
	}

	if s.CartName() == "" {
		t.Fatalf("expected a cart name")
	}
	if s.Mode() != agent.ModeManual || s.Auto() {
		t.Fatalf("unexpected initial mode/auto: %s/%v", s.Mode(), s.Auto())
	}

	frames := make(chan Frame, 4)
	states := make(chan StateUpdate, 4)
	s.SetFrameCallback(func(f Frame) {
		select {
		case frames <- f:
		default:
		}
	})
	s.SetStateCallback(func(u StateUpdate) {
		select {
		case states <- u:
		default:
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() { _ = s.Start(ctx) }()

	select {
	case f := <-frames:
		if len(f.PNG) == 0 {
			t.Fatalf("empty frame")
		}
	case <-time.After(3 * time.Second):
		t.Fatalf("no frame received")
	}
	select {
	case <-states:
	case <-time.After(3 * time.Second):
		t.Fatalf("no state received")
	}

	// Manual input then auto hand-off.
	s.Press(gb.ButtonRight)
	s.SetAuto(true)
	if !s.Auto() {
		t.Fatalf("expected auto on")
	}
	s.SetMode(agent.ModeRules)
	if s.Mode() != agent.ModeRules {
		t.Fatalf("expected rules mode")
	}
	s.SetAuto(false)
	if s.Auto() {
		t.Fatalf("expected auto off")
	}
	if err := s.SetPalette("bgb"); err != nil {
		t.Fatalf("SetPalette: %v", err)
	}
	if st := s.Status(); st["cart"] == nil || st["frames"] == nil {
		t.Fatalf("status missing fields: %+v", st)
	}
}

package gb

import (
	"os"
	"path/filepath"
	"testing"

	coregb "agigame/emulator/core/gb"
)

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

func TestHeadlessRun(t *testing.T) {
	rom := filepath.Join(t.TempDir(), "smoke.gb")
	if err := os.WriteFile(rom, makeSmokeROM(), 0o644); err != nil {
		t.Fatal(err)
	}

	c, err := New(rom)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if c.Name() != "gb" || c.Width() != Width || c.Height() != Height {
		t.Fatalf("bad console meta: %s %dx%d", c.Name(), c.Width(), c.Height())
	}

	for i := 0; i < 30; i++ {
		c.SetButton(coregb.ButtonRight, i%2 == 0)
		c.Step()
		c.DrainPCM()
	}
	if got := len(c.Pixels()); got != Width*Height*4 {
		t.Fatalf("pixels len = %d, want %d", got, Width*Height*4)
	}
	if c.FrameCount() == 0 {
		t.Fatalf("expected frames to advance")
	}
	if c.SampleRate() != SampleRate {
		t.Fatalf("sample rate = %d", c.SampleRate())
	}

	var audioBytes int
	for i := 0; i < 120; i++ {
		c.Step()
		audioBytes += len(c.DrainPCM())
	}
	if audioBytes == 0 {
		t.Fatalf("no PCM produced")
	}

	// Agent GameReader + debug snapshot must not panic.
	_ = c.ReadMemory(0xC000)
	if st := c.Snapshot(); st.Frame == 0 {
		t.Fatalf("snapshot frame = 0")
	}

	if err := c.SetPalette("bgb"); err != nil {
		t.Fatalf("SetPalette: %v", err)
	}
	if err := c.SetPalette("nope"); err == nil {
		t.Fatalf("expected error for unknown palette")
	}

	c.SetPaused(true)
	if !c.Paused() {
		t.Fatalf("expected paused")
	}
	c.SetPaused(false)
	c.ReleaseAll()
	c.Reset()
}

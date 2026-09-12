package gba

import (
	"os"
	"path/filepath"
	"testing"

	"agigame/emulator/core/gb"
)

// TestHeadlessRun boots the guac GBA core on a dummy ROM and steps it a few
// frames without a window, verifying the framebuffer contract and input API.
func TestHeadlessRun(t *testing.T) {
	rom := filepath.Join(t.TempDir(), "dummy.gba")
	if err := os.WriteFile(rom, make([]byte, 1<<20), 0o644); err != nil {
		t.Fatal(err)
	}

	c, err := New(rom)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	for i := 0; i < 30; i++ {
		c.SetButton(gb.ButtonRight, i%2 == 0)
		c.Step()
	}

	if got := len(c.Pixels()); got != Width*Height*4 {
		t.Fatalf("pixels len = %d, want %d", got, Width*Height*4)
	}
	if c.FrameCount() == 0 {
		t.Fatalf("expected frames to advance")
	}

	c.SetPaused(true)
	if !c.Paused() {
		t.Fatalf("expected paused")
	}
	c.SetPaused(false)
	c.ReleaseAll()
	c.Reset()
}

func TestReadTitle(t *testing.T) {
	rom := filepath.Join(t.TempDir(), "titled.gba")
	buf := make([]byte, 0xC0)
	copy(buf[0xA0:0xAC], []byte("TESTGAME"))
	if err := os.WriteFile(rom, buf, 0o644); err != nil {
		t.Fatal(err)
	}
	if got := readTitle(rom); got != "TESTGAME" {
		t.Fatalf("readTitle = %q", got)
	}
}

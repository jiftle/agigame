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

	// Drain audio every frame (required to avoid backpressure stalling) and
	// verify we get a plausible amount of PCM for the elapsed frames.
	var audioBytes int
	for i := 0; i < 120; i++ {
		c.Step()
		audioBytes += len(c.DrainPCM())
	}
	if audioBytes == 0 {
		t.Fatalf("no PCM produced")
	}
	fps := FPS
	expected := int(float64(120) / fps * SampleRate * 4)
	if audioBytes < expected*8/10 || audioBytes > expected*12/10 {
		t.Fatalf("PCM bytes = %d, expected ~%d", audioBytes, expected)
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

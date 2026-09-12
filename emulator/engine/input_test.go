package engine

import (
	"testing"

	"agigame/emulator/core/gb"
)

func TestInputStateDiff(t *testing.T) {
	in := NewInputState()
	in.Press(gb.ButtonA)
	if btns := in.ReadButtons(); len(btns.Pressed) != 1 || btns.Pressed[0] != gb.ButtonA {
		t.Errorf("expected A pressed once, got %+v", btns)
	}
	in.Set(gb.ButtonA, false)
	if btns := in.ReadButtons(); len(btns.Pressed) != 0 || len(btns.Released) != 1 {
		t.Errorf("expected A released once, got %+v", btns)
	}
}

func TestParseButton(t *testing.T) {
	if b, err := ParseButton("Right"); err != nil || b != gb.ButtonRight {
		t.Fatalf("ParseButton(Right) = %v, %v", b, err)
	}
	if _, err := ParseButton("Nope"); err == nil {
		t.Fatalf("expected error for unknown button")
	}
	if name := ButtonName(gb.ButtonStart); name != "Start" {
		t.Fatalf("ButtonName(Start) = %q", name)
	}
}

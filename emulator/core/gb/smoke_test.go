package gb

import (
	"testing"
)

// makeSmokeROM builds a minimal 32KB, no-MBC DMG ROM with a valid header. The
// entry point is an infinite NOP loop so the CPU executes real opcodes while the
// PPU keeps producing frames. Interrupts stay disabled so execution never jumps
// into header data.
func makeSmokeROM() []byte {
	rom := make([]byte, 0x8000)

	// Header title (0x134..0x143).
	copy(rom[0x134:0x13E], "SMOKE-TEST")

	// Cart type 0x00 = ROM only, no MBC. Size 0x00 = 32KB.
	rom[0x147] = 0x00
	rom[0x148] = 0x00

	// Entry point loop:
	//   0x100: NOP
	//   0x101: NOP
	//   0x102: NOP
	//   0x103: JR -3   (0x18 0xFD)  -> jumps back to 0x103
	rom[0x100] = 0x00
	rom[0x101] = 0x00
	rom[0x102] = 0x00
	rom[0x103] = 0x18
	rom[0x104] = 0xFD

	// Nintendo logo (0x0104..0x0133).
	logo := []byte{
		0xCE, 0xED, 0x66, 0x66, 0xCC, 0x0D, 0x00, 0x0B, 0x03, 0x73, 0x00, 0x83,
		0x00, 0x0C, 0x00, 0x0D, 0x00, 0x08, 0x11, 0x1F, 0x88, 0x89, 0x00, 0x0E,
		0xDC, 0xCC, 0x6E, 0xE6, 0xDD, 0xDD, 0xD9, 0x99, 0xBB, 0xBB, 0x67, 0x63,
		0x6E, 0x0E, 0xEC, 0xCC, 0xDD, 0xDC, 0x99, 0x9F, 0xBB, 0xB9, 0x33, 0x3E,
	}
	copy(rom[0x104:0x134], logo)

	// Header checksum (0x014D): such that 0x134..0x014C sum to zero.
	var sum int
	for i := 0x134; i <= 0x14C; i++ {
		sum = (sum - int(rom[i]) - 1) & 0xFF
	}
	rom[0x14D] = byte(sum)
	return rom
}

func TestSmokeRunFrames(t *testing.T) {
	gameboy := &Gameboy{}
	gameboy.setup()
	if _, err := gameboy.memory.LoadCartBytes(makeSmokeROM()); err != nil {
		t.Fatalf("load cart: %v", err)
	}

	frames := 0
	states := 0
	gameboy.SetFrameCallback(func(_ *[ScreenWidth][ScreenHeight][3]uint8, _ uint64) {
		frames++
	})
	gameboy.SetStateCallback(func(s State) {
		states++
		if s.CartName != "SMOKE-TEST" {
			t.Errorf("cart name = %q", s.CartName)
		}
	}, 5)

	for i := 0; i < 120; i++ {
		gameboy.Step()
	}

	if frames != 120 {
		t.Fatalf("expected 120 frames, got %d", frames)
	}
	if states != 24 {
		t.Errorf("expected 24 state snapshots (every 5 frames), got %d", states)
	}
	if gameboy.FrameNumber() != 120 {
		t.Errorf("frame number = %d", gameboy.FrameNumber())
	}
}

func TestSmokeInput(t *testing.T) {
	gameboy := &Gameboy{}
	gameboy.setup()
	if _, err := gameboy.memory.LoadCartBytes(makeSmokeROM()); err != nil {
		t.Fatalf("load cart: %v", err)
	}

	// Select the button group (bit 5 clear, bit 4 set) and read the joypad line.
	buttons := func() byte {
		return gameboy.joypadValue(0x10)
	}
	if v := buttons(); v != 0xDF {
		t.Fatalf("expected all buttons released (0xDF), got %#02x", v)
	}
	gameboy.pressButton(ButtonA)
	if v := buttons(); v != 0xDE {
		t.Errorf("pressing A did not clear the joypad bit: got %#02x", v)
	}
	gameboy.releaseButton(ButtonA)
	if v := buttons(); v != 0xDF {
		t.Errorf("releasing A did not restore the joypad: got %#02x", v)
	}
}

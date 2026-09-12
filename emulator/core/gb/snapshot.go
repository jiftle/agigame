package gb

// Registers is a snapshot of the GameBoy CPU registers.
type Registers struct {
	PC uint16 `json:"pc"`
	SP uint16 `json:"sp"`
	AF uint16 `json:"af"`
	BC uint16 `json:"bc"`
	DE uint16 `json:"de"`
	HL uint16 `json:"hl"`
}

// State is a snapshot of the emulator state. It is consumed by the web UI and
// by the agent layer for game state extraction.
type State struct {
	Frame    uint64    `json:"frame"`
	CartName string    `json:"cart"`
	CGB      bool      `json:"cgb"`
	Speed    int       `json:"speed"`
	Paused   bool      `json:"paused"`
	CPU      Registers `json:"cpu"`
	LY       byte      `json:"ly"`
	LCDC     byte      `json:"lcdc"`
	IF       byte      `json:"if"`
	IE       byte      `json:"ie"`
	Joypad   byte      `json:"joypad"`
	DIV      byte      `json:"div"`
	TIMA     byte      `json:"tima"`
	SCX      byte      `json:"scx"`
	SCY      byte      `json:"scy"`
}

// Snapshot returns the current emulator state.
func (gb *Gameboy) Snapshot() State {
	state := State{
		Frame:  gb.frameCount.Load(),
		CGB:    gb.IsCGB(),
		Speed:  gb.getSpeed(),
		Paused: gb.IsPaused(),
		CPU: Registers{
			PC: gb.cpu.PC,
			SP: gb.cpu.SP.HiLo(),
			AF: gb.cpu.AF.HiLo(),
			BC: gb.cpu.BC.HiLo(),
			DE: gb.cpu.DE.HiLo(),
			HL: gb.cpu.HL.HiLo(),
		},
		LY:     gb.memory.ReadHighRam(0xFF44),
		LCDC:   gb.memory.ReadHighRam(0xFF40),
		IF:     gb.memory.ReadHighRam(0xFF0F),
		IE:     gb.memory.ReadHighRam(0xFFFF),
		Joypad: gb.memory.HighRAM[0x00],
		DIV:    gb.memory.ReadHighRam(0xFF04),
		TIMA:   gb.memory.ReadHighRam(0xFF05),
		SCX:    gb.memory.ReadHighRam(0xFF43),
		SCY:    gb.memory.ReadHighRam(0xFF42),
	}
	if gb.IsCartLoaded() {
		state.CartName = gb.memory.Cart.GetName()
	}
	return state
}

// ReadMemory reads a byte from the emulated memory address space. Agent game
// plugins use this to extract game state from WRAM / HRAM addresses.
func (gb *Gameboy) ReadMemory(address uint16) byte {
	if gb.memory == nil {
		return 0xFF
	}
	return gb.memory.Read(address)
}

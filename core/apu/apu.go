package apu

import "log"

// APU is the GameBoy's audio processing unit. In the headless/server build the
// APU keeps the register and waveform RAM state (games poll these registers),
// but does not synthesise or play audio samples. A real audio backend can be
// added later behind the same public methods without touching core/gb.
type APU struct {
	playing bool

	memory      [52]byte
	waveformRam []byte
}

// Init the (headless) APU. The sound flag is accepted for API compatibility;
// headless builds never open an audio device.
func (a *APU) Init(sound bool) {
	a.playing = sound
	a.waveformRam = make([]byte, 0x20)

	// Sets waveform ram to:
	// 00 FF 00 FF  00 FF 00 FF  00 FF 00 FF  00 FF 00 FF
	for x := 0x0; x < 0x20; x++ {
		if x&2 == 0 {
			a.waveformRam[x] = 0x00
		} else {
			a.waveformRam[x] = 0xFF
		}
	}
}

// Buffer accepts CPU ticks for the audio pipeline. The headless APU consumes
// them without generating samples.
func (a *APU) Buffer(cpuTicks int, speed int) {}

// Read returns a value from the APU.
func (a *APU) Read(address uint16) byte {
	if address >= 0xFF30 {
		return a.waveformRam[address-0xFF30]
	}
	return a.memory[address-0xFF00] & soundMask[address-0xFF10]
}

// Write stores a value to the APU registers without side effects.
func (a *APU) Write(address uint16, value byte) {
	a.memory[address-0xFF00] = value
}

// WriteWaveform writes to channel 3 waveform RAM.
func (a *APU) WriteWaveform(address uint16, value byte) {
	a.waveformRam[address-0xFF30] = value
}

// ToggleSoundChannel is a no-op in the headless build.
func (a *APU) ToggleSoundChannel(channel int) {}

// LogSoundState logs the current sound register state.
func (a *APU) LogSoundState() {
	log.Printf("APU (headless): NR52=0x%02x NR50=0x%02x", a.memory[0x26], a.memory[0x24])
}

var soundMask = []byte{
	/* 0xFF10 */ 0xFF, 0xC0, 0xFF, 0x00, 0x40,
	/* 0xFF15 */ 0x00, 0xC0, 0xFF, 0x00, 0x40,
	/* 0xFF1A */ 0x80, 0x00, 0x60, 0x00, 0x40,
	/* 0xFF20 */ 0x00, 0x3F, 0xFF, 0xFF, 0x40,
	/* 0xFF24 */ 0xFF, 0xFF, 0x80,
}
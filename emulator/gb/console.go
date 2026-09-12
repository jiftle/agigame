// Package gb adapts the guac Game Boy / Game Boy Color core (aabalke/guac)
// into the agigame console layer. It runs headless (no window, no BIOS file,
// no audio device) and produces RGBA frames plus stereo s16le PCM.
package gb

import (
	"fmt"
	"image/color"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	guacconfig "github.com/aabalke/guac/config"
	guacgb "github.com/aabalke/guac/emu/gb"
	"github.com/aabalke/guac/utils"
	"github.com/hajimehoshi/ebiten/v2"

	coregb "agigame/emulator/core/gb"
)

const (
	// Width and Height are the DMG/CGB screen dimensions.
	Width  = 160
	Height = 144
	// FPS is the real DMG/CGB refresh rate.
	FPS = 59.727500569606
	// SampleRate is the PCM sample rate we capture (stereo s16le).
	SampleRate = 32768
	// audioBufferDur is the APU ring buffer size (also provides backpressure).
	audioBufferDur = 40 * time.Millisecond
)

var initOnce sync.Once

// ensureConfig initialises the guac global config once, without touching the
// filesystem (guac's file.Decode would write ./config.toml).
func ensureConfig() {
	initOnce.Do(func() {
		guacconfig.Conf.General.Headless = true
		guacconfig.Conf.General.Muted = true
		guacconfig.Conf.Gb.Bios.Direct = true // direct boot, no BIOS file needed
		guacconfig.Conf.Gb.Keyboard = guacconfig.EmulatorKeyboard{
			A:      []ebiten.Key{ebiten.KeyZ},
			B:      []ebiten.Key{ebiten.KeyX},
			Start:  []ebiten.Key{ebiten.KeyEnter},
			Select: []ebiten.Key{ebiten.KeyBackspace},
			Up:     []ebiten.Key{ebiten.KeyArrowUp},
			Down:   []ebiten.Key{ebiten.KeyArrowDown},
			Left:   []ebiten.Key{ebiten.KeyArrowLeft},
			Right:  []ebiten.Key{ebiten.KeyArrowRight},
		}
		setPaletteConfig("greyscale")
	})
}

// palettes maps our palette names to the 4 DMG shades (light -> dark).
var palettes = map[string][4]color.RGBA{
	"greyscale": {{0xFF, 0xFF, 0xFF, 0xFF}, {0xCC, 0xCC, 0xCC, 0xFF}, {0x77, 0x77, 0x77, 0xFF}, {0x00, 0x00, 0x00, 0xFF}},
	"original":  {{0x9B, 0xBC, 0x0F, 0xFF}, {0x8B, 0xAC, 0x0F, 0xFF}, {0x30, 0x62, 0x30, 0xFF}, {0x0F, 0x38, 0x0F, 0xFF}},
	"bgb":       {{0xE0, 0xF8, 0xD0, 0xFF}, {0x88, 0xC0, 0x70, 0xFF}, {0x34, 0x68, 0x56, 0xFF}, {0x08, 0x18, 0x20, 0xFF}},
}

func setPaletteConfig(name string) bool {
	rgba, ok := palettes[name]
	if !ok {
		return false
	}
	for i := 0; i < 4; i++ {
		guacconfig.Conf.Gb.Palette[i] = rgba[i]
	}
	return true
}

// configureAudio wires guac's APU to an in-memory stream and schedules the
// sample/frame-sequencer events so SoundClock produces PCM.
func configureAudio(g *guacgb.GameBoy) {
	g.Apu.Stream = utils.NewStream(audioBufferDur, SampleRate)
	g.CyclesPerSndGen = int64(guacgb.CPU_SPEED / SampleRate)
	g.Scheduler.Schedule(guacgb.EVENT_SND_FRAME_SEQ, 1, 0, g.ClockFrameSequencerEvent, nil)
	g.Scheduler.Schedule(guacgb.EVENT_SND_SAMPLE_GEN, 1, 0, g.AudioSampleEvent, nil)
}

// keyFor maps an agigame button to the ebiten key configured above.
func keyFor(b coregb.Button) (ebiten.Key, bool) {
	switch b {
	case coregb.ButtonA:
		return ebiten.KeyZ, true
	case coregb.ButtonB:
		return ebiten.KeyX, true
	case coregb.ButtonStart:
		return ebiten.KeyEnter, true
	case coregb.ButtonSelect:
		return ebiten.KeyBackspace, true
	case coregb.ButtonUp:
		return ebiten.KeyArrowUp, true
	case coregb.ButtonDown:
		return ebiten.KeyArrowDown, true
	case coregb.ButtonLeft:
		return ebiten.KeyArrowLeft, true
	case coregb.ButtonRight:
		return ebiten.KeyArrowRight, true
	default:
		return 0, false
	}
}

// Console is a running GB/GBC core.
type Console struct {
	mu       sync.Mutex
	g        *guacgb.GameBoy
	romPath  string
	title    string
	held     map[coregb.Button]bool
	frames   uint64
	audioAcc int64
	paused   bool
}

// New loads a GB/GBC ROM and boots the core.
func New(romPath string) (c *Console, err error) {
	ensureConfig()
	if _, statErr := os.Stat(romPath); statErr != nil {
		return nil, statErr
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("gb: boot failed: %v", r)
		}
	}()
	g := guacgb.NewGameBoy(nil, romPath, true)
	configureAudio(g)
	return &Console{
		g:       g,
		romPath: romPath,
		title:   readTitle(romPath),
		held:    make(map[coregb.Button]bool),
	}, nil
}

// Step advances exactly one frame, feeding the currently held buttons.
func (c *Console) Step() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.paused || c.g == nil {
		return
	}
	keys := make([]ebiten.Key, 0, len(c.held))
	for b, down := range c.held {
		if !down {
			continue
		}
		if k, ok := keyFor(b); ok {
			keys = append(keys, k)
		}
	}
	c.g.InputHandler(keys, nil)
	c.g.Update()
	c.frames++
}

// Pixels returns the framebuffer as RGBA bytes (Width*Height*4).
func (c *Console) Pixels() []byte { return c.g.Pixels }

// SampleRate returns the PCM sample rate of DrainPCM's output.
func (c *Console) SampleRate() int { return SampleRate }

// DrainPCM returns the stereo s16le PCM produced during the most recent Step.
// Must be called once per Step (the APU ring has backpressure).
func (c *Console) DrainPCM() []byte {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.g == nil || c.g.Apu == nil || c.g.Apu.Stream == nil {
		return nil
	}
	cyclesPerSample := c.g.CyclesPerSndGen
	if cyclesPerSample <= 0 {
		return nil
	}
	c.audioAcc += int64(guacgb.CYCLES_FRAME)
	frames := c.audioAcc / cyclesPerSample
	c.audioAcc -= frames * cyclesPerSample
	if frames <= 0 {
		return nil
	}
	buf := make([]byte, frames*4)
	_, _ = c.g.Apu.Stream.Read(buf)
	return buf
}

// SetButton sets a button's held state.
func (c *Console) SetButton(b coregb.Button, down bool) {
	c.mu.Lock()
	c.held[b] = down
	c.mu.Unlock()
}

// Set is an alias for SetButton (agent InputSink compatible).
func (c *Console) Set(b coregb.Button, down bool) { c.SetButton(b, down) }

// ReleaseAll releases every held button.
func (c *Console) ReleaseAll() {
	c.mu.Lock()
	c.held = make(map[coregb.Button]bool)
	c.mu.Unlock()
}

// Reset reboots the core with the same ROM.
func (c *Console) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	g, err := newCore(c.romPath)
	if err != nil {
		return
	}
	configureAudio(g)
	c.g = g
	c.frames = 0
	c.audioAcc = 0
}

// SetPaused pauses or resumes stepping.
func (c *Console) SetPaused(paused bool) {
	c.mu.Lock()
	c.paused = paused
	c.mu.Unlock()
}

// Paused reports the paused state.
func (c *Console) Paused() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.paused
}

// FrameCount returns the number of frames stepped.
func (c *Console) FrameCount() uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.frames
}

// CartName returns the ROM's internal title.
func (c *Console) CartName() string { return c.title }

// Name returns the console kind.
func (c *Console) Name() string { return "gb" }

// Width and Height report the screen size.
func (c *Console) Width() int  { return Width }
func (c *Console) Height() int { return Height }

// SetPalette switches the DMG palette by name (greyscale | original | bgb).
func (c *Console) SetPalette(name string) error {
	if !setPaletteConfig(name) {
		return fmt.Errorf("unknown palette %q", name)
	}
	return nil
}

// ReadMemory reads a byte from the emulated address space (agent GameReader).
func (c *Console) ReadMemory(address uint16) byte {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.g == nil {
		return 0xFF
	}
	return c.g.Read(address)
}

// Snapshot returns a debug snapshot for the UI.
func (c *Console) Snapshot() coregb.State {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.g == nil {
		return coregb.State{}
	}
	st := coregb.State{Frame: c.frames, Paused: c.paused, CartName: c.title}
	st.CPU.PC = c.g.Cpu.PC
	st.CPU.SP = c.g.Cpu.SP
	if c.g.Cpu.BC != nil {
		st.CPU.BC = *c.g.Cpu.BC
	}
	if c.g.Cpu.DE != nil {
		st.CPU.DE = *c.g.Cpu.DE
	}
	if c.g.Cpu.HL != nil {
		st.CPU.HL = *c.g.Cpu.HL
	}
	st.LY = c.g.Read(0xFF44)
	st.LCDC = c.g.Read(0xFF40)
	st.IF = c.g.Read(0xFF0F)
	st.IE = c.g.Read(0xFFFF)
	return st
}

// newCore boots a fresh core, recovering from guac panics.
func newCore(romPath string) (g *guacgb.GameBoy, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("gb: boot failed: %v", r)
		}
	}()
	return guacgb.NewGameBoy(nil, romPath, true), nil
}

// readTitle reads the DMG/CGB title at 0x134.
func readTitle(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()
	buf := make([]byte, 0x144)
	n, _ := io.ReadFull(f, buf)
	if n < 0x144 {
		return ""
	}
	s := string(buf[0x134:0x143]) // 0x143 is the CGB flag
	if i := strings.IndexByte(s, 0); i >= 0 {
		s = s[:i]
	}
	return strings.TrimSpace(s)
}

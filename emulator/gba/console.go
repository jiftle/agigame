// Package gba adapts the guac Game Boy Advance core (aabalke/guac) into the
// agigame console layer. It runs headless: the guac core needs no window, no
// BIOS file and no audio device (we pass a nil audio context and keep it muted).
package gba

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	guacconfig "github.com/aabalke/guac/config"
	guacgba "github.com/aabalke/guac/emu/gba"
	"github.com/aabalke/guac/utils"
	"github.com/hajimehoshi/ebiten/v2"

	"agigame/emulator/core/gb"
)

const (
	// Width and Height are the GBA screen dimensions.
	Width  = 240
	Height = 160
	// FPS is the GBA refresh rate.
	FPS = 59.727500569606
	// SampleRate is the PCM sample rate we capture from the GBA APU (stereo s16le).
	SampleRate = 32768
	// audioBufferDur is the GBA APU ring buffer size; it also provides
	// backpressure that keeps emulation synced to audio consumption.
	audioBufferDur = 40 * time.Millisecond
)

var initOnce sync.Once

// ensureConfig initialises the guac global config once, without touching the
// filesystem (guac's file.Decode would write ./config.toml).
func ensureConfig() {
	initOnce.Do(func() {
		guacconfig.Conf.General.Headless = true
		guacconfig.Conf.General.Muted = true
		guacconfig.Conf.Gba.Bios.Direct = true // direct boot, no BIOS file needed
		guacconfig.Conf.Gba.Keyboard = guacconfig.EmulatorKeyboard{
			A:      []ebiten.Key{ebiten.KeyZ},
			B:      []ebiten.Key{ebiten.KeyX},
			Start:  []ebiten.Key{ebiten.KeyEnter},
			Select: []ebiten.Key{ebiten.KeyBackspace},
			Up:     []ebiten.Key{ebiten.KeyArrowUp},
			Down:   []ebiten.Key{ebiten.KeyArrowDown},
			Left:   []ebiten.Key{ebiten.KeyArrowLeft},
			Right:  []ebiten.Key{ebiten.KeyArrowRight},
			L:      []ebiten.Key{ebiten.KeyQ},
			R:      []ebiten.Key{ebiten.KeyW},
		}
	})
}

// keyFor maps an agigame button to the ebiten key configured above.
func keyFor(b gb.Button) (ebiten.Key, bool) {
	switch b {
	case gb.ButtonA:
		return ebiten.KeyZ, true
	case gb.ButtonB:
		return ebiten.KeyX, true
	case gb.ButtonStart:
		return ebiten.KeyEnter, true
	case gb.ButtonSelect:
		return ebiten.KeyBackspace, true
	case gb.ButtonUp:
		return ebiten.KeyArrowUp, true
	case gb.ButtonDown:
		return ebiten.KeyArrowDown, true
	case gb.ButtonLeft:
		return ebiten.KeyArrowLeft, true
	case gb.ButtonRight:
		return ebiten.KeyArrowRight, true
	case gb.ButtonL:
		return ebiten.KeyQ, true
	case gb.ButtonR:
		return ebiten.KeyW, true
	default:
		return 0, false
	}
}

// configureAudio wires guac's APU to an in-memory stream (no audio device) and
// schedules the sample/frame-sequencer events so SoundClock produces PCM.
func configureAudio(g *guacgba.GBA) {
	g.Apu.Stream = utils.NewStream(audioBufferDur, SampleRate)
	g.CyclesPerSndGen = int64(guacgba.CPU_SPEED / SampleRate)
	g.Scheduler.Schedule(guacgba.EVENT_SND_FRAME_SEQ, 1, 0, g.ClockFrameSequencerEvent, nil)
	g.Scheduler.Schedule(guacgba.EVENT_SND_SAMPLE_GEN, 1, 0, g.AudioSampleEvent, nil)
}

// Console is a running GBA core.
type Console struct {
	mu       sync.Mutex
	g        *guacgba.GBA
	romPath  string
	title    string
	held     map[gb.Button]bool
	frames   uint64
	audioAcc int64
	paused   bool
}

// New loads a GBA ROM and boots the core. It returns an error (instead of
// panicking) if the ROM is invalid.
func New(romPath string) (c *Console, err error) {
	ensureConfig()
	if _, statErr := os.Stat(romPath); statErr != nil {
		return nil, statErr
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("gba: boot failed: %v", r)
		}
	}()
	g := guacgba.NewGBA(nil, romPath, true)
	configureAudio(g)
	return &Console{
		g:       g,
		romPath: romPath,
		title:   readTitle(romPath),
		held:    make(map[gb.Button]bool),
	}, nil
}

// Step advances exactly one frame, feeding the currently held buttons.
func (c *Console) Step() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.paused {
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
	c.g.InputHandler(nil, keys, nil, nil)
	c.g.Update()
	c.frames++
}

// Pixels returns the current framebuffer as RGBA bytes (Width*Height*4).
func (c *Console) Pixels() []byte { return c.g.Pixels }

// SampleRate returns the PCM sample rate of DrainPCM's output.
func (c *Console) SampleRate() int { return SampleRate }

// DrainPCM returns the stereo s16le PCM produced during the most recent Step.
// It must be called once per Step: the underlying ring has backpressure, so
// skipping it would eventually stall the emulator.
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
	c.audioAcc += int64(guacgba.CYCLES_FRAME)
	frames := c.audioAcc / cyclesPerSample
	c.audioAcc -= frames * cyclesPerSample
	if frames <= 0 {
		return nil
	}
	buf := make([]byte, frames*4) // stereo s16
	_, _ = c.g.Apu.Stream.Read(buf)
	return buf
}

// SetButton sets a button's held state.
func (c *Console) SetButton(b gb.Button, down bool) {
	c.mu.Lock()
	c.held[b] = down
	c.mu.Unlock()
}

// ReleaseAll releases every held button.
func (c *Console) ReleaseAll() {
	c.mu.Lock()
	c.held = make(map[gb.Button]bool)
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

// newCore boots a fresh core, recovering from guac panics.
func newCore(romPath string) (g *guacgba.GBA, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("gba: boot failed: %v", r)
		}
	}()
	return guacgba.NewGBA(nil, romPath, true), nil
}

// readTitle reads the 12-byte GBA ROM title at 0xA0.
func readTitle(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()
	buf := make([]byte, 0xC0)
	n, _ := io.ReadFull(f, buf)
	if n < 0xAC {
		return ""
	}
	s := string(buf[0xA0:0xAC])
	if i := strings.IndexByte(s, 0); i >= 0 {
		s = s[:i]
	}
	return strings.TrimSpace(s)
}

package engine

import (
	"bytes"
	"image"
	"image/png"

	"agigame/emulator/core/gb"
)

// Frame is a rendered emulator frame. PNG is always set (for the standalone
// web UI); RGBA/Width/Height are set for consoles that expose a raw
// framebuffer (GBA), letting hosts ship binary frames without PNG decoding.
type Frame struct {
	PNG    []byte
	Tick   uint64
	RGBA   []byte
	Width  int
	Height int
}

// Audio is a chunk of stereo s16le PCM produced by a console's APU.
type Audio struct {
	PCM        []byte
	SampleRate int
}

// StateUpdate is a periodic emulator state snapshot plus the agent summary
// (non-nil whenever auto mode is on).
type StateUpdate struct {
	State gb.State
	Auto  bool
	Agent map[string]any
	// Console/Width/Height describe the active handheld (gb or gba).
	Console string
	Width   int
	Height  int
}

// encodePNG converts an emulator frame buffer into a PNG-encoded image.
func encodePNG(frame *[gb.ScreenWidth][gb.ScreenHeight][3]uint8) []byte {
	img := image.NewRGBA(image.Rect(0, 0, gb.ScreenWidth, gb.ScreenHeight))
	for y := 0; y < gb.ScreenHeight; y++ {
		for x := 0; x < gb.ScreenWidth; x++ {
			i := img.PixOffset(x, y)
			img.Pix[i+0] = frame[x][y][0]
			img.Pix[i+1] = frame[x][y][1]
			img.Pix[i+2] = frame[x][y][2]
			img.Pix[i+3] = 0xFF
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil
	}
	return buf.Bytes()
}

// encodeRGBA converts an RGBA framebuffer (width*height*4 bytes, as produced by
// the GBA core) into a PNG-encoded image.
func encodeRGBA(pixels []byte, width, height int) []byte {
	if len(pixels) < width*height*4 {
		return nil
	}
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	copy(img.Pix, pixels[:width*height*4])
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil
	}
	return buf.Bytes()
}

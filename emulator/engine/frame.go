package engine

import (
	"bytes"
	"image"
	"image/png"

	"agigame/emulator/core/gb"
)

// Frame is a rendered, PNG-encoded emulator frame.
type Frame struct {
	PNG  []byte
	Tick uint64
}

// StateUpdate is a periodic emulator state snapshot plus the agent summary
// (non-nil whenever auto mode is on).
type StateUpdate struct {
	State gb.State
	Auto  bool
	Agent map[string]any
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

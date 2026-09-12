package engine

import (
	"agigame/emulator/core/gb"
)

// Frame is a rendered emulator frame as raw RGBA (width*height*4). Hosts ship
// it as a binary WebSocket message; no PNG/base64 encoding is done on the hot
// path.
type Frame struct {
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

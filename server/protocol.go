package server

import (
	"encoding/json"
	"fmt"

	"agigame/core/gb"
)

// WebSocket message protocol between server and the web UI / agent clients.

// Server → Client -------------------------------------------------------------

// FrameMsg carries a rendered frame as a base64-encoded PNG image.
type FrameMsg struct {
	Type string `json:"type"`
	Img  string `json:"img"` // base64 PNG
	Tick uint64 `json:"tick"`
}

// StateMsg carries a snapshot of the emulator state. In auto mode the Agent
// field carries the agent's live extract/reward summary.
type StateMsg struct {
	Type  string   `json:"type"`
	State gb.State `json:"state"`
	Agent any      `json:"agent,omitempty"`
}

// LogMsg is a server log entry shown in the UI log panel.
type LogMsg struct {
	Type  string `json:"type"`
	Level string `json:"level"`
	Msg   string `json:"msg"`
}

// HelloMsg is the first message sent after a client connects.
type HelloMsg struct {
	Type string `json:"type"`
	Cart string `json:"cart"`
	FPS  int    `json:"fps"`
	Game string `json:"game"`
}

// Client → Server -------------------------------------------------------------

// KeysMsg forwards button press/release events from the client.
type KeysMsg struct {
	Type     string   `json:"type"`
	Pressed  []string `json:"pressed"`
	Released []string `json:"released"`
}

// ControlMsg requests emulator control actions.
type ControlMsg struct {
	Type   string `json:"type"`
	Action string `json:"action"` // reset | pause | resume
}

// ConfigMsg selects the control mode and DMG palette.
type ConfigMsg struct {
	Type    string `json:"type"`
	Auto    bool   `json:"auto"`    // true = agent controls, false = manual
	Agent   string `json:"agent"`   // manual | rules | llm | hybrid
	Palette string `json:"palette"` // greyscale | original | bgb
}

// decodeMessage parses a raw client websocket message into a typed struct.
func decodeMessage(data []byte) (any, error) {
	var header struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &header); err != nil {
		return nil, fmt.Errorf("invalid message: %w", err)
	}
	switch header.Type {
	case "keys":
		var msg KeysMsg
		if err := json.Unmarshal(data, &msg); err != nil {
			return nil, err
		}
		return &msg, nil
	case "control":
		var msg ControlMsg
		if err := json.Unmarshal(data, &msg); err != nil {
			return nil, err
		}
		return &msg, nil
	case "config":
		var msg ConfigMsg
		if err := json.Unmarshal(data, &msg); err != nil {
			return nil, err
		}
		return &msg, nil
	default:
		return nil, fmt.Errorf("unknown message type %q", header.Type)
	}
}

// buttonNames maps the wire-format button names to gb buttons.
var buttonNames = map[string]gb.Button{
	"A":      gb.ButtonA,
	"B":      gb.ButtonB,
	"Start":  gb.ButtonStart,
	"Select": gb.ButtonSelect,
	"Up":     gb.ButtonUp,
	"Down":   gb.ButtonDown,
	"Left":   gb.ButtonLeft,
	"Right":  gb.ButtonRight,
}

// parseButton converts a wire button name into a gb button.
func parseButton(name string) (gb.Button, error) {
	if b, ok := buttonNames[name]; ok {
		return b, nil
	}
	return 0, fmt.Errorf("unknown button %q", name)
}

// paletteIndex maps a palette name to its gb palette index.
func paletteIndex(name string) (byte, error) {
	switch name {
	case "greyscale", "gray":
		return gb.PaletteGreyscale, nil
	case "original":
		return gb.PaletteOriginal, nil
	case "bgb":
		return gb.PaletteBGB, nil
	default:
		return 0, fmt.Errorf("unknown palette %q (want greyscale|original|bgb)", name)
	}
}

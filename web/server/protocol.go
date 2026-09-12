package server

import (
	"encoding/json"
	"fmt"

	"agigame/emulator/core/gb"
)

// WebSocket message protocol between server and the web UI / agent clients.

// Server → Client -------------------------------------------------------------

// StateMsg carries a snapshot of the emulator state. In auto mode the Agent
// field carries the agent's live extract/reward summary.
type StateMsg struct {
	Type    string   `json:"type"`
	State   gb.State `json:"state"`
	Agent   any      `json:"agent,omitempty"`
	Console string   `json:"console,omitempty"`
	Width   int      `json:"width,omitempty"`
	Height  int      `json:"height,omitempty"`
}

// LogMsg is a server log entry shown in the UI log panel.
type LogMsg struct {
	Type  string `json:"type"`
	Level string `json:"level"`
	Msg   string `json:"msg"`
}

// AudioMsg carries a chunk of stereo s16le PCM as base64.
type AudioMsg struct {
	Type string `json:"type"`
	PCM  string `json:"pcm"`
	Rate int    `json:"rate"`
}

// HelloMsg is the first message sent after a client connects.
type HelloMsg struct {
	Type    string `json:"type"`
	Cart    string `json:"cart"`
	FPS     int    `json:"fps"`
	Game    string `json:"game"`
	Console string `json:"console,omitempty"`
	Width   int    `json:"width,omitempty"`
	Height  int    `json:"height,omitempty"`
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

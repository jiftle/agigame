package server

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

// makeSmokeROM builds a minimal valid DMG ROM (NOP loop) in memory.
func makeSmokeROM() []byte {
	rom := make([]byte, 0x8000)
	copy(rom[0x134:0x13E], "SMOKE-TEST")
	rom[0x147] = 0x00 // ROM only, no MBC
	rom[0x148] = 0x00 // 32KB
	rom[0x100] = 0x00 // NOP
	rom[0x101] = 0x00
	rom[0x102] = 0x00
	rom[0x103] = 0x18 // JR
	rom[0x104] = 0xFD // JR -3
	logo := []byte{
		0xCE, 0xED, 0x66, 0x66, 0xCC, 0x0D, 0x00, 0x0B, 0x03, 0x73, 0x00, 0x83,
		0x00, 0x0C, 0x00, 0x0D, 0x00, 0x08, 0x11, 0x1F, 0x88, 0x89, 0x00, 0x0E,
		0xDC, 0xCC, 0x6E, 0xE6, 0xDD, 0xDD, 0xD9, 0x99, 0xBB, 0xBB, 0x67, 0x63,
		0x6E, 0x0E, 0xEC, 0xCC, 0xDD, 0xDC, 0x99, 0x9F, 0xBB, 0xB9, 0x33, 0x3E,
	}
	copy(rom[0x104:0x134], logo)
	return rom
}

func TestServerEndToEnd(t *testing.T) {
	romPath := filepath.Join(t.TempDir(), "smoke.gb")
	if err := os.WriteFile(romPath, makeSmokeROM(), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := DefaultConfig()
	cfg.ROM = romPath
	cfg.Emulator.FrameSkip = 4

	srv, err := New(cfg)
	if err != nil {
		t.Fatalf("new server: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() { _ = srv.Start(ctx) }()

	ts := newTestHTTPServer(srv)
	defer ts.Close()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http") + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial ws: %v", err)
	}
	defer conn.Close()

	// Must receive a hello message and soon after a frame.
	deadline := time.Now().Add(5 * time.Second)
	var sawHello, sawFrame, sawState bool
	for time.Now().Before(deadline) {
		_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		_, data, err := conn.ReadMessage()
		if err != nil {
			break
		}
		var msg struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(data, &msg); err != nil {
			t.Fatalf("bad msg: %v", err)
		}
		switch msg.Type {
		case "hello":
			sawHello = true
		case "frame":
			sawFrame = true
		case "state":
			sawState = true
		}
		if sawHello && sawFrame && sawState {
			break
		}
	}
	if !sawHello {
		t.Error("never received hello")
	}
	if !sawFrame {
		t.Error("never received a frame")
	}
	if !sawState {
		t.Error("never received a state snapshot")
	}
}

// newTestHTTPServer wraps the handler in an httptest server.
func newTestHTTPServer(srv *Server) *httptest.Server {
	return httptest.NewServer(srv.Handler())
}

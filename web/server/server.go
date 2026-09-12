package server

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"agigame/emulator/agent"
	"agigame/emulator/core/gb"
	"agigame/emulator/engine"
)

// Server is the standalone Web emulator: an HTTP/WebSocket adapter on top of
// an engine.Session. It owns the client hub; the emulator lifecycle, input
// aggregation and agent all live in the session.
type Server struct {
	cfg     *Config
	session *engine.Session
	hub     *Hub

	upgrader websocket.Upgrader
}

// New creates a Server and loads the configured ROM into a new session.
func New(cfg *Config) (*Server, error) {
	s := &Server{
		cfg: cfg,
		hub: NewHub(),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}

	session, err := engine.New(engine.Config{
		ROM:           cfg.ROM,
		Console:       cfg.Emulator.Console,
		Game:          cfg.Game,
		Palette:       cfg.Emulator.Palette,
		CGB:           cfg.Emulator.CGB,
		FrameSkip:     cfg.Emulator.FrameSkip,
		StateInterval: 5,
		Agent: agent.Config{
			Mode:            agent.Mode(cfg.Agent.Mode),
			LLMIntervalMs:   cfg.Agent.LLMIntervalMS,
			SafetyNetFrames: cfg.Agent.SafetyNetFrames,
		},
		Logf: func(level, format string, args ...any) {
			s.logf(level, format, args...)
		},
	})
	if err != nil {
		return nil, err
	}
	s.session = session
	return s, nil
}

// Gameboy returns the underlying emulator instance.
func (s *Server) Gameboy() *gb.Gameboy { return s.session.Gameboy() }

// Agent returns the decision agent.
func (s *Server) Agent() *agent.Agent { return s.session.Agent() }

// Session returns the underlying engine session.
func (s *Server) Session() *engine.Session { return s.session }

// Start wires the session callbacks and runs the emulation loop until ctx is
// cancelled. It blocks.
func (s *Server) Start(ctx context.Context) error {
	s.session.SetFrameCallback(s.onFrame)
	s.session.SetStateCallback(s.onState)
	s.logf("info", "loaded cart: %q", s.session.CartName())
	return s.session.Start(ctx)
}

// onFrame runs on the session encoder goroutine and broadcasts each frame.
func (s *Server) onFrame(f engine.Frame) {
	s.hub.BroadcastJSON(FrameMsg{
		Type: "frame",
		Img:  base64.StdEncoding.EncodeToString(f.PNG),
		Tick: f.Tick,
	})
}

// onState runs periodically on the emulator goroutine and broadcasts state.
func (s *Server) onState(u engine.StateUpdate) {
	s.hub.BroadcastJSON(StateMsg{
		Type:    "state",
		State:   u.State,
		Agent:   u.Agent,
		Console: u.Console,
		Width:   u.Width,
		Height:  u.Height,
	})
}

// handleWS upgrades and serves a websocket client connection.
func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket upgrade error: %v", err)
		return
	}

	client := s.hub.Add(conn)
	client.send <- mustJSON(HelloMsg{
		Type:    "hello",
		Cart:    s.session.CartName(),
		FPS:     60,
		Game:    s.cfg.Game,
		Console: s.session.Console(),
		Width:   s.session.Width(),
		Height:  s.session.Height(),
	})

	client.readPump(s.handleClientMsg)
}

// handleClientMsg routes decoded client messages to the session or config.
func (s *Server) handleClientMsg(data []byte) {
	msg, err := decodeMessage(data)
	if err != nil {
		s.logf("warn", "bad client message: %v", err)
		return
	}

	switch m := msg.(type) {
	case *KeysMsg:
		for _, name := range m.Pressed {
			if b, err := engine.ParseButton(name); err == nil {
				s.session.Press(b)
			}
		}
		for _, name := range m.Released {
			if b, err := engine.ParseButton(name); err == nil {
				s.session.Release(b)
			}
		}

	case *ControlMsg:
		switch m.Action {
		case "reset":
			s.logf("info", "reset requested")
			s.session.Reset()
		case "pause":
			s.session.SetPaused(true)
			s.logf("info", "paused")
		case "resume":
			s.session.SetPaused(false)
			s.logf("info", "resumed")
		}

	case *ConfigMsg:
		s.session.SetAuto(m.Auto)
		if m.Agent != "" {
			s.session.SetMode(agent.Mode(m.Agent))
		}
		if m.Palette != "" {
			if err := s.session.SetPalette(m.Palette); err != nil {
				s.logf("warn", "%v", err)
			} else {
				s.logf("info", "palette -> %s", m.Palette)
			}
		}
		s.logf("info", "config: auto=%v mode=%s", m.Auto, s.session.Mode())
	}
}

// Handler returns the root HTTP handler (websocket + static assets).
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", s.handleWS)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		out := s.session.Status()
		out["status"] = "ok"
		out["agent"] = string(s.session.Mode())
		_ = json.NewEncoder(w).Encode(out)
	})
	mux.Handle("/", http.FileServer(http.Dir(s.cfg.WebUI)))
	return logMiddleware(mux)
}

// logMiddleware sets permissive CORS headers.
func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// logf broadcasts a log line to websocket clients and the server log.
func (s *Server) logf(level, format string, args ...any) {
	msg := struct {
		Type  string `json:"type"`
		Level string `json:"level"`
		Msg   string `json:"msg"`
	}{Type: "log", Level: level, Msg: fmt.Sprintf(format, args...)}
	buf, _ := json.Marshal(msg)
	if s.hub != nil {
		s.hub.BroadcastRaw(buf)
	}
	log.Printf("[%s] %s", level, msg.Msg)
}

// mustJSON marshals v to JSON bytes.
func mustJSON(v any) []byte {
	buf, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return buf
}

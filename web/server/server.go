package server

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/gorilla/websocket"

	"agigame/emulator/agent"
	"agigame/emulator/agent/games"
	"agigame/emulator/core/gb"
)

// framePush carries a rendered frame added to the encode queue with its tick.
type framePush struct {
	data *[gb.ScreenWidth][gb.ScreenHeight][3]uint8
	tick uint64
}

// Server is the GoBoy-LLM HTTP/WebSocket service. It owns the GameBoy, the
// input aggregation, the client hub and the decision agent.
type Server struct {
	cfg *Config

	gb        *gb.Gameboy
	hub       *Hub
	input     *InputState
	agent     *agent.Agent
	frames    chan framePush
	frameSkip uint64

	auto atomic.Bool // auto/manual switch

	upgrader websocket.Upgrader
}

// New creates a Server, loading the configured ROM and building the agent for
// the configured game plugin.
func New(cfg *Config) (*Server, error) {
	var opts []gb.GameboyOption
	if cfg.Emulator.CGB {
		opts = append(opts, gb.WithCGBEnabled())
	}
	gameboy, err := gb.New(cfg.ROM, opts...)
	if err != nil {
		return nil, err
	}

	var plugin games.GamePlugin
	switch cfg.Game {
	case "sml", "":
		plugin = &games.SMLPlugin{}
	default:
		return nil, fmt.Errorf("unknown game plugin %q", cfg.Game)
	}

	s := &Server{
		cfg:       cfg,
		gb:        gameboy,
		hub:       NewHub(),
		input:     NewInputState(),
		frames:    make(chan framePush, 1),
		frameSkip: uint64(cfg.Emulator.FrameSkip),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
	if s.frameSkip == 0 {
		s.frameSkip = 1
	}
	if idx, err := paletteIndex(cfg.Emulator.Palette); err == nil {
		gb.SetDMGPalette(idx)
	} else {
		log.Printf("warn: %v (falling back to greyscale)", err)
		gb.SetDMGPalette(gb.PaletteGreyscale)
	}

	// P3 skeleton: always the stub provider; swap in a real LLM via
	// agent.Config.Provider when wiring up credentials.
	s.agent = agent.New(agent.Config{
		Mode:            agent.Mode(cfg.Agent.Mode),
		LLMIntervalMs:   cfg.Agent.LLMIntervalMS,
		SafetyNetFrames: cfg.Agent.SafetyNetFrames,
		Provider:        &agent.StubLLM{},
		Logf: func(format string, args ...any) {
			s.logf("agent", format, args...)
		},
	}, plugin, s.input)
	return s, nil
}

// Gameboy returns the underlying emulator instance.
func (s *Server) Gameboy() *gb.Gameboy {
	return s.gb
}

// Agent returns the decision agent (may be nil in future builds).
func (s *Server) Agent() *agent.Agent {
	return s.agent
}

// Start launches the emulator loop, the agent LLM loop and the frame encoder.
// It blocks until ctx is cancelled.
func (s *Server) Start(ctx context.Context) error {
	s.logf("info", "loaded cart: %q", s.gb.GetLoadedCart().GetName())
	s.gb.SetInputProvider(s.input)
	s.gb.SetPreFrameCallback(s.onPreFrame)
	s.gb.SetFrameCallback(s.onFrame)
	s.gb.SetStateCallback(s.onState, 5)

	go s.runEncoder(ctx)
	go s.agent.StartLLM(ctx)

	go func() {
		s.logf("info", "emulation started")
		s.gb.Run(ctx)
		s.logf("info", "emulation stopped")
	}()

	<-ctx.Done()
	return nil
}

// onPreFrame runs on the emulator goroutine before each frame advances. It
// fires the decision agent once per frame when auto mode is on.
func (s *Server) onPreFrame() {
	if !s.auto.Load() {
		return
	}
	if s.agent != nil {
		s.agent.Tick(s.gb)
	}
}

// onFrame runs on the emulator goroutine. Frames requested by the UI are copied
// and queued for the encoder (latest-wins, never blocks the emulator).
func (s *Server) onFrame(frame *[gb.ScreenWidth][gb.ScreenHeight][3]uint8, tick uint64) {
	if s.frameSkip == 0 || tick%s.frameSkip != 0 {
		return
	}
	copy := *frame
	push := framePush{data: &copy, tick: tick}
	select {
	case s.frames <- push:
	default:
		// Drop the queued frame in favour of the newest.
		select {
		case <-s.frames:
		default:
		}
		select {
		case s.frames <- push:
		default:
		}
	}
}

// onState runs periodically on the emulator goroutine. In auto mode it
// attaches the agent's latest extract / reward summary.
func (s *Server) onState(state gb.State) {
	msg := StateMsg{Type: "state", State: state}
	if s.auto.Load() && s.agent != nil {
		msg.Agent = s.agent.Status()
	}
	s.hub.BroadcastJSON(msg)
}

// runEncoder converts queued frames to PNG and broadcasts them.
func (s *Server) runEncoder(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case push := <-s.frames:
			img := encodePNG(push.data)
			s.hub.BroadcastJSON(FrameMsg{
				Type: "frame",
				Img:  base64.StdEncoding.EncodeToString(img),
				Tick: push.tick,
			})
		}
	}
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
		Type: "hello",
		Cart: s.gb.GetLoadedCart().GetName(),
		FPS:  60,
		Game: s.cfg.Game,
	})

	client.readPump(s.handleClientMsg)
}

// handleClientMsg routes decoded client messages to the emulator or config.
func (s *Server) handleClientMsg(data []byte) {
	msg, err := decodeMessage(data)
	if err != nil {
		s.logf("warn", "bad client message: %v", err)
		return
	}

	switch m := msg.(type) {
	case *KeysMsg:
		for _, name := range m.Pressed {
			if b, err := parseButton(name); err == nil {
				s.input.Press(b)
			}
		}
		for _, name := range m.Released {
			if b, err := parseButton(name); err == nil {
				s.input.Release(b)
			}
		}

	case *ControlMsg:
		switch m.Action {
		case "reset":
			s.logf("info", "reset requested")
			s.gb.RequestReset()
		case "pause":
			s.gb.SetPaused(true)
			s.logf("info", "paused")
		case "resume":
			s.gb.SetPaused(false)
			s.logf("info", "resumed")
		}

	case *ConfigMsg:
		s.auto.Store(m.Auto)
		if m.Agent != "" && s.agent != nil {
			s.agent.SetMode(agent.Mode(m.Agent))
		}
		if m.Palette != "" {
			if idx, err := paletteIndex(m.Palette); err == nil {
				gb.SetDMGPalette(idx)
				s.logf("info", "palette -> %s", m.Palette)
			} else {
				s.logf("warn", "%v", err)
			}
		}
		if !m.Auto {
			// Hand control back to the keyboard: release agent-held buttons.
			s.input.Clear()
			if s.agent != nil {
				s.agent.ReleaseAll()
			}
		}
		mode := "manual"
		if s.agent != nil {
			mode = string(s.agent.Mode())
		}
		s.logf("info", "config: auto=%v mode=%s", m.Auto, mode)
	}
}

// Handler returns the root HTTP handler (websocket + static assets).
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", s.handleWS)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		mode, stats := "manual", map[string]any{}
		if s.agent != nil {
			mode = string(s.agent.Mode())
			stats = s.agent.Status()
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status": "ok",
			"cart":   s.gb.GetLoadedCart().GetName(),
			"frames": s.gb.FrameNumber(),
			"auto":   s.auto.Load(),
			"agent":  mode,
			"stats":  stats,
		})
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

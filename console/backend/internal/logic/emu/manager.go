package emu

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/grand"

	"agigame/console/backend/internal/model"
	"agigame/console/backend/internal/service"
	"agigame/console/backend/utility/errcode"
	"agigame/emulator/agent"
	"agigame/emulator/engine"
)

// subBuffer is the per-subscriber outbound queue size; slow subscribers drop
// the oldest messages instead of blocking the emulator.
const subBuffer = 64

type sEmu struct {
	mu       sync.Mutex
	sessions map[string]*session
}

type session struct {
	id        string
	game      string
	rom       string
	createdAt time.Time
	engine    *engine.Session
	cancel    context.CancelFunc

	mu   sync.Mutex
	subs map[chan []byte]struct{}
}

func init() {
	service.RegisterEmu(New())
}

// New 创建模拟器会话管理器。
func New() *sEmu {
	return &sEmu{sessions: make(map[string]*session)}
}

func (m *sEmu) List(ctx context.Context) []*model.EmuSessionInfo {
	m.mu.Lock()
	list := make([]*session, 0, len(m.sessions))
	for _, s := range m.sessions {
		list = append(list, s)
	}
	m.mu.Unlock()

	sort.Slice(list, func(i, j int) bool { return list[i].createdAt.After(list[j].createdAt) })
	out := make([]*model.EmuSessionInfo, 0, len(list))
	for _, s := range list {
		out = append(out, s.info())
	}
	return out
}

func (m *sEmu) Start(ctx context.Context, in *model.EmuStartInput) (*model.EmuSessionInfo, error) {
	console := in.Console
	if console == "" {
		console = "gb"
	}
	romPath, err := resolveROM(ctx, console, in.Rom)
	if err != nil {
		return nil, err
	}

	game := in.Game
	if game == "" {
		game = engine.DefaultGame
	}
	palette := in.Palette
	if palette == "" {
		palette = g.Cfg().MustGet(ctx, "emulator.palette", "greyscale").String()
	}
	mode := agent.Mode(in.Mode)
	if mode == "" {
		mode = agent.ModeRules
	}
	frameSkip := g.Cfg().MustGet(ctx, "emulator.frameSkip", 4).Int()

	id := grand.S(8)
	sess := &session{
		id:        id,
		game:      game,
		rom:       romPath,
		createdAt: time.Now(),
		subs:      make(map[chan []byte]struct{}),
	}

	eng, err := engine.New(engine.Config{
		ROM:           romPath,
		Console:       console,
		Game:          game,
		Palette:       palette,
		FrameSkip:     frameSkip,
		StateInterval: 5,
		Agent: agent.Config{
			Mode:            mode,
			LLMIntervalMs:   500,
			SafetyNetFrames: 120,
		},
		Logf: sess.logf,
	})
	if err != nil {
		return nil, errcode.Business(fmt.Sprintf("启动模拟器失败: %v", err))
	}
	sess.engine = eng
	eng.SetFrameCallback(sess.onFrame)
	eng.SetStateCallback(sess.onState)

	runCtx, cancel := context.WithCancel(context.Background())
	sess.cancel = cancel
	go func() { _ = eng.Start(runCtx) }()

	m.mu.Lock()
	m.sessions[id] = sess
	m.mu.Unlock()

	g.Log().Infof(ctx, "emu session started id=%s console=%s rom=%s game=%s mode=%s", id, console, romPath, game, mode)
	return sess.info(), nil
}

func (m *sEmu) Get(ctx context.Context, id string) (*model.EmuSessionInfo, error) {
	s, err := m.get(id)
	if err != nil {
		return nil, err
	}
	return s.info(), nil
}

func (m *sEmu) Stop(ctx context.Context, id string) error {
	m.mu.Lock()
	s, ok := m.sessions[id]
	if ok {
		delete(m.sessions, id)
	}
	m.mu.Unlock()
	if !ok {
		return nil
	}
	s.cancel()
	s.closeSubs()
	g.Log().Infof(ctx, "emu session stopped id=%s", id)
	return nil
}

func (m *sEmu) Control(ctx context.Context, id, action string) error {
	s, err := m.get(id)
	if err != nil {
		return err
	}
	switch action {
	case "reset":
		s.engine.Reset()
	case "pause":
		s.engine.SetPaused(true)
	case "resume":
		s.engine.SetPaused(false)
	case "stop":
		return m.Stop(ctx, id)
	default:
		return errcode.BadRequest("未知操作: " + action)
	}
	return nil
}

func (m *sEmu) Config(ctx context.Context, id string, in *model.EmuConfigInput) (*model.EmuSessionInfo, error) {
	s, err := m.get(id)
	if err != nil {
		return nil, err
	}
	if in.Auto != nil {
		s.engine.SetAuto(*in.Auto)
	}
	if in.Mode != "" {
		s.engine.SetMode(agent.Mode(in.Mode))
	}
	if in.Palette != "" {
		if err := s.engine.SetPalette(in.Palette); err != nil {
			return nil, errcode.BadRequest(err.Error())
		}
	}
	return s.info(), nil
}

func (m *sEmu) Input(ctx context.Context, id string, pressed, released []string) error {
	s, err := m.get(id)
	if err != nil {
		return err
	}
	for _, name := range pressed {
		if b, err := engine.ParseButton(name); err == nil {
			s.engine.Press(b)
		}
	}
	for _, name := range released {
		if b, err := engine.ParseButton(name); err == nil {
			s.engine.Release(b)
		}
	}
	return nil
}

func (m *sEmu) Subscribe(ctx context.Context, id string) (*service.EmuSubscriber, error) {
	s, err := m.get(id)
	if err != nil {
		return nil, err
	}
	ch, closeFn := s.subscribe()
	return &service.EmuSubscriber{Messages: ch, Close: closeFn}, nil
}

func (m *sEmu) get(id string) (*session, error) {
	m.mu.Lock()
	s, ok := m.sessions[id]
	m.mu.Unlock()
	if !ok {
		return nil, errcode.NotFound("会话不存在或已结束: " + id)
	}
	return s, nil
}

// --- session ---------------------------------------------------------------

func (s *session) info() *model.EmuSessionInfo {
	return &model.EmuSessionInfo{
		Id:        s.id,
		Console:   s.engine.Console(),
		Cart:      s.engine.CartName(),
		Game:      s.game,
		Mode:      string(s.engine.Mode()),
		Auto:      s.engine.Auto(),
		Paused:    s.engine.IsPaused(),
		Frames:    s.engine.FrameNumber(),
		Width:     s.engine.Width(),
		Height:    s.engine.Height(),
		CreatedAt: s.createdAt.Format("2006-01-02 15:04:05"),
	}
}

func (s *session) onFrame(f engine.Frame) {
	if !s.hasSubs() {
		return
	}
	s.broadcastJSON(frameMsg{
		Type: "frame",
		Img:  base64.StdEncoding.EncodeToString(f.PNG),
		Tick: f.Tick,
	})
}

func (s *session) onState(u engine.StateUpdate) {
	if !s.hasSubs() {
		return
	}
	s.broadcastJSON(stateMsg{Type: "state", State: u.State, Agent: u.Agent})
}

func (s *session) logf(level, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	g.Log().Infof(context.Background(), "[emu %s] %s", s.id, msg)
	s.broadcastJSON(logMsg{Type: "log", Level: level, Msg: msg})
}

func (s *session) subscribe() (<-chan []byte, func()) {
	ch := make(chan []byte, subBuffer)
	s.mu.Lock()
	s.subs[ch] = struct{}{}
	s.mu.Unlock()

	var once sync.Once
	closeFn := func() {
		once.Do(func() {
			s.mu.Lock()
			if _, ok := s.subs[ch]; ok {
				delete(s.subs, ch)
				close(ch)
			}
			s.mu.Unlock()
		})
	}
	return ch, closeFn
}

func (s *session) closeSubs() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for ch := range s.subs {
		delete(s.subs, ch)
		close(ch)
	}
}

func (s *session) hasSubs() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.subs) > 0
}

func (s *session) broadcastJSON(v any) {
	data, err := json.Marshal(v)
	if err != nil {
		return
	}
	s.broadcast(data)
}

func (s *session) broadcast(data []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for ch := range s.subs {
		select {
		case ch <- data:
		default:
			// 慢订阅者：丢最旧一条，保证不阻塞仿真
			select {
			case <-ch:
			default:
			}
			select {
			case ch <- data:
			default:
			}
		}
	}
}

// resolveROM 把请求里的 ROM 名解析为绝对/相对路径，限制在 romDir 内。
func resolveROM(ctx context.Context, console, name string) (string, error) {
	if name == "" {
		key := "emulator.defaultRom"
		if console == "gba" {
			key = "emulator.defaultGbaRom"
		}
		name = g.Cfg().MustGet(ctx, key, "").String()
	}
	if name == "" {
		return "", errcode.BadRequest("未指定 ROM")
	}
	if filepath.IsAbs(name) {
		return name, nil
	}
	clean := filepath.Clean(name)
	if strings.HasPrefix(clean, "..") {
		return "", errcode.BadRequest("非法 ROM 路径")
	}
	dir := g.Cfg().MustGet(ctx, "emulator.romDir", "./roms").String()
	return filepath.Join(dir, clean), nil
}

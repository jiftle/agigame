package emu

import (
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gorilla/websocket"

	api "agigame/console/backend/api/v1/emu"
	"agigame/console/backend/internal/model"
	"agigame/console/backend/internal/service"
	"agigame/console/backend/utility/authz"
)

type cStream struct{}

// NewStream 创建会话实时流控制器。
func NewStream() *cStream {
	return &cStream{}
}

// streamClientMsg 客户端上行消息（keys / control / config）。
type streamClientMsg struct {
	Type     string   `json:"type"`
	Pressed  []string `json:"pressed"`
	Released []string `json:"released"`
	Action   string   `json:"action"`
	Auto     *bool    `json:"auto"`
	Agent    string   `json:"agent"`
	Palette  string   `json:"palette"`
}

// Stream 建立会话 WebSocket：先发 hello，再双向转发帧/状态与按键/控制。
func (c *cStream) Stream(ctx context.Context, req *api.StreamReq) (res *api.StreamRes, err error) {
	if err = authz.Check(ctx, permList); err != nil {
		return nil, err
	}

	info, err := service.Emu().Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	sub, err := service.Emu().Subscribe(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	defer sub.Close()

	ws, err := g.RequestFromCtx(ctx).WebSocket()
	if err != nil {
		return nil, err
	}

	hello, _ := json.Marshal(map[string]any{
		"type": "hello", "id": info.Id, "cart": info.Cart, "game": info.Game, "fps": 60,
	})
	if err = ws.WriteMessage(websocket.TextMessage, hello); err != nil {
		return nil, nil
	}

	// 下行：把订阅到的 JSON 消息写到 socket。
	go func() {
		for msg := range sub.Messages {
			if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		}
	}()

	// 上行：按键/控制/配置。
	for {
		_, data, readErr := ws.ReadMessage()
		if readErr != nil {
			break
		}
		var m streamClientMsg
		if json.Unmarshal(data, &m) != nil {
			continue
		}
		switch m.Type {
		case "keys":
			_ = service.Emu().Input(ctx, req.Id, m.Pressed, m.Released)
		case "control":
			_ = service.Emu().Control(ctx, req.Id, m.Action)
		case "config":
			_, _ = service.Emu().Config(ctx, req.Id, &model.EmuConfigInput{
				Auto:    m.Auto,
				Mode:    m.Agent,
				Palette: m.Palette,
			})
		}
	}
	return nil, nil
}

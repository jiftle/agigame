package emu

import (
	"github.com/gogf/gf/v2/frame/g"

	"agigame/console/backend/internal/model"
)

// ------------------------------ 会话 ------------------------------

type SessionListReq struct {
	g.Meta `path:"/emu/session/list" method:"get" tags:"模拟器" summary:"会话列表"`
}

type SessionListRes struct {
	List []*model.EmuSessionInfo `json:"list"`
}

type SessionStartReq struct {
	g.Meta  `path:"/emu/session/start" method:"post" tags:"模拟器" summary:"启动会话"`
	Rom     string `json:"rom"`
	Console string `json:"console"` // gb | gba
	Game    string `json:"game"`
	Mode    string `json:"mode"`
	Palette string `json:"palette"`
}

type SessionStartRes struct {
	Session *model.EmuSessionInfo `json:"session"`
}

type SessionGetReq struct {
	g.Meta `path:"/emu/session/{id}" method:"get" tags:"模拟器" summary:"会话详情"`
	Id     string `json:"id" in:"path" v:"required#会话ID不能为空"`
}

type SessionGetRes struct {
	Session *model.EmuSessionInfo `json:"session"`
}

type SessionStopReq struct {
	g.Meta `path:"/emu/session/{id}/stop" method:"post" tags:"模拟器" summary:"停止会话"`
	Id     string `json:"id" in:"path" v:"required#会话ID不能为空"`
}

type SessionStopRes struct{}

type SessionControlReq struct {
	g.Meta `path:"/emu/session/{id}/control" method:"post" tags:"模拟器" summary:"会话控制"`
	Id     string `json:"id" in:"path" v:"required#会话ID不能为空"`
	Action string `json:"action" v:"required#操作不能为空"` // reset | pause | resume
}

type SessionControlRes struct{}

type SessionConfigReq struct {
	g.Meta  `path:"/emu/session/{id}/config" method:"post" tags:"模拟器" summary:"会话配置"`
	Id      string `json:"id" in:"path" v:"required#会话ID不能为空"`
	Auto    *bool  `json:"auto"`
	Mode    string `json:"mode"`
	Palette string `json:"palette"`
}

type SessionConfigRes struct {
	Session *model.EmuSessionInfo `json:"session"`
}

// ------------------------------ WebSocket 流 ------------------------------

type StreamReq struct {
	g.Meta `path:"/emu/{id}/stream" method:"get" tags:"模拟器" summary:"会话实时流(WebSocket)"`
	Id     string `json:"id" in:"path" v:"required#会话ID不能为空"`
}

type StreamRes struct{}

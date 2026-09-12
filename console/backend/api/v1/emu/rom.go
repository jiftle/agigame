package emu

import (
	"github.com/gogf/gf/v2/frame/g"

	"agigame/console/backend/internal/model"
)

type RomListReq struct {
	g.Meta `path:"/emu/rom/list" method:"get" tags:"模拟器" summary:"可选 ROM 列表"`
}

type RomListRes struct {
	List []*model.EmuRomInfo `json:"list"`
}

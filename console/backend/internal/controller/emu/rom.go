package emu

import (
	"context"

	api "agigame/console/backend/api/v1/emu"
	"agigame/console/backend/internal/service"
	"agigame/console/backend/utility/authz"
)

type cRom struct{}

// NewRom 创建 ROM 控制器。
func NewRom() *cRom {
	return &cRom{}
}

// List 可选 ROM 列表
func (c *cRom) List(ctx context.Context, req *api.RomListReq) (res *api.RomListRes, err error) {
	if err = authz.Check(ctx, permList); err != nil {
		return nil, err
	}
	list, err := service.Emu().ListRoms(ctx)
	if err != nil {
		return nil, err
	}
	return &api.RomListRes{List: list}, nil
}

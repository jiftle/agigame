package emu

import (
	"context"

	api "agigame/console/backend/api/v1/emu"
	"agigame/console/backend/internal/model"
	"agigame/console/backend/internal/service"
	"agigame/console/backend/utility/authz"
)

// 权限码（需在菜单/角色里授予对应角色）
const (
	permList    = "emu:session:list"
	permControl = "emu:session:control"
)

type cSession struct{}

// NewSession 创建模拟器会话控制器。
func NewSession() *cSession {
	return &cSession{}
}

// List 会话列表
func (c *cSession) List(ctx context.Context, req *api.SessionListReq) (res *api.SessionListRes, err error) {
	if err = authz.Check(ctx, permList); err != nil {
		return nil, err
	}
	return &api.SessionListRes{List: service.Emu().List(ctx)}, nil
}

// Start 启动会话
func (c *cSession) Start(ctx context.Context, req *api.SessionStartReq) (res *api.SessionStartRes, err error) {
	if err = authz.Check(ctx, permControl); err != nil {
		return nil, err
	}
	info, err := service.Emu().Start(ctx, &model.EmuStartInput{
		Rom:     req.Rom,
		Console: req.Console,
		Game:    req.Game,
		Mode:    req.Mode,
		Palette: req.Palette,
	})
	if err != nil {
		return nil, err
	}
	return &api.SessionStartRes{Session: info}, nil
}

// Get 会话详情
func (c *cSession) Get(ctx context.Context, req *api.SessionGetReq) (res *api.SessionGetRes, err error) {
	if err = authz.Check(ctx, permList); err != nil {
		return nil, err
	}
	info, err := service.Emu().Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &api.SessionGetRes{Session: info}, nil
}

// Stop 停止会话
func (c *cSession) Stop(ctx context.Context, req *api.SessionStopReq) (res *api.SessionStopRes, err error) {
	if err = authz.Check(ctx, permControl); err != nil {
		return nil, err
	}
	if err = service.Emu().Stop(ctx, req.Id); err != nil {
		return nil, err
	}
	return &api.SessionStopRes{}, nil
}

// Control 会话控制（reset | pause | resume）
func (c *cSession) Control(ctx context.Context, req *api.SessionControlReq) (res *api.SessionControlRes, err error) {
	if err = authz.Check(ctx, permControl); err != nil {
		return nil, err
	}
	if err = service.Emu().Control(ctx, req.Id, req.Action); err != nil {
		return nil, err
	}
	return &api.SessionControlRes{}, nil
}

// Config 会话配置（auto/mode/palette）
func (c *cSession) Config(ctx context.Context, req *api.SessionConfigReq) (res *api.SessionConfigRes, err error) {
	if err = authz.Check(ctx, permControl); err != nil {
		return nil, err
	}
	info, err := service.Emu().Config(ctx, req.Id, &model.EmuConfigInput{
		Auto:    req.Auto,
		Mode:    req.Mode,
		Palette: req.Palette,
	})
	if err != nil {
		return nil, err
	}
	return &api.SessionConfigRes{Session: info}, nil
}

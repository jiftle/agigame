package system

import (
	"context"

	api "agigame/console/backend/api/v1/system"
	"agigame/console/backend/internal/model"
	"agigame/console/backend/internal/service"
	"agigame/console/backend/utility/authz"
)

type cLog struct{}

// NewLog 创建日志控制器
func NewLog() *cLog {
	return &cLog{}
}

// LoginList 登录日志列表
func (c *cLog) LoginList(ctx context.Context, req *api.LoginLogListReq) (res *api.LoginLogListRes, err error) {
	if err = authz.Check(ctx, "system:loginlog:list"); err != nil {
		return nil, err
	}
	total, list, err := service.Log().LoginLogList(ctx, &model.LoginLogQueryInput{
		PageInput: model.PageInput{PageNum: req.PageNum, PageSize: req.PageSize},
		Username:  req.Username,
		Status:    req.Status,
		BeginTime: req.BeginTime,
		EndTime:   req.EndTime,
	})
	if err != nil {
		return nil, err
	}
	return &api.LoginLogListRes{Total: total, List: list}, nil
}

// LoginDelete 删除登录日志
func (c *cLog) LoginDelete(ctx context.Context, req *api.LoginLogDeleteReq) (res *api.LoginLogDeleteRes, err error) {
	if err = authz.Check(ctx, "system:loginlog:remove"); err != nil {
		return nil, err
	}
	if err = service.Log().LoginLogDelete(ctx, req.Ids); err != nil {
		return nil, err
	}
	return &api.LoginLogDeleteRes{}, nil
}

// LoginClear 清空登录日志
func (c *cLog) LoginClear(ctx context.Context, req *api.LoginLogClearReq) (res *api.LoginLogClearRes, err error) {
	if err = authz.Check(ctx, "system:loginlog:remove"); err != nil {
		return nil, err
	}
	if err = service.Log().LoginLogClear(ctx); err != nil {
		return nil, err
	}
	return &api.LoginLogClearRes{}, nil
}

// OperList 操作日志列表
func (c *cLog) OperList(ctx context.Context, req *api.OperLogListReq) (res *api.OperLogListRes, err error) {
	if err = authz.Check(ctx, "system:operlog:list"); err != nil {
		return nil, err
	}
	total, list, err := service.Log().OperLogList(ctx, &model.OperLogQueryInput{
		PageInput: model.PageInput{PageNum: req.PageNum, PageSize: req.PageSize},
		Title:     req.Title,
		OperName:  req.OperName,
		Status:    req.Status,
		BeginTime: req.BeginTime,
		EndTime:   req.EndTime,
	})
	if err != nil {
		return nil, err
	}
	return &api.OperLogListRes{Total: total, List: list}, nil
}

// OperDelete 删除操作日志
func (c *cLog) OperDelete(ctx context.Context, req *api.OperLogDeleteReq) (res *api.OperLogDeleteRes, err error) {
	if err = authz.Check(ctx, "system:operlog:remove"); err != nil {
		return nil, err
	}
	if err = service.Log().OperLogDelete(ctx, req.Ids); err != nil {
		return nil, err
	}
	return &api.OperLogDeleteRes{}, nil
}

// OperClear 清空操作日志
func (c *cLog) OperClear(ctx context.Context, req *api.OperLogClearReq) (res *api.OperLogClearRes, err error) {
	if err = authz.Check(ctx, "system:operlog:remove"); err != nil {
		return nil, err
	}
	if err = service.Log().OperLogClear(ctx); err != nil {
		return nil, err
	}
	return &api.OperLogClearRes{}, nil
}

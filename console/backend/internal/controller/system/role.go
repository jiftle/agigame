package system

import (
	"context"

	api "agigame/console/backend/api/v1/system"
	"agigame/console/backend/internal/model"
	"agigame/console/backend/internal/service"
	"agigame/console/backend/utility/authz"
)

type cRole struct{}

// NewRole 创建角色控制器
func NewRole() *cRole {
	return &cRole{}
}

// List 角色列表
func (c *cRole) List(ctx context.Context, req *api.RoleListReq) (res *api.RoleListRes, err error) {
	if err = authz.Check(ctx, "system:role:list"); err != nil {
		return nil, err
	}
	total, list, err := service.Role().List(ctx, &model.RoleQueryInput{
		PageInput: model.PageInput{PageNum: req.PageNum, PageSize: req.PageSize},
		Name:      req.Name,
		Code:      req.Code,
		Status:    req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &api.RoleListRes{Total: total, List: list}, nil
}

// Get 角色详情
func (c *cRole) Get(ctx context.Context, req *api.RoleGetReq) (res *api.RoleGetRes, err error) {
	if err = authz.Check(ctx, "system:role:list"); err != nil {
		return nil, err
	}
	role, err := service.Role().GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	menuIds, err := service.Role().GetMenuIds(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &api.RoleGetRes{Role: role, MenuIds: menuIds}, nil
}

// Create 新增角色
func (c *cRole) Create(ctx context.Context, req *api.RoleCreateReq) (res *api.RoleCreateRes, err error) {
	if err = authz.Check(ctx, "system:role:add"); err != nil {
		return nil, err
	}
	err = service.Role().Create(ctx, &model.RoleSaveInput{
		Name:      req.Name,
		Code:      req.Code,
		Sort:      req.Sort,
		DataScope: req.DataScope,
		Status:    req.Status,
		Remark:    req.Remark,
		MenuIds:   req.MenuIds,
	})
	if err != nil {
		return nil, err
	}
	return &api.RoleCreateRes{}, nil
}

// Update 修改角色
func (c *cRole) Update(ctx context.Context, req *api.RoleUpdateReq) (res *api.RoleUpdateRes, err error) {
	if err = authz.Check(ctx, "system:role:edit"); err != nil {
		return nil, err
	}
	err = service.Role().Update(ctx, &model.RoleSaveInput{
		Id:        req.Id,
		Name:      req.Name,
		Code:      req.Code,
		Sort:      req.Sort,
		DataScope: req.DataScope,
		Status:    req.Status,
		Remark:    req.Remark,
		MenuIds:   req.MenuIds,
	})
	if err != nil {
		return nil, err
	}
	return &api.RoleUpdateRes{}, nil
}

// Delete 删除角色
func (c *cRole) Delete(ctx context.Context, req *api.RoleDeleteReq) (res *api.RoleDeleteRes, err error) {
	if err = authz.Check(ctx, "system:role:remove"); err != nil {
		return nil, err
	}
	if err = service.Role().Delete(ctx, req.Ids); err != nil {
		return nil, err
	}
	return &api.RoleDeleteRes{}, nil
}

package system

import (
	"context"

	api "agigame/console/backend/api/v1/system"
	"agigame/console/backend/internal/model"
	"agigame/console/backend/internal/service"
	"agigame/console/backend/utility/authz"
)

type cMenu struct{}

// NewMenu 创建菜单控制器
func NewMenu() *cMenu {
	return &cMenu{}
}

// List 菜单列表
func (c *cMenu) List(ctx context.Context, req *api.MenuListReq) (res *api.MenuListRes, err error) {
	if err = authz.CheckAny(ctx, "system:menu:list", "system:role:list"); err != nil {
		return nil, err
	}
	list, err := service.Menu().GetTree(ctx, &model.MenuQueryInput{Title: req.Title, Status: req.Status})
	if err != nil {
		return nil, err
	}
	return &api.MenuListRes{List: list}, nil
}

// Get 菜单详情
func (c *cMenu) Get(ctx context.Context, req *api.MenuGetReq) (res *api.MenuGetRes, err error) {
	if err = authz.Check(ctx, "system:menu:list"); err != nil {
		return nil, err
	}
	menu, err := service.Menu().GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &api.MenuGetRes{Menu: menu}, nil
}

// Create 新增菜单
func (c *cMenu) Create(ctx context.Context, req *api.MenuCreateReq) (res *api.MenuCreateRes, err error) {
	if err = authz.Check(ctx, "system:menu:add"); err != nil {
		return nil, err
	}
	err = service.Menu().Create(ctx, &model.MenuSaveInput{
		ParentId:  req.ParentId,
		Title:     req.Title,
		Name:      req.Name,
		Path:      req.Path,
		Component: req.Component,
		Icon:      req.Icon,
		Type:      req.Type,
		Perms:     req.Perms,
		Sort:      req.Sort,
		Visible:   req.Visible,
		Status:    req.Status,
		Redirect:  req.Redirect,
		IsFrame:   req.IsFrame,
		IsCache:   req.IsCache,
	})
	if err != nil {
		return nil, err
	}
	return &api.MenuCreateRes{}, nil
}

// Update 修改菜单
func (c *cMenu) Update(ctx context.Context, req *api.MenuUpdateReq) (res *api.MenuUpdateRes, err error) {
	if err = authz.Check(ctx, "system:menu:edit"); err != nil {
		return nil, err
	}
	err = service.Menu().Update(ctx, &model.MenuSaveInput{
		Id:        req.Id,
		ParentId:  req.ParentId,
		Title:     req.Title,
		Name:      req.Name,
		Path:      req.Path,
		Component: req.Component,
		Icon:      req.Icon,
		Type:      req.Type,
		Perms:     req.Perms,
		Sort:      req.Sort,
		Visible:   req.Visible,
		Status:    req.Status,
		Redirect:  req.Redirect,
		IsFrame:   req.IsFrame,
		IsCache:   req.IsCache,
	})
	if err != nil {
		return nil, err
	}
	return &api.MenuUpdateRes{}, nil
}

// Delete 删除菜单
func (c *cMenu) Delete(ctx context.Context, req *api.MenuDeleteReq) (res *api.MenuDeleteRes, err error) {
	if err = authz.Check(ctx, "system:menu:remove"); err != nil {
		return nil, err
	}
	if err = service.Menu().Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &api.MenuDeleteRes{}, nil
}

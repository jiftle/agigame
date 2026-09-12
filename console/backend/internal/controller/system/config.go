package system

import (
	"context"

	api "agigame/console/backend/api/v1/system"
	"agigame/console/backend/internal/model"
	"agigame/console/backend/internal/service"
	"agigame/console/backend/utility/authz"
)

type cConfig struct{}

// NewConfig 创建参数控制器
func NewConfig() *cConfig {
	return &cConfig{}
}

// List 参数列表
func (c *cConfig) List(ctx context.Context, req *api.ConfigListReq) (res *api.ConfigListRes, err error) {
	if err = authz.Check(ctx, "system:config:list"); err != nil {
		return nil, err
	}
	total, list, err := service.Config().List(ctx, &model.ConfigQueryInput{
		PageInput:  model.PageInput{PageNum: req.PageNum, PageSize: req.PageSize},
		ConfigName: req.ConfigName,
		ConfigKey:  req.ConfigKey,
	})
	if err != nil {
		return nil, err
	}
	return &api.ConfigListRes{Total: total, List: list}, nil
}

// Get 参数详情
func (c *cConfig) Get(ctx context.Context, req *api.ConfigGetReq) (res *api.ConfigGetRes, err error) {
	if err = authz.Check(ctx, "system:config:list"); err != nil {
		return nil, err
	}
	item, err := service.Config().GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &api.ConfigGetRes{Config: item}, nil
}

// Create 新增参数
func (c *cConfig) Create(ctx context.Context, req *api.ConfigCreateReq) (res *api.ConfigCreateRes, err error) {
	if err = authz.Check(ctx, "system:config:add"); err != nil {
		return nil, err
	}
	err = service.Config().Create(ctx, &model.ConfigSaveInput{
		ConfigName: req.ConfigName, ConfigKey: req.ConfigKey, ConfigValue: req.ConfigValue,
		ConfigType: req.ConfigType, Remark: req.Remark,
	})
	if err != nil {
		return nil, err
	}
	return &api.ConfigCreateRes{}, nil
}

// Update 修改参数
func (c *cConfig) Update(ctx context.Context, req *api.ConfigUpdateReq) (res *api.ConfigUpdateRes, err error) {
	if err = authz.Check(ctx, "system:config:edit"); err != nil {
		return nil, err
	}
	err = service.Config().Update(ctx, &model.ConfigSaveInput{
		Id: req.Id, ConfigName: req.ConfigName, ConfigKey: req.ConfigKey, ConfigValue: req.ConfigValue,
		ConfigType: req.ConfigType, Remark: req.Remark,
	})
	if err != nil {
		return nil, err
	}
	return &api.ConfigUpdateRes{}, nil
}

// Delete 删除参数
func (c *cConfig) Delete(ctx context.Context, req *api.ConfigDeleteReq) (res *api.ConfigDeleteRes, err error) {
	if err = authz.Check(ctx, "system:config:remove"); err != nil {
		return nil, err
	}
	if err = service.Config().Delete(ctx, req.Ids); err != nil {
		return nil, err
	}
	return &api.ConfigDeleteRes{}, nil
}

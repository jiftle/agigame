package system

import (
	"context"

	api "agigame/console/backend/api/v1/system"
	"agigame/console/backend/internal/model"
	"agigame/console/backend/internal/service"
	"agigame/console/backend/utility/authz"
)

type cDict struct{}

// NewDict 创建字典控制器
func NewDict() *cDict {
	return &cDict{}
}

// TypeList 字典类型列表
func (c *cDict) TypeList(ctx context.Context, req *api.DictTypeListReq) (res *api.DictTypeListRes, err error) {
	if err = authz.Check(ctx, "system:dict:list"); err != nil {
		return nil, err
	}
	total, list, err := service.Dict().TypeList(ctx, &model.DictTypeQueryInput{
		PageInput: model.PageInput{PageNum: req.PageNum, PageSize: req.PageSize},
		Name:      req.Name,
		Type:      req.Type,
		Status:    req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &api.DictTypeListRes{Total: total, List: list}, nil
}

// TypeGet 字典类型详情
func (c *cDict) TypeGet(ctx context.Context, req *api.DictTypeGetReq) (res *api.DictTypeGetRes, err error) {
	if err = authz.Check(ctx, "system:dict:list"); err != nil {
		return nil, err
	}
	item, err := service.Dict().TypeGetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &api.DictTypeGetRes{DictType: item}, nil
}

// TypeCreate 新增字典类型
func (c *cDict) TypeCreate(ctx context.Context, req *api.DictTypeCreateReq) (res *api.DictTypeCreateRes, err error) {
	if err = authz.Check(ctx, "system:dict:add"); err != nil {
		return nil, err
	}
	err = service.Dict().TypeCreate(ctx, &model.DictTypeSaveInput{
		Name: req.Name, Type: req.Type, Status: req.Status, Remark: req.Remark,
	})
	if err != nil {
		return nil, err
	}
	return &api.DictTypeCreateRes{}, nil
}

// TypeUpdate 修改字典类型
func (c *cDict) TypeUpdate(ctx context.Context, req *api.DictTypeUpdateReq) (res *api.DictTypeUpdateRes, err error) {
	if err = authz.Check(ctx, "system:dict:edit"); err != nil {
		return nil, err
	}
	err = service.Dict().TypeUpdate(ctx, &model.DictTypeSaveInput{
		Id: req.Id, Name: req.Name, Type: req.Type, Status: req.Status, Remark: req.Remark,
	})
	if err != nil {
		return nil, err
	}
	return &api.DictTypeUpdateRes{}, nil
}

// TypeDelete 删除字典类型
func (c *cDict) TypeDelete(ctx context.Context, req *api.DictTypeDeleteReq) (res *api.DictTypeDeleteRes, err error) {
	if err = authz.Check(ctx, "system:dict:remove"); err != nil {
		return nil, err
	}
	if err = service.Dict().TypeDelete(ctx, req.Ids); err != nil {
		return nil, err
	}
	return &api.DictTypeDeleteRes{}, nil
}

// DataList 字典数据列表
func (c *cDict) DataList(ctx context.Context, req *api.DictDataListReq) (res *api.DictDataListRes, err error) {
	if err = authz.Check(ctx, "system:dict:list"); err != nil {
		return nil, err
	}
	total, list, err := service.Dict().DataList(ctx, &model.DictDataQueryInput{
		PageInput: model.PageInput{PageNum: req.PageNum, PageSize: req.PageSize},
		DictType:  req.DictType,
		DictLabel: req.DictLabel,
		Status:    req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &api.DictDataListRes{Total: total, List: list}, nil
}

// DataByType 按类型获取字典数据
func (c *cDict) DataByType(ctx context.Context, req *api.DictDataByTypeReq) (res *api.DictDataByTypeRes, err error) {
	if err = authz.Check(ctx, "system:dict:list"); err != nil {
		return nil, err
	}
	list, err := service.Dict().DataByType(ctx, req.DictType)
	if err != nil {
		return nil, err
	}
	return &api.DictDataByTypeRes{List: list}, nil
}

// DataGet 字典数据详情
func (c *cDict) DataGet(ctx context.Context, req *api.DictDataGetReq) (res *api.DictDataGetRes, err error) {
	if err = authz.Check(ctx, "system:dict:list"); err != nil {
		return nil, err
	}
	item, err := service.Dict().DataGetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &api.DictDataGetRes{DictData: item}, nil
}

// DataCreate 新增字典数据
func (c *cDict) DataCreate(ctx context.Context, req *api.DictDataCreateReq) (res *api.DictDataCreateRes, err error) {
	if err = authz.Check(ctx, "system:dict:add"); err != nil {
		return nil, err
	}
	err = service.Dict().DataCreate(ctx, &model.DictDataSaveInput{
		DictSort: req.DictSort, DictLabel: req.DictLabel, DictValue: req.DictValue,
		DictType: req.DictType, IsDefault: req.IsDefault, Status: req.Status, Remark: req.Remark,
	})
	if err != nil {
		return nil, err
	}
	return &api.DictDataCreateRes{}, nil
}

// DataUpdate 修改字典数据
func (c *cDict) DataUpdate(ctx context.Context, req *api.DictDataUpdateReq) (res *api.DictDataUpdateRes, err error) {
	if err = authz.Check(ctx, "system:dict:edit"); err != nil {
		return nil, err
	}
	err = service.Dict().DataUpdate(ctx, &model.DictDataSaveInput{
		Id: req.Id, DictSort: req.DictSort, DictLabel: req.DictLabel, DictValue: req.DictValue,
		DictType: req.DictType, IsDefault: req.IsDefault, Status: req.Status, Remark: req.Remark,
	})
	if err != nil {
		return nil, err
	}
	return &api.DictDataUpdateRes{}, nil
}

// DataDelete 删除字典数据
func (c *cDict) DataDelete(ctx context.Context, req *api.DictDataDeleteReq) (res *api.DictDataDeleteRes, err error) {
	if err = authz.Check(ctx, "system:dict:remove"); err != nil {
		return nil, err
	}
	if err = service.Dict().DataDelete(ctx, req.Ids); err != nil {
		return nil, err
	}
	return &api.DictDataDeleteRes{}, nil
}

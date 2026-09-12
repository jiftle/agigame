package service

import (
	"context"

	"agigame/console/backend/internal/model"
	"agigame/console/backend/internal/model/entity"
)

// IDict 字典服务
type IDict interface {
	TypeList(ctx context.Context, in *model.DictTypeQueryInput) (total int, list []*entity.DictType, err error)
	TypeAll(ctx context.Context) ([]*entity.DictType, error)
	TypeGetById(ctx context.Context, id int) (*entity.DictType, error)
	TypeCreate(ctx context.Context, in *model.DictTypeSaveInput) error
	TypeUpdate(ctx context.Context, in *model.DictTypeSaveInput) error
	TypeDelete(ctx context.Context, ids []int) error

	DataList(ctx context.Context, in *model.DictDataQueryInput) (total int, list []*entity.DictData, err error)
	DataByType(ctx context.Context, dictType string) ([]*entity.DictData, error)
	DataGetById(ctx context.Context, id int) (*entity.DictData, error)
	DataCreate(ctx context.Context, in *model.DictDataSaveInput) error
	DataUpdate(ctx context.Context, in *model.DictDataSaveInput) error
	DataDelete(ctx context.Context, ids []int) error
}

var localDict IDict

// Dict 获取字典服务
func Dict() IDict {
	if localDict == nil {
		panic("implement not found for interface IDict, forgot register?")
	}
	return localDict
}

// RegisterDict 注册字典服务
func RegisterDict(i IDict) {
	localDict = i
}

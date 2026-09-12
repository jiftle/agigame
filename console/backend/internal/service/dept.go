package service

import (
	"context"

	"agigame/console/backend/internal/model"
	"agigame/console/backend/internal/model/entity"
)

// IDept 部门服务
type IDept interface {
	List(ctx context.Context, in *model.DeptQueryInput) ([]*entity.Dept, error)
	GetById(ctx context.Context, id int) (*entity.Dept, error)
	Create(ctx context.Context, in *model.DeptSaveInput) error
	Update(ctx context.Context, in *model.DeptSaveInput) error
	Delete(ctx context.Context, id int) error
	Tree(ctx context.Context) ([]*model.DeptNode, error)
}

var localDept IDept

// Dept 获取部门服务
func Dept() IDept {
	if localDept == nil {
		panic("implement not found for interface IDept, forgot register?")
	}
	return localDept
}

// RegisterDept 注册部门服务
func RegisterDept(i IDept) {
	localDept = i
}

package service

import (
	"context"

	"agigame/console/backend/internal/model"
	"agigame/console/backend/internal/model/entity"
)

// IRole 角色服务
type IRole interface {
	List(ctx context.Context, in *model.RoleQueryInput) (total int, list []*entity.Role, err error)
	GetById(ctx context.Context, id int) (*entity.Role, error)
	Create(ctx context.Context, in *model.RoleSaveInput) error
	Update(ctx context.Context, in *model.RoleSaveInput) error
	Delete(ctx context.Context, ids []int) error
	GetMenuIds(ctx context.Context, roleId int) ([]int, error)
	GetByUserId(ctx context.Context, userId int) ([]*entity.Role, error)
}

var localRole IRole

// Role 获取角色服务
func Role() IRole {
	if localRole == nil {
		panic("implement not found for interface IRole, forgot register?")
	}
	return localRole
}

// RegisterRole 注册角色服务
func RegisterRole(i IRole) {
	localRole = i
}

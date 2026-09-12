package service

import (
	"context"

	"agigame/console/backend/internal/model"
	"agigame/console/backend/internal/model/entity"
)

// IMenu 菜单服务
type IMenu interface {
	List(ctx context.Context, in *model.MenuQueryInput) ([]*entity.Menu, error)
	GetTree(ctx context.Context, in *model.MenuQueryInput) ([]*model.MenuNode, error)
	GetById(ctx context.Context, id int) (*entity.Menu, error)
	Create(ctx context.Context, in *model.MenuSaveInput) error
	Update(ctx context.Context, in *model.MenuSaveInput) error
	Delete(ctx context.Context, id int) error
	GetPermsByUserId(ctx context.Context, userId int) ([]string, error)
	GetMenuTreeByUserId(ctx context.Context, userId int, isSuper bool) ([]*model.MenuNode, error)
	GetAllTree(ctx context.Context) ([]*model.MenuNode, error)
}

var localMenu IMenu

// Menu 获取菜单服务
func Menu() IMenu {
	if localMenu == nil {
		panic("implement not found for interface IMenu, forgot register?")
	}
	return localMenu
}

// RegisterMenu 注册菜单服务
func RegisterMenu(i IMenu) {
	localMenu = i
}

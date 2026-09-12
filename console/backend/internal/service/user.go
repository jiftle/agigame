package service

import (
	"context"

	"agigame/console/backend/internal/model"
	"agigame/console/backend/internal/model/entity"
)

// IUser 用户服务
type IUser interface {
	List(ctx context.Context, in *model.UserQueryInput) (total int, list []*entity.User, err error)
	GetById(ctx context.Context, id int) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	Create(ctx context.Context, in *model.UserSaveInput) error
	Update(ctx context.Context, in *model.UserSaveInput) error
	Delete(ctx context.Context, ids []int) error
	ResetPwd(ctx context.Context, id int, password string) error
	ChangeStatus(ctx context.Context, id, status int) error
	UpdateProfile(ctx context.Context, id int, in *model.ProfileUpdateInput) error
	GetRoleIds(ctx context.Context, userId int) ([]int, error)
	IsUsernameExist(ctx context.Context, username string, excludeId int) (bool, error)
	UpdateLoginInfo(ctx context.Context, userId int, ip string) error
}

var localUser IUser

// User 获取用户服务
func User() IUser {
	if localUser == nil {
		panic("implement not found for interface IUser, forgot register?")
	}
	return localUser
}

// RegisterUser 注册用户服务
func RegisterUser(i IUser) {
	localUser = i
}

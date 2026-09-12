package service

import (
	"context"

	"agigame/console/backend/internal/model"
)

// IAuth 认证服务
type IAuth interface {
	Login(ctx context.Context, in *model.LoginInput) (*model.LoginOutput, error)
	RefreshToken(ctx context.Context, refreshToken string) (*model.LoginOutput, error)
	GetContextUser(ctx context.Context, userId int) (*model.ContextUser, error)
	GetUserInfo(ctx context.Context, userId int) (*model.UserInfoOutput, error)
	ChangePassword(ctx context.Context, userId int, oldPassword, newPassword string) error
	UpdateProfile(ctx context.Context, userId int, in *model.ProfileUpdateInput) error
}

var localAuth IAuth

// Auth 获取认证服务
func Auth() IAuth {
	if localAuth == nil {
		panic("implement not found for interface IAuth, forgot register?")
	}
	return localAuth
}

// RegisterAuth 注册认证服务
func RegisterAuth(i IAuth) {
	localAuth = i
}

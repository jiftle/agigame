package auth

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"

	api "agigame/console/backend/api/v1/auth"
	"agigame/console/backend/internal/middleware"
	"agigame/console/backend/internal/model"
	"agigame/console/backend/internal/service"
	"agigame/console/backend/utility/ctxuser"
	"agigame/console/backend/utility/errcode"
	jwtutil "agigame/console/backend/utility/jwtutil"
)

// cLogin 公开认证接口
type cLogin struct{}

// NewLogin 创建公开认证控制器
func NewLogin() *cLogin {
	return &cLogin{}
}

// Login 登录
func (c *cLogin) Login(ctx context.Context, req *api.LoginReq) (res *api.LoginRes, err error) {
	r := g.RequestFromCtx(ctx)
	browser, os := parseUA(r.Header.Get("User-Agent"))
	out, err := service.Auth().Login(ctx, &model.LoginInput{
		Username: req.Username,
		Password: req.Password,
		Ip:       r.GetClientIp(),
		Browser:  browser,
		Os:       os,
	})
	if err != nil {
		return nil, err
	}
	return &api.LoginRes{
		Token:        out.Token,
		RefreshToken: out.RefreshToken,
		ExpiresIn:    out.ExpiresIn,
		TokenType:    out.TokenType,
	}, nil
}

// Refresh 刷新令牌
func (c *cLogin) Refresh(ctx context.Context, req *api.RefreshReq) (res *api.RefreshRes, err error) {
	out, err := service.Auth().RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}
	return &api.RefreshRes{
		Token:        out.Token,
		RefreshToken: out.RefreshToken,
		ExpiresIn:    out.ExpiresIn,
		TokenType:    out.TokenType,
	}, nil
}

// cProfile 受保护的用户接口
type cProfile struct{}

// NewProfile 创建用户信息控制器
func NewProfile() *cProfile {
	return &cProfile{}
}

// UserInfo 当前用户信息
func (c *cProfile) UserInfo(ctx context.Context, req *api.UserInfoReq) (res *api.UserInfoRes, err error) {
	user := ctxuser.Get(ctx)
	if user == nil {
		return nil, errcode.Unauthorized("未登录或登录已过期")
	}
	info, err := service.Auth().GetUserInfo(ctx, user.Id)
	if err != nil {
		return nil, err
	}
	return &api.UserInfoRes{
		User:  info.User,
		Roles: info.Roles,
		Perms: info.Perms,
		Menus: info.Menus,
	}, nil
}

// Password 修改当前用户密码
func (c *cProfile) Password(ctx context.Context, req *api.PasswordReq) (res *api.PasswordRes, err error) {
	user := ctxuser.Get(ctx)
	if user == nil {
		return nil, errcode.Unauthorized("未登录或登录已过期")
	}
	if err = service.Auth().ChangePassword(ctx, user.Id, req.OldPassword, req.NewPassword); err != nil {
		return nil, err
	}
	return &api.PasswordRes{}, nil
}

// ProfileUpdate 更新当前用户个人资料
func (c *cProfile) ProfileUpdate(ctx context.Context, req *api.ProfileUpdateReq) (res *api.ProfileUpdateRes, err error) {
	user := ctxuser.Get(ctx)
	if user == nil {
		return nil, errcode.Unauthorized("未登录或登录已过期")
	}
	err = service.Auth().UpdateProfile(ctx, user.Id, &model.ProfileUpdateInput{
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Avatar:   req.Avatar,
	})
	if err != nil {
		return nil, err
	}
	return &api.ProfileUpdateRes{}, nil
}

// Logout 退出登录
func (c *cProfile) Logout(ctx context.Context, req *api.LogoutReq) (res *api.LogoutRes, err error) {
	if req.RefreshToken != "" {
		jwtutil.Revoke(ctx, req.RefreshToken)
	}
	if r := g.RequestFromCtx(ctx); r != nil {
		jwtutil.Revoke(ctx, middleware.ExtractToken(r))
	}
	return &api.LogoutRes{}, nil
}

func parseUA(ua string) (browser, os string) {
	ua = strings.ToLower(ua)
	switch {
	case strings.Contains(ua, "edg"):
		browser = "Edge"
	case strings.Contains(ua, "chrome"):
		browser = "Chrome"
	case strings.Contains(ua, "firefox"):
		browser = "Firefox"
	case strings.Contains(ua, "safari"):
		browser = "Safari"
	default:
		browser = "Unknown"
	}
	switch {
	case strings.Contains(ua, "windows"):
		os = "Windows"
	case strings.Contains(ua, "mac os"):
		os = "macOS"
	case strings.Contains(ua, "android"):
		os = "Android"
	case strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad"):
		os = "iOS"
	case strings.Contains(ua, "linux"):
		os = "Linux"
	default:
		os = "Unknown"
	}
	return
}

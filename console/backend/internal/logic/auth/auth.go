package auth

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	"agigame/console/backend/internal/consts"
	"agigame/console/backend/internal/model"
	"agigame/console/backend/internal/service"
	"agigame/console/backend/utility/appcfg"
	"agigame/console/backend/utility/errcode"
	jwtutil "agigame/console/backend/utility/jwtutil"
	"agigame/console/backend/utility/password"
)

type sAuth struct{}

func init() {
	service.RegisterAuth(New())
}

// New 创建认证服务
func New() *sAuth {
	return &sAuth{}
}

// Login 登录
func (s *sAuth) Login(ctx context.Context, in *model.LoginInput) (*model.LoginOutput, error) {
	user, err := service.User().GetByUsername(ctx, in.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		_ = service.Log().LoginLogCreate(ctx, &model.LoginLogCreateInput{
			Username: in.Username, Ip: in.Ip, Browser: in.Browser, Os: in.Os,
			Status: consts.StatusDisabled, Msg: "用户不存在",
		})
		return nil, errcode.BadRequest("用户名或密码错误")
	}
	if !password.Verify(user.Password, in.Password) {
		_ = service.Log().LoginLogCreate(ctx, &model.LoginLogCreateInput{
			Username: in.Username, Ip: in.Ip, Browser: in.Browser, Os: in.Os,
			Status: consts.StatusDisabled, Msg: "密码错误",
		})
		return nil, errcode.BadRequest("用户名或密码错误")
	}
	if user.Status != consts.StatusEnabled {
		_ = service.Log().LoginLogCreate(ctx, &model.LoginLogCreateInput{
			Username: in.Username, Ip: in.Ip, Browser: in.Browser, Os: in.Os,
			Status: consts.StatusDisabled, Msg: "账号已停用",
		})
		return nil, errcode.BadRequest("账号已被停用，请联系管理员")
	}

	out, err := s.issueToken(ctx, user.Id, user.Username)
	if err != nil {
		return nil, err
	}

	_ = service.User().UpdateLoginInfo(ctx, user.Id, in.Ip)
	_ = service.Log().LoginLogCreate(ctx, &model.LoginLogCreateInput{
		Username: in.Username, Ip: in.Ip, Browser: in.Browser, Os: in.Os,
		Status: consts.StatusEnabled, Msg: "登录成功",
	})
	return out, nil
}

// RefreshToken 刷新令牌
func (s *sAuth) RefreshToken(ctx context.Context, refreshToken string) (*model.LoginOutput, error) {
	claims, err := jwtutil.Parse(ctx, refreshToken)
	if err != nil || claims.Type != jwtutil.TypeRefresh {
		return nil, errcode.TokenExpired("刷新令牌已失效，请重新登录")
	}
	user, err := service.User().GetById(ctx, claims.UserId)
	if err != nil || user == nil || user.Status != consts.StatusEnabled {
		return nil, errcode.Unauthorized("用户不存在或已停用")
	}
	out, err := s.issueToken(ctx, user.Id, user.Username)
	if err != nil {
		return nil, err
	}
	jwtutil.Revoke(ctx, refreshToken)
	return out, nil
}

// GetContextUser 获取登录上下文用户
func (s *sAuth) GetContextUser(ctx context.Context, userId int) (*model.ContextUser, error) {
	user, err := service.User().GetById(ctx, userId)
	if err != nil || user == nil {
		return nil, errcode.Unauthorized("用户不存在")
	}
	if user.Status != consts.StatusEnabled {
		return nil, errcode.Forbidden("账号已被停用")
	}

	roles, err := service.Role().GetByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	roleIds := make([]int, 0, len(roles))
	roleCodes := make([]string, 0, len(roles))
	isSuper := userId == appcfg.SuperAdminId(ctx)
	for _, r := range roles {
		roleIds = append(roleIds, r.Id)
		roleCodes = append(roleCodes, r.Code)
		if r.Code == consts.SuperRoleCode {
			isSuper = true
		}
	}

	perms, err := service.Menu().GetPermsByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	if isSuper {
		perms = []string{"*:*:*"}
	}

	return &model.ContextUser{
		Id:        user.Id,
		Username:  user.Username,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		DeptId:    user.DeptId,
		IsSuper:   isSuper,
		RoleIds:   roleIds,
		RoleCodes: roleCodes,
		Perms:     perms,
	}, nil
}

// GetUserInfo 获取当前用户信息
func (s *sAuth) GetUserInfo(ctx context.Context, userId int) (*model.UserInfoOutput, error) {
	user, err := service.User().GetById(ctx, userId)
	if err != nil || user == nil {
		return nil, errcode.Unauthorized("用户不存在")
	}
	user.Password = ""

	roles, err := service.Role().GetByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	roleCodes := make([]string, 0, len(roles))
	for _, r := range roles {
		roleCodes = append(roleCodes, r.Code)
	}

	ctxUser, err := s.GetContextUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	menus, err := service.Menu().GetMenuTreeByUserId(ctx, userId, ctxUser.IsSuper)
	if err != nil {
		return nil, err
	}
	return &model.UserInfoOutput{
		User:  user,
		Roles: roleCodes,
		Perms: ctxUser.Perms,
		Menus: menus,
	}, nil
}

// ChangePassword 修改当前用户密码
func (s *sAuth) ChangePassword(ctx context.Context, userId int, oldPassword, newPassword string) error {
	user, err := service.User().GetById(ctx, userId)
	if err != nil {
		return err
	}
	if !password.Verify(user.Password, oldPassword) {
		return errcode.BadRequest("原密码错误")
	}
	if oldPassword == newPassword {
		return errcode.BadRequest("新密码不能与原密码相同")
	}
	return service.User().ResetPwd(ctx, userId, newPassword)
}

// UpdateProfile 更新当前用户个人资料
func (s *sAuth) UpdateProfile(ctx context.Context, userId int, in *model.ProfileUpdateInput) error {
	return service.User().UpdateProfile(ctx, userId, in)
}

func (s *sAuth) issueToken(ctx context.Context, userId int, username string) (*model.LoginOutput, error) {
	expireSec := g.Cfg().MustGet(ctx, "jwt.expire", 86400).Int64()
	refreshSec := g.Cfg().MustGet(ctx, "jwt.refreshExpire", 604800).Int64()

	access, err := jwtutil.Generate(ctx, userId, username, jwtutil.TypeAccess, time.Duration(expireSec)*time.Second)
	if err != nil {
		return nil, err
	}
	refresh, err := jwtutil.Generate(ctx, userId, username, jwtutil.TypeRefresh, time.Duration(refreshSec)*time.Second)
	if err != nil {
		return nil, err
	}
	return &model.LoginOutput{
		Token:        access,
		RefreshToken: refresh,
		ExpiresIn:    expireSec,
		TokenType:    "Bearer",
	}, nil
}

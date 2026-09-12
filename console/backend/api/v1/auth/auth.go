package auth

import (
	"github.com/gogf/gf/v2/frame/g"

	"agigame/console/backend/internal/model"
	"agigame/console/backend/internal/model/entity"
)

// LoginReq 登录请求
type LoginReq struct {
	g.Meta   `path:"/auth/login" method:"post" tags:"认证" summary:"用户登录"`
	Username string `json:"username" v:"required#请输入用户名"`
	Password string `json:"password" v:"required#请输入密码"`
}

// LoginRes 登录响应
type LoginRes struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
	TokenType    string `json:"tokenType"`
}

// RefreshReq 刷新令牌请求
type RefreshReq struct {
	g.Meta       `path:"/auth/refresh" method:"post" tags:"认证" summary:"刷新令牌"`
	RefreshToken string `json:"refreshToken" v:"required#缺少刷新令牌"`
}

// RefreshRes 刷新令牌响应
type RefreshRes struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
	TokenType    string `json:"tokenType"`
}

// UserInfoReq 当前用户信息请求
type UserInfoReq struct {
	g.Meta `path:"/auth/user-info" method:"get" tags:"认证" summary:"当前用户信息"`
}

// UserInfoRes 当前用户信息响应
type UserInfoRes struct {
	User  *entity.User      `json:"user"`
	Roles []string          `json:"roles"`
	Perms []string          `json:"perms"`
	Menus []*model.MenuNode `json:"menus"`
}

// LogoutReq 退出登录请求
type LogoutReq struct {
	g.Meta       `path:"/auth/logout" method:"post" tags:"认证" summary:"退出登录"`
	RefreshToken string `json:"refreshToken"`
}

// LogoutRes 退出登录响应
type LogoutRes struct{}

// PasswordReq 修改当前用户密码请求
type PasswordReq struct {
	g.Meta      `path:"/auth/password" method:"put" tags:"认证" summary:"修改当前用户密码"`
	OldPassword string `json:"oldPassword" v:"required#请输入原密码"`
	NewPassword string `json:"newPassword" v:"required|length:6,64#请输入新密码|新密码长度需为6-64位"`
}

// PasswordRes 修改当前用户密码响应
type PasswordRes struct{}

// ProfileUpdateReq 更新个人资料请求
type ProfileUpdateReq struct {
	g.Meta   `path:"/auth/profile" method:"put" tags:"认证" summary:"更新个人资料"`
	Nickname string `json:"nickname" v:"required#请输入昵称"`
	Email    string `json:"email" v:"email#邮箱格式不正确"`
	Phone    string `json:"phone" v:"phone#手机号格式不正确"`
	Avatar   string `json:"avatar"`
}

// ProfileUpdateRes 更新个人资料响应
type ProfileUpdateRes struct{}

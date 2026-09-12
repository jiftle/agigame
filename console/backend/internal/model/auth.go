package model

import "agigame/console/backend/internal/model/entity"

// LoginInput 登录入参
type LoginInput struct {
	Username string
	Password string
	Ip       string
	Browser  string
	Os       string
}

// LoginOutput 登录出参
type LoginOutput struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
	TokenType    string `json:"tokenType"`
}

// UserInfoOutput 当前用户信息
type UserInfoOutput struct {
	User  *entity.User `json:"user"`
	Roles []string     `json:"roles"`
	Perms []string     `json:"perms"`
	Menus []*MenuNode  `json:"menus"`
}

// ProfileUpdateInput 个人资料更新
type ProfileUpdateInput struct {
	Nickname string
	Email    string
	Phone    string
	Avatar   string
}

// MenuNode 菜单树节点
type MenuNode struct {
	Id        int         `json:"id"`
	ParentId  int         `json:"parentId"`
	Title     string      `json:"title"`
	Name      string      `json:"name"`
	Path      string      `json:"path"`
	Component string      `json:"component"`
	Icon      string      `json:"icon"`
	Type      string      `json:"type"`
	Perms     string      `json:"perms"`
	Sort      int         `json:"sort"`
	Visible   int         `json:"visible"`
	Redirect  string      `json:"redirect"`
	IsFrame   int         `json:"isFrame"`
	IsCache   int         `json:"isCache"`
	Children  []*MenuNode `json:"children,omitempty"`
}

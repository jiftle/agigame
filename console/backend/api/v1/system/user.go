package system

import (
	"github.com/gogf/gf/v2/frame/g"

	"agigame/console/backend/internal/model/entity"
)

// ------------------------------ 用户 ------------------------------

type UserListReq struct {
	g.Meta   `path:"/system/user/list" method:"get" tags:"用户管理" summary:"用户列表"`
	PageNum  int    `json:"pageNum" d:"1"`
	PageSize int    `json:"pageSize" d:"10"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
	Status   *int   `json:"status"`
	DeptId   int    `json:"deptId"`
}

type UserListRes struct {
	Total int            `json:"total"`
	List  []*entity.User `json:"list"`
}

type UserGetReq struct {
	g.Meta `path:"/system/user/{id}" method:"get" tags:"用户管理" summary:"用户详情"`
	Id     int `json:"id" in:"path" v:"required#用户ID不能为空"`
}

type UserGetRes struct {
	User    *entity.User `json:"user"`
	RoleIds []int        `json:"roleIds"`
}

type UserCreateReq struct {
	g.Meta   `path:"/system/user" method:"post" tags:"用户管理" summary:"新增用户"`
	DeptId   int    `json:"deptId"`
	Username string `json:"username" v:"required#请输入用户名"`
	Nickname string `json:"nickname" v:"required#请输入昵称"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Sex      int    `json:"sex"`
	Status   int    `json:"status" d:"1"`
	Remark   string `json:"remark"`
	RoleIds  []int  `json:"roleIds"`
}

type UserCreateRes struct{}

type UserUpdateReq struct {
	g.Meta   `path:"/system/user" method:"put" tags:"用户管理" summary:"修改用户"`
	Id       int    `json:"id" v:"required#用户ID不能为空"`
	DeptId   int    `json:"deptId"`
	Username string `json:"username" v:"required#请输入用户名"`
	Nickname string `json:"nickname" v:"required#请输入昵称"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Sex      int    `json:"sex"`
	Status   int    `json:"status"`
	Remark   string `json:"remark"`
	RoleIds  []int  `json:"roleIds"`
}

type UserUpdateRes struct{}

type UserDeleteReq struct {
	g.Meta `path:"/system/user" method:"delete" tags:"用户管理" summary:"删除用户"`
	Ids    []int `json:"ids" v:"required#请选择要删除的用户"`
}

type UserDeleteRes struct{}

type UserResetPwdReq struct {
	g.Meta   `path:"/system/user/resetPwd" method:"put" tags:"用户管理" summary:"重置密码"`
	Id       int    `json:"id" v:"required#用户ID不能为空"`
	Password string `json:"password"`
}

type UserResetPwdRes struct{}

type UserStatusReq struct {
	g.Meta `path:"/system/user/status" method:"put" tags:"用户管理" summary:"修改用户状态"`
	Id     int `json:"id" v:"required#用户ID不能为空"`
	Status int `json:"status" v:"required#状态不能为空"`
}

type UserStatusRes struct{}

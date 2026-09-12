package system

import (
	"github.com/gogf/gf/v2/frame/g"

	"agigame/console/backend/internal/model/entity"
)

// ------------------------------ 角色 ------------------------------

type RoleListReq struct {
	g.Meta   `path:"/system/role/list" method:"get" tags:"角色管理" summary:"角色列表"`
	PageNum  int    `json:"pageNum" d:"1"`
	PageSize int    `json:"pageSize" d:"10"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	Status   *int   `json:"status"`
}

type RoleListRes struct {
	Total int            `json:"total"`
	List  []*entity.Role `json:"list"`
}

type RoleGetReq struct {
	g.Meta `path:"/system/role/{id}" method:"get" tags:"角色管理" summary:"角色详情"`
	Id     int `json:"id" in:"path" v:"required#角色ID不能为空"`
}

type RoleGetRes struct {
	Role    *entity.Role `json:"role"`
	MenuIds []int        `json:"menuIds"`
}

type RoleCreateReq struct {
	g.Meta    `path:"/system/role" method:"post" tags:"角色管理" summary:"新增角色"`
	Name      string `json:"name" v:"required#请输入角色名称"`
	Code      string `json:"code" v:"required#请输入角色编码"`
	Sort      int    `json:"sort"`
	DataScope int    `json:"dataScope" d:"1"`
	Status    int    `json:"status" d:"1"`
	Remark    string `json:"remark"`
	MenuIds   []int  `json:"menuIds"`
}

type RoleCreateRes struct{}

type RoleUpdateReq struct {
	g.Meta    `path:"/system/role" method:"put" tags:"角色管理" summary:"修改角色"`
	Id        int    `json:"id" v:"required#角色ID不能为空"`
	Name      string `json:"name" v:"required#请输入角色名称"`
	Code      string `json:"code" v:"required#请输入角色编码"`
	Sort      int    `json:"sort"`
	DataScope int    `json:"dataScope"`
	Status    int    `json:"status"`
	Remark    string `json:"remark"`
	MenuIds   []int  `json:"menuIds"`
}

type RoleUpdateRes struct{}

type RoleDeleteReq struct {
	g.Meta `path:"/system/role" method:"delete" tags:"角色管理" summary:"删除角色"`
	Ids    []int `json:"ids" v:"required#请选择要删除的角色"`
}

type RoleDeleteRes struct{}

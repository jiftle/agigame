package system

import (
	"github.com/gogf/gf/v2/frame/g"

	"agigame/console/backend/internal/model"
	"agigame/console/backend/internal/model/entity"
)

// ------------------------------ 部门 ------------------------------

type DeptListReq struct {
	g.Meta `path:"/system/dept/list" method:"get" tags:"部门管理" summary:"部门列表(树)"`
	Name   string `json:"name"`
	Status *int   `json:"status"`
}

type DeptListRes struct {
	List []*model.DeptNode `json:"list"`
}

type DeptGetReq struct {
	g.Meta `path:"/system/dept/{id}" method:"get" tags:"部门管理" summary:"部门详情"`
	Id     int `json:"id" in:"path" v:"required#部门ID不能为空"`
}

type DeptGetRes struct {
	Dept *entity.Dept `json:"dept"`
}

type DeptCreateReq struct {
	g.Meta   `path:"/system/dept" method:"post" tags:"部门管理" summary:"新增部门"`
	ParentId int    `json:"parentId"`
	Name     string `json:"name" v:"required#请输入部门名称"`
	Leader   string `json:"leader"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Sort     int    `json:"sort"`
	Status   int    `json:"status" d:"1"`
	Remark   string `json:"remark"`
}

type DeptCreateRes struct{}

type DeptUpdateReq struct {
	g.Meta   `path:"/system/dept" method:"put" tags:"部门管理" summary:"修改部门"`
	Id       int    `json:"id" v:"required#部门ID不能为空"`
	ParentId int    `json:"parentId"`
	Name     string `json:"name" v:"required#请输入部门名称"`
	Leader   string `json:"leader"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Sort     int    `json:"sort"`
	Status   int    `json:"status"`
	Remark   string `json:"remark"`
}

type DeptUpdateRes struct{}

type DeptDeleteReq struct {
	g.Meta `path:"/system/dept/{id}" method:"delete" tags:"部门管理" summary:"删除部门"`
	Id     int `json:"id" in:"path" v:"required#部门ID不能为空"`
}

type DeptDeleteRes struct{}

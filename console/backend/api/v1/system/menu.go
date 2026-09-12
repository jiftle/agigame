package system

import (
	"github.com/gogf/gf/v2/frame/g"

	"agigame/console/backend/internal/model"
	"agigame/console/backend/internal/model/entity"
)

// ------------------------------ 菜单 ------------------------------

type MenuListReq struct {
	g.Meta `path:"/system/menu/list" method:"get" tags:"菜单管理" summary:"菜单列表(树)"`
	Title  string `json:"title"`
	Status *int   `json:"status"`
}

type MenuListRes struct {
	List []*model.MenuNode `json:"list"`
}

type MenuGetReq struct {
	g.Meta `path:"/system/menu/{id}" method:"get" tags:"菜单管理" summary:"菜单详情"`
	Id     int `json:"id" in:"path" v:"required#菜单ID不能为空"`
}

type MenuGetRes struct {
	Menu *entity.Menu `json:"menu"`
}

type MenuCreateReq struct {
	g.Meta    `path:"/system/menu" method:"post" tags:"菜单管理" summary:"新增菜单"`
	ParentId  int    `json:"parentId"`
	Title     string `json:"title" v:"required#请输入菜单标题"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	Component string `json:"component"`
	Icon      string `json:"icon"`
	Type      string `json:"type" d:"C"`
	Perms     string `json:"perms"`
	Sort      int    `json:"sort"`
	Visible   int    `json:"visible" d:"1"`
	Status    int    `json:"status" d:"1"`
	Redirect  string `json:"redirect"`
	IsFrame   int    `json:"isFrame"`
	IsCache   int    `json:"isCache" d:"1"`
}

type MenuCreateRes struct{}

type MenuUpdateReq struct {
	g.Meta    `path:"/system/menu" method:"put" tags:"菜单管理" summary:"修改菜单"`
	Id        int    `json:"id" v:"required#菜单ID不能为空"`
	ParentId  int    `json:"parentId"`
	Title     string `json:"title" v:"required#请输入菜单标题"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	Component string `json:"component"`
	Icon      string `json:"icon"`
	Type      string `json:"type"`
	Perms     string `json:"perms"`
	Sort      int    `json:"sort"`
	Visible   int    `json:"visible"`
	Status    int    `json:"status"`
	Redirect  string `json:"redirect"`
	IsFrame   int    `json:"isFrame"`
	IsCache   int    `json:"isCache"`
}

type MenuUpdateRes struct{}

type MenuDeleteReq struct {
	g.Meta `path:"/system/menu/{id}" method:"delete" tags:"菜单管理" summary:"删除菜单"`
	Id     int `json:"id" in:"path" v:"required#菜单ID不能为空"`
}

type MenuDeleteRes struct{}

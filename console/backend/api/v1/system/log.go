package system

import (
	"github.com/gogf/gf/v2/frame/g"

	"agigame/console/backend/internal/model/entity"
)

// ------------------------------ 登录日志 ------------------------------

type LoginLogListReq struct {
	g.Meta    `path:"/system/log/login/list" method:"get" tags:"日志管理" summary:"登录日志列表"`
	PageNum   int    `json:"pageNum" d:"1"`
	PageSize  int    `json:"pageSize" d:"10"`
	Username  string `json:"username"`
	Status    *int   `json:"status"`
	BeginTime string `json:"beginTime"`
	EndTime   string `json:"endTime"`
}

type LoginLogListRes struct {
	Total int                `json:"total"`
	List  []*entity.LoginLog `json:"list"`
}

type LoginLogDeleteReq struct {
	g.Meta `path:"/system/log/login" method:"delete" tags:"日志管理" summary:"删除登录日志"`
	Ids    []int `json:"ids" v:"required#请选择要删除的日志"`
}

type LoginLogDeleteRes struct{}

type LoginLogClearReq struct {
	g.Meta `path:"/system/log/login/clear" method:"delete" tags:"日志管理" summary:"清空登录日志"`
}

type LoginLogClearRes struct{}

// ------------------------------ 操作日志 ------------------------------

type OperLogListReq struct {
	g.Meta    `path:"/system/log/oper/list" method:"get" tags:"日志管理" summary:"操作日志列表"`
	PageNum   int    `json:"pageNum" d:"1"`
	PageSize  int    `json:"pageSize" d:"10"`
	Title     string `json:"title"`
	OperName  string `json:"operName"`
	Status    *int   `json:"status"`
	BeginTime string `json:"beginTime"`
	EndTime   string `json:"endTime"`
}

type OperLogListRes struct {
	Total int               `json:"total"`
	List  []*entity.OperLog `json:"list"`
}

type OperLogDeleteReq struct {
	g.Meta `path:"/system/log/oper" method:"delete" tags:"日志管理" summary:"删除操作日志"`
	Ids    []int `json:"ids" v:"required#请选择要删除的日志"`
}

type OperLogDeleteRes struct{}

type OperLogClearReq struct {
	g.Meta `path:"/system/log/oper/clear" method:"delete" tags:"日志管理" summary:"清空操作日志"`
}

type OperLogClearRes struct{}

package system

import (
	"github.com/gogf/gf/v2/frame/g"

	"agigame/console/backend/internal/model/entity"
)

// ------------------------------ 参数配置 ------------------------------

type ConfigListReq struct {
	g.Meta     `path:"/system/config/list" method:"get" tags:"参数管理" summary:"参数列表"`
	PageNum    int    `json:"pageNum" d:"1"`
	PageSize   int    `json:"pageSize" d:"10"`
	ConfigName string `json:"configName"`
	ConfigKey  string `json:"configKey"`
}

type ConfigListRes struct {
	Total int              `json:"total"`
	List  []*entity.Config `json:"list"`
}

type ConfigGetReq struct {
	g.Meta `path:"/system/config/{id}" method:"get" tags:"参数管理" summary:"参数详情"`
	Id     int `json:"id" in:"path" v:"required#参数ID不能为空"`
}

type ConfigGetRes struct {
	Config *entity.Config `json:"config"`
}

type ConfigCreateReq struct {
	g.Meta      `path:"/system/config" method:"post" tags:"参数管理" summary:"新增参数"`
	ConfigName  string `json:"configName" v:"required#请输入参数名称"`
	ConfigKey   string `json:"configKey" v:"required#请输入参数键名"`
	ConfigValue string `json:"configValue" v:"required#请输入参数键值"`
	ConfigType  int    `json:"configType"`
	Remark      string `json:"remark"`
}

type ConfigCreateRes struct{}

type ConfigUpdateReq struct {
	g.Meta      `path:"/system/config" method:"put" tags:"参数管理" summary:"修改参数"`
	Id          int    `json:"id" v:"required#参数ID不能为空"`
	ConfigName  string `json:"configName" v:"required#请输入参数名称"`
	ConfigKey   string `json:"configKey" v:"required#请输入参数键名"`
	ConfigValue string `json:"configValue" v:"required#请输入参数键值"`
	ConfigType  int    `json:"configType"`
	Remark      string `json:"remark"`
}

type ConfigUpdateRes struct{}

type ConfigDeleteReq struct {
	g.Meta `path:"/system/config" method:"delete" tags:"参数管理" summary:"删除参数"`
	Ids    []int `json:"ids" v:"required#请选择要删除的参数"`
}

type ConfigDeleteRes struct{}

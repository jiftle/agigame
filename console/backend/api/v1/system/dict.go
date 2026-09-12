package system

import (
	"github.com/gogf/gf/v2/frame/g"

	"agigame/console/backend/internal/model/entity"
)

// ------------------------------ 字典类型 ------------------------------

type DictTypeListReq struct {
	g.Meta   `path:"/system/dict/type/list" method:"get" tags:"字典管理" summary:"字典类型列表"`
	PageNum  int    `json:"pageNum" d:"1"`
	PageSize int    `json:"pageSize" d:"10"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Status   *int   `json:"status"`
}

type DictTypeListRes struct {
	Total int                `json:"total"`
	List  []*entity.DictType `json:"list"`
}

type DictTypeGetReq struct {
	g.Meta `path:"/system/dict/type/{id}" method:"get" tags:"字典管理" summary:"字典类型详情"`
	Id     int `json:"id" in:"path" v:"required#字典类型ID不能为空"`
}

type DictTypeGetRes struct {
	DictType *entity.DictType `json:"dictType"`
}

type DictTypeCreateReq struct {
	g.Meta `path:"/system/dict/type" method:"post" tags:"字典管理" summary:"新增字典类型"`
	Name   string `json:"name" v:"required#请输入字典名称"`
	Type   string `json:"type" v:"required#请输入字典类型"`
	Status int    `json:"status" d:"1"`
	Remark string `json:"remark"`
}

type DictTypeCreateRes struct{}

type DictTypeUpdateReq struct {
	g.Meta `path:"/system/dict/type" method:"put" tags:"字典管理" summary:"修改字典类型"`
	Id     int    `json:"id" v:"required#字典类型ID不能为空"`
	Name   string `json:"name" v:"required#请输入字典名称"`
	Type   string `json:"type" v:"required#请输入字典类型"`
	Status int    `json:"status"`
	Remark string `json:"remark"`
}

type DictTypeUpdateRes struct{}

type DictTypeDeleteReq struct {
	g.Meta `path:"/system/dict/type" method:"delete" tags:"字典管理" summary:"删除字典类型"`
	Ids    []int `json:"ids" v:"required#请选择要删除的字典类型"`
}

type DictTypeDeleteRes struct{}

// ------------------------------ 字典数据 ------------------------------

type DictDataListReq struct {
	g.Meta    `path:"/system/dict/data/list" method:"get" tags:"字典管理" summary:"字典数据列表"`
	PageNum   int    `json:"pageNum" d:"1"`
	PageSize  int    `json:"pageSize" d:"10"`
	DictType  string `json:"dictType"`
	DictLabel string `json:"dictLabel"`
	Status    *int   `json:"status"`
}

type DictDataListRes struct {
	Total int                `json:"total"`
	List  []*entity.DictData `json:"list"`
}

type DictDataByTypeReq struct {
	g.Meta   `path:"/system/dict/data/type/{dictType}" method:"get" tags:"字典管理" summary:"按类型获取字典数据"`
	DictType string `json:"dictType" in:"path" v:"required#字典类型不能为空"`
}

type DictDataByTypeRes struct {
	List []*entity.DictData `json:"list"`
}

type DictDataGetReq struct {
	g.Meta `path:"/system/dict/data/{id}" method:"get" tags:"字典管理" summary:"字典数据详情"`
	Id     int `json:"id" in:"path" v:"required#字典数据ID不能为空"`
}

type DictDataGetRes struct {
	DictData *entity.DictData `json:"dictData"`
}

type DictDataCreateReq struct {
	g.Meta    `path:"/system/dict/data" method:"post" tags:"字典管理" summary:"新增字典数据"`
	DictSort  int    `json:"dictSort"`
	DictLabel string `json:"dictLabel" v:"required#请输入字典标签"`
	DictValue string `json:"dictValue" v:"required#请输入字典键值"`
	DictType  string `json:"dictType" v:"required#请选择字典类型"`
	IsDefault int    `json:"isDefault"`
	Status    int    `json:"status" d:"1"`
	Remark    string `json:"remark"`
}

type DictDataCreateRes struct{}

type DictDataUpdateReq struct {
	g.Meta    `path:"/system/dict/data" method:"put" tags:"字典管理" summary:"修改字典数据"`
	Id        int    `json:"id" v:"required#字典数据ID不能为空"`
	DictSort  int    `json:"dictSort"`
	DictLabel string `json:"dictLabel" v:"required#请输入字典标签"`
	DictValue string `json:"dictValue" v:"required#请输入字典键值"`
	DictType  string `json:"dictType" v:"required#请选择字典类型"`
	IsDefault int    `json:"isDefault"`
	Status    int    `json:"status"`
	Remark    string `json:"remark"`
}

type DictDataUpdateRes struct{}

type DictDataDeleteReq struct {
	g.Meta `path:"/system/dict/data" method:"delete" tags:"字典管理" summary:"删除字典数据"`
	Ids    []int `json:"ids" v:"required#请选择要删除的字典数据"`
}

type DictDataDeleteRes struct{}

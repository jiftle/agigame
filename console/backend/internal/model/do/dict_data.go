// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-09-10 22:09:07
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DictData is the golang structure of table sys_dict_data for DAO operations like Where/Data.
type DictData struct {
	g.Meta    `orm:"table:sys_dict_data, do:true"`
	Id        any         //
	DictSort  any         //
	DictLabel any         //
	DictValue any         //
	DictType  any         //
	IsDefault any         //
	Status    any         //
	Remark    any         //
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
	DeletedAt *gtime.Time //
}

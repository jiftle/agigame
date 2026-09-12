// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-09-10 22:09:07
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DictType is the golang structure of table sys_dict_type for DAO operations like Where/Data.
type DictType struct {
	g.Meta    `orm:"table:sys_dict_type, do:true"`
	Id        any         //
	Name      any         //
	Type      any         //
	Status    any         //
	Remark    any         //
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
	DeletedAt *gtime.Time //
}

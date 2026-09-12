// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-09-10 22:09:07
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Dept is the golang structure of table sys_dept for DAO operations like Where/Data.
type Dept struct {
	g.Meta    `orm:"table:sys_dept, do:true"`
	Id        any         //
	ParentId  any         //
	Ancestors any         //
	Name      any         //
	Leader    any         //
	Phone     any         //
	Email     any         //
	Sort      any         //
	Status    any         //
	Remark    any         //
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
	DeletedAt *gtime.Time //
}

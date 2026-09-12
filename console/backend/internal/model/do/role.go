// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-09-10 22:09:07
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Role is the golang structure of table sys_role for DAO operations like Where/Data.
type Role struct {
	g.Meta    `orm:"table:sys_role, do:true"`
	Id        any         //
	Name      any         //
	Code      any         //
	Sort      any         //
	DataScope any         //
	Status    any         //
	Remark    any         //
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
	DeletedAt *gtime.Time //
}

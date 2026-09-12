// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-09-10 22:09:07
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Menu is the golang structure of table sys_menu for DAO operations like Where/Data.
type Menu struct {
	g.Meta    `orm:"table:sys_menu, do:true"`
	Id        any         //
	ParentId  any         //
	Title     any         //
	Name      any         //
	Path      any         //
	Component any         //
	Icon      any         //
	Type      any         //
	Perms     any         //
	Sort      any         //
	Visible   any         //
	Status    any         //
	Redirect  any         //
	IsFrame   any         //
	IsCache   any         //
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
	DeletedAt *gtime.Time //
}

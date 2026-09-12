// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-09-10 22:09:07
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// RoleMenu is the golang structure of table sys_role_menu for DAO operations like Where/Data.
type RoleMenu struct {
	g.Meta    `orm:"table:sys_role_menu, do:true"`
	Id        any         //
	RoleId    any         //
	MenuId    any         //
	CreatedAt *gtime.Time //
}

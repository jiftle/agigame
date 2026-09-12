// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-09-10 22:09:07
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserRole is the golang structure of table sys_user_role for DAO operations like Where/Data.
type UserRole struct {
	g.Meta    `orm:"table:sys_user_role, do:true"`
	Id        any         //
	UserId    any         //
	RoleId    any         //
	CreatedAt *gtime.Time //
}

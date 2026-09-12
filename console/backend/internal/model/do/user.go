// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-09-10 22:09:07
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// User is the golang structure of table sys_user for DAO operations like Where/Data.
type User struct {
	g.Meta    `orm:"table:sys_user, do:true"`
	Id        any         //
	DeptId    any         //
	Username  any         //
	Password  any         //
	Nickname  any         //
	Email     any         //
	Phone     any         //
	Sex       any         //
	Avatar    any         //
	Status    any         //
	LoginIp   any         //
	LoginDate *gtime.Time //
	Remark    any         //
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
	DeletedAt *gtime.Time //
}

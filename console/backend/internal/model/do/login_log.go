// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-09-10 22:09:07
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// LoginLog is the golang structure of table sys_login_log for DAO operations like Where/Data.
type LoginLog struct {
	g.Meta    `orm:"table:sys_login_log, do:true"`
	Id        any         //
	Username  any         //
	Ip        any         //
	Location  any         //
	Browser   any         //
	Os        any         //
	Status    any         //
	Msg       any         //
	LoginTime *gtime.Time //
}

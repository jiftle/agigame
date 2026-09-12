// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-09-10 22:09:07
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// OperLog is the golang structure of table sys_oper_log for DAO operations like Where/Data.
type OperLog struct {
	g.Meta        `orm:"table:sys_oper_log, do:true"`
	Id            any         //
	Title         any         //
	BusinessType  any         //
	Method        any         //
	RequestMethod any         //
	OperName      any         //
	OperUrl       any         //
	OperIp        any         //
	OperParam     any         //
	JsonResult    any         //
	Status        any         //
	ErrorMsg      any         //
	Cost          any         //
	OperTime      *gtime.Time //
}

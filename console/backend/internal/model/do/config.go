// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-09-10 22:09:07
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Config is the golang structure of table sys_config for DAO operations like Where/Data.
type Config struct {
	g.Meta      `orm:"table:sys_config, do:true"`
	Id          any         //
	ConfigName  any         //
	ConfigKey   any         //
	ConfigValue any         //
	ConfigType  any         //
	Remark      any         //
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
	DeletedAt   *gtime.Time //
}

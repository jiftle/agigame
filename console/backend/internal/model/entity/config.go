// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-09-10 22:09:07
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Config is the golang structure for table config.
type Config struct {
	Id          int         `json:"id"          orm:"id"           description:""` //
	ConfigName  string      `json:"configName"  orm:"config_name"  description:""` //
	ConfigKey   string      `json:"configKey"   orm:"config_key"   description:""` //
	ConfigValue string      `json:"configValue" orm:"config_value" description:""` //
	ConfigType  int         `json:"configType"  orm:"config_type"  description:""` //
	Remark      string      `json:"remark"      orm:"remark"       description:""` //
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   description:""` //
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"   description:""` //
	DeletedAt   *gtime.Time `json:"deletedAt"   orm:"deleted_at"   description:""` //
}

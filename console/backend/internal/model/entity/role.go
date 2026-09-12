// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-09-10 22:09:07
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Role is the golang structure for table role.
type Role struct {
	Id        int         `json:"id"        orm:"id"         description:""` //
	Name      string      `json:"name"      orm:"name"       description:""` //
	Code      string      `json:"code"      orm:"code"       description:""` //
	Sort      int         `json:"sort"      orm:"sort"       description:""` //
	DataScope int         `json:"dataScope" orm:"data_scope" description:""` //
	Status    int         `json:"status"    orm:"status"     description:""` //
	Remark    string      `json:"remark"    orm:"remark"     description:""` //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""` //
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:""` //
	DeletedAt *gtime.Time `json:"deletedAt" orm:"deleted_at" description:""` //
}

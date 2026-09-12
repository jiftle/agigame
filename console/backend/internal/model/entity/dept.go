// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-09-10 22:09:07
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Dept is the golang structure for table dept.
type Dept struct {
	Id        int         `json:"id"        orm:"id"         description:""` //
	ParentId  int         `json:"parentId"  orm:"parent_id"  description:""` //
	Ancestors string      `json:"ancestors" orm:"ancestors"  description:""` //
	Name      string      `json:"name"      orm:"name"       description:""` //
	Leader    string      `json:"leader"    orm:"leader"     description:""` //
	Phone     string      `json:"phone"     orm:"phone"      description:""` //
	Email     string      `json:"email"     orm:"email"      description:""` //
	Sort      int         `json:"sort"      orm:"sort"       description:""` //
	Status    int         `json:"status"    orm:"status"     description:""` //
	Remark    string      `json:"remark"    orm:"remark"     description:""` //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""` //
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:""` //
	DeletedAt *gtime.Time `json:"deletedAt" orm:"deleted_at" description:""` //
}

// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-09-10 22:09:07
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// DictData is the golang structure for table dict_data.
type DictData struct {
	Id        int         `json:"id"        orm:"id"         description:""` //
	DictSort  int         `json:"dictSort"  orm:"dict_sort"  description:""` //
	DictLabel string      `json:"dictLabel" orm:"dict_label" description:""` //
	DictValue string      `json:"dictValue" orm:"dict_value" description:""` //
	DictType  string      `json:"dictType"  orm:"dict_type"  description:""` //
	IsDefault int         `json:"isDefault" orm:"is_default" description:""` //
	Status    int         `json:"status"    orm:"status"     description:""` //
	Remark    string      `json:"remark"    orm:"remark"     description:""` //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""` //
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:""` //
	DeletedAt *gtime.Time `json:"deletedAt" orm:"deleted_at" description:""` //
}

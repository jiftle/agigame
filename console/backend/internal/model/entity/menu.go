// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-09-10 22:09:07
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Menu is the golang structure for table menu.
type Menu struct {
	Id        int         `json:"id"        orm:"id"         description:""` //
	ParentId  int         `json:"parentId"  orm:"parent_id"  description:""` //
	Title     string      `json:"title"     orm:"title"      description:""` //
	Name      string      `json:"name"      orm:"name"       description:""` //
	Path      string      `json:"path"      orm:"path"       description:""` //
	Component string      `json:"component" orm:"component"  description:""` //
	Icon      string      `json:"icon"      orm:"icon"       description:""` //
	Type      string      `json:"type"      orm:"type"       description:""` //
	Perms     string      `json:"perms"     orm:"perms"      description:""` //
	Sort      int         `json:"sort"      orm:"sort"       description:""` //
	Visible   int         `json:"visible"   orm:"visible"    description:""` //
	Status    int         `json:"status"    orm:"status"     description:""` //
	Redirect  string      `json:"redirect"  orm:"redirect"   description:""` //
	IsFrame   int         `json:"isFrame"   orm:"is_frame"   description:""` //
	IsCache   int         `json:"isCache"   orm:"is_cache"   description:""` //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""` //
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:""` //
	DeletedAt *gtime.Time `json:"deletedAt" orm:"deleted_at" description:""` //
}

// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-09-10 22:09:07
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// RoleMenu is the golang structure for table role_menu.
type RoleMenu struct {
	Id        int         `json:"id"        orm:"id"         description:""` //
	RoleId    int         `json:"roleId"    orm:"role_id"    description:""` //
	MenuId    int         `json:"menuId"    orm:"menu_id"    description:""` //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""` //
}

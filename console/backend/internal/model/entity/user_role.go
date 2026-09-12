// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-09-10 22:09:07
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserRole is the golang structure for table user_role.
type UserRole struct {
	Id        int         `json:"id"        orm:"id"         description:""` //
	UserId    int         `json:"userId"    orm:"user_id"    description:""` //
	RoleId    int         `json:"roleId"    orm:"role_id"    description:""` //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""` //
}

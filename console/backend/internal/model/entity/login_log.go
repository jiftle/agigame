// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-09-10 22:09:07
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// LoginLog is the golang structure for table login_log.
type LoginLog struct {
	Id        int         `json:"id"        orm:"id"         description:""` //
	Username  string      `json:"username"  orm:"username"   description:""` //
	Ip        string      `json:"ip"        orm:"ip"         description:""` //
	Location  string      `json:"location"  orm:"location"   description:""` //
	Browser   string      `json:"browser"   orm:"browser"    description:""` //
	Os        string      `json:"os"        orm:"os"         description:""` //
	Status    int         `json:"status"    orm:"status"     description:""` //
	Msg       string      `json:"msg"       orm:"msg"        description:""` //
	LoginTime *gtime.Time `json:"loginTime" orm:"login_time" description:""` //
}

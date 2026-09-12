package model

// ContextUser 登录用户上下文
type ContextUser struct {
	Id        int      `json:"id"`
	Username  string   `json:"username"`
	Nickname  string   `json:"nickname"`
	Avatar    string   `json:"avatar"`
	DeptId    int      `json:"deptId"`
	IsSuper   bool     `json:"isSuper"`
	RoleIds   []int    `json:"roleIds"`
	RoleCodes []string `json:"roleCodes"`
	Perms     []string `json:"perms"`
}

// HasPerm 判断是否拥有指定权限
func (u *ContextUser) HasPerm(perm string) bool {
	if u == nil {
		return false
	}
	if u.IsSuper {
		return true
	}
	for _, p := range u.Perms {
		if p == perm {
			return true
		}
	}
	return false
}

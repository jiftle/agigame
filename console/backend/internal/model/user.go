package model

// UserQueryInput 用户查询
type UserQueryInput struct {
	PageInput
	Username string
	Nickname string
	Phone    string
	Status   *int
	DeptId   int
}

// UserSaveInput 用户新增/修改
type UserSaveInput struct {
	Id       int
	DeptId   int
	Username string
	Nickname string
	Email    string
	Phone    string
	Sex      int
	Status   int
	Remark   string
	Password string
	RoleIds  []int
}

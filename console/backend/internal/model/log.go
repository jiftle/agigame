package model

// LoginLogQueryInput 登录日志查询
type LoginLogQueryInput struct {
	PageInput
	Username  string
	Status    *int
	BeginTime string
	EndTime   string
}

// OperLogQueryInput 操作日志查询
type OperLogQueryInput struct {
	PageInput
	Title     string
	OperName  string
	Status    *int
	BeginTime string
	EndTime   string
}

// OperLogCreateInput 操作日志写入
type OperLogCreateInput struct {
	Title         string
	BusinessType  string
	Method        string
	RequestMethod string
	OperName      string
	OperUrl       string
	OperIp        string
	OperParam     string
	JsonResult    string
	Status        int
	ErrorMsg      string
	Cost          int
}

// LoginLogCreateInput 登录日志写入
type LoginLogCreateInput struct {
	Username string
	Ip       string
	Location string
	Browser  string
	Os       string
	Status   int
	Msg      string
}

package model

// ConfigQueryInput 参数查询
type ConfigQueryInput struct {
	PageInput
	ConfigName string
	ConfigKey  string
}

// ConfigSaveInput 参数新增/修改
type ConfigSaveInput struct {
	Id          int
	ConfigName  string
	ConfigKey   string
	ConfigValue string
	ConfigType  int
	Remark      string
}

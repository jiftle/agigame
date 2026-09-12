package model

// EmuSessionInfo 是控制台里一个运行中模拟器会话的摘要。
type EmuSessionInfo struct {
	Id        string `json:"id"`
	Console   string `json:"console"`
	Cart      string `json:"cart"`
	Game      string `json:"game"`
	Mode      string `json:"mode"`
	Auto      bool   `json:"auto"`
	Paused    bool   `json:"paused"`
	Frames    uint64 `json:"frames"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	CreatedAt string `json:"createdAt"`
}

// EmuStartInput 启动一个新的模拟器会话。
type EmuStartInput struct {
	Rom     string `json:"rom"`
	Console string `json:"console"` // gb | gba
	Game    string `json:"game"`
	Mode    string `json:"mode"`
	Palette string `json:"palette"`
}

// EmuConfigInput 更新运行中的会话配置。
type EmuConfigInput struct {
	Auto    *bool  `json:"auto"`
	Mode    string `json:"mode"`
	Palette string `json:"palette"`
}

// EmuRomInfo 描述 romDir 里一个可选 ROM。
type EmuRomInfo struct {
	Name    string `json:"name"`    // 文件名（启动会话时传它）
	Title   string `json:"title"`   // ROM 内部标题
	Console string `json:"console"` // gb | gba
	Ext     string `json:"ext"`
	Size    int64  `json:"size"`
}

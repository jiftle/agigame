package model

// EmuSessionInfo 是控制台里一个运行中模拟器会话的摘要。
type EmuSessionInfo struct {
	Id        string `json:"id"`
	Cart      string `json:"cart"`
	Game      string `json:"game"`
	Mode      string `json:"mode"`
	Auto      bool   `json:"auto"`
	Paused    bool   `json:"paused"`
	Frames    uint64 `json:"frames"`
	CreatedAt string `json:"createdAt"`
}

// EmuStartInput 启动一个新的模拟器会话。
type EmuStartInput struct {
	Rom     string `json:"rom"`
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

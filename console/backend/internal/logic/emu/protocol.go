package emu

// WebSocket 消息协议（与控制台前端约定，形状与 agigame/web 一致）。

type frameMsg struct {
	Type string `json:"type"`
	Img  string `json:"img"` // base64 PNG
	Tick uint64 `json:"tick"`
}

type stateMsg struct {
	Type  string         `json:"type"`
	State any            `json:"state"`
	Agent map[string]any `json:"agent,omitempty"`
}

type helloMsg struct {
	Type string `json:"type"`
	Id   string `json:"id"`
	Cart string `json:"cart"`
	Game string `json:"game"`
	FPS  int    `json:"fps"`
}

type logMsg struct {
	Type  string `json:"type"`
	Level string `json:"level"`
	Msg   string `json:"msg"`
}

// audioMsg 携带一段立体声 s16le PCM（base64）。
type audioMsg struct {
	Type string `json:"type"`
	PCM  string `json:"pcm"`
	Rate int    `json:"rate"`
}

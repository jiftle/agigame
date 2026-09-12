package service

import (
	"context"

	"agigame/console/backend/internal/model"
)

// EmuSubscriber 是对一个模拟器会话实时消息流的订阅句柄。
// Messages 里是已序列化好的 JSON 消息（hello/frame/state/log）。
type EmuSubscriber struct {
	Messages <-chan []byte
	Close    func()
}

// IEmu 模拟器会话服务。
type IEmu interface {
	List(ctx context.Context) []*model.EmuSessionInfo
	ListRoms(ctx context.Context) ([]*model.EmuRomInfo, error)
	Start(ctx context.Context, in *model.EmuStartInput) (*model.EmuSessionInfo, error)
	Get(ctx context.Context, id string) (*model.EmuSessionInfo, error)
	Stop(ctx context.Context, id string) error
	Control(ctx context.Context, id, action string) error
	Config(ctx context.Context, id string, in *model.EmuConfigInput) (*model.EmuSessionInfo, error)
	Subscribe(ctx context.Context, id string) (*EmuSubscriber, error)
	Input(ctx context.Context, id string, pressed, released []string) error
}

var localEmu IEmu

// Emu 获取模拟器会话服务。
func Emu() IEmu {
	if localEmu == nil {
		panic("implement not found for interface IEmu, forgot register?")
	}
	return localEmu
}

// RegisterEmu 注册模拟器会话服务。
func RegisterEmu(i IEmu) {
	localEmu = i
}

package service

import (
	"context"

	"agigame/console/backend/internal/model"
	"agigame/console/backend/internal/model/entity"
)

// IConfig 参数配置服务
type IConfig interface {
	List(ctx context.Context, in *model.ConfigQueryInput) (total int, list []*entity.Config, err error)
	GetById(ctx context.Context, id int) (*entity.Config, error)
	GetValue(ctx context.Context, configKey string) (string, error)
	Create(ctx context.Context, in *model.ConfigSaveInput) error
	Update(ctx context.Context, in *model.ConfigSaveInput) error
	Delete(ctx context.Context, ids []int) error
}

var localConfig IConfig

// Config 获取参数配置服务
func Config() IConfig {
	if localConfig == nil {
		panic("implement not found for interface IConfig, forgot register?")
	}
	return localConfig
}

// RegisterConfig 注册参数配置服务
func RegisterConfig(i IConfig) {
	localConfig = i
}

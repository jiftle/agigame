package service

import (
	"context"

	"agigame/console/backend/internal/model"
	"agigame/console/backend/internal/model/entity"
)

// ILog 日志服务
type ILog interface {
	LoginLogList(ctx context.Context, in *model.LoginLogQueryInput) (total int, list []*entity.LoginLog, err error)
	LoginLogCreate(ctx context.Context, in *model.LoginLogCreateInput) error
	LoginLogDelete(ctx context.Context, ids []int) error
	LoginLogClear(ctx context.Context) error

	OperLogList(ctx context.Context, in *model.OperLogQueryInput) (total int, list []*entity.OperLog, err error)
	OperLogCreate(ctx context.Context, in *model.OperLogCreateInput) error
	OperLogDelete(ctx context.Context, ids []int) error
	OperLogClear(ctx context.Context) error
}

var localLog ILog

// Log 获取日志服务
func Log() ILog {
	if localLog == nil {
		panic("implement not found for interface ILog, forgot register?")
	}
	return localLog
}

// RegisterLog 注册日志服务
func RegisterLog(i ILog) {
	localLog = i
}

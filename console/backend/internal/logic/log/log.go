package log

import (
	"context"

	"github.com/gogf/gf/v2/os/gtime"

	"agigame/console/backend/internal/dao"
	"agigame/console/backend/internal/model"
	"agigame/console/backend/internal/model/do"
	"agigame/console/backend/internal/model/entity"
	"agigame/console/backend/internal/service"
)

type sLog struct{}

func init() {
	service.RegisterLog(New())
}

// New 创建日志服务
func New() *sLog {
	return &sLog{}
}

// LoginLogList 登录日志列表
func (s *sLog) LoginLogList(ctx context.Context, in *model.LoginLogQueryInput) (total int, list []*entity.LoginLog, err error) {
	cols := dao.LoginLog.Columns()
	m := dao.LoginLog.Ctx(ctx)
	if in.Username != "" {
		m = m.WhereLike(cols.Username, "%"+in.Username+"%")
	}
	if in.Status != nil {
		m = m.Where(cols.Status, *in.Status)
	}
	if in.BeginTime != "" {
		m = m.WhereGTE(cols.LoginTime, in.BeginTime)
	}
	if in.EndTime != "" {
		m = m.WhereLTE(cols.LoginTime, in.EndTime)
	}
	total, err = m.Count()
	if err != nil || total == 0 {
		return 0, nil, err
	}
	err = m.Page(in.PageNum, in.PageSize).OrderDesc(cols.Id).Scan(&list)
	return
}

// LoginLogCreate 写入登录日志
func (s *sLog) LoginLogCreate(ctx context.Context, in *model.LoginLogCreateInput) error {
	_, err := dao.LoginLog.Ctx(ctx).Data(do.LoginLog{
		Username:  in.Username,
		Ip:        in.Ip,
		Location:  in.Location,
		Browser:   in.Browser,
		Os:        in.Os,
		Status:    in.Status,
		Msg:       in.Msg,
		LoginTime: gtime.Now(),
	}).Insert()
	return err
}

// LoginLogDelete 删除登录日志
func (s *sLog) LoginLogDelete(ctx context.Context, ids []int) error {
	_, err := dao.LoginLog.Ctx(ctx).WhereIn(dao.LoginLog.Columns().Id, ids).Delete()
	return err
}

// LoginLogClear 清空登录日志
func (s *sLog) LoginLogClear(ctx context.Context) error {
	_, err := dao.LoginLog.Ctx(ctx).Where("1=1").Delete()
	return err
}

// OperLogList 操作日志列表
func (s *sLog) OperLogList(ctx context.Context, in *model.OperLogQueryInput) (total int, list []*entity.OperLog, err error) {
	cols := dao.OperLog.Columns()
	m := dao.OperLog.Ctx(ctx)
	if in.Title != "" {
		m = m.WhereLike(cols.Title, "%"+in.Title+"%")
	}
	if in.OperName != "" {
		m = m.WhereLike(cols.OperName, "%"+in.OperName+"%")
	}
	if in.Status != nil {
		m = m.Where(cols.Status, *in.Status)
	}
	if in.BeginTime != "" {
		m = m.WhereGTE(cols.OperTime, in.BeginTime)
	}
	if in.EndTime != "" {
		m = m.WhereLTE(cols.OperTime, in.EndTime)
	}
	total, err = m.Count()
	if err != nil || total == 0 {
		return 0, nil, err
	}
	err = m.Page(in.PageNum, in.PageSize).OrderDesc(cols.Id).Scan(&list)
	return
}

// OperLogCreate 写入操作日志
func (s *sLog) OperLogCreate(ctx context.Context, in *model.OperLogCreateInput) error {
	_, err := dao.OperLog.Ctx(ctx).Data(do.OperLog{
		Title:         in.Title,
		BusinessType:  in.BusinessType,
		Method:        in.Method,
		RequestMethod: in.RequestMethod,
		OperName:      in.OperName,
		OperUrl:       in.OperUrl,
		OperIp:        in.OperIp,
		OperParam:     in.OperParam,
		JsonResult:    in.JsonResult,
		Status:        in.Status,
		ErrorMsg:      in.ErrorMsg,
		Cost:          in.Cost,
		OperTime:      gtime.Now(),
	}).Insert()
	return err
}

// OperLogDelete 删除操作日志
func (s *sLog) OperLogDelete(ctx context.Context, ids []int) error {
	_, err := dao.OperLog.Ctx(ctx).WhereIn(dao.OperLog.Columns().Id, ids).Delete()
	return err
}

// OperLogClear 清空操作日志
func (s *sLog) OperLogClear(ctx context.Context) error {
	_, err := dao.OperLog.Ctx(ctx).Where("1=1").Delete()
	return err
}

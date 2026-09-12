package config

import (
	"context"

	"agigame/console/backend/internal/dao"
	"agigame/console/backend/internal/model"
	"agigame/console/backend/internal/model/do"
	"agigame/console/backend/internal/model/entity"
	"agigame/console/backend/internal/service"
	"agigame/console/backend/utility/errcode"
)

type sConfig struct{}

func init() {
	service.RegisterConfig(New())
}

// New 创建参数配置服务
func New() *sConfig {
	return &sConfig{}
}

// List 参数列表
func (s *sConfig) List(ctx context.Context, in *model.ConfigQueryInput) (total int, list []*entity.Config, err error) {
	cols := dao.Config.Columns()
	m := dao.Config.Ctx(ctx)
	if in.ConfigName != "" {
		m = m.WhereLike(cols.ConfigName, "%"+in.ConfigName+"%")
	}
	if in.ConfigKey != "" {
		m = m.WhereLike(cols.ConfigKey, "%"+in.ConfigKey+"%")
	}
	total, err = m.Count()
	if err != nil || total == 0 {
		return 0, nil, err
	}
	err = m.Page(in.PageNum, in.PageSize).OrderDesc(cols.Id).Scan(&list)
	return
}

// GetById 获取参数
func (s *sConfig) GetById(ctx context.Context, id int) (*entity.Config, error) {
	var item *entity.Config
	if err := dao.Config.Ctx(ctx).Where(dao.Config.Columns().Id, id).Scan(&item); err != nil {
		return nil, err
	}
	if item == nil {
		return nil, errcode.NotFound("参数不存在")
	}
	return item, nil
}

// GetValue 按键获取参数值
func (s *sConfig) GetValue(ctx context.Context, configKey string) (string, error) {
	value, err := dao.Config.Ctx(ctx).Where(dao.Config.Columns().ConfigKey, configKey).Value(dao.Config.Columns().ConfigValue)
	if err != nil {
		return "", err
	}
	return value.String(), nil
}

// Create 新增参数
func (s *sConfig) Create(ctx context.Context, in *model.ConfigSaveInput) error {
	exist, err := dao.Config.Ctx(ctx).Where(dao.Config.Columns().ConfigKey, in.ConfigKey).Count()
	if err != nil {
		return err
	}
	if exist > 0 {
		return errcode.BadRequest("参数键名已存在")
	}
	_, err = dao.Config.Ctx(ctx).Data(do.Config{
		ConfigName:  in.ConfigName,
		ConfigKey:   in.ConfigKey,
		ConfigValue: in.ConfigValue,
		ConfigType:  in.ConfigType,
		Remark:      in.Remark,
	}).Insert()
	return err
}

// Update 修改参数
func (s *sConfig) Update(ctx context.Context, in *model.ConfigSaveInput) error {
	exist, err := dao.Config.Ctx(ctx).
		Where(dao.Config.Columns().ConfigKey, in.ConfigKey).
		WhereNot(dao.Config.Columns().Id, in.Id).
		Count()
	if err != nil {
		return err
	}
	if exist > 0 {
		return errcode.BadRequest("参数键名已存在")
	}
	_, err = dao.Config.Ctx(ctx).Where(dao.Config.Columns().Id, in.Id).Data(do.Config{
		ConfigName:  in.ConfigName,
		ConfigKey:   in.ConfigKey,
		ConfigValue: in.ConfigValue,
		ConfigType:  in.ConfigType,
		Remark:      in.Remark,
	}).Update()
	return err
}

// Delete 删除参数
func (s *sConfig) Delete(ctx context.Context, ids []int) error {
	_, err := dao.Config.Ctx(ctx).WhereIn(dao.Config.Columns().Id, ids).Delete()
	return err
}

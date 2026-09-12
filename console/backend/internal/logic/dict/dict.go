package dict

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"

	"agigame/console/backend/internal/consts"
	"agigame/console/backend/internal/dao"
	"agigame/console/backend/internal/model"
	"agigame/console/backend/internal/model/do"
	"agigame/console/backend/internal/model/entity"
	"agigame/console/backend/internal/service"
	"agigame/console/backend/utility/errcode"
)

type sDict struct{}

func init() {
	service.RegisterDict(New())
}

// New 创建字典服务
func New() *sDict {
	return &sDict{}
}

// TypeList 字典类型列表
func (s *sDict) TypeList(ctx context.Context, in *model.DictTypeQueryInput) (total int, list []*entity.DictType, err error) {
	cols := dao.DictType.Columns()
	m := dao.DictType.Ctx(ctx)
	if in.Name != "" {
		m = m.WhereLike(cols.Name, "%"+in.Name+"%")
	}
	if in.Type != "" {
		m = m.WhereLike(cols.Type, "%"+in.Type+"%")
	}
	if in.Status != nil {
		m = m.Where(cols.Status, *in.Status)
	}
	total, err = m.Count()
	if err != nil || total == 0 {
		return 0, nil, err
	}
	err = m.Page(in.PageNum, in.PageSize).OrderDesc(cols.Id).Scan(&list)
	return
}

// TypeAll 全部字典类型
func (s *sDict) TypeAll(ctx context.Context) ([]*entity.DictType, error) {
	var list []*entity.DictType
	err := dao.DictType.Ctx(ctx).Where(dao.DictType.Columns().Status, consts.StatusEnabled).OrderAsc(dao.DictType.Columns().Id).Scan(&list)
	return list, err
}

// TypeGetById 获取字典类型
func (s *sDict) TypeGetById(ctx context.Context, id int) (*entity.DictType, error) {
	var item *entity.DictType
	if err := dao.DictType.Ctx(ctx).Where(dao.DictType.Columns().Id, id).Scan(&item); err != nil {
		return nil, err
	}
	if item == nil {
		return nil, errcode.NotFound("字典类型不存在")
	}
	return item, nil
}

// TypeCreate 新增字典类型
func (s *sDict) TypeCreate(ctx context.Context, in *model.DictTypeSaveInput) error {
	exist, err := dao.DictType.Ctx(ctx).Where(dao.DictType.Columns().Type, in.Type).Count()
	if err != nil {
		return err
	}
	if exist > 0 {
		return errcode.BadRequest("字典类型已存在")
	}
	_, err = dao.DictType.Ctx(ctx).Data(do.DictType{
		Name:   in.Name,
		Type:   in.Type,
		Status: in.Status,
		Remark: in.Remark,
	}).Insert()
	return err
}

// TypeUpdate 修改字典类型
func (s *sDict) TypeUpdate(ctx context.Context, in *model.DictTypeSaveInput) error {
	old, err := s.TypeGetById(ctx, in.Id)
	if err != nil {
		return err
	}
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err := tx.Model(dao.DictType.Table()).Ctx(ctx).Where(dao.DictType.Columns().Id, in.Id).Data(do.DictType{
			Name:   in.Name,
			Type:   in.Type,
			Status: in.Status,
			Remark: in.Remark,
		}).Update()
		if err != nil {
			return err
		}
		// 同步字典数据的类型标识
		if old.Type != in.Type {
			_, err = tx.Model(dao.DictData.Table()).Ctx(ctx).
				Where(dao.DictData.Columns().DictType, old.Type).
				Data(do.DictData{DictType: in.Type}).Update()
		}
		return err
	})
}

// TypeDelete 删除字典类型
func (s *sDict) TypeDelete(ctx context.Context, ids []int) error {
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		values, err := tx.Model(dao.DictType.Table()).Ctx(ctx).WhereIn(dao.DictType.Columns().Id, ids).Array(dao.DictType.Columns().Type)
		if err != nil {
			return err
		}
		types := make([]string, 0, len(values))
		for _, v := range values {
			types = append(types, v.String())
		}
		if _, err = tx.Model(dao.DictType.Table()).Ctx(ctx).WhereIn(dao.DictType.Columns().Id, ids).Delete(); err != nil {
			return err
		}
		if len(types) > 0 {
			_, err = tx.Model(dao.DictData.Table()).Ctx(ctx).WhereIn(dao.DictData.Columns().DictType, types).Delete()
			return err
		}
		return nil
	})
}

// DataList 字典数据列表
func (s *sDict) DataList(ctx context.Context, in *model.DictDataQueryInput) (total int, list []*entity.DictData, err error) {
	cols := dao.DictData.Columns()
	m := dao.DictData.Ctx(ctx)
	if in.DictType != "" {
		m = m.Where(cols.DictType, in.DictType)
	}
	if in.DictLabel != "" {
		m = m.WhereLike(cols.DictLabel, "%"+in.DictLabel+"%")
	}
	if in.Status != nil {
		m = m.Where(cols.Status, *in.Status)
	}
	total, err = m.Count()
	if err != nil || total == 0 {
		return 0, nil, err
	}
	err = m.Page(in.PageNum, in.PageSize).OrderAsc(cols.DictSort).OrderAsc(cols.Id).Scan(&list)
	return
}

// DataByType 按类型获取字典数据
func (s *sDict) DataByType(ctx context.Context, dictType string) ([]*entity.DictData, error) {
	var list []*entity.DictData
	err := dao.DictData.Ctx(ctx).
		Where(dao.DictData.Columns().DictType, dictType).
		Where(dao.DictData.Columns().Status, consts.StatusEnabled).
		OrderAsc(dao.DictData.Columns().DictSort).OrderAsc(dao.DictData.Columns().Id).
		Scan(&list)
	return list, err
}

// DataGetById 获取字典数据
func (s *sDict) DataGetById(ctx context.Context, id int) (*entity.DictData, error) {
	var item *entity.DictData
	if err := dao.DictData.Ctx(ctx).Where(dao.DictData.Columns().Id, id).Scan(&item); err != nil {
		return nil, err
	}
	if item == nil {
		return nil, errcode.NotFound("字典数据不存在")
	}
	return item, nil
}

// DataCreate 新增字典数据
func (s *sDict) DataCreate(ctx context.Context, in *model.DictDataSaveInput) error {
	_, err := dao.DictData.Ctx(ctx).Data(do.DictData{
		DictSort:  in.DictSort,
		DictLabel: in.DictLabel,
		DictValue: in.DictValue,
		DictType:  in.DictType,
		IsDefault: in.IsDefault,
		Status:    in.Status,
		Remark:    in.Remark,
	}).Insert()
	return err
}

// DataUpdate 修改字典数据
func (s *sDict) DataUpdate(ctx context.Context, in *model.DictDataSaveInput) error {
	_, err := dao.DictData.Ctx(ctx).Where(dao.DictData.Columns().Id, in.Id).Data(do.DictData{
		DictSort:  in.DictSort,
		DictLabel: in.DictLabel,
		DictValue: in.DictValue,
		DictType:  in.DictType,
		IsDefault: in.IsDefault,
		Status:    in.Status,
		Remark:    in.Remark,
	}).Update()
	return err
}

// DataDelete 删除字典数据
func (s *sDict) DataDelete(ctx context.Context, ids []int) error {
	_, err := dao.DictData.Ctx(ctx).WhereIn(dao.DictData.Columns().Id, ids).Delete()
	return err
}

package dept

import (
	"context"
	"fmt"
	"strings"

	"agigame/console/backend/internal/dao"
	"agigame/console/backend/internal/model"
	"agigame/console/backend/internal/model/do"
	"agigame/console/backend/internal/model/entity"
	"agigame/console/backend/internal/service"
	"agigame/console/backend/utility/errcode"
)

type sDept struct{}

func init() {
	service.RegisterDept(New())
}

// New 创建部门服务
func New() *sDept {
	return &sDept{}
}

// List 部门列表
func (s *sDept) List(ctx context.Context, in *model.DeptQueryInput) ([]*entity.Dept, error) {
	cols := dao.Dept.Columns()
	m := dao.Dept.Ctx(ctx)
	if in.Name != "" {
		m = m.WhereLike(cols.Name, "%"+in.Name+"%")
	}
	if in.Status != nil {
		m = m.Where(cols.Status, *in.Status)
	}
	var list []*entity.Dept
	err := m.OrderAsc(cols.Sort).OrderAsc(cols.Id).Scan(&list)
	return list, err
}

// GetById 获取部门
func (s *sDept) GetById(ctx context.Context, id int) (*entity.Dept, error) {
	var dept *entity.Dept
	if err := dao.Dept.Ctx(ctx).Where(dao.Dept.Columns().Id, id).Scan(&dept); err != nil {
		return nil, err
	}
	if dept == nil {
		return nil, errcode.NotFound("部门不存在")
	}
	return dept, nil
}

// Tree 部门树
func (s *sDept) Tree(ctx context.Context) ([]*model.DeptNode, error) {
	var list []*entity.Dept
	if err := dao.Dept.Ctx(ctx).OrderAsc(dao.Dept.Columns().Sort).OrderAsc(dao.Dept.Columns().Id).Scan(&list); err != nil {
		return nil, err
	}
	return buildTree(list, 0), nil
}

// Create 新增部门
func (s *sDept) Create(ctx context.Context, in *model.DeptSaveInput) error {
	ancestors, err := s.buildAncestors(ctx, in.ParentId)
	if err != nil {
		return err
	}
	_, err = dao.Dept.Ctx(ctx).Data(do.Dept{
		ParentId:  in.ParentId,
		Ancestors: ancestors,
		Name:      in.Name,
		Leader:    in.Leader,
		Phone:     in.Phone,
		Email:     in.Email,
		Sort:      in.Sort,
		Status:    in.Status,
		Remark:    in.Remark,
	}).Insert()
	return err
}

// Update 修改部门
func (s *sDept) Update(ctx context.Context, in *model.DeptSaveInput) error {
	if in.Id == in.ParentId {
		return errcode.BadRequest("上级部门不能是自己")
	}
	old, err := s.GetById(ctx, in.Id)
	if err != nil {
		return err
	}
	ancestors, err := s.buildAncestors(ctx, in.ParentId)
	if err != nil {
		return err
	}
	if strings.Contains(ancestors, fmt.Sprintf(",%d", in.Id)) || ancestors == fmt.Sprintf("%d", in.Id) {
		return errcode.BadRequest("上级部门不能是自己的子部门")
	}
	_, err = dao.Dept.Ctx(ctx).Where(dao.Dept.Columns().Id, in.Id).Data(do.Dept{
		ParentId:  in.ParentId,
		Ancestors: ancestors,
		Name:      in.Name,
		Leader:    in.Leader,
		Phone:     in.Phone,
		Email:     in.Email,
		Sort:      in.Sort,
		Status:    in.Status,
		Remark:    in.Remark,
	}).Update()
	if err != nil {
		return err
	}
	// 父级变更时同步子部门 ancestors
	if old.ParentId != in.ParentId {
		oldPrefix := old.Ancestors
		if oldPrefix == "" {
			oldPrefix = "0"
		}
		var children []*entity.Dept
		if err = dao.Dept.Ctx(ctx).WhereLike(dao.Dept.Columns().Ancestors, oldPrefix+","+fmt.Sprint(in.Id)+"%").Scan(&children); err != nil {
			return err
		}
		for _, child := range children {
			newAnc := strings.Replace(child.Ancestors, oldPrefix+","+fmt.Sprint(in.Id), ancestors+","+fmt.Sprint(in.Id), 1)
			if _, err = dao.Dept.Ctx(ctx).Where(dao.Dept.Columns().Id, child.Id).Data(do.Dept{Ancestors: newAnc}).Update(); err != nil {
				return err
			}
		}
	}
	return nil
}

// Delete 删除部门
func (s *sDept) Delete(ctx context.Context, id int) error {
	count, err := dao.Dept.Ctx(ctx).Where(dao.Dept.Columns().ParentId, id).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return errcode.BadRequest("存在下级部门，不允许删除")
	}
	userCount, err := dao.User.Ctx(ctx).Where(dao.User.Columns().DeptId, id).Count()
	if err != nil {
		return err
	}
	if userCount > 0 {
		return errcode.BadRequest("部门下存在用户，不允许删除")
	}
	_, err = dao.Dept.Ctx(ctx).Where(dao.Dept.Columns().Id, id).Delete()
	return err
}

func (s *sDept) buildAncestors(ctx context.Context, parentId int) (string, error) {
	if parentId == 0 {
		return "0", nil
	}
	parent, err := s.GetById(ctx, parentId)
	if err != nil {
		return "", err
	}
	if parent.Ancestors == "" {
		return fmt.Sprint(parentId), nil
	}
	return parent.Ancestors + "," + fmt.Sprint(parentId), nil
}

func buildTree(list []*entity.Dept, parentId int) []*model.DeptNode {
	nodes := make([]*model.DeptNode, 0)
	for _, d := range list {
		if d.ParentId != parentId {
			continue
		}
		node := &model.DeptNode{
			Id:       d.Id,
			ParentId: d.ParentId,
			Name:     d.Name,
			Leader:   d.Leader,
			Phone:    d.Phone,
			Email:    d.Email,
			Sort:     d.Sort,
			Status:   d.Status,
		}
		node.Children = buildTree(list, d.Id)
		nodes = append(nodes, node)
	}
	return nodes
}

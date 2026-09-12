package role

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

type sRole struct{}

func init() {
	service.RegisterRole(New())
}

// New 创建角色服务
func New() *sRole {
	return &sRole{}
}

// List 角色列表
func (s *sRole) List(ctx context.Context, in *model.RoleQueryInput) (total int, list []*entity.Role, err error) {
	cols := dao.Role.Columns()
	m := dao.Role.Ctx(ctx)
	if in.Name != "" {
		m = m.WhereLike(cols.Name, "%"+in.Name+"%")
	}
	if in.Code != "" {
		m = m.WhereLike(cols.Code, "%"+in.Code+"%")
	}
	if in.Status != nil {
		m = m.Where(cols.Status, *in.Status)
	}
	total, err = m.Count()
	if err != nil || total == 0 {
		return 0, nil, err
	}
	err = m.Page(in.PageNum, in.PageSize).OrderAsc(cols.Sort).OrderAsc(cols.Id).Scan(&list)
	return
}

// GetById 获取角色
func (s *sRole) GetById(ctx context.Context, id int) (*entity.Role, error) {
	var role *entity.Role
	if err := dao.Role.Ctx(ctx).Where(dao.Role.Columns().Id, id).Scan(&role); err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errcode.NotFound("角色不存在")
	}
	return role, nil
}

// Create 新增角色
func (s *sRole) Create(ctx context.Context, in *model.RoleSaveInput) error {
	if err := s.checkCodeUnique(ctx, in.Code, 0); err != nil {
		return err
	}
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		id, err := tx.Model(dao.Role.Table()).Ctx(ctx).Data(do.Role{
			Name:      in.Name,
			Code:      in.Code,
			Sort:      in.Sort,
			DataScope: in.DataScope,
			Status:    in.Status,
			Remark:    in.Remark,
		}).InsertAndGetId()
		if err != nil {
			return err
		}
		return s.saveMenus(ctx, tx, int(id), in.MenuIds)
	})
}

// Update 修改角色
func (s *sRole) Update(ctx context.Context, in *model.RoleSaveInput) error {
	old, err := s.GetById(ctx, in.Id)
	if err != nil {
		return err
	}
	if old.Code == consts.SuperRoleCode {
		return errcode.BadRequest("内置超级管理员角色不允许修改")
	}
	if err := s.checkCodeUnique(ctx, in.Code, in.Id); err != nil {
		return err
	}
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err := tx.Model(dao.Role.Table()).Ctx(ctx).Where(dao.Role.Columns().Id, in.Id).Data(do.Role{
			Name:      in.Name,
			Code:      in.Code,
			Sort:      in.Sort,
			DataScope: in.DataScope,
			Status:    in.Status,
			Remark:    in.Remark,
		}).Update()
		if err != nil {
			return err
		}
		if _, err = tx.Model(dao.RoleMenu.Table()).Ctx(ctx).Where(dao.RoleMenu.Columns().RoleId, in.Id).Delete(); err != nil {
			return err
		}
		return s.saveMenus(ctx, tx, in.Id, in.MenuIds)
	})
}

// Delete 删除角色
func (s *sRole) Delete(ctx context.Context, ids []int) error {
	for _, id := range ids {
		role, err := s.GetById(ctx, id)
		if err != nil {
			return err
		}
		if role.Code == consts.SuperRoleCode {
			return errcode.BadRequest("内置超级管理员角色不允许删除")
		}
		count, err := dao.UserRole.Ctx(ctx).Where(dao.UserRole.Columns().RoleId, id).Count()
		if err != nil {
			return err
		}
		if count > 0 {
			return errcode.BadRequest("角色已分配给用户，不允许删除")
		}
	}
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(dao.Role.Table()).Ctx(ctx).WhereIn(dao.Role.Columns().Id, ids).Delete(); err != nil {
			return err
		}
		_, err := tx.Model(dao.RoleMenu.Table()).Ctx(ctx).WhereIn(dao.RoleMenu.Columns().RoleId, ids).Delete()
		return err
	})
}

// GetMenuIds 获取角色菜单ID
func (s *sRole) GetMenuIds(ctx context.Context, roleId int) ([]int, error) {
	values, err := dao.RoleMenu.Ctx(ctx).Where(dao.RoleMenu.Columns().RoleId, roleId).Array(dao.RoleMenu.Columns().MenuId)
	if err != nil {
		return nil, err
	}
	ids := make([]int, 0, len(values))
	for _, v := range values {
		ids = append(ids, v.Int())
	}
	return ids, nil
}

// GetByUserId 获取用户角色
func (s *sRole) GetByUserId(ctx context.Context, userId int) ([]*entity.Role, error) {
	values, err := dao.UserRole.Ctx(ctx).Where(dao.UserRole.Columns().UserId, userId).Array(dao.UserRole.Columns().RoleId)
	if err != nil {
		return nil, err
	}
	if len(values) == 0 {
		return nil, nil
	}
	roleIds := make([]int, 0, len(values))
	for _, v := range values {
		roleIds = append(roleIds, v.Int())
	}
	var list []*entity.Role
	err = dao.Role.Ctx(ctx).
		WhereIn(dao.Role.Columns().Id, roleIds).
		Where(dao.Role.Columns().Status, consts.StatusEnabled).
		OrderAsc(dao.Role.Columns().Sort).
		Scan(&list)
	return list, err
}

func (s *sRole) checkCodeUnique(ctx context.Context, code string, excludeId int) error {
	m := dao.Role.Ctx(ctx).Where(dao.Role.Columns().Code, code)
	if excludeId > 0 {
		m = m.WhereNot(dao.Role.Columns().Id, excludeId)
	}
	count, err := m.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return errcode.BadRequest("角色编码已存在")
	}
	return nil
}

func (s *sRole) saveMenus(ctx context.Context, tx gdb.TX, roleId int, menuIds []int) error {
	if len(menuIds) == 0 {
		return nil
	}
	data := make([]g.Map, 0, len(menuIds))
	for _, menuId := range menuIds {
		data = append(data, g.Map{
			dao.RoleMenu.Columns().RoleId: roleId,
			dao.RoleMenu.Columns().MenuId: menuId,
		})
	}
	_, err := tx.Model(dao.RoleMenu.Table()).Ctx(ctx).Data(data).Insert()
	return err
}

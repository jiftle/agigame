package menu

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	"agigame/console/backend/internal/consts"
	"agigame/console/backend/internal/dao"
	"agigame/console/backend/internal/model"
	"agigame/console/backend/internal/model/do"
	"agigame/console/backend/internal/model/entity"
	"agigame/console/backend/internal/service"
	"agigame/console/backend/utility/errcode"
)

type sMenu struct{}

func init() {
	service.RegisterMenu(New())
}

// New 创建菜单服务
func New() *sMenu {
	return &sMenu{}
}

// List 菜单列表
func (s *sMenu) List(ctx context.Context, in *model.MenuQueryInput) ([]*entity.Menu, error) {
	m := dao.Menu.Ctx(ctx)
	if in.Title != "" {
		m = m.WhereLike(dao.Menu.Columns().Title, "%"+in.Title+"%")
	}
	if in.Status != nil {
		m = m.Where(dao.Menu.Columns().Status, *in.Status)
	}
	var list []*entity.Menu
	err := m.OrderAsc(dao.Menu.Columns().Sort).OrderAsc(dao.Menu.Columns().Id).Scan(&list)
	return list, err
}

// GetTree 菜单树(带筛选)
func (s *sMenu) GetTree(ctx context.Context, in *model.MenuQueryInput) ([]*model.MenuNode, error) {
	list, err := s.List(ctx, in)
	if err != nil {
		return nil, err
	}
	return buildTree(list, 0), nil
}

// GetById 获取菜单
func (s *sMenu) GetById(ctx context.Context, id int) (*entity.Menu, error) {
	var menu *entity.Menu
	err := dao.Menu.Ctx(ctx).Where(dao.Menu.Columns().Id, id).Scan(&menu)
	if err != nil {
		return nil, err
	}
	if menu == nil {
		return nil, errcode.NotFound("菜单不存在")
	}
	return menu, nil
}

// Create 新增菜单
func (s *sMenu) Create(ctx context.Context, in *model.MenuSaveInput) error {
	_, err := dao.Menu.Ctx(ctx).Data(do.Menu{
		ParentId:  in.ParentId,
		Title:     in.Title,
		Name:      in.Name,
		Path:      in.Path,
		Component: in.Component,
		Icon:      in.Icon,
		Type:      in.Type,
		Perms:     in.Perms,
		Sort:      in.Sort,
		Visible:   in.Visible,
		Status:    in.Status,
		Redirect:  in.Redirect,
		IsFrame:   in.IsFrame,
		IsCache:   in.IsCache,
	}).Insert()
	return err
}

// Update 修改菜单
func (s *sMenu) Update(ctx context.Context, in *model.MenuSaveInput) error {
	if in.Id == in.ParentId {
		return errcode.BadRequest("上级菜单不能是自己")
	}
	_, err := dao.Menu.Ctx(ctx).Where(dao.Menu.Columns().Id, in.Id).Data(do.Menu{
		ParentId:  in.ParentId,
		Title:     in.Title,
		Name:      in.Name,
		Path:      in.Path,
		Component: in.Component,
		Icon:      in.Icon,
		Type:      in.Type,
		Perms:     in.Perms,
		Sort:      in.Sort,
		Visible:   in.Visible,
		Status:    in.Status,
		Redirect:  in.Redirect,
		IsFrame:   in.IsFrame,
		IsCache:   in.IsCache,
	}).Update()
	return err
}

// Delete 删除菜单
func (s *sMenu) Delete(ctx context.Context, id int) error {
	count, err := dao.Menu.Ctx(ctx).Where(dao.Menu.Columns().ParentId, id).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return errcode.BadRequest("存在子菜单，不允许删除")
	}
	roleCount, err := dao.RoleMenu.Ctx(ctx).Where(dao.RoleMenu.Columns().MenuId, id).Count()
	if err != nil {
		return err
	}
	if roleCount > 0 {
		return errcode.BadRequest("菜单已分配角色，不允许删除")
	}
	_, err = dao.Menu.Ctx(ctx).Where(dao.Menu.Columns().Id, id).Delete()
	return err
}

// GetPermsByUserId 获取用户权限标识
func (s *sMenu) GetPermsByUserId(ctx context.Context, userId int) ([]string, error) {
	sql := `SELECT DISTINCT m.perms AS perms
		FROM sys_menu m
		JOIN sys_role_menu rm ON rm.menu_id = m.id
		JOIN sys_role r ON r.id = rm.role_id AND r.status = 1 AND r.deleted_at IS NULL
		JOIN sys_user_role ur ON ur.role_id = rm.role_id
		WHERE ur.user_id = ? AND m.status = 1 AND m.perms <> '' AND m.deleted_at IS NULL`
	result, err := g.DB().GetAll(ctx, sql, userId)
	if err != nil {
		return nil, err
	}
	perms := make([]string, 0, len(result))
	for _, r := range result {
		perms = append(perms, r["perms"].String())
	}
	return perms, nil
}

// GetMenuTreeByUserId 获取用户菜单树
func (s *sMenu) GetMenuTreeByUserId(ctx context.Context, userId int, isSuper bool) ([]*model.MenuNode, error) {
	var list []*entity.Menu
	if isSuper {
		err := dao.Menu.Ctx(ctx).
			Where(dao.Menu.Columns().Status, consts.StatusEnabled).
			WhereIn(dao.Menu.Columns().Type, []string{consts.MenuTypeDir, consts.MenuTypeMenu}).
			OrderAsc(dao.Menu.Columns().Sort).OrderAsc(dao.Menu.Columns().Id).
			Scan(&list)
		if err != nil {
			return nil, err
		}
	} else {
		sql := `SELECT DISTINCT m.*
			FROM sys_menu m
			JOIN sys_role_menu rm ON rm.menu_id = m.id
			JOIN sys_role r ON r.id = rm.role_id AND r.status = 1 AND r.deleted_at IS NULL
			JOIN sys_user_role ur ON ur.role_id = rm.role_id
			WHERE ur.user_id = ? AND m.status = 1 AND m.type IN ('M','C') AND m.deleted_at IS NULL
			ORDER BY m.sort ASC, m.id ASC`
		result, err := g.DB().GetAll(ctx, sql, userId)
		if err != nil {
			return nil, err
		}
		if err = result.Structs(&list); err != nil {
			return nil, err
		}
	}
	return buildTree(list, 0), nil
}

// GetAllTree 获取全部菜单树
func (s *sMenu) GetAllTree(ctx context.Context) ([]*model.MenuNode, error) {
	var list []*entity.Menu
	err := dao.Menu.Ctx(ctx).
		OrderAsc(dao.Menu.Columns().Sort).OrderAsc(dao.Menu.Columns().Id).
		Scan(&list)
	if err != nil {
		return nil, err
	}
	return buildTree(list, 0), nil
}

func buildTree(list []*entity.Menu, parentId int) []*model.MenuNode {
	nodes := make([]*model.MenuNode, 0)
	for _, m := range list {
		if m.ParentId != parentId {
			continue
		}
		node := &model.MenuNode{
			Id:        m.Id,
			ParentId:  m.ParentId,
			Title:     m.Title,
			Name:      m.Name,
			Path:      m.Path,
			Component: m.Component,
			Icon:      m.Icon,
			Type:      m.Type,
			Perms:     m.Perms,
			Sort:      m.Sort,
			Visible:   m.Visible,
			Redirect:  m.Redirect,
			IsFrame:   m.IsFrame,
			IsCache:   m.IsCache,
		}
		node.Children = buildTree(list, m.Id)
		nodes = append(nodes, node)
	}
	return nodes
}

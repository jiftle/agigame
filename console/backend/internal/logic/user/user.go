package user

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"agigame/console/backend/internal/consts"
	"agigame/console/backend/internal/dao"
	"agigame/console/backend/internal/model"
	"agigame/console/backend/internal/model/do"
	"agigame/console/backend/internal/model/entity"
	"agigame/console/backend/internal/service"
	"agigame/console/backend/utility/appcfg"
	"agigame/console/backend/utility/errcode"
	"agigame/console/backend/utility/password"
)

type sUser struct{}

func init() {
	service.RegisterUser(New())
}

// New 创建用户服务
func New() *sUser {
	return &sUser{}
}

// List 用户列表
func (s *sUser) List(ctx context.Context, in *model.UserQueryInput) (total int, list []*entity.User, err error) {
	cols := dao.User.Columns()
	m := dao.User.Ctx(ctx)
	if in.Username != "" {
		m = m.WhereLike(cols.Username, "%"+in.Username+"%")
	}
	if in.Nickname != "" {
		m = m.WhereLike(cols.Nickname, "%"+in.Nickname+"%")
	}
	if in.Phone != "" {
		m = m.WhereLike(cols.Phone, "%"+in.Phone+"%")
	}
	if in.Status != nil {
		m = m.Where(cols.Status, *in.Status)
	}
	if in.DeptId > 0 {
		m = m.Where(cols.DeptId, in.DeptId)
	}
	total, err = m.Count()
	if err != nil || total == 0 {
		return 0, nil, err
	}
	err = m.Page(in.PageNum, in.PageSize).OrderDesc(cols.Id).Scan(&list)
	if err != nil {
		return
	}
	for _, u := range list {
		u.Password = ""
	}
	return
}

// GetById 获取用户
func (s *sUser) GetById(ctx context.Context, id int) (*entity.User, error) {
	var user *entity.User
	if err := dao.User.Ctx(ctx).Where(dao.User.Columns().Id, id).Scan(&user); err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errcode.NotFound("用户不存在")
	}
	return user, nil
}

// GetByUsername 按用户名获取用户
func (s *sUser) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user *entity.User
	if err := dao.User.Ctx(ctx).Where(dao.User.Columns().Username, username).Scan(&user); err != nil {
		return nil, err
	}
	return user, nil
}

// Create 新增用户
func (s *sUser) Create(ctx context.Context, in *model.UserSaveInput) error {
	if exist, err := s.IsUsernameExist(ctx, in.Username, 0); err != nil {
		return err
	} else if exist {
		return errcode.BadRequest("用户名已存在")
	}
	plain := in.Password
	if plain == "" {
		plain, _ = service.Config().GetValue(ctx, "sys.user.initPwd")
	}
	if plain == "" {
		plain = appcfg.DefaultPassword(ctx)
	}
	hash, err := password.Encrypt(plain)
	if err != nil {
		return err
	}
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		id, err := tx.Model(dao.User.Table()).Ctx(ctx).Data(do.User{
			DeptId:   in.DeptId,
			Username: in.Username,
			Password: hash,
			Nickname: in.Nickname,
			Email:    in.Email,
			Phone:    in.Phone,
			Sex:      in.Sex,
			Status:   in.Status,
			Remark:   in.Remark,
		}).InsertAndGetId()
		if err != nil {
			return err
		}
		return s.saveRoles(ctx, tx, int(id), in.RoleIds)
	})
}

// Update 修改用户
func (s *sUser) Update(ctx context.Context, in *model.UserSaveInput) error {
	if in.Id == appcfg.SuperAdminId(ctx) {
		old, err := s.GetById(ctx, in.Id)
		if err != nil {
			return err
		}
		if in.Status != consts.StatusEnabled {
			return errcode.BadRequest("内置超级管理员不允许停用")
		}
		if old.Username != in.Username {
			return errcode.BadRequest("内置超级管理员用户名不允许修改")
		}
	}
	if exist, err := s.IsUsernameExist(ctx, in.Username, in.Id); err != nil {
		return err
	} else if exist {
		return errcode.BadRequest("用户名已存在")
	}
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err := tx.Model(dao.User.Table()).Ctx(ctx).Where(dao.User.Columns().Id, in.Id).Data(do.User{
			DeptId:   in.DeptId,
			Username: in.Username,
			Nickname: in.Nickname,
			Email:    in.Email,
			Phone:    in.Phone,
			Sex:      in.Sex,
			Status:   in.Status,
			Remark:   in.Remark,
		}).Update()
		if err != nil {
			return err
		}
		if _, err = tx.Model(dao.UserRole.Table()).Ctx(ctx).Where(dao.UserRole.Columns().UserId, in.Id).Delete(); err != nil {
			return err
		}
		return s.saveRoles(ctx, tx, in.Id, in.RoleIds)
	})
}

// Delete 删除用户
func (s *sUser) Delete(ctx context.Context, ids []int) error {
	for _, id := range ids {
		if id == appcfg.SuperAdminId(ctx) {
			return errcode.BadRequest("内置超级管理员不允许删除")
		}
	}
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(dao.User.Table()).Ctx(ctx).WhereIn(dao.User.Columns().Id, ids).Delete(); err != nil {
			return err
		}
		_, err := tx.Model(dao.UserRole.Table()).Ctx(ctx).WhereIn(dao.UserRole.Columns().UserId, ids).Delete()
		return err
	})
}

// ResetPwd 重置密码
func (s *sUser) ResetPwd(ctx context.Context, id int, plain string) error {
	if plain == "" {
		plain, _ = service.Config().GetValue(ctx, "sys.user.initPwd")
	}
	if plain == "" {
		plain = appcfg.DefaultPassword(ctx)
	}
	hash, err := password.Encrypt(plain)
	if err != nil {
		return err
	}
	_, err = dao.User.Ctx(ctx).Where(dao.User.Columns().Id, id).Data(do.User{Password: hash}).Update()
	return err
}

// ChangeStatus 修改状态
func (s *sUser) ChangeStatus(ctx context.Context, id, status int) error {
	if id == appcfg.SuperAdminId(ctx) && status == consts.StatusDisabled {
		return errcode.BadRequest("内置超级管理员不允许停用")
	}
	_, err := dao.User.Ctx(ctx).Where(dao.User.Columns().Id, id).Data(do.User{Status: status}).Update()
	return err
}

// UpdateProfile 更新个人资料
func (s *sUser) UpdateProfile(ctx context.Context, id int, in *model.ProfileUpdateInput) error {
	_, err := dao.User.Ctx(ctx).Where(dao.User.Columns().Id, id).Data(do.User{
		Nickname: in.Nickname,
		Email:    in.Email,
		Phone:    in.Phone,
		Avatar:   in.Avatar,
	}).Update()
	return err
}

// GetRoleIds 获取用户角色ID
func (s *sUser) GetRoleIds(ctx context.Context, userId int) ([]int, error) {
	values, err := dao.UserRole.Ctx(ctx).Where(dao.UserRole.Columns().UserId, userId).Array(dao.UserRole.Columns().RoleId)
	if err != nil {
		return nil, err
	}
	ids := make([]int, 0, len(values))
	for _, v := range values {
		ids = append(ids, v.Int())
	}
	return ids, nil
}

// IsUsernameExist 用户名是否存在
func (s *sUser) IsUsernameExist(ctx context.Context, username string, excludeId int) (bool, error) {
	m := dao.User.Ctx(ctx).Where(dao.User.Columns().Username, username)
	if excludeId > 0 {
		m = m.WhereNot(dao.User.Columns().Id, excludeId)
	}
	count, err := m.Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// UpdateLoginInfo 更新登录信息 (供认证模块调用)
func (s *sUser) UpdateLoginInfo(ctx context.Context, userId int, ip string) error {
	_, err := dao.User.Ctx(ctx).Where(dao.User.Columns().Id, userId).Data(g.Map{
		dao.User.Columns().LoginIp:   ip,
		dao.User.Columns().LoginDate: gtime.Now(),
	}).Update()
	return err
}

func (s *sUser) saveRoles(ctx context.Context, tx gdb.TX, userId int, roleIds []int) error {
	if len(roleIds) == 0 {
		return nil
	}
	data := make([]g.Map, 0, len(roleIds))
	for _, roleId := range roleIds {
		data = append(data, g.Map{
			dao.UserRole.Columns().UserId: userId,
			dao.UserRole.Columns().RoleId: roleId,
		})
	}
	_, err := tx.Model(dao.UserRole.Table()).Ctx(ctx).Data(data).Insert()
	return err
}

package system

import (
	"context"

	api "agigame/console/backend/api/v1/system"
	"agigame/console/backend/internal/model"
	"agigame/console/backend/internal/service"
	"agigame/console/backend/utility/authz"
)

type cUser struct{}

// NewUser 创建用户控制器
func NewUser() *cUser {
	return &cUser{}
}

// List 用户列表
func (c *cUser) List(ctx context.Context, req *api.UserListReq) (res *api.UserListRes, err error) {
	if err = authz.Check(ctx, "system:user:list"); err != nil {
		return nil, err
	}
	total, list, err := service.User().List(ctx, &model.UserQueryInput{
		PageInput: model.PageInput{PageNum: req.PageNum, PageSize: req.PageSize},
		Username:  req.Username,
		Nickname:  req.Nickname,
		Phone:     req.Phone,
		Status:    req.Status,
		DeptId:    req.DeptId,
	})
	if err != nil {
		return nil, err
	}
	return &api.UserListRes{Total: total, List: list}, nil
}

// Get 用户详情
func (c *cUser) Get(ctx context.Context, req *api.UserGetReq) (res *api.UserGetRes, err error) {
	if err = authz.Check(ctx, "system:user:list"); err != nil {
		return nil, err
	}
	user, err := service.User().GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	roleIds, err := service.User().GetRoleIds(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	user.Password = ""
	return &api.UserGetRes{User: user, RoleIds: roleIds}, nil
}

// Create 新增用户
func (c *cUser) Create(ctx context.Context, req *api.UserCreateReq) (res *api.UserCreateRes, err error) {
	if err = authz.Check(ctx, "system:user:add"); err != nil {
		return nil, err
	}
	err = service.User().Create(ctx, &model.UserSaveInput{
		DeptId:   req.DeptId,
		Username: req.Username,
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Sex:      req.Sex,
		Status:   req.Status,
		Remark:   req.Remark,
		Password: req.Password,
		RoleIds:  req.RoleIds,
	})
	if err != nil {
		return nil, err
	}
	return &api.UserCreateRes{}, nil
}

// Update 修改用户
func (c *cUser) Update(ctx context.Context, req *api.UserUpdateReq) (res *api.UserUpdateRes, err error) {
	if err = authz.Check(ctx, "system:user:edit"); err != nil {
		return nil, err
	}
	err = service.User().Update(ctx, &model.UserSaveInput{
		Id:       req.Id,
		DeptId:   req.DeptId,
		Username: req.Username,
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Sex:      req.Sex,
		Status:   req.Status,
		Remark:   req.Remark,
		RoleIds:  req.RoleIds,
	})
	if err != nil {
		return nil, err
	}
	return &api.UserUpdateRes{}, nil
}

// Delete 删除用户
func (c *cUser) Delete(ctx context.Context, req *api.UserDeleteReq) (res *api.UserDeleteRes, err error) {
	if err = authz.Check(ctx, "system:user:remove"); err != nil {
		return nil, err
	}
	if err = service.User().Delete(ctx, req.Ids); err != nil {
		return nil, err
	}
	return &api.UserDeleteRes{}, nil
}

// ResetPwd 重置密码
func (c *cUser) ResetPwd(ctx context.Context, req *api.UserResetPwdReq) (res *api.UserResetPwdRes, err error) {
	if err = authz.Check(ctx, "system:user:resetPwd"); err != nil {
		return nil, err
	}
	if err = service.User().ResetPwd(ctx, req.Id, req.Password); err != nil {
		return nil, err
	}
	return &api.UserResetPwdRes{}, nil
}

// Status 修改用户状态
func (c *cUser) Status(ctx context.Context, req *api.UserStatusReq) (res *api.UserStatusRes, err error) {
	if err = authz.Check(ctx, "system:user:edit"); err != nil {
		return nil, err
	}
	if err = service.User().ChangeStatus(ctx, req.Id, req.Status); err != nil {
		return nil, err
	}
	return &api.UserStatusRes{}, nil
}

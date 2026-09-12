package system

import (
	"context"

	api "agigame/console/backend/api/v1/system"
	"agigame/console/backend/internal/model"
	"agigame/console/backend/internal/service"
	"agigame/console/backend/utility/authz"
)

type cDept struct{}

// NewDept 创建部门控制器
func NewDept() *cDept {
	return &cDept{}
}

// List 部门列表
func (c *cDept) List(ctx context.Context, req *api.DeptListReq) (res *api.DeptListRes, err error) {
	if err = authz.CheckAny(ctx, "system:dept:list", "system:user:list"); err != nil {
		return nil, err
	}
	list, err := service.Dept().Tree(ctx)
	if err != nil {
		return nil, err
	}
	return &api.DeptListRes{List: list}, nil
}

// Get 部门详情
func (c *cDept) Get(ctx context.Context, req *api.DeptGetReq) (res *api.DeptGetRes, err error) {
	if err = authz.Check(ctx, "system:dept:list"); err != nil {
		return nil, err
	}
	dept, err := service.Dept().GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &api.DeptGetRes{Dept: dept}, nil
}

// Create 新增部门
func (c *cDept) Create(ctx context.Context, req *api.DeptCreateReq) (res *api.DeptCreateRes, err error) {
	if err = authz.Check(ctx, "system:dept:add"); err != nil {
		return nil, err
	}
	err = service.Dept().Create(ctx, &model.DeptSaveInput{
		ParentId: req.ParentId,
		Name:     req.Name,
		Leader:   req.Leader,
		Phone:    req.Phone,
		Email:    req.Email,
		Sort:     req.Sort,
		Status:   req.Status,
		Remark:   req.Remark,
	})
	if err != nil {
		return nil, err
	}
	return &api.DeptCreateRes{}, nil
}

// Update 修改部门
func (c *cDept) Update(ctx context.Context, req *api.DeptUpdateReq) (res *api.DeptUpdateRes, err error) {
	if err = authz.Check(ctx, "system:dept:edit"); err != nil {
		return nil, err
	}
	err = service.Dept().Update(ctx, &model.DeptSaveInput{
		Id:       req.Id,
		ParentId: req.ParentId,
		Name:     req.Name,
		Leader:   req.Leader,
		Phone:    req.Phone,
		Email:    req.Email,
		Sort:     req.Sort,
		Status:   req.Status,
		Remark:   req.Remark,
	})
	if err != nil {
		return nil, err
	}
	return &api.DeptUpdateRes{}, nil
}

// Delete 删除部门
func (c *cDept) Delete(ctx context.Context, req *api.DeptDeleteReq) (res *api.DeptDeleteRes, err error) {
	if err = authz.Check(ctx, "system:dept:remove"); err != nil {
		return nil, err
	}
	if err = service.Dept().Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &api.DeptDeleteRes{}, nil
}

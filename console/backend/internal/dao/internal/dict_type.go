// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-09-10 22:09:07
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DictTypeDao is the data access object for the table sys_dict_type.
type DictTypeDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  DictTypeColumns    // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// DictTypeColumns defines and stores column names for the table sys_dict_type.
type DictTypeColumns struct {
	Id        string //
	Name      string //
	Type      string //
	Status    string //
	Remark    string //
	CreatedAt string //
	UpdatedAt string //
	DeletedAt string //
}

// dictTypeColumns holds the columns for the table sys_dict_type.
var dictTypeColumns = DictTypeColumns{
	Id:        "id",
	Name:      "name",
	Type:      "type",
	Status:    "status",
	Remark:    "remark",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
	DeletedAt: "deleted_at",
}

// NewDictTypeDao creates and returns a new DAO object for table data access.
func NewDictTypeDao(handlers ...gdb.ModelHandler) *DictTypeDao {
	return &DictTypeDao{
		group:    "default",
		table:    "sys_dict_type",
		columns:  dictTypeColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DictTypeDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DictTypeDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DictTypeDao) Columns() DictTypeColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DictTypeDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DictTypeDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *DictTypeDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}

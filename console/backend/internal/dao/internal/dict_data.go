// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-09-10 22:09:07
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DictDataDao is the data access object for the table sys_dict_data.
type DictDataDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  DictDataColumns    // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// DictDataColumns defines and stores column names for the table sys_dict_data.
type DictDataColumns struct {
	Id        string //
	DictSort  string //
	DictLabel string //
	DictValue string //
	DictType  string //
	IsDefault string //
	Status    string //
	Remark    string //
	CreatedAt string //
	UpdatedAt string //
	DeletedAt string //
}

// dictDataColumns holds the columns for the table sys_dict_data.
var dictDataColumns = DictDataColumns{
	Id:        "id",
	DictSort:  "dict_sort",
	DictLabel: "dict_label",
	DictValue: "dict_value",
	DictType:  "dict_type",
	IsDefault: "is_default",
	Status:    "status",
	Remark:    "remark",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
	DeletedAt: "deleted_at",
}

// NewDictDataDao creates and returns a new DAO object for table data access.
func NewDictDataDao(handlers ...gdb.ModelHandler) *DictDataDao {
	return &DictDataDao{
		group:    "default",
		table:    "sys_dict_data",
		columns:  dictDataColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DictDataDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DictDataDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DictDataDao) Columns() DictDataColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DictDataDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DictDataDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *DictDataDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}

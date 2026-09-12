// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-09-10 22:09:07
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// OperLogDao is the data access object for the table sys_oper_log.
type OperLogDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  OperLogColumns     // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// OperLogColumns defines and stores column names for the table sys_oper_log.
type OperLogColumns struct {
	Id            string //
	Title         string //
	BusinessType  string //
	Method        string //
	RequestMethod string //
	OperName      string //
	OperUrl       string //
	OperIp        string //
	OperParam     string //
	JsonResult    string //
	Status        string //
	ErrorMsg      string //
	Cost          string //
	OperTime      string //
}

// operLogColumns holds the columns for the table sys_oper_log.
var operLogColumns = OperLogColumns{
	Id:            "id",
	Title:         "title",
	BusinessType:  "business_type",
	Method:        "method",
	RequestMethod: "request_method",
	OperName:      "oper_name",
	OperUrl:       "oper_url",
	OperIp:        "oper_ip",
	OperParam:     "oper_param",
	JsonResult:    "json_result",
	Status:        "status",
	ErrorMsg:      "error_msg",
	Cost:          "cost",
	OperTime:      "oper_time",
}

// NewOperLogDao creates and returns a new DAO object for table data access.
func NewOperLogDao(handlers ...gdb.ModelHandler) *OperLogDao {
	return &OperLogDao{
		group:    "default",
		table:    "sys_oper_log",
		columns:  operLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *OperLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *OperLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *OperLogDao) Columns() OperLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *OperLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *OperLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *OperLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}

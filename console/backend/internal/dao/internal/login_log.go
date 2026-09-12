// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-09-10 22:09:07
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// LoginLogDao is the data access object for the table sys_login_log.
type LoginLogDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  LoginLogColumns    // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// LoginLogColumns defines and stores column names for the table sys_login_log.
type LoginLogColumns struct {
	Id        string //
	Username  string //
	Ip        string //
	Location  string //
	Browser   string //
	Os        string //
	Status    string //
	Msg       string //
	LoginTime string //
}

// loginLogColumns holds the columns for the table sys_login_log.
var loginLogColumns = LoginLogColumns{
	Id:        "id",
	Username:  "username",
	Ip:        "ip",
	Location:  "location",
	Browser:   "browser",
	Os:        "os",
	Status:    "status",
	Msg:       "msg",
	LoginTime: "login_time",
}

// NewLoginLogDao creates and returns a new DAO object for table data access.
func NewLoginLogDao(handlers ...gdb.ModelHandler) *LoginLogDao {
	return &LoginLogDao{
		group:    "default",
		table:    "sys_login_log",
		columns:  loginLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *LoginLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *LoginLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *LoginLogDao) Columns() LoginLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *LoginLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *LoginLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *LoginLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}

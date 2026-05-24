// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

type UserAccountsDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  UserAccountsColumns    // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

type UserAccountsColumns struct {
	Id            string
	UserId        string
	Balance       string
	FrozenAmount  string
	TotalRecharge string
	TotalConsume  string
	Version       string
	CreatedAt     string
	UpdatedAt     string
}

var userAccountsColumns = UserAccountsColumns{
	Id:            "id",
	UserId:        "user_id",
	Balance:       "balance",
	FrozenAmount:  "frozen_amount",
	TotalRecharge: "total_recharge",
	TotalConsume:  "total_consume",
	Version:       "version",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
}

func NewUserAccountsDao(handlers ...gdb.ModelHandler) *UserAccountsDao {
	return &UserAccountsDao{
		group:    "default",
		table:    "user_accounts",
		columns:  userAccountsColumns,
		handlers: handlers,
	}
}

func (dao *UserAccountsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

func (dao *UserAccountsDao) Table() string {
	return dao.table
}

func (dao *UserAccountsDao) Columns() UserAccountsColumns {
	return dao.columns
}

func (dao *UserAccountsDao) Group() string {
	return dao.group
}

func (dao *UserAccountsDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

func (dao *UserAccountsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}

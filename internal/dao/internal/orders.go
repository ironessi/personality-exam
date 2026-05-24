// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// OrdersDao is the data access object for the table orders.
type OrdersDao struct {
	table    string
	group    string
	columns  OrdersColumns
	handlers []gdb.ModelHandler
}

// OrdersColumns defines and stores column names for the table orders.
type OrdersColumns struct {
	Id           string
	UserId       string
	TestTypeId   string
	TestRecordId string
	OrderNo      string
	TradeNo      string
	Amount       string
	Status       string
	PayType      string
	PayTime      string
	CreatedAt    string
	UpdatedAt    string
}

var ordersColumns = OrdersColumns{
	Id:           "id",
	UserId:       "user_id",
	TestTypeId:   "test_type_id",
	TestRecordId: "test_record_id",
	OrderNo:      "order_no",
	TradeNo:      "trade_no",
	Amount:       "amount",
	Status:       "status",
	PayType:      "pay_type",
	PayTime:      "pay_time",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
}

func NewOrdersDao(handlers ...gdb.ModelHandler) *OrdersDao {
	return &OrdersDao{
		group:    "default",
		table:    "orders",
		columns:  ordersColumns,
		handlers: handlers,
	}
}

func (dao *OrdersDao) DB() gdb.DB {
	return g.DB(dao.group)
}

func (dao *OrdersDao) Table() string {
	return dao.table
}

func (dao *OrdersDao) Columns() OrdersColumns {
	return dao.columns
}

func (dao *OrdersDao) Group() string {
	return dao.group
}

func (dao *OrdersDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

func (dao *OrdersDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}

// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserAnswersDao is the data access object for the table user_answers.
type UserAnswersDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  UserAnswersColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// UserAnswersColumns defines and stores column names for the table user_answers.
type UserAnswersColumns struct {
	Id          string // 答案ID，自增主键
	UserId      string // 关联用户ID
	TestTypeId  string // 测试类型ID
	QuestionId  string // 问题ID
	OptionId    string // 选择的选项ID
	CreatedTime string // 创建时间
	Answer      string // 用户作答选择
	DeletedTime string // 软删除标记
}

// userAnswersColumns holds the columns for the table user_answers.
var userAnswersColumns = UserAnswersColumns{
	Id:          "id",
	UserId:      "user_id",
	TestTypeId:  "test_type_id",
	QuestionId:  "question_id",
	OptionId:    "option_id",
	CreatedTime: "created_time",
	Answer:      "answer",
	DeletedTime: "deleted_time",
}

// NewUserAnswersDao creates and returns a new DAO object for table data access.
func NewUserAnswersDao(handlers ...gdb.ModelHandler) *UserAnswersDao {
	return &UserAnswersDao{
		group:    "default",
		table:    "user_answers",
		columns:  userAnswersColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UserAnswersDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UserAnswersDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UserAnswersDao) Columns() UserAnswersColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UserAnswersDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UserAnswersDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UserAnswersDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}

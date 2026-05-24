// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserAnswers is the golang structure of table user_answers for DAO operations like Where/Data.
type UserAnswers struct {
	g.Meta      `orm:"table:user_answers, do:true"`
	Id          any         // 答案ID，自增主键
	UserId      any         // 关联用户ID
	TestTypeId  any         // 测试类型ID
	QuestionId  any         // 问题ID
	OptionId    any         // 选择的选项ID
	CreatedTime *gtime.Time // 创建时间
	Answer      any         // 用户作答选择
	DeletedTime *gtime.Time // 软删除标记
}

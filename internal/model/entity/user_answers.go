// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserAnswers is the golang structure for table user_answers.
type UserAnswers struct {
	Id          int64       `json:"id"          orm:"id"           description:"答案ID，自增主键"` // 答案ID，自增主键
	UserId      int64       `json:"userId"      orm:"user_id"      description:"关联用户ID"`    // 关联用户ID
	TestTypeId  int64       `json:"testTypeId"  orm:"test_type_id" description:"测试类型ID"`    // 测试类型ID
	QuestionId  int64       `json:"questionId"  orm:"question_id"  description:"问题ID"`      // 问题ID
	OptionId    int64       `json:"optionId"    orm:"option_id"    description:"选择的选项ID"`   // 选择的选项ID
	CreatedTime *gtime.Time `json:"createdTime" orm:"created_time" description:"创建时间"`      // 创建时间
	Answer      string      `json:"answer"      orm:"answer"       description:"用户作答选择"`    // 用户作答选择
	DeletedTime *gtime.Time `json:"deletedTime" orm:"deleted_time" description:"软删除标记"`     // 软删除标记
}

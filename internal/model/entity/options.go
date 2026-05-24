// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Options is the golang structure for table options.
type Options struct {
	Id          int64       `json:"id"          orm:"id"           description:"选项ID，自增主键"` // 选项ID，自增主键
	QuestionId  int64       `json:"questionId"  orm:"question_id"  description:"关联的问题ID"`   // 关联的问题ID
	Content     string      `json:"content"     orm:"content"      description:"选项内容文本"`    // 选项内容文本
	Score       int         `json:"score"       orm:"score"        description:"选项对应的分数"`   // 选项对应的分数
	CreatedTime *gtime.Time `json:"createdTime" orm:"created_time" description:"创建时间"`      // 创建时间
	UpdatedTime *gtime.Time `json:"updatedTime" orm:"updated_time" description:"更新时间"`      // 更新时间
}

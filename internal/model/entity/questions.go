// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Questions is the golang structure for table questions.
type Questions struct {
	Id          int64       `json:"id"          orm:"id"           description:"问题ID，自增主键"` // 问题ID，自增主键
	TestTypeId  int64       `json:"testTypeId"  orm:"test_type_id" description:"关联的测试类型ID"` // 关联的测试类型ID
	Content     string      `json:"content"     orm:"content"      description:"问题内容文本"`    // 问题内容文本
	Sequence    int         `json:"sequence"    orm:"sequence"     description:"问题顺序号"`     // 问题顺序号
	CreatedTime *gtime.Time `json:"createdTime" orm:"created_time" description:"创建时间"`      // 创建时间
	UpdatedTime *gtime.Time `json:"updatedTime" orm:"updated_time" description:"更新时间"`      // 更新时间
}

// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TestResults is the golang structure for table test_results.
type TestResults struct {
	Id          int64       `json:"id"          orm:"id"           description:"结果ID，自增主键"` // 结果ID，自增主键
	TestTypeId  int64       `json:"testTypeId"  orm:"test_type_id" description:"关联的测试类型ID"` // 关联的测试类型ID
	Description string      `json:"description" orm:"description"  description:"结果详细描述"`    // 结果详细描述
	CreateTime  *gtime.Time `json:"createTime"  orm:"create_time"  description:"创建时间"`      // 创建时间
	UpdateTime  *gtime.Time `json:"updateTime"  orm:"update_time"  description:"更新时间"`      // 更新时间
}

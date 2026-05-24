// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TestRecords is the golang structure for table test_records.
type TestRecords struct {
	Id          int64       `json:"id"          orm:"id"           description:"记录ID，自增主键"`          // 记录ID，自增主键
	UserId      int64       `json:"userId"      orm:"user_id"      description:"关联用户ID"`             // 关联用户ID
	TestTypeId  int         `json:"testTypeId"  orm:"test_type_id" description:"测试类型ID"`             // 测试类型ID
	Status      string      `json:"status"      orm:"status"       description:"测试状态(进行中/已完成)"`      // 测试状态(进行中/已完成)
	StartTime   *gtime.Time `json:"startTime"   orm:"start_time"   description:"测试开始时间"`             // 测试开始时间
	EndTime     *gtime.Time `json:"endTime"     orm:"end_time"     description:"测试结束时间"`             // 测试结束时间
	CreatedTime *gtime.Time `json:"createdTime" orm:"created_time" description:"记录创建时间"`             // 记录创建时间
	UpdatedTime *gtime.Time `json:"updatedTime" orm:"updated_time" description:"记录更新时间"`             // 记录更新时间
	Deleted     int         `json:"deleted"     orm:"deleted"      description:"软删除标记(0-未删除,1-已删除)"` // 软删除标记(0-未删除,1-已删除)
}

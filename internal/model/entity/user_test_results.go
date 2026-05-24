// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserTestResults is the golang structure for table user_test_results.
type UserTestResults struct {
	Id          int64       `json:"id"          orm:"id"           description:"结果ID，自增主键"`           // 结果ID，自增主键
	UserId      int64       `json:"userId"      orm:"user_id"      description:"用户ID"`                // 用户ID
	TestTypeId  int64       `json:"testTypeId"  orm:"test_type_id" description:"测试类型ID"`              // 测试类型ID
	ResultData  string      `json:"resultData"  orm:"result_data"  description:"结果数据(包含总分、分析结果、建议等)"` // 结果数据(包含总分、分析结果、建议等)
	TestDate    *gtime.Time `json:"testDate"    orm:"test_date"    description:"测试时间"`                // 测试时间
	CreatedTime *gtime.Time `json:"createdTime" orm:"created_time" description:"创建时间"`                // 创建时间
}

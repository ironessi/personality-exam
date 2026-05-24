// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserTestHistory is the golang structure for table user_test_history.
type UserTestHistory struct {
	Id          int64       `json:"id"          orm:"id"           description:"历史记录ID，自增主键"` // 历史记录ID，自增主键
	UserId      int64       `json:"userId"      orm:"user_id"      description:"关联用户ID"`      // 关联用户ID
	TestTypeId  int64       `json:"testTypeId"  orm:"test_type_id" description:"测试类型ID"`      // 测试类型ID
	ResultId    int64       `json:"resultId"    orm:"result_id"    description:"关联的测试结果ID"`   // 关联的测试结果ID
	ResultType  string      `json:"resultType"  orm:"result_type"  description:"结果类型代码"`      // 结果类型代码
	ResultTitle string      `json:"resultTitle" orm:"result_title" description:"结果标题"`        // 结果标题
	WingType    string      `json:"wingType"    orm:"wing_type"    description:"翼型(九型人格专用)"`  // 翼型(九型人格专用)
	CreateTime  *gtime.Time `json:"createTime"  orm:"create_time"  description:"创建时间"`        // 创建时间
}

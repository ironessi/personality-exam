// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserTestResults is the golang structure of table user_test_results for DAO operations like Where/Data.
type UserTestResults struct {
	g.Meta      `orm:"table:user_test_results, do:true"`
	Id          any         // 结果ID，自增主键
	UserId      any         // 用户ID
	TestTypeId  any         // 测试类型ID
	ResultData  any         // 结果数据(包含总分、分析结果、建议等)
	TestDate    *gtime.Time // 测试时间
	CreatedTime *gtime.Time // 创建时间
}

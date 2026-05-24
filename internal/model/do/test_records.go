// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// TestRecords is the golang structure of table test_records for DAO operations like Where/Data.
type TestRecords struct {
	g.Meta      `orm:"table:test_records, do:true"`
	Id          any         // 记录ID，自增主键
	UserId      any         // 关联用户ID
	TestTypeId  any         // 测试类型ID
	Status      any         // 测试状态(进行中/已完成)
	StartTime   *gtime.Time // 测试开始时间
	EndTime     *gtime.Time // 测试结束时间
	CreatedTime *gtime.Time // 记录创建时间
	UpdatedTime *gtime.Time // 记录更新时间
	Deleted     any         // 软删除标记(0-未删除,1-已删除)
}

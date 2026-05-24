// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// TestResults is the golang structure of table test_results for DAO operations like Where/Data.
type TestResults struct {
	g.Meta      `orm:"table:test_results, do:true"`
	Id          any         // 结果ID，自增主键
	TestTypeId  any         // 关联的测试类型ID
	Description any         // 结果详细描述
	CreateTime  *gtime.Time // 创建时间
	UpdateTime  *gtime.Time // 更新时间
}

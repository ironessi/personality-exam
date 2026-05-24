// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserTestHistory is the golang structure of table user_test_history for DAO operations like Where/Data.
type UserTestHistory struct {
	g.Meta      `orm:"table:user_test_history, do:true"`
	Id          any         // 历史记录ID，自增主键
	UserId      any         // 关联用户ID
	TestTypeId  any         // 测试类型ID
	ResultId    any         // 关联的测试结果ID
	ResultType  any         // 结果类型代码
	ResultTitle any         // 结果标题
	WingType    any         // 翼型(九型人格专用)
	CreateTime  *gtime.Time // 创建时间
}

// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Questions is the golang structure of table questions for DAO operations like Where/Data.
type Questions struct {
	g.Meta      `orm:"table:questions, do:true"`
	Id          any         // 问题ID，自增主键
	TestTypeId  any         // 关联的测试类型ID
	Content     any         // 问题内容文本
	Sequence    any         // 问题顺序号
	CreatedTime *gtime.Time // 创建时间
	UpdatedTime *gtime.Time // 更新时间
}

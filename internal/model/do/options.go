// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Options is the golang structure of table options for DAO operations like Where/Data.
type Options struct {
	g.Meta      `orm:"table:options, do:true"`
	Id          any         // 选项ID，自增主键
	QuestionId  any         // 关联的问题ID
	Content     any         // 选项内容文本
	Score       any         // 选项对应的分数
	CreatedTime *gtime.Time // 创建时间
	UpdatedTime *gtime.Time // 更新时间
}

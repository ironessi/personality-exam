// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// IdempotencyTokens is the golang structure of table idempotency_tokens for DAO operations like Where/Data.
type IdempotencyTokens struct {
	g.Meta       `orm:"table:idempotency_tokens, do:true"`
	Token        any         //
	BusinessType any         // 业务类型
	BusinessId   any         // 业务单号
	Status       any         // 0-处理中 1-成功 2-失败
	ExpireAt     any         //
	CreatedAt    *gtime.Time //
}

// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// VerificationCodes is the golang structure of table verification_codes for DAO operations like Where/Data.
type VerificationCodes struct {
	g.Meta    `orm:"table:verification_codes, do:true"`
	Id        any         // 验证码ID
	Mobile    any         // 手机号码
	Code      any         // 验证码
	Type      any         // 业务类型(LOGIN/REGISTER/...)
	Status    any         // 状态(0-未使用, 1-已使用, 2-已过期)
	ExpireAt  any         // 过期时间
	ClientIp  any         // 客户端IP
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
}

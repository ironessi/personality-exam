// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// VerificationCodes is the golang structure for table verification_codes.
type VerificationCodes struct {
	Id        int64       `json:"id"        orm:"id"         description:"验证码ID"`                    // 验证码ID
	Mobile    string      `json:"mobile"    orm:"mobile"     description:"手机号码"`                     // 手机号码
	Code      string      `json:"code"      orm:"code"       description:"验证码"`                      // 验证码
	Type      string      `json:"type"      orm:"type"       description:"业务类型(LOGIN/REGISTER/...)"` // 业务类型(LOGIN/REGISTER/...)
	Status    int         `json:"status"    orm:"status"     description:"状态(0-未使用, 1-已使用, 2-已过期)"`  // 状态(0-未使用, 1-已使用, 2-已过期)
	ExpireAt  int64       `json:"expireAt"  orm:"expire_at"  description:"过期时间"`                     // 过期时间
	ClientIp  string      `json:"clientIp"  orm:"client_ip"  description:"客户端IP"`                    // 客户端IP
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`                     // 创建时间
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:"更新时间"`                     // 更新时间
}

// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// IdempotencyTokens is the golang structure for table idempotency_tokens.
type IdempotencyTokens struct {
	Token        string      `json:"token"        orm:"token"         description:""`                //
	BusinessType string      `json:"businessType" orm:"business_type" description:"业务类型"`            // 业务类型
	BusinessId   string      `json:"businessId"   orm:"business_id"   description:"业务单号"`            // 业务单号
	Status       int         `json:"status"       orm:"status"        description:"0-处理中 1-成功 2-失败"` // 0-处理中 1-成功 2-失败
	ExpireAt     int64       `json:"expireAt"     orm:"expire_at"     description:""`                //
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    description:""`                //
}

// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package user

import (
	"context"

	"PersonalityTest/api/user/v1"
)

type IUserV1 interface {
	GetUserInfo(ctx context.Context, req *v1.GetUserInfoReq) (res *v1.GetUserInfoRes, err error)
	UpdateUserInfo(ctx context.Context, req *v1.UpdateUserInfoReq) (res *v1.UpdateUserInfoRes, err error)
}

package user

import (
	"context"

	"PersonalityTest/api/user/v1"
)

func (c *ControllerV1) GetUserInfo(ctx context.Context, req *v1.GetUserInfoReq) (res *v1.GetUserInfoRes, err error) {
	user, err := c.user.GetUserInfo(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserInfoRes{
		Id:            user.Id,
		Username:      user.Username,
		Nickname:      user.Nickname,
		Gender:        user.Gender,
		LastLoginTime: user.LastLoginTime.Time,
	}, nil
}
func (c *ControllerV1) UpdateUserInfo(ctx context.Context, req *v1.UpdateUserInfoReq) (res *v1.UpdateUserInfoRes, err error) {
	err = c.user.UpdateUserInfo(ctx, req.Id, req.Nickname, req.Gender)
	if err != nil {
		return nil, err
	}
	return &v1.UpdateUserInfoRes{}, err
}

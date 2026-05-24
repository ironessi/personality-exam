package auth

import (
	"context"

	v1 "PersonalityTest/api/auth/v1"
)

func (c *ControllerV1) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	token, err := c.auth.Login(ctx, req.Account, req.Password)
	if err != nil {
		return nil, err
	}
	return &v1.LoginRes{
		Token: token,
	}, nil
}
func (c *ControllerV1) Register(ctx context.Context, req *v1.RegisterReq) (res *v1.RegisterRes, err error) {
	return nil, c.auth.Register(ctx, req.Mobile, req.Password, req.Code)
}

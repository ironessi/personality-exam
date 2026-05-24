package sms

import (
	"PersonalityTest/internal/model"
	"PersonalityTest/internal/service"
	"context"
)

type ControllerV1 struct{}

func NewV1() *ControllerV1 {
	return &ControllerV1{}
}

func (c *ControllerV1) Send(ctx context.Context, req *model.SmsSendReq) (res *model.SmsSendRes, err error) {
	res, err = service.Sms().SendCode(ctx, *req)
	return
}

func (c *ControllerV1) Verify(ctx context.Context, req *model.SmsVerifyReq) (res *model.SmsVerifyRes, err error) {
	res, err = service.Sms().VerifyCode(ctx, *req)
	return
}

func (c *ControllerV1) Status(ctx context.Context, req *model.SmsStatusReq) (res *model.SmsStatusRes, err error) {
	res, err = service.Sms().GetStatus(ctx, req.Mobile)
	return
}

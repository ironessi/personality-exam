package service

import (
	"PersonalityTest/internal/model"
	"context"
)

type (
	ISms interface {
		SendCode(ctx context.Context, in model.SmsSendReq) (res *model.SmsSendRes, err error)
		VerifyCode(ctx context.Context, in model.SmsVerifyReq) (res *model.SmsVerifyRes, err error)
		CheckCode(ctx context.Context, mobile, code, typeVal string) error
		GetStatus(ctx context.Context, mobile string) (res *model.SmsStatusRes, err error)
	}
)

var (
	localSms ISms
)

func Sms() ISms {
	if localSms == nil {
		panic("implement not found for interface ISms, forgot register?")
	}
	return localSms
}

func RegisterSms(i ISms) {
	localSms = i
}

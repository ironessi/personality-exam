package sms

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/logic/user"
	"PersonalityTest/internal/model"
	"PersonalityTest/internal/model/do"
	"PersonalityTest/internal/model/entity"
	"PersonalityTest/internal/service"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dypnsapi "github.com/alibabacloud-go/dypnsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/golang-jwt/jwt/v5"
)

type sSms struct{}

func init() {
	service.RegisterSms(New())
}

func New() *sSms {
	return &sSms{}
}

// Redis Keys
const (
	DriverSmsCode        = "sms:code:%s:%s" // sms:code:{mobile}:{type}
	DriverSmsLastSend    = "sms:last:%s"    // sms:last:{mobile}
	DriverSmsCountMobile = "sms:cnt:m:%s"   // sms:cnt:m:{mobile} (Daily)
	DriverSmsCountIp     = "sms:cnt:i:%s"   // sms:cnt:i:{ip} (Hourly)

	smsCodeExpireSeconds = 300
	smsSendCooldown      = 60
	smsDailyMobileLimit  = 20
	smsHourlyIpLimit     = 10
)

// getRedisGroup 返回 Redis 组名（固定使用 default）
func (s *sSms) getRedisGroup(ctx context.Context) string {
	return "default"
}

func (s *sSms) SendCode(ctx context.Context, in model.SmsSendReq) (res *model.SmsSendRes, err error) {
	// 1. Check Rate Limits
	clientIp := g.RequestFromCtx(ctx).GetClientIp()
	if err = s.checkRateLimit(ctx, in.Mobile, clientIp); err != nil {
		return nil, err
	}

	// 2. Generate Code
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	// 3. Save to Redis & DB
	// DB Audit
	_, err = dao.VerificationCodes.Ctx(ctx).Data(do.VerificationCodes{
		Mobile:   in.Mobile,
		Code:     code,
		Type:     in.Type,
		Status:   0, // Unused
		ExpireAt: gtime.Timestamp() + smsCodeExpireSeconds,
		ClientIp: clientIp,
	}).Insert()
	if err != nil {
		return nil, err
	}

	// Redis Cache (5 minutes)
	redisKey := fmt.Sprintf(DriverSmsCode, in.Mobile, in.Type)
	_, err = g.Redis(s.getRedisGroup(ctx)).Do(ctx, "SETEX", redisKey, smsCodeExpireSeconds, code)
	if err != nil {
		// Redis 失败不影响主流程，仅记录日志
		g.Log().Warningf(ctx, "Redis cache failed (ignored): %v", err)
		err = nil // 清除错误，继续执行
	}

	// 4. Send SMS via Alibaba Cloud
	// If credentials are missing, we might want to mock it (or just fail).
	// For now, we try to send if config exists, else log and return success (dev mode).
	ak := g.Cfg().MustGet(ctx, "sms.accessKeyId").String()
	sk := g.Cfg().MustGet(ctx, "sms.accessKeySecret").String()

	if ak == "" || sk == "" {
		g.Log().Warningf(ctx, "SMS credentials not configured. Mock sending code %s to %s", code, in.Mobile)
		// Mock success for development if no creds
		return &model.SmsSendRes{
			RequestId: "MOCK-REQUEST-ID",
			ExpireIn:  smsCodeExpireSeconds,
		}, nil
	}

	// Real Send
	err = s.sendAliSms(ctx, ak, sk, in.Mobile, code)
	if err != nil {
		// SMS 发送失败，但验证码已保存到 Redis/DB，记录警告但不阻止流程
		// 开发环境可以从数据库查看验证码进行测试
		g.Log().Warningf(ctx, "SMS send failed (code saved): %v", err)
	}

	// Update Rate Limit Counters (无论 SMS 是否成功都要更新)
	s.updateRateLimit(ctx, in.Mobile, clientIp)

	return &model.SmsSendRes{
		RequestId: "REQUEST-ID",
		ExpireIn:  smsCodeExpireSeconds,
	}, nil
}

func (s *sSms) VerifyCode(ctx context.Context, in model.SmsVerifyReq) (res *model.SmsVerifyRes, err error) {
	// 1. Verify Code from Redis (如果可用)
	redisKey := fmt.Sprintf(DriverSmsCode, in.Mobile, in.Type)
	v, err := g.Redis(s.getRedisGroup(ctx)).Do(ctx, "GET", redisKey)

	var codeValid bool
	if err != nil || v.IsNil() {
		// Redis 不可用或没有缓存，从数据库验证
		g.Log().Warningf(ctx, "Redis unavailable, checking database for verification code")
		var dbCode entity.VerificationCodes
		err = dao.VerificationCodes.Ctx(ctx).Where(do.VerificationCodes{
			Mobile: in.Mobile,
			Code:   in.Code,
			Type:   in.Type,
			Status: 0,
		}).Where("expire_at > ?", gtime.Timestamp()).Order("created_at desc").Limit(1).Scan(&dbCode)

		if err != nil || dbCode.Id == 0 {
			return nil, gerror.New("验证码错误或已过期")
		}
		codeValid = true
	} else {
		// 从 Redis 验证
		if v.String() != in.Code {
			return nil, gerror.New("验证码错误或已过期")
		}
		codeValid = true
	}

	if !codeValid {
		return nil, gerror.New("验证码错误或已过期")
	}

	// 2. Clear Code
	g.Redis(s.getRedisGroup(ctx)).Do(ctx, "DEL", redisKey) // 忽略错误
	// Update DB Status
	dao.VerificationCodes.Ctx(ctx).Where(do.VerificationCodes{
		Mobile: in.Mobile,
		Code:   in.Code,
		Type:   in.Type,
		Status: 0,
	}).Order("created_at desc").Limit(1).Data(do.VerificationCodes{
		Status:    1, // Used
		UpdatedAt: gtime.Now(),
	}).Update()

	// 3. User Logic (Login/Register)
	var userEntity entity.Users
	err = dao.Users.Ctx(ctx).Where(do.Users{Mobile: in.Mobile}).Scan(&userEntity)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if userEntity.Id == 0 {
		// Register
		// 使用 do.Users 插入，确保 OpenId 为 NULL，避免空字符串导致的唯一索引冲突
		insertData := do.Users{
			Mobile:      in.Mobile,
			Username:    in.Mobile, // Default username = mobile
			Nickname:    "用户" + in.Mobile[len(in.Mobile)-4:],
			UserType:    0,
			CreatedTime: gtime.Now(),
			UpdatedTime: gtime.Now(),
		}
		// Insert
		id, err := dao.Users.Ctx(ctx).Data(insertData).InsertAndGetId()
		if err != nil {
			return nil, gerror.Wrap(err, "注册用户失败")
		}
		// Re-assign to userEntity
		userEntity = entity.Users{
			Id:          id,
			Mobile:      in.Mobile,
			Username:    in.Mobile,
			Nickname:    "用户" + in.Mobile[len(in.Mobile)-4:],
			UserType:    0,
			CreatedTime: gtime.Now(),
			UpdatedTime: gtime.Now(),
		}
	}

	// 4. Generate Token
	token, err := s.generateToken(ctx, &userEntity)
	if err != nil {
		return nil, err
	}

	return &model.SmsVerifyRes{
		Token: token,
		User: &model.UserDetail{
			Id:       userEntity.Id,
			Mobile:   userEntity.Mobile,
			Nickname: userEntity.Nickname,
		},
	}, nil
}

func (s *sSms) CheckCode(ctx context.Context, mobile, code, typeVal string) error {
	redisKey := fmt.Sprintf(DriverSmsCode, mobile, typeVal)
	v, err := g.Redis(s.getRedisGroup(ctx)).Do(ctx, "GET", redisKey)
	if err != nil {
		return err
	}
	if v.IsNil() || v.String() != code {
		return gerror.New("验证码错误或已过期")
	}

	// 验证通过后删除验证码（或是可选行为，这里为了安全选择删除）
	g.Redis(s.getRedisGroup(ctx)).Do(ctx, "DEL", redisKey)

	// 更新 DB 状态
	dao.VerificationCodes.Ctx(ctx).Where(do.VerificationCodes{
		Mobile: mobile,
		Code:   code,
		Type:   typeVal,
		Status: 0,
	}).Order("created_at desc").Limit(1).Data(do.VerificationCodes{
		Status:    1, // Used
		UpdatedAt: gtime.Now(),
	}).Update()

	return nil
}

func (s *sSms) GetStatus(ctx context.Context, mobile string) (res *model.SmsStatusRes, err error) {
	// Implementation for status check if needed
	// For simplicity, returning empty or basic info

	// Check last send time from Redis
	lastSendKey := fmt.Sprintf(DriverSmsLastSend, mobile)
	lastSendTime, _ := g.Redis(s.getRedisGroup(ctx)).Do(ctx, "GET", lastSendKey) // timestamp

	// remaining count...

	res = &model.SmsStatusRes{
		HasSent: !lastSendTime.IsNil(),
	}
	return
}

// -------------------------------------------------------------------------
// Private Helpers
// -------------------------------------------------------------------------

func (s *sSms) checkRateLimit(ctx context.Context, mobile, ip string) error {
	// 1. Check 60s cooldown
	lastKey := fmt.Sprintf(DriverSmsLastSend, mobile)
	v, err := g.Redis(s.getRedisGroup(ctx)).Do(ctx, "GET", lastKey)
	if err != nil {
		g.Log().Warningf(ctx, "Redis rate limit check failed, falling back to database: %v", err)
		return s.checkDatabaseRateLimit(ctx, mobile, ip)
	}
	if !v.IsNil() {
		return gerror.New("发送过于频繁,请稍后再试")
	}

	// Redis 中没有冷却标记时，继续用数据库兜底。
	// 这可以覆盖上一次 Redis 写入失败、Redis 重启或缓存被清理的情况。
	return s.checkDatabaseRateLimit(ctx, mobile, ip)
}

func (s *sSms) updateRateLimit(ctx context.Context, mobile, ip string) {
	// Set 60s cooldown
	lastKey := fmt.Sprintf(DriverSmsLastSend, mobile)
	if _, err := g.Redis(s.getRedisGroup(ctx)).Do(ctx, "SETEX", lastKey, smsSendCooldown, gtime.Timestamp()); err != nil {
		g.Log().Warningf(ctx, "Redis rate limit update failed: %v", err)
	}
}

func (s *sSms) checkDatabaseRateLimit(ctx context.Context, mobile, ip string) error {
	now := time.Now()
	cooldownSince := now.Add(-time.Duration(smsSendCooldown) * time.Second)
	hourSince := now.Add(-time.Hour)
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	recentCount, err := dao.VerificationCodes.Ctx(ctx).
		Where(do.VerificationCodes{Mobile: mobile}).
		Where("created_at >= ?", cooldownSince).
		Count()
	if err != nil {
		return gerror.Wrap(err, "短信限流检查失败")
	}
	if recentCount > 0 {
		return gerror.New("发送过于频繁,请稍后再试")
	}

	mobileDailyCount, err := dao.VerificationCodes.Ctx(ctx).
		Where(do.VerificationCodes{Mobile: mobile}).
		Where("created_at >= ?", dayStart).
		Count()
	if err != nil {
		return gerror.Wrap(err, "短信限流检查失败")
	}
	if mobileDailyCount >= smsDailyMobileLimit {
		return gerror.New("该手机号今日验证码发送次数已达上限")
	}

	if ip == "" {
		return nil
	}
	ipHourlyCount, err := dao.VerificationCodes.Ctx(ctx).
		Where(do.VerificationCodes{ClientIp: ip}).
		Where("created_at >= ?", hourSince).
		Count()
	if err != nil {
		return gerror.Wrap(err, "短信限流检查失败")
	}
	if ipHourlyCount >= smsHourlyIpLimit {
		return gerror.New("当前IP验证码发送过于频繁,请稍后再试")
	}

	return nil
}

func (s *sSms) sendAliSms(ctx context.Context, ak, sk, mobile, code string) error {
	// Ref: https://next.api.aliyun.com/api/Dypnsapi/2017-05-25/SendSmsVerifyCode
	config := &openapi.Config{
		AccessKeyId:     tea.String(ak),
		AccessKeySecret: tea.String(sk),
		Endpoint:        tea.String("dypnsapi.aliyuncs.com"),
	}
	client, err := dypnsapi.NewClient(config)
	if err != nil {
		return err
	}

	// 构建模板参数 - 使用 map 然后 JSON marshal
	templateParam := map[string]string{
		"code": code,
		"min":  "5", // 添加 min 参数以匹配模板要求
	}
	templateParamBytes, _ := json.Marshal(templateParam)

	req := &dypnsapi.SendSmsVerifyCodeRequest{
		PhoneNumber:   tea.String(mobile),
		SignName:      tea.String(g.Cfg().MustGet(ctx, "sms.signName").String()),
		TemplateCode:  tea.String(g.Cfg().MustGet(ctx, "sms.templateCode").String()),
		TemplateParam: tea.String(string(templateParamBytes)),
	}

	runtime := &util.RuntimeOptions{}
	resp, err := client.SendSmsVerifyCodeWithOptions(req, runtime)
	if err != nil {
		return err
	}

	if *resp.Body.Code != "OK" {
		return gerror.Newf("短信发送失败: %s", *resp.Body.Message)
	}

	return nil
}

func (s *sSms) generateToken(ctx context.Context, u *entity.Users) (string, error) {
	uc := &user.JwtClaims{
		Id:       u.Id,
		Username: u.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(720 * time.Hour)), // 30 days
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	// Using hardcoded key from consts as seen in auth.go
	return token.SignedString([]byte(consts.JwtKey))

}

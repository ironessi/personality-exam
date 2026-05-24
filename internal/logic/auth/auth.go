package auth

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	userLogic "PersonalityTest/internal/logic/user"
	"PersonalityTest/internal/model/do"
	"PersonalityTest/internal/model/entity"
	"PersonalityTest/internal/service"
	"context"
	"time"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/golang-jwt/jwt/v5"
)

type Auth struct {
	u *userLogic.User
}

func New() *Auth {
	return &Auth{
		u: userLogic.New(),
	}
}

func (a *Auth) EncryptPassword(password string) string {
	return gmd5.MustEncryptString(password)
}

func (a *Auth) Register(ctx context.Context, mobile, password, code string) (err error) {
	// 1. Verify Code
	if err = service.Sms().CheckCode(ctx, mobile, code, "REGISTER"); err != nil {
		return err
	}

	// 2. Check if mobile already exists
	count, err := dao.Users.Ctx(ctx).Where("mobile=?", mobile).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return gerror.New("该手机号已注册")
	}

	// 3. Create User (手机号同时作为用户名)
	_, err = dao.Users.Ctx(ctx).Data(&do.Users{
		Username:    mobile,
		Password:    a.EncryptPassword(password),
		Mobile:      mobile,
		Nickname:    "用户" + mobile[len(mobile)-4:],
		CreatedTime: gtime.Now(),
		UpdatedTime: gtime.Now(),
	}).Insert()

	return err
}

func (a *Auth) Login(ctx context.Context, account, password string) (tokenString string, err error) {
	var userEntity entity.Users
	// Login with username OR mobile
	err = dao.Users.Ctx(ctx).Where("username=? OR mobile=?", account, account).Scan(&userEntity)
	if err != nil {
		return "", gerror.New("用户名或密码错误")
	}

	if userEntity.Id == 0 {
		return "", gerror.New("用户不存在")
	}

	if userEntity.Password != a.u.EncryptPassword(password) {
		return "", gerror.New("用户名或密码错误")
	}

	uc := &userLogic.JwtClaims{
		Id:       userEntity.Id,
		Username: userEntity.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(999999 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	tokenString, err = token.SignedString([]byte(consts.JwtKey))
	if err != nil {
		return "", gerror.New("生成Token失败")
	}

	_, err = dao.Users.Ctx(ctx).Where("id", userEntity.Id).Data(do.Users{
		LastLoginTime: gtime.Now(),
	}).Update()
	if err != nil {
		return "", gerror.New("更新用户信息失败")
	}
	return tokenString, nil
}

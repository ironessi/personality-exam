package user

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"PersonalityTest/internal/model/entity"
	"context"
	"errors"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/golang-jwt/jwt/v5"
)

type User struct {
}

func New() *User {
	return &User{}
}

type JwtClaims struct {
	Id       int64
	Username string
	jwt.RegisteredClaims
}

func (u *User) EncryptPassword(password string) string {
	return gmd5.MustEncryptString(password)
}

func (u *User) Info(ctx context.Context) (user *entity.Users, err error) {
	tokenString := g.RequestFromCtx(ctx).Request.Header.Get("authorization")
	// 去掉 Bearer 前缀
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	tokenClaim, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(consts.JwtKey), nil
	})

	if err != nil || tokenClaim == nil || !tokenClaim.Valid {
		return nil, errors.New("invalid token")
	}

	if claims, ok := tokenClaim.Claims.(*JwtClaims); ok && tokenClaim.Valid {
		err = dao.Users.Ctx(ctx).Where("username", claims.Username).Scan(&user)
		if err == nil && user == nil {
			return nil, errors.New("用户不存在或已被删除")
		}
	}
	return
}

func (u *User) GetUserInfo(ctx context.Context, id int64) (user *entity.Users, err error) {
	info, err := u.Info(ctx)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, errors.New("无法获取当前登录用户信息")
	}
	if info.Id != id {
		return nil, errors.New("无权限访问")
	}

	err = dao.Users.Ctx(ctx).Where("id", id).Scan(&user)
	if err == nil && user == nil {
		return nil, errors.New("未找到指定用户的数据")
	}
	return user, err
}

func (u *User) UpdateUserInfo(ctx context.Context, id int64, nickname string, gender int) (err error) {
	info, err := u.Info(ctx)
	if err != nil {
		return err
	}
	if info == nil {
		return errors.New("无法获取当前登录用户信息")
	}
	if info.Id != id {
		return errors.New("无权限访问")
	}
	_, err = dao.Users.Ctx(ctx).Where("id", id).Data(&do.Users{
		Nickname: nickname,
		Gender:   gender,
	}).Update()
	return err
}

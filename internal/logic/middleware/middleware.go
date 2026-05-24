package middleware

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/logic/user"
	"net/http"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(r *ghttp.Request) {
	var tokenString = r.Header.Get("Authorization")
	// 去掉 Bearer 前缀（如果有）
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	tokenClaim, err := jwt.ParseWithClaims(tokenString, &user.JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(consts.JwtKey), nil
	})

	if err != nil || !tokenClaim.Valid {
		r.Response.WriteStatus(http.StatusForbidden)
		r.Exit()
	}

	if claims, ok := tokenClaim.Claims.(*user.JwtClaims); ok && tokenClaim.Valid {
		r.SetCtxVar("uid", claims.Id)
	} else {
		r.Response.WriteStatus(http.StatusForbidden)
		r.Exit()
	}

	r.Middleware.Next()
}

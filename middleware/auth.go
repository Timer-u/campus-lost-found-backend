package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"campus-lost-found-backend/pkg/auth"
	"campus-lost-found-backend/pkg/response"
)

// JWT 校验通过后写入上下文的 key，controller 通过 c.Get 取值
const (
	CtxUserID   = "userID"
	CtxUsername = "username"
	CtxRole     = "role"
)

// JWTAuth 登录校验：请求头 Authorization: Bearer <token>，失败统一返回 10002
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		claims, err := auth.ParseToken(token)
		if err != nil {
			response.Fail(c, response.ErrUnauthorized)
			c.Abort()
			return
		}
		setClaims(c, claims)
		c.Next()
	}
}

// OptionalJWT 尝试解析登录态但从不拦截：token 有效时写入上下文，匿名请求以未登录身份继续
func OptionalJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		if claims, err := auth.ParseToken(token); err == nil {
			setClaims(c, claims)
		}
		c.Next()
	}
}

func setClaims(c *gin.Context, claims *auth.Claims) {
	c.Set(CtxUserID, claims.UserID)
	c.Set(CtxUsername, claims.Username)
	c.Set(CtxRole, claims.Role)
}

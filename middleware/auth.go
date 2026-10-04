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
		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxUsername, claims.Username)
		c.Set(CtxRole, claims.Role)
		c.Next()
	}
}

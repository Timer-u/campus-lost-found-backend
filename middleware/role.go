package middleware

import (
	"github.com/gin-gonic/gin"

	"campus-lost-found-backend/pkg/response"
)

// RequireRole 角色校验，必须挂在 JWTAuth 之后使用
func RequireRole(allowed ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get(CtxRole)
		current, _ := role.(string)

		for _, a := range allowed {
			if a == current {
				c.Next()
				return
			}
		}

		// 仅 system_admin 可访问的接口，lost_admin 访问时提示需要系统管理员权限
		if len(allowed) == 1 && allowed[0] == "system_admin" && current == "lost_admin" {
			response.Fail(c, response.ErrNeedSystemAdmin)
		} else {
			response.Fail(c, response.ErrNeedLostAdmin)
		}
		c.Abort()
	}
}

// RequireLostAdmin 放行 lost_admin 与 system_admin，用于物品/认领审核类接口
func RequireLostAdmin() gin.HandlerFunc {
	return RequireRole("lost_admin", "system_admin")
}

// RequireSystemAdmin 仅放行 system_admin，用于用户管理/公告/统计类接口
func RequireSystemAdmin() gin.HandlerFunc {
	return RequireRole("system_admin")
}

package middleware

import (
	"log"
	"runtime/debug"

	"github.com/gin-gonic/gin"

	"campus-lost-found-backend/pkg/response"
)

// Recovery panic 恢复中间件：未捕获异常统一返回 10000，不向外暴露堆栈信息
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[panic] %v\n%s", err, debug.Stack())
				response.Fail(c, response.ErrInternal)
				c.Abort()
			}
		}()
		c.Next()
	}
}

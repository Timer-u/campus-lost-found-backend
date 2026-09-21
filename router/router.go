package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	// 创建一个默认 Gin 引擎
	r := gin.Default()

	// 注册路由映射
	// 当有人用 GET 方法访问 "/health" 网址时，执行后面的匿名函数
	r.GET("/health", func(c *gin.Context) {
		// c *gin.Context 是 Gin 的上下文，代表这次 HTTP 请求的所有信息
		// c.JSON 表示给客户端返回一个 JSON 格式的数据
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "success",
			"data": gin.H{
				"status": "ok",
			},
		})
	})

	return r
}

package router

import (
	"github.com/gin-contrib/cors" // 跨域中间件
	"github.com/gin-gonic/gin"

	"campus-lost-found-backend/controller"
	"campus-lost-found-backend/middleware"
	"campus-lost-found-backend/pkg/response"
)

func SetupRouter() *gin.Engine {
	r := gin.New()

	// 全局中间件：日志、panic 恢复（统一返回 10000）、跨域
	r.Use(gin.Logger(), middleware.Recovery(), cors.Default())

	// 健康检查接口
	r.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{
			"status": "ok",
		})
	})

	// 业务路由：/api/v1，与 docs/openapi.yaml 保持一致
	apiV1 := r.Group("/api/v1")
	{
		// 认证：注册、登录
		authGroup := apiV1.Group("/auth")
		{
			authGroup.POST("/register", controller.Register)
			authGroup.POST("/login", controller.Login)
		}

		// 物品：公开查询
		apiV1.GET("/items", controller.GetItemList)
		apiV1.GET("/items/:itemId", controller.GetItemDetail)
	}

	// 需鉴权接口组骨架（JWT 中间件已就绪，对应接口交付时启用）：
	// authed := apiV1.Group("")
	// authed.Use(middleware.JWTAuth())
	// {
	// 	// 普通登录用户接口
	// }
	// admin := apiV1.Group("/admin")
	// admin.Use(middleware.JWTAuth(), middleware.RequireLostAdmin())
	// {
	// 	// 失物招领管理员接口；用户/公告/统计类接口用 middleware.RequireSystemAdmin()
	// }

	return r
}

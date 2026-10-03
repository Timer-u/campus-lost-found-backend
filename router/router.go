package router

import (
	//"net/http" // 原有的导入，健康检查接口会用到

	"github.com/gin-contrib/cors" // 新增：跨域中间件
	"github.com/gin-gonic/gin"

	"campus-lost-found-backend/controller"   // 新增：导入控制器
	"campus-lost-found-backend/pkg/response" // 新增：统一响应工具包，暂未使用先注释
)

func SetupRouter() *gin.Engine {
	// 创建一个默认 Gin 引擎
	r := gin.Default()

	// ===== 新增：全局跨域中间件 =====
	r.Use(cors.Default())

	// 原健康检查接口代码（恢复启用，保证编译运行）
	// 注册路由映射
	// 当有人用 GET 方法访问 "/health" 网址时，执行后面的匿名函数
	// r.GET("/health", func(c *gin.Context) {
	// 	// c *gin.Context 是 Gin 的上下文，代表这次 HTTP 请求的所有信息
	// 	// c.JSON 表示给客户端返回一个 JSON 格式的数据
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"code": 0,
	// 		"msg":  "success",
	// 		"data": gin.H{
	// 			"status": "ok",
	// 		},
	// 	})
	// })

	// 新统一响应写法（暂注释，等确认response包函数名再启用）
	// 健康检查接口：使用统一响应格式，返回结构和原代码完全一致
	r.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{
			"status": "ok",
		})
	})

	//新增：业务路由分组预留骨架
	// 公开接口组：不需要登录就能访问
	publicGroup := r.Group("/api/public")
	{
		// 用户相关
		publicGroup.POST("/register", controller.Register) // 用户注册
		publicGroup.POST("/login", controller.Login)       // 用户登录

		// 物品相关
		publicGroup.GET("/items", controller.GetItemList)      // 获取物品列表
		publicGroup.GET("/item/:id", controller.GetItemDetail) // 获取物品详情
	}

	//需鉴权接口组：必须登录才能访问（后续写完JWT中间件再取消注释启用）
	// authGroup := r.Group("/api")
	// authGroup.Use(middleware.JWTAuth())
	// {
	// 	// 用户相关接口
	// 	userGroup := authGroup.Group("/user")
	// 	{
	// 		// userGroup.GET("/info", controller.GetUserInfo)
	// 	}

	// 	// 物品相关接口
	// 	itemGroup := authGroup.Group("/item")
	// 	{
	// 		// itemGroup.POST("/publish", controller.PublishItem)
	// 		// itemGroup.PUT("/:id", controller.UpdateItem)
	// 		// itemGroup.DELETE("/:id", controller.DeleteItem)
	// 	}
	// }

	return r
}

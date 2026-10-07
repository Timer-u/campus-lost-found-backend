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

	// 图片静态托管：上传接口返回的 /uploads/xxx 直接由 Gin 提供
	r.Static("/uploads", "uploads")

	// 健康检查接口
	r.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{
			"status": "ok",
		})
	})

	// 业务路由：/api/v1，与 docs/openapi.yaml 保持一致
	apiV1 := r.Group("/api/v1")
	{
		// 认证：注册、登录（公开）；me、logout（需登录）
		authGroup := apiV1.Group("/auth")
		{
			authGroup.POST("/register", controller.Register)
			authGroup.POST("/login", controller.Login)

			authed := authGroup.Group("")
			authed.Use(middleware.JWTAuth())
			{
				authed.GET("/me", controller.GetCurrentUser)
				authed.POST("/logout", controller.Logout)
				authed.PATCH("/profile", controller.UpdateProfile)
				authed.DELETE("/account", controller.DeleteAccount)
			}
		}

		// 物品：公开查询（详情带可选登录态，发布者本人可见自己的待审核物品）
		apiV1.GET("/items", controller.GetItemList)
		apiV1.GET("/items/:itemId", middleware.OptionalJWT(), controller.GetItemDetail)

		// 物品：发布/编辑/删除/我的发布与图片上传（需登录）
		authedItems := apiV1.Group("")
		authedItems.Use(middleware.JWTAuth())
		{
			authedItems.POST("/items", controller.CreateItem)
			authedItems.PATCH("/items/:itemId", controller.UpdateItem)
			authedItems.DELETE("/items/:itemId", controller.DeleteItem)
			authedItems.GET("/me/items", controller.ListMyItems)
			authedItems.POST("/uploads/images", controller.UploadImage)
		}

		// 认领申请：提交/查看/取消（需登录）
		authedClaims := apiV1.Group("")
		authedClaims.Use(middleware.JWTAuth())
		{
			authedClaims.POST("/items/:itemId/claims", controller.CreateClaim)
			authedClaims.GET("/items/:itemId/claims", controller.ListItemClaims)
			authedClaims.GET("/me/claims", controller.ListMyClaims)
			authedClaims.POST("/claims/:claimId/cancel", controller.CancelClaim)
		}
	}

	// 管理员接口：物品与认领审核（lost_admin / system_admin）
	admin := apiV1.Group("/admin")
	admin.Use(middleware.JWTAuth(), middleware.RequireLostAdmin())
	{
		admin.GET("/items", controller.AdminListItems)
		admin.PATCH("/items/:itemId/review", controller.AdminReviewItem)
		admin.PATCH("/items/:itemId/status", controller.AdminUpdateItemStatus)
		admin.PATCH("/claims/:claimId/status", controller.AdminReviewClaim)
	}

	// 系统管理员接口：用户管理（system_admin）
	sysAdmin := apiV1.Group("/admin/users")
	sysAdmin.Use(middleware.JWTAuth(), middleware.RequireSystemAdmin())
	{
		sysAdmin.GET("", controller.AdminListUsers)
		sysAdmin.PATCH("/:userId/role", controller.AdminUpdateUserRole)
		sysAdmin.PATCH("/:userId/status", controller.AdminUpdateUserStatus)
	}

	// 系统管理员接口：公告管理（system_admin）
	annAdmin := apiV1.Group("/admin/announcements")
	annAdmin.Use(middleware.JWTAuth(), middleware.RequireSystemAdmin())
	{
		annAdmin.GET("", controller.AdminListAnnouncements)
		annAdmin.POST("", controller.AdminCreateAnnouncement)
		annAdmin.PATCH("/:announcementId", controller.AdminUpdateAnnouncement)
		annAdmin.DELETE("/:announcementId", controller.AdminDeleteAnnouncement)
	}

	// 统计总览（lost_admin / system_admin）
	apiV1.GET("/admin/statistics/overview",
		middleware.JWTAuth(), middleware.RequireLostAdmin(), controller.GetStatisticsOverview)

	// 公开公告列表
	apiV1.GET("/announcements", controller.ListPublicAnnouncements)

	return r
}

package controller

import (
	"github.com/gin-gonic/gin"

	"campus-lost-found-backend/pkg/response"
	"campus-lost-found-backend/service"
)

// GetStatisticsOverview 系统统计总览（lost_admin / system_admin 均可查看）
func GetStatisticsOverview(c *gin.Context) {
	stats, eno := service.GetStatisticsOverview()
	if eno != nil {
		response.Fail(c, eno)
		return
	}

	response.Success(c, stats)
}

package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"campus-lost-found-backend/pkg/response"
	"campus-lost-found-backend/pkg/util"
	"campus-lost-found-backend/service"
)

// AdminListUsers 用户列表（keyword 按学号或姓名搜索）
func AdminListUsers(c *gin.Context) {
	page, pageSize := util.ParsePage(c.Query("page"), c.Query("pageSize"))

	users, meta, eno := service.AdminListUsers(service.AdminUserListQuery{
		Page:     page,
		PageSize: pageSize,
		Keyword:  c.Query("keyword"),
		Role:     c.Query("role"),
		Status:   c.Query("status"),
	})
	if eno != nil {
		response.Fail(c, eno)
		return
	}

	response.Success(c, gin.H{"users": users, "meta": meta})
}

// adminUpdateUserRoleRequest 角色修改参数
type adminUpdateUserRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=student lost_admin system_admin"`
}

// AdminUpdateUserRole 修改用户角色
func AdminUpdateUserRole(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		response.Fail(c, response.ErrInvalidParams.WithMsg("无效的用户ID"))
		return
	}

	var req adminUpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrInvalidParams)
		return
	}

	user, eno := service.AdminUpdateUserRole(uint(userID), req.Role)
	if eno != nil {
		response.Fail(c, eno)
		return
	}

	response.Success(c, user)
}

// adminUpdateUserStatusRequest 账号状态参数
type adminUpdateUserStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active disabled"`
}

// AdminUpdateUserStatus 启用或禁用用户
func AdminUpdateUserStatus(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		response.Fail(c, response.ErrInvalidParams.WithMsg("无效的用户ID"))
		return
	}

	var req adminUpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrInvalidParams)
		return
	}

	user, eno := service.AdminUpdateUserStatus(uint(userID), req.Status)
	if eno != nil {
		response.Fail(c, eno)
		return
	}

	response.Success(c, user)
}

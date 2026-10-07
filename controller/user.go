package controller

import (
	"github.com/gin-gonic/gin"

	"campus-lost-found-backend/middleware"
	"campus-lost-found-backend/pkg/response"
	"campus-lost-found-backend/service"
)

// registerRequest 注册参数，校验规则与 docs/openapi.yaml 的 RegisterRequest 一致
type registerRequest struct {
	Username string `json:"username" binding:"required,numeric"`
	Name     string `json:"name" binding:"required,max=50"`
	Password string `json:"password" binding:"required,min=8,max=32"`
}

// Register 用户注册，成功返回 201 与用户信息
func Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrInvalidParams)
		return
	}

	user, eno := service.Register(req.Username, req.Name, req.Password)
	if eno != nil {
		response.Fail(c, eno)
		return
	}

	response.Created(c, user)
}

// loginRequest 登录参数
type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 用户登录，签发 JWT
func Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrInvalidParams)
		return
	}

	result, eno := service.Login(req.Username, req.Password)
	if eno != nil {
		response.Fail(c, eno)
		return
	}

	response.Success(c, result)
}

// GetCurrentUser 获取当前登录用户
func GetCurrentUser(c *gin.Context) {
	id, _ := c.Get(middleware.CtxUserID)
	userID, _ := id.(uint)

	user, eno := service.GetByID(userID)
	if eno != nil {
		response.Fail(c, eno)
		return
	}

	response.Success(c, user)
}

// Logout 退出登录：无状态 JWT，前端收到成功后删除本地 Token 即可
func Logout(c *gin.Context) {
	response.Success(c, nil)
}

// updateProfileRequest 修改个人信息参数，约束与 docs/openapi.yaml 的 UpdateProfileRequest 一致
type updateProfileRequest struct {
	Name string `json:"name" binding:"required,max=50"`
}

// UpdateProfile 修改当前登录用户的个人信息
func UpdateProfile(c *gin.Context) {
	var req updateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrInvalidParams)
		return
	}

	user, eno := service.UpdateProfile(currentUserID(c), req.Name)
	if eno != nil {
		response.Fail(c, eno)
		return
	}

	response.Success(c, user)
}

// deleteAccountRequest 注销账号参数（危险操作，需验证当前密码）
type deleteAccountRequest struct {
	Password string `json:"password" binding:"required"`
}

// DeleteAccount 注销当前登录账号
func DeleteAccount(c *gin.Context) {
	var req deleteAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrInvalidParams)
		return
	}

	if eno := service.DeleteAccount(currentUserID(c), req.Password); eno != nil {
		response.Fail(c, eno)
		return
	}

	response.Success(c, nil)
}

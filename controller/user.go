package controller

import (
	"time"

	"github.com/gin-gonic/gin"

	"campus-lost-found-backend/model"
	"campus-lost-found-backend/pkg/response"
)

// 模拟用户数据
var userList = []model.User{
	{
		Base:         model.Base{ID: 1, CreatedAt: time.Now()},
		Username:     "20260001",
		PasswordHash: "123456",
		Name:         "测试用户",
		Role:         "student",
		Status:       "active",
	},
}

// Register 用户注册
func Register(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	// 绑定前端传的JSON数据
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrInvalidParams)
		return
	}

	// 简单校验
	if req.Username == "" || req.Password == "" {
		response.Fail(c, response.ErrInvalidParams.WithMsg("用户名和密码不能为空"))
		return
	}

	// 创建新用户（模拟）
	newUser := model.User{
		Base:         model.Base{ID: uint(len(userList) + 1), CreatedAt: time.Now()},
		Username:     req.Username,
		PasswordHash: req.Password,
		Name:         req.Name,
	}
	userList = append(userList, newUser)

	response.Success(c, "注册成功")
}

// Login 用户登录
func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrInvalidParams)
		return
	}

	// 校验用户名密码
	for _, user := range userList {
		if user.Username == req.Username && user.PasswordHash == req.Password {
			// 后续这里会返回JWT令牌，现在先返回用户信息
			response.Success(c, gin.H{
				"user":  user,
				"token": "模拟token_xxxxxx",
			})
			return
		}
	}

	response.Fail(c, response.ErrInvalidCredentials)
}

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
		ID:        1,
		RealName:  "test",
		Password:  "123456",
		Username:  "测试用户",
		Phone:     "138****1234",
		CreatedAt: time.Now(),
	},
}

// Register 用户注册
func Register(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		RealName string `json:"real_name"`
	}

	// 绑定前端传的JSON数据
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, "参数错误")
		return
	}

	// 简单校验
	if req.Username == "" || req.Password == "" {
		response.Error(c, "用户名和密码不能为空")
		return
	}

	// 创建新用户（模拟）
	newUser := model.User{
		ID:        uint(len(userList) + 1),
		Username:  req.Username,
		Password:  req.Password,
		RealName:  req.RealName,
		CreatedAt: time.Now(),
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
		response.ParamError(c, "参数错误")
		return
	}

	// 校验用户名密码
	for _, user := range userList {
		if user.Username == req.Username && user.Password == req.Password {
			// 后续这里会返回JWT令牌，现在先返回用户信息
			response.Success(c, gin.H{
				"user":  user,
				"token": "模拟token_xxxxxx",
			})
			return
		}
	}

	response.Error(c, "用户名或密码错误")
}

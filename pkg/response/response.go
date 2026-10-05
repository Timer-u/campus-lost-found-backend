package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Success 成功响应，自动填充 code=0 和 msg="success"
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": data,
	})
}

// Created 201 成功响应（如注册接口）
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, gin.H{
		"code": 0,
		"msg":  "success",
		"data": data,
	})
}

// Fail 按错误码返回统一失败响应，HTTP 状态码与业务码同时传递
func Fail(c *gin.Context, e *Errno) {
	c.JSON(e.HTTP, gin.H{
		"code": e.Code,
		"msg":  e.Msg,
		"data": nil,
	})
}

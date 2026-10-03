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

// Error 普通业务失败响应，传入错误提示信息
func Error(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  msg,
		"data": nil,
	})
}

// ErrorWithCode 自定义错误码的失败响应
func ErrorWithCode(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
		"data": nil,
	})
}

// ParamError 参数错误响应，返回 400 状态码
func ParamError(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"code": 400,
		"msg":  msg,
		"data": nil,
	})
}

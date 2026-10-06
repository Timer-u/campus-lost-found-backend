package controller

import (
	"github.com/gin-gonic/gin"

	"campus-lost-found-backend/pkg/response"
	"campus-lost-found-backend/service"
)

// UploadImage 上传物品图片，返回可直接引用的图片 URL
func UploadImage(c *gin.Context) {
	fh, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, response.ErrInvalidParams.WithMsg("缺少 file 字段"))
		return
	}

	url, eno := service.SaveImage(fh)
	if eno != nil {
		response.Fail(c, eno)
		return
	}

	response.Created(c, gin.H{"url": url})
}

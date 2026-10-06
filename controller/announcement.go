package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"campus-lost-found-backend/pkg/response"
	"campus-lost-found-backend/pkg/util"
	"campus-lost-found-backend/service"
)

// ListPublicAnnouncements 公开公告列表（仅已发布）
func ListPublicAnnouncements(c *gin.Context) {
	page, pageSize := util.ParsePage(c.Query("page"), c.Query("pageSize"))

	list, meta, eno := service.ListPublicAnnouncements(page, pageSize)
	if eno != nil {
		response.Fail(c, eno)
		return
	}

	response.Success(c, gin.H{"announcements": list, "meta": meta})
}

// AdminListAnnouncements 管理员公告列表（published=true/false 筛选）
func AdminListAnnouncements(c *gin.Context) {
	page, pageSize := util.ParsePage(c.Query("page"), c.Query("pageSize"))

	list, meta, eno := service.AdminListAnnouncements(c.Query("published"), page, pageSize)
	if eno != nil {
		response.Fail(c, eno)
		return
	}

	response.Success(c, gin.H{"announcements": list, "meta": meta})
}

// adminCreateAnnouncementRequest 创建公告参数，约束与 docs/openapi.yaml 一致
type adminCreateAnnouncementRequest struct {
	Title     string `json:"title" binding:"required,max=100"`
	Content   string `json:"content" binding:"required,max=10000"`
	Published bool   `json:"published"`
}

// AdminCreateAnnouncement 创建公告
func AdminCreateAnnouncement(c *gin.Context) {
	var req adminCreateAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrInvalidParams)
		return
	}

	announcement, eno := service.CreateAnnouncement(currentUserID(c), service.CreateAnnouncementInput{
		Title:     req.Title,
		Content:   req.Content,
		Published: req.Published,
	})
	if eno != nil {
		response.Fail(c, eno)
		return
	}

	response.Created(c, announcement)
}

// adminUpdateAnnouncementRequest 修改公告参数：指针非 nil 表示该字段需要更新
type adminUpdateAnnouncementRequest struct {
	Title     *string `json:"title" binding:"omitempty,max=100"`
	Content   *string `json:"content" binding:"omitempty,max=10000"`
	Published *bool   `json:"published"`
}

// AdminUpdateAnnouncement 修改公告
func AdminUpdateAnnouncement(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("announcementId"), 10, 64)
	if err != nil {
		response.Fail(c, response.ErrInvalidParams.WithMsg("无效的公告ID"))
		return
	}

	var req adminUpdateAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrInvalidParams)
		return
	}
	if req.Title == nil && req.Content == nil && req.Published == nil {
		response.Fail(c, response.ErrInvalidParams.WithMsg("至少提供一个修改字段"))
		return
	}

	announcement, eno := service.UpdateAnnouncement(uint(id), service.UpdateAnnouncementInput{
		Title:     req.Title,
		Content:   req.Content,
		Published: req.Published,
	})
	if eno != nil {
		response.Fail(c, eno)
		return
	}

	response.Success(c, announcement)
}

// AdminDeleteAnnouncement 删除公告
func AdminDeleteAnnouncement(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("announcementId"), 10, 64)
	if err != nil {
		response.Fail(c, response.ErrInvalidParams.WithMsg("无效的公告ID"))
		return
	}

	if eno := service.DeleteAnnouncement(uint(id)); eno != nil {
		response.Fail(c, eno)
		return
	}

	response.Success(c, nil)
}

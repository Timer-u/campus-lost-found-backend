package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"campus-lost-found-backend/pkg/response"
	"campus-lost-found-backend/pkg/util"
	"campus-lost-found-backend/service"
)

// AdminListItems 管理员物品列表：全部状态可见，待审核优先
func AdminListItems(c *gin.Context) {
	page, pageSize := util.ParsePage(c.Query("page"), c.Query("pageSize"))

	items, meta, eno := service.AdminListItems(service.AdminItemListQuery{
		Page:         page,
		PageSize:     pageSize,
		ReviewStatus: c.Query("reviewStatus"),
		ItemStatus:   c.Query("itemStatus"),
		Keyword:      c.Query("keyword"),
	})
	if eno != nil {
		response.Fail(c, eno)
		return
	}

	response.Success(c, gin.H{"items": items, "meta": meta})
}

// adminReviewItemRequest 审核参数，约束与 docs/openapi.yaml 的 ReviewItemRequest 一致
type adminReviewItemRequest struct {
	ReviewStatus string `json:"reviewStatus" binding:"required,oneof=approved rejected offline"`
	RejectReason string `json:"rejectReason" binding:"omitempty,max=500"`
}

// AdminReviewItem 审核物品发布
func AdminReviewItem(c *gin.Context) {
	itemID, err := strconv.ParseUint(c.Param("itemId"), 10, 64)
	if err != nil {
		response.Fail(c, response.ErrInvalidParams.WithMsg("无效的物品ID"))
		return
	}

	var req adminReviewItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrInvalidParams)
		return
	}

	item, eno := service.AdminReviewItem(uint(itemID), req.ReviewStatus, req.RejectReason)
	if eno != nil {
		response.Fail(c, eno)
		return
	}

	response.Success(c, item)
}

// adminUpdateItemStatusRequest 物品状态调整参数
type adminUpdateItemStatusRequest struct {
	ItemStatus string `json:"itemStatus" binding:"required,oneof=open claimed resolved closed"`
}

// AdminUpdateItemStatus 修改物品状态
func AdminUpdateItemStatus(c *gin.Context) {
	itemID, err := strconv.ParseUint(c.Param("itemId"), 10, 64)
	if err != nil {
		response.Fail(c, response.ErrInvalidParams.WithMsg("无效的物品ID"))
		return
	}

	var req adminUpdateItemStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrInvalidParams)
		return
	}

	item, eno := service.AdminUpdateItemStatus(uint(itemID), req.ItemStatus)
	if eno != nil {
		response.Fail(c, eno)
		return
	}

	response.Success(c, item)
}

// adminReviewClaimRequest 认领审核参数，约束与 docs/openapi.yaml 的 ReviewClaimRequest 一致
type adminReviewClaimRequest struct {
	Status       string `json:"status" binding:"required,oneof=approved rejected"`
	ReviewReason string `json:"reviewReason" binding:"omitempty,max=500"`
}

// AdminReviewClaim 审核认领申请（通过时事务内同步物品状态）
func AdminReviewClaim(c *gin.Context) {
	claimID, err := strconv.ParseUint(c.Param("claimId"), 10, 64)
	if err != nil {
		response.Fail(c, response.ErrInvalidParams.WithMsg("无效的申请ID"))
		return
	}

	var req adminReviewClaimRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrInvalidParams)
		return
	}

	claim, eno := service.AdminReviewClaim(currentUserID(c), uint(claimID), req.Status, req.ReviewReason)
	if eno != nil {
		response.Fail(c, eno)
		return
	}

	response.Success(c, claim)
}

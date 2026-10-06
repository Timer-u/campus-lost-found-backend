package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"campus-lost-found-backend/middleware"
	"campus-lost-found-backend/pkg/response"
	"campus-lost-found-backend/pkg/util"
	"campus-lost-found-backend/service"
)

// createClaimRequest 提交认领申请参数，约束与 docs/openapi.yaml 的 CreateClaimRequest 一致
type createClaimRequest struct {
	Description string `json:"description" binding:"required,max=1000"`
	Contact     string `json:"contact" binding:"required,max=100"`
}

// CreateClaim 提交认领申请
func CreateClaim(c *gin.Context) {
	itemID, err := strconv.ParseUint(c.Param("itemId"), 10, 64)
	if err != nil {
		response.Fail(c, response.ErrInvalidParams.WithMsg("无效的物品ID"))
		return
	}

	var req createClaimRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrInvalidParams)
		return
	}

	claim, eno := service.CreateClaim(currentUserID(c), uint(itemID), service.CreateClaimInput{
		Description: req.Description,
		Contact:     req.Contact,
	})
	if eno != nil {
		response.Fail(c, eno)
		return
	}

	response.Created(c, claim)
}

// ListItemClaims 查询物品的认领申请（仅发布者或管理员）
func ListItemClaims(c *gin.Context) {
	itemID, err := strconv.ParseUint(c.Param("itemId"), 10, 64)
	if err != nil {
		response.Fail(c, response.ErrInvalidParams.WithMsg("无效的物品ID"))
		return
	}
	role, _ := c.Get(middleware.CtxRole)
	page, pageSize := util.ParsePage(c.Query("page"), c.Query("pageSize"))

	claims, meta, eno := service.ListItemClaims(currentUserID(c), role.(string), uint(itemID), page, pageSize)
	if eno != nil {
		response.Fail(c, eno)
		return
	}

	response.Success(c, gin.H{"claims": claims, "meta": meta})
}

// ListMyClaims 我的认领申请
func ListMyClaims(c *gin.Context) {
	page, pageSize := util.ParsePage(c.Query("page"), c.Query("pageSize"))

	claims, meta, eno := service.ListMyClaims(currentUserID(c), c.Query("status"), page, pageSize)
	if eno != nil {
		response.Fail(c, eno)
		return
	}

	response.Success(c, gin.H{"claims": claims, "meta": meta})
}

// CancelClaim 取消自己的认领申请
func CancelClaim(c *gin.Context) {
	claimID, err := strconv.ParseUint(c.Param("claimId"), 10, 64)
	if err != nil {
		response.Fail(c, response.ErrInvalidParams.WithMsg("无效的申请ID"))
		return
	}

	claim, eno := service.CancelClaim(currentUserID(c), uint(claimID))
	if eno != nil {
		response.Fail(c, eno)
		return
	}

	response.Success(c, claim)
}

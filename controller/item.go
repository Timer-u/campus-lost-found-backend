package controller

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"campus-lost-found-backend/middleware"
	"campus-lost-found-backend/pkg/response"
	"campus-lost-found-backend/pkg/util"
	"campus-lost-found-backend/service"
)

// GetItemList 公开物品列表：筛选、模糊搜索、排序与分页
func GetItemList(c *gin.Context) {
	page, pageSize := util.ParsePage(c.Query("page"), c.Query("pageSize"))

	items, meta, eno := service.ListPublicItems(service.ItemListQuery{
		Page:       page,
		PageSize:   pageSize,
		Keyword:    c.Query("keyword"),
		Type:       c.Query("type"),
		Category:   c.Query("category"),
		ItemStatus: c.Query("itemStatus"),
		Location:   c.Query("location"),
		Sort:       c.Query("sort"),
	})
	if eno != nil {
		response.Fail(c, eno)
		return
	}

	response.Success(c, gin.H{"items": items, "meta": meta})
}

// GetItemDetail 公开物品详情
func GetItemDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("itemId"), 10, 64)
	if err != nil {
		response.Fail(c, response.ErrInvalidParams.WithMsg("无效的物品ID"))
		return
	}

	item, eno := service.GetPublicItem(uint(id))
	if eno != nil {
		response.Fail(c, eno)
		return
	}

	response.Success(c, item)
}

// createItemRequest 发布参数，约束与 docs/openapi.yaml 的 CreateItemRequest 一致
type createItemRequest struct {
	Type        string     `json:"type" binding:"required,oneof=lost found"`
	Title       string     `json:"title" binding:"required,max=100"`
	Description string     `json:"description" binding:"required,max=2000"`
	Category    string     `json:"category" binding:"required,oneof=id_card wallet phone computer book clothing key daily other"`
	Location    string     `json:"location" binding:"required,max=100"`
	LostAt      *time.Time `json:"lostAt"`
	Contact     string     `json:"contact" binding:"required,max=100"`
	ImageURLs   []string   `json:"imageUrls" binding:"omitempty,max=6,dive,max=255"`
}

// CreateItem 发布失物/招领，默认进入待审核
func CreateItem(c *gin.Context) {
	userID := currentUserID(c)

	var req createItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrInvalidParams)
		return
	}

	item, eno := service.CreateItem(userID, service.CreateItemInput{
		Type:        req.Type,
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		Location:    req.Location,
		LostAt:      req.LostAt,
		Contact:     req.Contact,
		ImageURLs:   req.ImageURLs,
	})
	if eno != nil {
		response.Fail(c, eno)
		return
	}

	response.Created(c, item)
}

// updateItemRequest 编辑参数：指针非 nil 表示该字段需要更新，全部为 nil 视为空请求
type updateItemRequest struct {
	Title       *string    `json:"title" binding:"omitempty,max=100"`
	Description *string    `json:"description" binding:"omitempty,max=2000"`
	Category    *string    `json:"category" binding:"omitempty,oneof=id_card wallet phone computer book clothing key daily other"`
	Location    *string    `json:"location" binding:"omitempty,max=100"`
	LostAt      *time.Time `json:"lostAt"`
	Contact     *string    `json:"contact" binding:"omitempty,max=100"`
	ImageURLs   *[]string  `json:"imageUrls" binding:"omitempty,max=6,dive,max=255"`
}

// UpdateItem 编辑自己的发布
func UpdateItem(c *gin.Context) {
	userID := currentUserID(c)

	itemID, err := strconv.ParseUint(c.Param("itemId"), 10, 64)
	if err != nil {
		response.Fail(c, response.ErrInvalidParams.WithMsg("无效的物品ID"))
		return
	}

	var req updateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrInvalidParams)
		return
	}
	if req.Title == nil && req.Description == nil && req.Category == nil && req.Location == nil &&
		req.LostAt == nil && req.Contact == nil && req.ImageURLs == nil {
		response.Fail(c, response.ErrInvalidParams.WithMsg("至少提供一个修改字段"))
		return
	}

	item, eno := service.UpdateItem(userID, uint(itemID), service.UpdateItemInput{
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		Location:    req.Location,
		LostAt:      req.LostAt,
		Contact:     req.Contact,
		ImageURLs:   req.ImageURLs,
	})
	if eno != nil {
		response.Fail(c, eno)
		return
	}

	response.Success(c, item)
}

// DeleteItem 删除自己的发布
func DeleteItem(c *gin.Context) {
	userID := currentUserID(c)

	itemID, err := strconv.ParseUint(c.Param("itemId"), 10, 64)
	if err != nil {
		response.Fail(c, response.ErrInvalidParams.WithMsg("无效的物品ID"))
		return
	}

	if eno := service.DeleteItem(userID, uint(itemID)); eno != nil {
		response.Fail(c, eno)
		return
	}

	response.Success(c, nil)
}

// ListMyItems 我的发布：reviewStatus/itemStatus 筛选 + 分页
func ListMyItems(c *gin.Context) {
	userID := currentUserID(c)
	page, pageSize := util.ParsePage(c.Query("page"), c.Query("pageSize"))

	items, meta, eno := service.ListMyItems(userID, c.Query("reviewStatus"), c.Query("itemStatus"), page, pageSize)
	if eno != nil {
		response.Fail(c, eno)
		return
	}

	response.Success(c, gin.H{"items": items, "meta": meta})
}

// currentUserID 从 JWT 中间件写入的上下文中取当前用户 ID
func currentUserID(c *gin.Context) uint {
	id, _ := c.Get(middleware.CtxUserID)
	userID, _ := id.(uint)
	return userID
}

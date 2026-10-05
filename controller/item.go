package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"

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

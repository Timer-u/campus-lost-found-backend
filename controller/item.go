package controller

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"campus-lost-found-backend/model"
	"campus-lost-found-backend/pkg/response"
)

// 模拟数据：临时用内存存储，后面对接数据库再替换
var itemList = []model.Item{
	{
		Base:         model.Base{ID: 1, CreatedAt: time.Now()},
		Title:        "黑色钱包",
		Type:         "lost",
		Description:  "内有身份证和校园卡，黑色皮质",
		Category:     "wallet",
		Location:     "食堂一楼",
		Contact:      "138****1234",
		ReviewStatus: "approved",
		ItemStatus:   "open",
		UserID:       1,
	},
	{
		Base:         model.Base{ID: 2, CreatedAt: time.Now()},
		Title:        "捡到一把雨伞",
		Type:         "found",
		Description:  "蓝色长柄雨伞",
		Category:     "daily",
		Location:     "教学楼A座门口",
		Contact:      "微信:test123",
		ReviewStatus: "approved",
		ItemStatus:   "open",
		UserID:       2,
	},
}

// GetItemList 获取物品列表
func GetItemList(c *gin.Context) {
	// 后续可以加分页、筛选类型，这里先返回全部
	response.Success(c, itemList)
}

// GetItemDetail 获取物品详情
func GetItemDetail(c *gin.Context) {
	// 从路径参数里获取id
	idStr := c.Param("itemId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Fail(c, response.ErrInvalidParams.WithMsg("无效的物品ID"))
		return
	}

	// 查找对应物品
	for _, item := range itemList {
		if int(item.ID) == id {
			response.Success(c, item)
			return
		}
	}

	response.Fail(c, response.ErrItemNotFound)
}

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
		ID:          1,
		Title:       "黑色钱包",
		Type:        "lost",
		Description: "内有身份证和校园卡，黑色皮质",
		Location:    "食堂一楼",
		Phone:       "138****1234",
		Status:      "未认领",
		UserID:      1,
		CreatedAt:   time.Now(),
	},
	{
		ID:          2,
		Title:       "捡到一把雨伞",
		Type:        "found",
		Description: "蓝色长柄雨伞",
		Location:    "教学楼A座门口",
		Phone:       "微信:test123",
		Status:      "未认领",
		UserID:      2,
		CreatedAt:   time.Now(),
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
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.ParamError(c, "无效的物品ID")
		return
	}

	// 查找对应物品
	for _, item := range itemList {
		if int(item.ID) == id {
			response.Success(c, item)
			return
		}
	}

	response.Error(c, "物品不存在")
}

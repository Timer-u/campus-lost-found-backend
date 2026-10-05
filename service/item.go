package service

import (
	"errors"

	"gorm.io/gorm"

	"campus-lost-found-backend/model"
	"campus-lost-found-backend/pkg/response"
	"campus-lost-found-backend/pkg/util"
)

// 物品公开接口允许的枚举值，与 docs/openapi.yaml 保持一致
var (
	itemTypeAllowed     = map[string]bool{"lost": true, "found": true}
	itemCategoryAllowed = map[string]bool{
		"id_card": true, "wallet": true, "phone": true, "computer": true, "book": true,
		"clothing": true, "key": true, "daily": true, "other": true,
	}
	itemStatusAllowed = map[string]bool{"open": true, "claimed": true, "resolved": true, "closed": true}
)

// ItemListQuery 物品列表查询参数
type ItemListQuery struct {
	Page       int
	PageSize   int
	Keyword    string
	Type       string
	Category   string
	ItemStatus string
	Location   string
	Sort       string
}

// ListPublicItems 公开物品列表：只返回审核通过的物品，支持筛选、模糊搜索、排序与分页
func ListPublicItems(q ItemListQuery) ([]model.Item, util.PageMeta, *response.Errno) {
	// 枚举参数非法时按参数错误处理（文档定义列表接口返回 400）
	if q.Type != "" && !itemTypeAllowed[q.Type] {
		return nil, util.PageMeta{}, response.ErrInvalidParams
	}
	if q.Category != "" && !itemCategoryAllowed[q.Category] {
		return nil, util.PageMeta{}, response.ErrInvalidParams
	}
	if q.ItemStatus != "" && !itemStatusAllowed[q.ItemStatus] {
		return nil, util.PageMeta{}, response.ErrInvalidParams
	}
	if q.Sort != "" && q.Sort != "latest" && q.Sort != "oldest" {
		return nil, util.PageMeta{}, response.ErrInvalidParams
	}

	query := model.DB.Model(&model.Item{}).Where("review_status = ?", "approved")
	if q.Keyword != "" {
		kw := "%" + q.Keyword + "%"
		query = query.Where("title LIKE ? OR description LIKE ? OR location LIKE ?", kw, kw, kw)
	}
	if q.Type != "" {
		query = query.Where("type = ?", q.Type)
	}
	if q.Category != "" {
		query = query.Where("category = ?", q.Category)
	}
	if q.ItemStatus != "" {
		query = query.Where("item_status = ?", q.ItemStatus)
	}
	if q.Location != "" {
		query = query.Where("location = ?", q.Location)
	}

	base := query.Session(&gorm.Session{})
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, util.PageMeta{}, response.ErrInternal
	}

	order := "created_at DESC, id DESC" // 默认 latest；id 作次级键保证同秒创建时顺序稳定
	if q.Sort == "oldest" {
		order = "created_at ASC, id ASC"
	}

	var items []model.Item
	if err := base.Preload("User").Order(order).
		Offset((q.Page - 1) * q.PageSize).Limit(q.PageSize).Find(&items).Error; err != nil {
		return nil, util.PageMeta{}, response.ErrInternal
	}
	for i := range items {
		items[i].PublisherName = items[i].User.Name
	}

	return items, util.NewPageMeta(int64(q.Page), int64(q.PageSize), total), nil
}

// GetPublicItem 物品详情：未审核通过的物品对公众不可见
func GetPublicItem(id uint) (*model.Item, *response.Errno) {
	var item model.Item
	if err := model.DB.Preload("User").
		Where("review_status = ?", "approved").First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.ErrItemNotFound
		}
		return nil, response.ErrInternal
	}
	item.PublisherName = item.User.Name
	return &item, nil
}

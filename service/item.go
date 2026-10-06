package service

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"campus-lost-found-backend/model"
	"campus-lost-found-backend/pkg/response"
	"campus-lost-found-backend/pkg/util"
)

// 物品接口允许的枚举值，与 docs/openapi.yaml 保持一致
var (
	itemTypeAllowed     = map[string]bool{"lost": true, "found": true}
	itemCategoryAllowed = map[string]bool{
		"id_card": true, "wallet": true, "phone": true, "computer": true, "book": true,
		"clothing": true, "key": true, "daily": true, "other": true,
	}
	itemStatusAllowed   = map[string]bool{"open": true, "claimed": true, "resolved": true, "closed": true}
	reviewStatusAllowed = map[string]bool{"pending": true, "approved": true, "rejected": true, "offline": true}
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

// CreateItemInput 发布物品参数，字段约束由 controller binding 校验
type CreateItemInput struct {
	Type        string
	Title       string
	Description string
	Category    string
	Location    string
	LostAt      *time.Time
	Contact     string
	ImageURLs   []string
}

// CreateItem 发布失物/招领：默认进入待审核，等待管理员审核通过后公开
func CreateItem(userID uint, in CreateItemInput) (*model.Item, *response.Errno) {
	if len(in.ImageURLs) > 6 {
		return nil, response.ErrInvalidParams.WithMsg("每个物品最多引用 6 张图片")
	}

	item := &model.Item{
		Type:         in.Type,
		Title:        in.Title,
		Description:  in.Description,
		Category:     in.Category,
		Location:     in.Location,
		LostAt:       in.LostAt,
		Contact:      in.Contact,
		ImageURLs:    model.ImageList(in.ImageURLs),
		ReviewStatus: "pending",
		ItemStatus:   "open",
		UserID:       userID,
	}
	if item.ImageURLs == nil {
		item.ImageURLs = model.ImageList{}
	}
	if err := model.DB.Create(item).Error; err != nil {
		return nil, response.ErrInternal
	}
	if user, eno := GetByID(userID); eno == nil {
		item.PublisherName = user.Name
	}
	return item, nil
}

// UpdateItemInput 编辑物品参数：指针非 nil 表示该字段需要更新
type UpdateItemInput struct {
	Title       *string
	Description *string
	Category    *string
	Location    *string
	LostAt      *time.Time
	Contact     *string
	ImageURLs   *[]string
}

// UpdateItem 编辑自己的发布：仅发布者可操作；修改后的物品统一回到待审核重新走审核流程
func UpdateItem(userID, itemID uint, in UpdateItemInput) (*model.Item, *response.Errno) {
	var item model.Item
	if err := model.DB.First(&item, itemID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.ErrItemNotFound
		}
		return nil, response.ErrInternal
	}
	if item.UserID != userID {
		return nil, response.ErrForbidden
	}
	if item.ItemStatus == "closed" {
		return nil, response.ErrItemStatusNotAllowed
	}

	updates := map[string]any{}
	if in.Title != nil {
		updates["title"] = *in.Title
	}
	if in.Description != nil {
		updates["description"] = *in.Description
	}
	if in.Category != nil {
		updates["category"] = *in.Category
	}
	if in.Location != nil {
		updates["location"] = *in.Location
	}
	if in.LostAt != nil {
		updates["lost_at"] = *in.LostAt
	}
	if in.Contact != nil {
		updates["contact"] = *in.Contact
	}
	if in.ImageURLs != nil {
		if len(*in.ImageURLs) > 6 {
			return nil, response.ErrInvalidParams.WithMsg("每个物品最多引用 6 张图片")
		}
		updates["image_urls"] = model.ImageList(*in.ImageURLs)
	}

	if len(updates) > 0 {
		updates["review_status"] = "pending"
		if err := model.DB.Model(&item).Updates(updates).Error; err != nil {
			return nil, response.ErrInternal
		}
	}

	if err := model.DB.Preload("User").First(&item, itemID).Error; err != nil {
		return nil, response.ErrInternal
	}
	item.PublisherName = item.User.Name
	return &item, nil
}

// DeleteItem 删除自己的发布（软删除）
func DeleteItem(userID, itemID uint) *response.Errno {
	var item model.Item
	if err := model.DB.First(&item, itemID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ErrItemNotFound
		}
		return response.ErrInternal
	}
	if item.UserID != userID {
		return response.ErrForbidden
	}
	if item.ItemStatus == "closed" {
		return response.ErrItemStatusNotAllowed
	}
	if err := model.DB.Delete(&item).Error; err != nil {
		return response.ErrInternal
	}
	return nil
}

// ListMyItems 我的发布：reviewStatus/itemStatus 筛选 + 分页
func ListMyItems(userID uint, reviewStatus, itemStatus string, page, pageSize int) ([]model.Item, util.PageMeta, *response.Errno) {
	if reviewStatus != "" && !reviewStatusAllowed[reviewStatus] {
		return nil, util.PageMeta{}, response.ErrInvalidParams
	}
	if itemStatus != "" && !itemStatusAllowed[itemStatus] {
		return nil, util.PageMeta{}, response.ErrInvalidParams
	}

	query := model.DB.Model(&model.Item{}).Where("user_id = ?", userID)
	if reviewStatus != "" {
		query = query.Where("review_status = ?", reviewStatus)
	}
	if itemStatus != "" {
		query = query.Where("item_status = ?", itemStatus)
	}

	base := query.Session(&gorm.Session{})
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, util.PageMeta{}, response.ErrInternal
	}

	var items []model.Item
	if err := base.Preload("User").Order("created_at DESC, id DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, util.PageMeta{}, response.ErrInternal
	}
	for i := range items {
		items[i].PublisherName = items[i].User.Name
	}

	return items, util.NewPageMeta(int64(page), int64(pageSize), total), nil
}

package service

import (
	"errors"
	"strconv"

	"gorm.io/gorm"

	"campus-lost-found-backend/model"
	"campus-lost-found-backend/pkg/response"
	"campus-lost-found-backend/pkg/util"
)

// ListPublicAnnouncements 公开公告列表：只返回已发布的公告
func ListPublicAnnouncements(page, pageSize int) ([]model.Announcement, util.PageMeta, *response.Errno) {
	base := model.DB.Model(&model.Announcement{}).
		Where("published = ?", true).Session(&gorm.Session{})

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, util.PageMeta{}, response.ErrInternal
	}

	var list []model.Announcement
	if err := base.Order("created_at DESC, id DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, util.PageMeta{}, response.ErrInternal
	}

	return list, util.NewPageMeta(int64(page), int64(pageSize), total), nil
}

// AdminListAnnouncements 管理员公告列表：published 为 "true"/"false" 时过滤，留空返回全部
func AdminListAnnouncements(published string, page, pageSize int) ([]model.Announcement, util.PageMeta, *response.Errno) {
	base := model.DB.Model(&model.Announcement{}).Session(&gorm.Session{})
	if published != "" {
		p, err := strconv.ParseBool(published)
		if err != nil {
			return nil, util.PageMeta{}, response.ErrInvalidParams
		}
		base = base.Where("published = ?", p)
	}

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, util.PageMeta{}, response.ErrInternal
	}

	var list []model.Announcement
	if err := base.Order("created_at DESC, id DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, util.PageMeta{}, response.ErrInternal
	}

	return list, util.NewPageMeta(int64(page), int64(pageSize), total), nil
}

// CreateAnnouncementInput 创建公告参数，字段约束由 controller binding 校验
type CreateAnnouncementInput struct {
	Title     string
	Content   string
	Published bool
}

// CreateAnnouncement 创建公告，作者为当前系统管理员
func CreateAnnouncement(authorID uint, in CreateAnnouncementInput) (*model.Announcement, *response.Errno) {
	announcement := &model.Announcement{
		Title:     in.Title,
		Content:   in.Content,
		Published: in.Published,
		AuthorID:  authorID,
	}
	if err := model.DB.Create(announcement).Error; err != nil {
		return nil, response.ErrInternal
	}
	return announcement, nil
}

// UpdateAnnouncementInput 修改公告参数：指针非 nil 表示该字段需要更新
type UpdateAnnouncementInput struct {
	Title     *string
	Content   *string
	Published *bool
}

// UpdateAnnouncement 修改公告（部分更新）
func UpdateAnnouncement(id uint, in UpdateAnnouncementInput) (*model.Announcement, *response.Errno) {
	var announcement model.Announcement
	if err := model.DB.First(&announcement, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.ErrAnnouncementNotFound
		}
		return nil, response.ErrInternal
	}

	updates := map[string]any{}
	if in.Title != nil {
		updates["title"] = *in.Title
	}
	if in.Content != nil {
		updates["content"] = *in.Content
	}
	if in.Published != nil {
		updates["published"] = *in.Published
	}

	if len(updates) > 0 {
		if err := model.DB.Model(&announcement).Updates(updates).Error; err != nil {
			return nil, response.ErrInternal
		}
	}
	return &announcement, nil
}

// DeleteAnnouncement 删除公告（软删除）
func DeleteAnnouncement(id uint) *response.Errno {
	var announcement model.Announcement
	if err := model.DB.First(&announcement, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ErrAnnouncementNotFound
		}
		return response.ErrInternal
	}
	if err := model.DB.Delete(&announcement).Error; err != nil {
		return response.ErrInternal
	}
	return nil
}

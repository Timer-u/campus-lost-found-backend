package service

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"campus-lost-found-backend/model"
	"campus-lost-found-backend/pkg/response"
	"campus-lost-found-backend/pkg/util"
)

// AdminItemListQuery 管理员物品列表查询参数
type AdminItemListQuery struct {
	Page         int
	PageSize     int
	ReviewStatus string
	ItemStatus   string
	Keyword      string
}

// AdminListItems 管理员物品列表：可见全部状态，默认待审核优先、按更新时间倒序
func AdminListItems(q AdminItemListQuery) ([]model.Item, util.PageMeta, *response.Errno) {
	if q.ReviewStatus != "" && !reviewStatusAllowed[q.ReviewStatus] {
		return nil, util.PageMeta{}, response.ErrInvalidParams
	}
	if q.ItemStatus != "" && !itemStatusAllowed[q.ItemStatus] {
		return nil, util.PageMeta{}, response.ErrInvalidParams
	}

	query := model.DB.Model(&model.Item{})
	if q.ReviewStatus != "" {
		query = query.Where("review_status = ?", q.ReviewStatus)
	}
	if q.ItemStatus != "" {
		query = query.Where("item_status = ?", q.ItemStatus)
	}
	if q.Keyword != "" {
		kw := "%" + q.Keyword + "%"
		query = query.Where("title LIKE ? OR description LIKE ? OR location LIKE ?", kw, kw, kw)
	}

	base := query.Session(&gorm.Session{})
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, util.PageMeta{}, response.ErrInternal
	}

	// 待审核排在最前，其余按更新时间倒序
	var items []model.Item
	if err := base.Preload("User").
		Order("review_status = 'pending' DESC, updated_at DESC, id DESC").
		Offset((q.Page - 1) * q.PageSize).Limit(q.PageSize).Find(&items).Error; err != nil {
		return nil, util.PageMeta{}, response.ErrInternal
	}
	for i := range items {
		items[i].PublisherName = items[i].User.Name
	}

	return items, util.NewPageMeta(int64(q.Page), int64(q.PageSize), total), nil
}

// AdminReviewItem 审核物品：approved 公开、rejected 需填原因、offline 下架；重复审核返回 409
func AdminReviewItem(itemID uint, reviewStatus, rejectReason string) (*model.Item, *response.Errno) {
	if !reviewStatusAllowed[reviewStatus] || reviewStatus == "pending" {
		return nil, response.ErrInvalidParams
	}

	var item model.Item
	if err := model.DB.First(&item, itemID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.ErrItemNotFound
		}
		return nil, response.ErrInternal
	}
	if item.ReviewStatus == reviewStatus {
		return nil, response.ErrItemStatusNotAllowed
	}

	updates := map[string]any{"review_status": reviewStatus, "reject_reason": rejectReason}
	if err := model.DB.Model(&item).Updates(updates).Error; err != nil {
		return nil, response.ErrInternal
	}

	if err := model.DB.Preload("User").First(&item, itemID).Error; err != nil {
		return nil, response.ErrInternal
	}
	item.PublisherName = item.User.Name
	return &item, nil
}

// AdminUpdateItemStatus 管理员调整物品状态：仅审核通过的物品可调整
func AdminUpdateItemStatus(itemID uint, itemStatus string) (*model.Item, *response.Errno) {
	if !itemStatusAllowed[itemStatus] {
		return nil, response.ErrInvalidParams
	}

	var item model.Item
	if err := model.DB.First(&item, itemID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.ErrItemNotFound
		}
		return nil, response.ErrInternal
	}
	if item.ReviewStatus != "approved" {
		return nil, response.ErrItemStatusNotAllowed
	}
	if item.ItemStatus == itemStatus {
		return nil, response.ErrItemStatusNotAllowed
	}

	if err := model.DB.Model(&item).Update("item_status", itemStatus).Error; err != nil {
		return nil, response.ErrInternal
	}

	if err := model.DB.Preload("User").First(&item, itemID).Error; err != nil {
		return nil, response.ErrInternal
	}
	item.PublisherName = item.User.Name
	return &item, nil
}

// AdminReviewClaim 审核认领申请：通过与拒绝仅限 pending 的申请。
// 通过时在事务中锁定申请并检查同物品是否已有通过的申请，防止一物多通过；
// 通过后物品状态同步置为 claimed。
func AdminReviewClaim(adminID, claimID uint, status, reviewReason string) (*model.Claim, *response.Errno) {
	if status != "approved" && status != "rejected" {
		return nil, response.ErrInvalidParams
	}

	var claim model.Claim
	err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 锁定申请行，避免并发审核同一申请
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&claim, claimID).Error; err != nil {
			return err
		}
		if claim.Status != "pending" {
			return response.ErrClaimStatusNotAllowed
		}

		if status == "approved" {
			// 同一物品不允许出现多个通过的申请
			var approvedCount int64
			if err := tx.Model(&model.Claim{}).
				Where("item_id = ? AND status = ? AND id <> ?", claim.ItemID, "approved", claimID).
				Count(&approvedCount).Error; err != nil {
				return err
			}
			if approvedCount > 0 {
				return response.ErrClaimStatusNotAllowed
			}

			if err := tx.Model(&model.Item{}).Where("id = ?", claim.ItemID).
				Update("item_status", "claimed").Error; err != nil {
				return err
			}
		}

		return tx.Model(&claim).Updates(map[string]any{
			"status":        status,
			"review_reason": reviewReason,
			"reviewed_by":   adminID,
		}).Error
	})
	if err != nil {
		var eno *response.Errno
		if errors.As(err, &eno) {
			return nil, eno
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.ErrClaimNotFound
		}
		return nil, response.ErrInternal
	}

	if err := model.DB.Preload("Applicant").First(&claim, claimID).Error; err != nil {
		return nil, response.ErrInternal
	}
	claim.ApplicantName = claim.Applicant.Name
	return &claim, nil
}

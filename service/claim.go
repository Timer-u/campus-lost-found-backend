package service

import (
	"errors"

	"gorm.io/gorm"

	"campus-lost-found-backend/model"
	"campus-lost-found-backend/pkg/response"
	"campus-lost-found-backend/pkg/util"
)

var claimStatusAllowed = map[string]bool{"pending": true, "approved": true, "rejected": true, "cancelled": true}

// CreateClaimInput 提交认领申请参数，字段约束由 controller binding 校验
type CreateClaimInput struct {
	Description string
	Contact     string
}

// CreateClaim 提交认领申请：仅审核通过且开放中的招领可以申请，
// 发布者本人、重复申请与状态不符统一返回 13001（文档该接口未定义 403/其他 409 码）
func CreateClaim(applicantID, itemID uint, in CreateClaimInput) (*model.Claim, *response.Errno) {
	var item model.Item
	if err := model.DB.First(&item, itemID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.ErrItemNotFound
		}
		return nil, response.ErrInternal
	}
	if item.Type != "found" || item.ReviewStatus != "approved" || item.ItemStatus != "open" {
		return nil, response.ErrClaimNotAllowed
	}
	if item.UserID == applicantID {
		return nil, response.ErrClaimNotAllowed
	}

	// 同一用户对同一物品只允许存在一个有效申请（已取消的可以重新提交）
	var count int64
	if err := model.DB.Model(&model.Claim{}).
		Where("item_id = ? AND applicant_id = ? AND status <> ?", itemID, applicantID, "cancelled").
		Count(&count).Error; err != nil {
		return nil, response.ErrInternal
	}
	if count > 0 {
		return nil, response.ErrClaimNotAllowed
	}

	claim := &model.Claim{
		ItemID:      itemID,
		ApplicantID: applicantID,
		Description: in.Description,
		Contact:     in.Contact,
		Status:      "pending",
	}
	if err := model.DB.Create(claim).Error; err != nil {
		return nil, response.ErrInternal
	}
	if user, eno := GetByID(applicantID); eno == nil {
		claim.ApplicantName = user.Name
	}
	return claim, nil
}

// ListItemClaims 物品的认领申请列表：仅发布者或管理员可见
func ListItemClaims(viewerID uint, viewerRole string, itemID uint, page, pageSize int) ([]model.Claim, util.PageMeta, *response.Errno) {
	var item model.Item
	if err := model.DB.Select("id", "user_id").First(&item, itemID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.PageMeta{}, response.ErrItemNotFound
		}
		return nil, util.PageMeta{}, response.ErrInternal
	}
	if item.UserID != viewerID && viewerRole != "lost_admin" && viewerRole != "system_admin" {
		return nil, util.PageMeta{}, response.ErrForbidden
	}

	base := model.DB.Model(&model.Claim{}).Where("item_id = ?", itemID).Session(&gorm.Session{})
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, util.PageMeta{}, response.ErrInternal
	}

	var claims []model.Claim
	if err := base.Preload("Applicant").Order("created_at ASC, id ASC").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&claims).Error; err != nil {
		return nil, util.PageMeta{}, response.ErrInternal
	}
	for i := range claims {
		claims[i].ApplicantName = claims[i].Applicant.Name
	}

	return claims, util.NewPageMeta(int64(page), int64(pageSize), total), nil
}

// ListMyClaims 我的认领申请：status 筛选 + 分页
func ListMyClaims(applicantID uint, status string, page, pageSize int) ([]model.Claim, util.PageMeta, *response.Errno) {
	if status != "" && !claimStatusAllowed[status] {
		return nil, util.PageMeta{}, response.ErrInvalidParams
	}

	query := model.DB.Model(&model.Claim{}).Where("applicant_id = ?", applicantID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	base := query.Session(&gorm.Session{})
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, util.PageMeta{}, response.ErrInternal
	}

	var claims []model.Claim
	if err := base.Preload("Applicant").Order("created_at DESC, id DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&claims).Error; err != nil {
		return nil, util.PageMeta{}, response.ErrInternal
	}
	for i := range claims {
		claims[i].ApplicantName = claims[i].Applicant.Name
	}

	return claims, util.NewPageMeta(int64(page), int64(pageSize), total), nil
}

// CancelClaim 取消自己的认领申请：仅 pending 状态可以取消
func CancelClaim(applicantID, claimID uint) (*model.Claim, *response.Errno) {
	var claim model.Claim
	if err := model.DB.First(&claim, claimID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.ErrClaimNotFound
		}
		return nil, response.ErrInternal
	}
	if claim.ApplicantID != applicantID {
		return nil, response.ErrForbidden
	}
	if claim.Status != "pending" {
		return nil, response.ErrClaimStatusNotAllowed
	}

	claim.Status = "cancelled"
	if err := model.DB.Model(&claim).Update("status", "cancelled").Error; err != nil {
		return nil, response.ErrInternal
	}
	return &claim, nil
}

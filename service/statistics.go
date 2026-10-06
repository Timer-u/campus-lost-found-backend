package service

import (
	"campus-lost-found-backend/model"
	"campus-lost-found-backend/pkg/response"
)

// Statistics 统计总览，字段与 docs/openapi.yaml 的 StatisticsResponse 一致
type Statistics struct {
	UserCount         int64            `json:"userCount"`
	ItemCount         int64            `json:"itemCount"`
	PendingItemCount  int64            `json:"pendingItemCount"`
	ClaimCount        int64            `json:"claimCount"`
	ResolvedItemCount int64            `json:"resolvedItemCount"`
	ItemsByType       map[string]int64 `json:"itemsByType"`
}

// GetStatisticsOverview 系统统计总览
func GetStatisticsOverview() (*Statistics, *response.Errno) {
	stats := &Statistics{ItemsByType: map[string]int64{}}

	if err := model.DB.Model(&model.User{}).Count(&stats.UserCount).Error; err != nil {
		return nil, response.ErrInternal
	}
	if err := model.DB.Model(&model.Item{}).Count(&stats.ItemCount).Error; err != nil {
		return nil, response.ErrInternal
	}
	if err := model.DB.Model(&model.Item{}).Where("review_status = ?", "pending").
		Count(&stats.PendingItemCount).Error; err != nil {
		return nil, response.ErrInternal
	}
	if err := model.DB.Model(&model.Claim{}).Count(&stats.ClaimCount).Error; err != nil {
		return nil, response.ErrInternal
	}
	if err := model.DB.Model(&model.Item{}).Where("item_status = ?", "resolved").
		Count(&stats.ResolvedItemCount).Error; err != nil {
		return nil, response.ErrInternal
	}

	var byType []struct {
		Type  string
		Count int64
	}
	if err := model.DB.Model(&model.Item{}).
		Select("type, COUNT(*) AS count").Group("type").Scan(&byType).Error; err != nil {
		return nil, response.ErrInternal
	}
	for _, row := range byType {
		stats.ItemsByType[row.Type] = row.Count
	}

	return stats, nil
}

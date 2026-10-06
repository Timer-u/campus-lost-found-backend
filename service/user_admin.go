package service

import (
	"errors"

	"gorm.io/gorm"

	"campus-lost-found-backend/model"
	"campus-lost-found-backend/pkg/response"
	"campus-lost-found-backend/pkg/util"
)

var (
	userRoleAllowed   = map[string]bool{"student": true, "lost_admin": true, "system_admin": true}
	userStatusAllowed = map[string]bool{"active": true, "disabled": true}
)

// AdminUserListQuery 管理员用户列表查询参数
type AdminUserListQuery struct {
	Page     int
	PageSize int
	Keyword  string
	Role     string
	Status   string
}

// AdminListUsers 用户列表：keyword 按学号或姓名模糊搜索，role/status 筛选
func AdminListUsers(q AdminUserListQuery) ([]model.User, util.PageMeta, *response.Errno) {
	if q.Role != "" && !userRoleAllowed[q.Role] {
		return nil, util.PageMeta{}, response.ErrInvalidParams
	}
	if q.Status != "" && !userStatusAllowed[q.Status] {
		return nil, util.PageMeta{}, response.ErrInvalidParams
	}

	query := model.DB.Model(&model.User{})
	if q.Keyword != "" {
		kw := "%" + q.Keyword + "%"
		query = query.Where("username LIKE ? OR name LIKE ?", kw, kw)
	}
	if q.Role != "" {
		query = query.Where("role = ?", q.Role)
	}
	if q.Status != "" {
		query = query.Where("status = ?", q.Status)
	}

	base := query.Session(&gorm.Session{})
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, util.PageMeta{}, response.ErrInternal
	}

	var users []model.User
	if err := base.Order("created_at ASC, id ASC").
		Offset((q.Page - 1) * q.PageSize).Limit(q.PageSize).Find(&users).Error; err != nil {
		return nil, util.PageMeta{}, response.ErrInternal
	}

	return users, util.NewPageMeta(int64(q.Page), int64(q.PageSize), total), nil
}

// AdminUpdateUserRole 修改用户角色
func AdminUpdateUserRole(targetUserID uint, role string) (*model.User, *response.Errno) {
	if !userRoleAllowed[role] {
		return nil, response.ErrInvalidParams
	}

	var user model.User
	if err := model.DB.First(&user, targetUserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.ErrUserNotFound
		}
		return nil, response.ErrInternal
	}

	if err := model.DB.Model(&user).Update("role", role).Error; err != nil {
		return nil, response.ErrInternal
	}
	user.Role = role
	return &user, nil
}

// AdminUpdateUserStatus 启用或禁用用户
func AdminUpdateUserStatus(targetUserID uint, status string) (*model.User, *response.Errno) {
	if !userStatusAllowed[status] {
		return nil, response.ErrInvalidParams
	}

	var user model.User
	if err := model.DB.First(&user, targetUserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.ErrUserNotFound
		}
		return nil, response.ErrInternal
	}

	if err := model.DB.Model(&user).Update("status", status).Error; err != nil {
		return nil, response.ErrInternal
	}
	user.Status = status
	return &user, nil
}

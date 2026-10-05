package service

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"campus-lost-found-backend/model"
	"campus-lost-found-backend/pkg/auth"
	"campus-lost-found-backend/pkg/response"
)

// LoginResult 登录成功后的令牌与用户信息，字段名与 docs/openapi.yaml 的 LoginData 一致
type LoginResult struct {
	Token     string      `json:"accessToken"`
	TokenType string      `json:"tokenType"`
	ExpiresIn int64       `json:"expiresIn"`
	User      *model.User `json:"user"`
}

// Register 注册学生账号：学号查重，密码 bcrypt 加密后入库
func Register(username, name, password string) (*model.User, *response.Errno) {
	var count int64
	if err := model.DB.Model(&model.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return nil, response.ErrInternal
	}
	if count > 0 {
		return nil, response.ErrUsernameExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, response.ErrInternal
	}

	user := &model.User{
		Username:     username,
		Name:         name,
		PasswordHash: string(hash),
		Role:         "student",
		Status:       "active",
	}
	if err := model.DB.Create(user).Error; err != nil {
		// 查重与写入之间存在并发窗口，唯一索引兜底
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, response.ErrUsernameExists
		}
		return nil, response.ErrInternal
	}
	return user, nil
}

// Login 校验学号密码并签发 JWT；用户不存在与密码错误返回同一错误码，避免暴露账号是否存在
func Login(username, password string) (*LoginResult, *response.Errno) {
	var user model.User
	if err := model.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.ErrInvalidCredentials
		}
		return nil, response.ErrInternal
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, response.ErrInvalidCredentials
	}

	if user.Status != "active" {
		return nil, response.ErrAccountDisabled
	}

	token, expiresIn, err := auth.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, response.ErrInternal
	}

	return &LoginResult{Token: token, TokenType: "Bearer", ExpiresIn: expiresIn, User: &user}, nil
}

// GetByID 按主键查询用户
func GetByID(id uint) (*model.User, *response.Errno) {
	var user model.User
	if err := model.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.ErrUserNotFound
		}
		return nil, response.ErrInternal
	}
	return &user, nil
}

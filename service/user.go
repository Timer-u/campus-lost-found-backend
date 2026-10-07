package service

import (
	"errors"
	"fmt"
	"regexp"

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

// 学号固定为 12 位纯数字（如 302026315155），具体分段含义不做校验
var studentIDRegexp = regexp.MustCompile(`^[0-9]{12}$`)

// Register 注册学生账号：学号查重，密码 bcrypt 加密后入库
func Register(username, name, password string) (*model.User, *response.Errno) {
	if !studentIDRegexp.MatchString(username) {
		return nil, response.ErrInvalidParams.WithMsg("学号必须为 12 位数字")
	}

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

// UpdateProfile 修改当前登录用户的个人信息（当前支持修改姓名）
func UpdateProfile(userID uint, name string) (*model.User, *response.Errno) {
	var user model.User
	if err := model.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.ErrUserNotFound
		}
		return nil, response.ErrInternal
	}

	if err := model.DB.Model(&user).Update("name", name).Error; err != nil {
		return nil, response.ErrInternal
	}
	user.Name = name
	return &user, nil
}

// DeleteAccount 注销账号：验证密码后软删除当前账号，级联下架名下物品并取消待审核的认领申请。
// 管理员账号不允许自我注销；注销时改写学号以释放唯一索引，允许该学号日后重新注册。
func DeleteAccount(userID uint, password string) *response.Errno {
	var user model.User
	if err := model.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ErrUserNotFound
		}
		return response.ErrInternal
	}

	// 危险操作：验证当前密码，防止设备被他人滥用
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return response.ErrInvalidCredentials
	}
	if user.Role != "student" {
		return response.ErrAdminUndeletable
	}

	err := model.DB.Transaction(func(tx *gorm.DB) error {
		// 改写学号，释放唯一索引
		if err := tx.Model(&user).Update("username",
			fmt.Sprintf("%s#deleted%d", user.Username, user.ID)).Error; err != nil {
			return err
		}
		// 级联下架名下物品（软删除）
		if err := tx.Where("user_id = ?", userID).Delete(&model.Item{}).Error; err != nil {
			return err
		}
		// 待审核的认领申请自动取消
		if err := tx.Model(&model.Claim{}).
			Where("applicant_id = ? AND status = ?", userID, "pending").
			Update("status", "cancelled").Error; err != nil {
			return err
		}
		// 软删除账号
		return tx.Delete(&user).Error
	})
	if err != nil {
		return response.ErrInternal
	}
	return nil
}

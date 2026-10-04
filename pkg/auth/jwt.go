package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"campus-lost-found-backend/config"
)

// ErrInvalidToken token 缺失、过期或签名不符时统一返回
var ErrInvalidToken = errors.New("无效的登录凭证")

// Claims 登录凭证负载：身份信息 + 标准声明
type Claims struct {
	UserID   uint   `json:"userId"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// expireDuration 有效期取配置值，未配置时兜底 2 小时（文档示例 expiresIn=7200 秒）
func expireDuration() time.Duration {
	if config.GlobalConfig.JWT.ExpireHours > 0 {
		return time.Duration(config.GlobalConfig.JWT.ExpireHours) * time.Hour
	}
	return 2 * time.Hour
}

// GenerateToken 签发 JWT，返回 token 字符串和有效期（秒）
func GenerateToken(userID uint, username, role string) (string, int64, error) {
	d := expireDuration()
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(d)),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   username,
		},
	}
	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(config.GlobalConfig.JWT.Secret))
	if err != nil {
		return "", 0, err
	}
	return signed, int64(d.Seconds()), nil
}

// ParseToken 校验并解析 JWT，只接受 HS256 签名
func ParseToken(tokenString string) (*Claims, error) {
	parsed, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(*jwt.Token) (any, error) {
		return []byte(config.GlobalConfig.JWT.Secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, ErrInvalidToken
	}
	claims, ok := parsed.Claims.(*Claims)
	if !ok || !parsed.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

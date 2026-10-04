package auth

import (
	"testing"

	"campus-lost-found-backend/config"
)

func TestGenerateAndParseToken(t *testing.T) {
	config.GlobalConfig.JWT = config.JWTConfig{Secret: "test-secret", ExpireHours: 1}

	token, expiresIn, err := GenerateToken(7, "20260001", "student")
	if err != nil {
		t.Fatalf("签发失败: %v", err)
	}
	if token == "" {
		t.Fatal("token 不应为空")
	}
	if expiresIn != 3600 {
		t.Fatalf("有效期应为 3600 秒，实际 %d", expiresIn)
	}

	claims, err := ParseToken(token)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}
	if claims.UserID != 7 || claims.Username != "20260001" || claims.Role != "student" {
		t.Fatalf("负载不符: %+v", claims)
	}
}

func TestParseTokenWrongSecret(t *testing.T) {
	config.GlobalConfig.JWT = config.JWTConfig{Secret: "secret-a", ExpireHours: 1}
	token, _, err := GenerateToken(1, "20260001", "student")
	if err != nil {
		t.Fatal(err)
	}

	config.GlobalConfig.JWT.Secret = "secret-b"
	if _, err := ParseToken(token); err == nil {
		t.Fatal("密钥不一致时应当校验失败")
	}
}

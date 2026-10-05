package service

import (
	"encoding/json"
	"strings"
	"testing"

	"campus-lost-found-backend/model"
)

// LoginResult 序列化后的字段名必须与 docs/openapi.yaml 的 LoginData 一致
func TestLoginResultJSONShape(t *testing.T) {
	result := LoginResult{
		Token:     "header.payload.signature",
		TokenType: "Bearer",
		ExpiresIn: 7200,
		User:      &model.User{Username: "302026315155", Name: "张三", Role: "student", Status: "active"},
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("序列化失败: %v", err)
	}
	body := string(data)

	for _, key := range []string{`"accessToken"`, `"tokenType":"Bearer"`, `"expiresIn":7200`, `"user"`} {
		if !strings.Contains(body, key) {
			t.Fatalf("缺少契约字段 %s: %s", key, body)
		}
	}
	if strings.Contains(body, "PasswordHash") || strings.Contains(body, "password") && strings.Contains(body, "hash") {
		t.Fatalf("密码字段不应出现在登录结果中: %s", body)
	}
}

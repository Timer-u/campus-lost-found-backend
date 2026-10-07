package response

import "testing"

// 所有错误码必须唯一，且 msg 与 HTTP 状态码不能为空/越界
func TestErrnoCodesUnique(t *testing.T) {
	seen := make(map[int]string)

	for _, e := range []*Errno{
		ErrInvalidParams, ErrUnauthorized, ErrForbidden, ErrAdminUndeletable,
		ErrUsernameExists, ErrInvalidCredentials, ErrAccountDisabled,
		ErrNeedLostAdmin, ErrNeedSystemAdmin,
		ErrItemNotFound, ErrItemStatusNotAllowed,
		ErrClaimNotAllowed, ErrClaimNotFound, ErrClaimStatusNotAllowed,
		ErrUserNotFound, ErrAnnouncementNotFound, ErrImageTooLarge, ErrInternal,
	} {
		if prev, ok := seen[e.Code]; ok {
			t.Fatalf("错误码 %d 重复使用: %s / %s", e.Code, prev, e.Msg)
		}
		seen[e.Code] = e.Msg

		if e.Msg == "" {
			t.Fatalf("错误码 %d 的 msg 不能为空", e.Code)
		}
		if e.HTTP < 400 || e.HTTP > 599 {
			t.Fatalf("错误码 %d 的 HTTP 状态码 %d 越界", e.Code, e.HTTP)
		}
	}
}

func TestWithMsgKeepsCodeAndHTTP(t *testing.T) {
	e := ErrInvalidParams.WithMsg("无效的物品ID")
	if e.Code != ErrInvalidParams.Code || e.HTTP != ErrInvalidParams.HTTP {
		t.Fatal("WithMsg 不应改变 Code 和 HTTP 状态码")
	}
	if e.Msg != "无效的物品ID" {
		t.Fatalf("Msg 应被替换，实际 %q", e.Msg)
	}
}

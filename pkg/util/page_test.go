package util

import "testing"

func TestParsePageDefaults(t *testing.T) {
	page, size := ParsePage("", "")
	if page != DefaultPage || size != DefaultPageSize {
		t.Fatalf("缺省值应为 1/10，实际 %d/%d", page, size)
	}
}

func TestParsePageClamp(t *testing.T) {
	page, size := ParsePage("0", "999")
	if page != 1 || size != MaxPageSize {
		t.Fatalf("page 应回落为 1，size 应截断为 50，实际 %d/%d", page, size)
	}

	page, size = ParsePage("abc", "-5")
	if page != 1 || size != DefaultPageSize {
		t.Fatalf("非法值应回落默认，实际 %d/%d", page, size)
	}
}

func TestNewPageMeta(t *testing.T) {
	meta := NewPageMeta(2, 10, 25)
	if meta.TotalPages != 3 {
		t.Fatalf("25 条按 10 条分页应为 3 页，实际 %d", meta.TotalPages)
	}

	meta = NewPageMeta(1, 10, 0)
	if meta.TotalPages != 0 {
		t.Fatalf("无数据时 totalPages 应为 0，实际 %d", meta.TotalPages)
	}
}

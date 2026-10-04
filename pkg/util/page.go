package util

import "strconv"

const (
	DefaultPage     = 1
	DefaultPageSize = 10
	MaxPageSize     = 50
)

// PageMeta 分页元信息，结构与 docs/openapi.yaml 的 PageMeta 一致
type PageMeta struct {
	Page       int64 `json:"page"`
	PageSize   int64 `json:"pageSize"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"totalPages"`
}

// ParsePage 解析分页参数：缺省或非法值回落默认，pageSize 夹在 [1, 50]
func ParsePage(pageStr, sizeStr string) (page, pageSize int) {
	page = DefaultPage
	if v, err := strconv.Atoi(pageStr); err == nil && v >= 1 {
		page = v
	}

	pageSize = DefaultPageSize
	if v, err := strconv.Atoi(sizeStr); err == nil {
		if v < 1 {
			v = DefaultPageSize
		}
		if v > MaxPageSize {
			v = MaxPageSize
		}
		pageSize = v
	}
	return page, pageSize
}

// NewPageMeta 由记录总数计算分页元信息
func NewPageMeta(page, pageSize, total int64) PageMeta {
	var totalPages int64
	if pageSize > 0 && total > 0 {
		totalPages = (total + pageSize - 1) / pageSize
	}
	return PageMeta{Page: page, PageSize: pageSize, Total: total, TotalPages: totalPages}
}

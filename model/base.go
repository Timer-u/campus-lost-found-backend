package model

import (
	"time"

	"gorm.io/gorm"
)

// Base 公共字段，json 名与 docs/openapi.yaml 保持一致（createdAt 等）
type Base struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

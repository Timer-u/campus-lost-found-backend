package model

import (
	"gorm.io/gorm"
)

type Item struct {
	gorm.Model

	Title       string `gorm:"type:varchar(100);not null" json:"title,omitempty"`
	Type        string `gorm:"type:varchar(20);not null" json:"type,omitempty"`
	Category    string `gorm:"type:varchar(20);not null" json:"category,omitempty"`
	Location    string `gorm:"type:varchar(100);not null" json:"location,omitempty"`
	Description string `gorm:"type:text" json:"description,omitempty"`
	ImageURL    string `gorm:"type:varchar(255)" json:"image_url,omitempty"`
	Status      string `gorm:"type:varchar(20);default:'published'" json:"status,omitempty"`
	UserID      uint   `gorm:"not null" json:"user_id,omitempty"`
}

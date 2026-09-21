package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string `gorm:"type:varchar(32);unique;not null" json:"username,omitempty"`
	PasswordHash string `gorm:"type:varchar(255);not null" json:"-"`
	StudentID    string `gorm:"type:varchar(32);not null;unique" json:"student_id,omitempty"`
	RealName     string `gorm:"type:varchar(32)" json:"real_name,omitempty"`
	Phone        string `gorm:"type:varchar(20)" json:"phone,omitempty"`
	Role         string `gorm:"type:varchar(20);default:'student'" json:"role,omitempty"`
}

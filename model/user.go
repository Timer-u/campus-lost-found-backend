package model

// User 用户表，username 即学号（12 位纯数字），见 docs/openapi.yaml
type User struct {
	Base
	Username     string `gorm:"type:varchar(32);unique;not null" json:"username"`
	PasswordHash string `gorm:"type:varchar(255);not null" json:"-"`
	Name         string `gorm:"type:varchar(50)" json:"name"`
	Role         string `gorm:"type:varchar(20);default:'student'" json:"role"`
	Status       string `gorm:"type:varchar(20);default:'active'" json:"status"`
}

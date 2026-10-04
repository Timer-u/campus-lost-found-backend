package model

// Announcement 公告表
type Announcement struct {
	Base
	Title     string `gorm:"type:varchar(100);not null" json:"title"`
	Content   string `gorm:"type:text;not null" json:"content"`
	Published bool   `gorm:"default:false" json:"published"`
	AuthorID  uint   `gorm:"not null" json:"authorId"`
}

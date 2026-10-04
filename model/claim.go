package model

// Claim 认领申请表
type Claim struct {
	Base
	ItemID       uint   `gorm:"not null;index" json:"itemId"`
	ApplicantID  uint   `gorm:"not null" json:"applicantId"`
	Description  string `gorm:"type:varchar(1000);not null" json:"description"`
	Contact      string `gorm:"type:varchar(100);not null" json:"contact"`
	Status       string `gorm:"type:varchar(20);default:'pending'" json:"status"`
	ReviewReason string `gorm:"type:varchar(255)" json:"reviewReason,omitempty"`
	ReviewedBy   *uint  `json:"reviewedBy,omitempty"`
}

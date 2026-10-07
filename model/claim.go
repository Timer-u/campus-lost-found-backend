package model

// Claim 认领申请表
type Claim struct {
	Base
	ItemID uint `gorm:"not null;index" json:"itemId"`
	// ItemTitle 所属物品标题，列表接口由 service 联表填充
	ItemTitle    string `gorm:"-" json:"itemTitle,omitempty"`
	ApplicantID  uint   `gorm:"not null" json:"applicantId"`
	Description  string `gorm:"type:varchar(1000);not null" json:"description"`
	Contact      string `gorm:"type:varchar(100);not null" json:"contact"`
	Status       string `gorm:"type:varchar(20);default:'pending'" json:"status"`
	ReviewReason string `gorm:"type:varchar(255)" json:"reviewReason,omitempty"`
	ReviewedBy   *uint  `json:"reviewedBy,omitempty"`
	// ApplicantName 申请人姓名，查询时由 service 通过 Preload("Applicant") 填充
	ApplicantName string `gorm:"-" json:"applicantName"`
	Applicant     User   `gorm:"foreignKey:ApplicantID" json:"-"`
}

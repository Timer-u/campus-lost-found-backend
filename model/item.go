package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// ImageList 物品图片列表（最多 6 张），数据库中以 JSON 数组存储
type ImageList []string

func (l ImageList) Value() (driver.Value, error) {
	if len(l) == 0 {
		return "[]", nil
	}
	b, err := json.Marshal(l)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

func (l *ImageList) Scan(v any) error {
	if v == nil {
		*l = ImageList{}
		return nil
	}
	var bytes []byte
	switch data := v.(type) {
	case []byte:
		bytes = data
	case string:
		bytes = []byte(data)
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
	if err := json.Unmarshal(bytes, l); err != nil {
		return errors.New("解析图片列表失败: " + err.Error())
	}
	return nil
}

// Item 失物/招领信息表，审核状态与物品状态字段见 docs/openapi.yaml
type Item struct {
	Base
	Type         string     `gorm:"type:varchar(20);not null" json:"type"`
	Title        string     `gorm:"type:varchar(100);not null" json:"title"`
	Description  string     `gorm:"type:text;not null" json:"description"`
	Category     string     `gorm:"type:varchar(20);not null" json:"category"`
	Location     string     `gorm:"type:varchar(100);not null" json:"location"`
	LostAt       *time.Time `gorm:"type:datetime" json:"lostAt,omitempty"`
	Contact      string     `gorm:"type:varchar(100);not null" json:"contact"`
	ImageURLs    ImageList  `gorm:"type:json" json:"imageUrls"`
	ReviewStatus string     `gorm:"type:varchar(20);default:'pending'" json:"reviewStatus"`
	RejectReason string     `gorm:"type:varchar(255)" json:"rejectReason,omitempty"`
	ItemStatus   string     `gorm:"type:varchar(20);default:'open'" json:"itemStatus"`
	UserID       uint       `gorm:"not null" json:"publisherId"`
	// PublisherName 发布者姓名，查询时由 service 通过 Preload("User") 填充
	PublisherName string `gorm:"-" json:"publisherName"`
	User          User   `gorm:"foreignKey:UserID" json:"-"`
}

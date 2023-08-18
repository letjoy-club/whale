package models

import "time"

type MatchingOfferSummary struct {
	MatchingID string `gorm:"primaryKey;type:varchar(32)"`
	UserID     string `gorm:"index;type:varchar(32)"`

	// 用户收到的匹配意向数量
	InOfferNum int
	// 未处理的匹配意向数量
	// UnprocessedInOfferNum int
	// 用户收到的匹配意向
	// InMatchingOfferIDs []string `gorm:"type:json;serializer:json"`

	// 用户发出的匹配意向数量
	OutOfferNum int
	// // 对方未处理的匹配意向数量
	// UnprocessedOutOfferNum int
	// 用户发出的匹配意向
	// OutMatchingOfferIDs []string `gorm:"type:json;serializer:json"`

	BasicQuota  int
	RemainQuota int

	Active bool

	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type MatchingOfferRecord struct {
	ID int `gorm:"primaryKey"`

	ToMatchingID string `gorm:"type:varchar(32);index"`
	ToUserID     string `gorm:"type:varchar(32);index"`

	MatchingID string `gorm:"index"`
	UserID     string `gorm:"type:varchar(32);index"`

	// Unprocessed, Accepted, Rejected, Canceled
	State     string `gorm:"type:varchar(32)"`
	ReactedAt *time.Time

	Remark string `gorm:"text"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
}

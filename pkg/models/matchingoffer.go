package models

import "time"

type MatchingOfferSummary struct {
	MatchingID string `gorm:"primaryKey"`
	UserID     string `gorm:"index"`

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

	ToMatchingID string `gorm:"index"`
	ToUserID     string

	MatchingID string `gorm:"index"`
	UserID     string

	// Unprocessed, Accepted, Rejected, Canceled
	State     string
	ReactedAt *time.Time

	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type MatchingOffer struct {
	ID string `json:"id"`
	// Unprocessed, Accepted, Rejected, Canceled
	State string `json:"state"`
}

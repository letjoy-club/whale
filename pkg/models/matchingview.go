package models

import "time"

type UserViewMatching struct {
	ID     int    `gorm:"primaryKey"`
	UserID string `gorm:"index"`

	ViewedMatchingIDs []string  `gorm:"type:json"`
	CreatedAt         time.Time `gorm:"autoCreateTime"`
}

type MatchingView struct {
	MatchingID string `gorm:"primaryKey"`

	// 查看次数
	ViewCount int
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type MatchingViewHistory struct {
	ID int `gorm:"primaryKey"`

	ViewedMatchingID string `gorm:"index"`
	UserID           string `gorm:"index"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
}

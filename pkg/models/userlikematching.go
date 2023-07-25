package models

import "time"

type UserLikeMatching struct {
	ID           int    `gorm:"primaryKey"`
	ToMatchingID string `gorm:"index:user_like_matching,unique,priority:2"`
	ToUserID     string

	UserID    string    `gorm:"index:user_like_matching,unique,priority:1"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type MatchingReceiveLike struct {
	MatchingID    string `gorm:"primaryKey"`
	LikeNum       int
	RecentUserIDs []string `gorm:"type:json"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

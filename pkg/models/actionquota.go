package models

import "time"

type ActionQuota struct {
	UserID string `gorm:"type:varchar(32);primaryKey;uniqueIndex:user_day"`

	CreationCount int
	ReactionCount int

	Day string `gorm:"type:varchar(16);uniqueIndex:user_day"`

	UpdatedAt time.Time
	CreatedAt time.Time
}

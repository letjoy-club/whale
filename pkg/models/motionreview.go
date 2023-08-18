package models

import "time"

type MotionReview struct {
	ID int `gorm:"primaryKey"`

	MatchingResultID int
	MotionID         string `gorm:"index:from_motion_review;type:varchar(32)"`
	UserID           string `gorm:"index:from_motion_review;type:varchar(32)"`

	ToMotionID string `gorm:"index:to_motion_review;type:varchar(32)"`
	ToUserID   string `gorm:"index:to_motion_review;type:varchar(32)"`
	TopicID    string `gorm:"type:varchar(32)"`

	Score     int
	Comment   string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

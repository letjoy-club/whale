package models

import "time"

type MotionReview struct {
	ID            int    `gorm:"primaryKey"`
	MotionOfferID int    `gorm:"index:idx_motion_offer_id"`
	ReviewerID    string `gorm:"index:idx_reviewer_id;type:varchar(32)"`
	ToUserID      string `gorm:"index:idx_to_user_id;type:varchar(32)"`
	TopicID       string `gorm:"type:varchar(32)"`
	Score         int
	Comment       string    `gorm:"type:text"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
}

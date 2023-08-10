package models

import "time"

type UserLikeMotion struct {
	ID         int    `gorm:"primaryKey"`
	ToMotionID string `gorm:"index:user_like_motion,unique,priority:2"`
	ToUserID   string `gorm:"type:varchar(32)"`

	UserID    string    `gorm:"type:varchar(32);index:user_like_motion,unique,priority:1"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type UserThumbsUpMotion struct {
	ID         int    `gorm:"primaryKey"`
	ToMotionID string `gorm:"index:user_thumbs_up_motion,unique,priority:2"`
	ToUserID   string `gorm:"type:varchar(32)"`

	UserID    string    `gorm:"type:varchar(32);index:user_thumbs_up_motion,unique,priority:1"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

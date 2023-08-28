package models

import "time"

type DurationConstraint struct {
	ID     int    `gorm:"primaryKey"`
	UserID string `gorm:"type:varchar(64);index"`

	TotalMotionQuota  int
	RemainMotionQuota int

	StartDate time.Time
	StopDate  time.Time

	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

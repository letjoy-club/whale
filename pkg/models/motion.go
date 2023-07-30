package models

import "time"

type Motion struct {
	ID      string   `gorm:"primaryKey"`
	UserID  string   `gorm:"index"`
	TopicID string   `gorm:"index"`
	CityID  string   `gorm:"index"`
	AreaIDs []string `gorm:"type:json;serializer:json"`

	Properties []MotionProperty `gorm:"type:json;serializer:json"`
	Active     bool
	Remark     string

	MyGender string
	Gender   string

	// 特定日期区间，格式 20060102
	DayRange []string `gorm:"serializer:json;type:json"`
	// 优先时间段
	PreferredPeriods []string `gorm:"serializer:json;type:json"`

	ViewCount int `gorm:"default:0"`
	LikeCount int `gorm:"default:0"`

	InOfferNum  int `gorm:"default:0"`
	OutOfferNum int `gorm:"default:0"`

	PendingInNum  int `gorm:"default:0"`
	PendingOutNum int `gorm:"default:0"`
	ActiveNum     int `gorm:"default:0"`

	Discoverable bool

	BasicQuota  int `gorm:"default:0"`
	RemainQuota int `gorm:"default:0"`

	Deadline  time.Time
	UpdatedAt time.Time
	CreatedAt time.Time
}

type RecentLikeMotion struct {
	MotionID string
	UserIDs  []string
}

type LikeMotion struct {
	ID        int    `gorm:"primaryKey"`
	UserID    string `gorm:"index"`
	MotionID  string `gorm:"index"`
	CreatedAt time.Time
}

type MotionProperty struct {
	ID     string   `json:"id"`
	Values []string `json:"values"`
}

type MotionOfferRecord struct {
	ID         int    `gorm:"primaryKey"`
	MotionID   string `gorm:"index:motion_to_motion,unique"`
	ToMotionID string `gorm:"index:motion_to_motion,unique"`

	UserID string
	// Pending, Accepted, Rejected, Timeout, Closed
	State       string
	ChatGroupID string

	ReactAt *time.Time
	Remark  string

	ExpiredAt time.Time
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type MotionViewHistory struct {
	ID        int      `gorm:"primaryKey"`
	UserID    string   `gorm:"index"`
	MotionIDs []string `gorm:"type:json;serializer:json"`
	CreatedAt time.Time
}

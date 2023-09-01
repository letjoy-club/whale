package models

import "time"

type Motion struct {
	ID      string   `gorm:"primaryKey"`
	UserID  string   `gorm:"index;type:varchar(32)"`
	TopicID string   `gorm:"index;type:varchar(32)"`
	CityID  string   `gorm:"index;type:varchar(32)"`
	AreaIDs []string `gorm:"type:json;serializer:json"`

	Properties []MotionProperty `gorm:"type:json;serializer:json"`
	Active     bool
	Remark     string `gorm:"type:varchar(255)"`

	MyGender string `gorm:"type:varchar(32)"`
	Gender   string `gorm:"type:varchar(32)"`

	// 特定日期区间，格式 20060102
	DayRange []string `gorm:"serializer:json;type:json"`
	// 优先时间段
	PreferredPeriods []string `gorm:"serializer:json;type:json"`

	ViewCount     int `gorm:"default:0"`
	LikeCount     int `gorm:"default:0"`
	ThumbsUpCount int `gorm:"default:0"`

	InOfferNum  int `gorm:"default:0"`
	OutOfferNum int `gorm:"default:0"`

	PendingInNum  int `gorm:"default:0"`
	PendingOutNum int `gorm:"default:0"`
	ActiveNum     int `gorm:"default:0"`

	Level int `gorm:"default:0"`

	Discoverable bool

	RelatedMatchingID string `gorm:"type:varchar(32)"`

	Deadline  time.Time
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type RecentLikeMotion struct {
	MotionID string
	UserIDs  []string
}

type MotionProperty struct {
	ID     string   `json:"id"`
	Values []string `json:"values"`
}

type MotionOfferRecord struct {
	ID         int    `gorm:"primaryKey"`
	MotionID   string `gorm:"index:motion_to_motion,unique"`
	ToMotionID string `gorm:"index:motion_to_motion,unique"`

	UserID   string `gorm:"index;type:varchar(32)"`
	ToUserID string `gorm:"index;type:varchar(32)"`

	// Pending, Accepted, Rejected, Canceled, Timeout, Finished
	State       string `gorm:"type:varchar(32)"`
	ChatGroupID string `gorm:"type:varchar(32)"`

	ReactAt *time.Time
	Remark  string `gorm:"type:varchar(255)"`

	ChatChance int `gorm:"default:0"`

	ExpiredAt time.Time
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type MotionViewHistory struct {
	ID        int       `gorm:"primaryKey"`
	UserID    string    `gorm:"index;type:varchar(32)"`
	MotionIDs []string  `gorm:"type:json;serializer:json"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

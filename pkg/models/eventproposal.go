package models

import "time"

type EventProposal struct {
	ID        string   `gorm:"primaryKey"`
	TopicID   string   `gorm:"index;type:varchar(32)"`
	CityID    string   `gorm:"index;type:varchar(32)"`
	Address   string   `gorm:"type:varchar(255)"`
	Images    []string `gorm:"type:json;serializer:json"`
	Latitude  float64
	Longitude float64

	MaxNum  int
	JoinNum int

	Title       string `gorm:"type:varchar(255)"`
	Desc        string `gorm:"type:text"`
	ChatGroupID string `gorm:"type:varchar(32)"`

	UserID string `gorm:"index;type:varchar(32)"`
	Active bool

	Deadline time.Time
	StartAt  time.Time

	CreatedAt time.Time
}

type UserJoinEventProposal struct {
	ID int `gorm:"primaryKey"`

	EventID string `gorm:"index;type:varchar(32)"`
	UserID  string `gorm:"index;type:varchar(32)"`

	// "joined" | "left" | "kickedOut"
	State string

	LeftAt    *time.Time
	CreatedAt time.Time
}

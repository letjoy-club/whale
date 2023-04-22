package models

import "time"

type Matching struct {
	ID string `gorm:"primaryKey"`

	TopicID string   `gorm:"index"`
	UserID  string   `gorm:"index;type:varchar(64)"`
	AreaIDs []string `gorm:"serializer:json;type:json"`
	CityID  string   `gorm:"index;type:varchar(64)"`

	Gender string `gorm:"type:varchar(4)"`

	RejectedUserIDs []string `gorm:"serializer:json;type:json"`
	InChatGroup     bool
	State           string `gorm:"type:varchar(64)"`
	ChatGroupState  string `gorm:"type:varchar(64)"`
	ResultID        int
	Remark          string `gorm:"type:varchar(64)"`

	Deadline  time.Time
	CreatedAt time.Time `gorm:"autoCreateTime;index"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (Matching) IsEntity() {}

type MatchingResult struct {
	ID          int      `gorm:"primaryKey"`
	MatchingIDs []string `gorm:"serializer:json;type:json"`
	TopicID     string   `gorm:"index"`
	UserIDs     []string `gorm:"serializer:json;type:json"`

	ConfirmStates  []string `gorm:"serializer:json;type:json"`
	ChatGroupState string   `gorm:"type:varchar(64)"`
	ChatGroupID    string   `gorm:"type:varchar(64)"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (MatchingResult) IsEntity() {}

type MatchingResultConfirmAction struct {
	ID               int    `gorm:"primaryKey"`
	MatchingResultID int    `gorm:"index"`
	UserID           string `gorm:"index;type:varchar(64)"`
	Confirmed        bool

	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type MatchingQuota struct {
	UserID string `gorm:"primaryKey"`

	Remain int
	Total  int

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (MatchingQuota) IsEntity() {}

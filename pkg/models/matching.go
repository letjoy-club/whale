package models

import (
	"time"

	"gorm.io/gorm"
)

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

type MatchingReview struct {
	ID int `gorm:"primaryKey"`

	MatchingResultID int
	MatchingID       string `gorm:"index:from_matching_user;type:varchar(64)"`
	UserID           string `gorm:"index:from_matching_user;type:varchar(64)"`

	ToMatchingID string `gorm:"index:to_matching_user;type:varchar(64)"`
	ToUserID     string `gorm:"index:to_matching_user;type:varchar(64)"`
	TopicID      string `gorm:"type:varchar(64)"`

	Score      int
	Comment    string    `gorm:"type:varchar(1024)"`
	CreateTime time.Time `gorm:"autoCreateTime"`
}

func (m *Matching) BeforeFind(db *gorm.DB) error {
	if m.RejectedUserIDs == nil {
		m.RejectedUserIDs = []string{}
	}
	return nil
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
	Closed         bool

	CreatedBy string `gorm:"type:varchar(64)"`

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

type MatchingPreview struct {
	UserID    string
	Remark    string
	CreatedAt time.Time
}

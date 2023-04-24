package models

import (
	"time"

	"gorm.io/gorm"
)

type MatchingInvitation struct {
	ID        string `gorm:"primaryKey"`
	UserID    string `gorm:"index;type:varchar(64)"`
	InviteeID string `gorm:"index;type:varchar(64)"`

	TopicID string `gorm:"type:varchar(64)"`

	Remark string `gorm:"type:varchar(128)"`

	MatchingResultId int

	MatchingIds []string `gorm:"serializer:json;type:json"`
	CityID      string   `gorm:"type:varchar(64)"`
	AreaIDs     []string `gorm:"serializer:json;type:json"`

	Closed       bool
	ConfirmState string `gorm:"index;type:varchar(64)"`
	ConfirmedAt  *time.Time
	CreatedAt    time.Time `gorm:"autoCreateTime;index"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

func (m *MatchingInvitation) BeforeFind(db *gorm.DB) error {
	if m.AreaIDs == nil {
		m.AreaIDs = []string{}
	}
	return nil
}

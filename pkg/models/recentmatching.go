package models

import "time"

type RecentMatching struct {
	ID      string `gorm:"primaryKey;type:varchar(64)"`
	CityID  string `gorm:"uniqueIndex:city_topic;type:varchar(32)"`
	TopicID string `gorm:"uniqueIndex:city_topic;type:varchar(32)"`

	MatchingIDs []string `gorm:"serializer:json;type:json"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

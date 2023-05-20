package models

import "time"

type CityTopics struct {
	CityID string `gorm:"primaryKey"`

	TopicIDs []string `gorm:"serializer:json;type:json"`

	UpdatedAt time.Time
}

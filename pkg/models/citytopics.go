package models

import "time"

type CityTopics struct {
	CityID string `gorm:"primaryKey;type:varchar(32)"`

	TopicIDs []string `gorm:"serializer:json;type:json"`

	UpdatedAt time.Time `gorm:"autoUpdatetime"`
}

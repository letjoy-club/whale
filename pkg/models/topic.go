package models

import "time"

type HotTopic struct {
	TopicID string `gorm:"uniqueIndex:topic_area;type:varchar(64)"`
	AreaID  string `gorm:"uniqueIndex:topic_area;type:varchar(64)"`

	InMatchingNum int
	MatchedNum    int
}

type TopicMetrics struct {
	ID       string `json:"id"`
	Matching int    `json:"matching"`
	Matched  int    `json:"matched"`
}

func (tm TopicMetrics) Total() int {
	return tm.Matching + tm.Matched
}

type HotTopicsInArea struct {
	CityID       string         `gorm:"primaryKey;type:varchar(64)"`
	TopicMetrics []TopicMetrics `gorm:"type:text;serializer:json"`

	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

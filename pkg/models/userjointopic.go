package models

type UserJoinTopic struct {
	ID               int    `gorm:"primaryKey"`
	TopicID          string `gorm:"uniqueIndex:topic_city_user;type:varchar(32)"`
	CityID           string `gorm:"uniqueIndex:topic_city_user;type:varchar(32)"`
	UserID           string `gorm:"uniqueIndex:topic_city_user;type:varchar(32)"`
	LatestMatchingID string `gorm:"type:varchar(32)"`
	Times            int
	UpdatedAt        int `gorm:"autoUpdateTime;index"`
	CreatedAt        int `gorm:"autoCreateTime"`
}

package models

import (
	"encoding/json"
	"gorm.io/plugin/soft_delete"
	"time"
)

type WhaleConfig struct {
	ID        int                   `json:"id" gorm:"primaryKey"`
	Name      string                `json:"name" gorm:"type:varchar(64)"`  // 配置名，唯一标识一种配置
	Desc      string                `json:"desc" gorm:"type:varchar(100)"` // 配置的中文说明
	Enable    bool                  `json:"enable"`                        // 配置是否启用
	StartAt   time.Time             `json:"startAt"`                       // 配置生效时间
	EndAt     *time.Time            `json:"endAt"`                         // 配置结束时间，为Null时表示永久有效
	Content   json.RawMessage       `json:"content" gorm:"type:text"`      // 配置内容，json字符串
	CreatedAt time.Time             `json:"createdAt" gorm:"index:idx_created_at;autoCreateTime"`
	UpdatedAt time.Time             `json:"updatedAt" gorm:"autoUpdateTime"`
	IsDeleted soft_delete.DeletedAt `json:"isDeleted" gorm:"softDelete:flag;default:0"`
}

var (
	ConfigLevelRights = "LevelRights" // 等级权益
)

// 等级权益

type UserLevelConfig struct {
	Rights []*LevelRights `json:"rights"`
}

type LevelRights struct {
	Level                      int `json:"level"`
	MotionQuota                int `json:"motionQuota"`
	OfferQuota                 int `json:"offerQuota"`
	MatchingQuota              int `json:"matchingQuota"`
	MatchingDurationConstraint int `json:"matchingDurationConstraint"`
}

// IsEntity GraphQL @key
func (LevelRights) IsEntity() {}

func (c *UserLevelConfig) Parse(config *WhaleConfig) error {
	return json.Unmarshal(config.Content, c)
}

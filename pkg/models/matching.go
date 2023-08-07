package models

import (
	"time"
	"whale/pkg/whalecode"

	"gorm.io/gorm"
)

type Matching struct {
	ID string `gorm:"primaryKey"`

	TopicID string   `gorm:"index"`
	UserID  string   `gorm:"index;type:varchar(64)"`
	AreaIDs []string `gorm:"serializer:json;type:json"`
	CityID  string   `gorm:"index;type:varchar(64)"`

	// 期望的性别
	Gender string `gorm:"type:varchar(4)"`

	// 创建者的性别
	MyGender string `gorm:"type:varchar(4)"`

	RejectedUserIDs []string `gorm:"serializer:json;type:json"`
	InChatGroup     bool
	State           string `gorm:"type:varchar(64)"`
	ChatGroupState  string `gorm:"type:varchar(64)"`
	ResultID        int
	Remark          string `gorm:"type:varchar(64)"`

	// 特定日期区间，格式 20060102
	DayRange []string `gorm:"serializer:json;type:json"`
	// 优先时间段
	PreferredPeriods []string           `gorm:"serializer:json;type:json"`
	Properties       []MatchingProperty `gorm:"serializer:json;type:json"`

	// 真正开始匹配的时间
	StartMatchingAt *time.Time

	Discoverable bool

	FinishedAt *time.Time
	MatchedAt  *time.Time

	Deadline  time.Time
	CreatedAt time.Time `gorm:"autoCreateTime;index"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (m *Matching) BeforeSave(db *gorm.DB) error {
	if m.DayRange != nil {
		if len(m.DayRange) == 0 {
			return nil
		}
		if len(m.DayRange) != 2 {
			return whalecode.ErrDayRangeNumInvalid
		}
		if m.DayRange[0] > m.DayRange[1] {
			return whalecode.ErrDayRangeNumInvalid
		}
		if _, err := time.Parse("20060102", m.DayRange[0]); err != nil {
			return whalecode.ErrDayRangeDateFormatInvalid
		}
		if _, err := time.Parse("20060102", m.DayRange[1]); err != nil {
			return whalecode.ErrDayRangeDateFormatInvalid
		}
	}
	return nil
}

func (m *Matching) HasSpecificDay(day string) bool {
	if len(m.DayRange) != 2 {
		return false
	}
	if day >= m.DayRange[0] && day <= m.DayRange[1] {
		return true
	}
	return false
}

func (m *Matching) AfterFind(tx *gorm.DB) error {
	if m.RejectedUserIDs == nil {
		m.RejectedUserIDs = []string{}
	}
	if m.Properties == nil {
		m.Properties = []MatchingProperty{}
	}
	if m.PreferredPeriods == nil {
		m.PreferredPeriods = []string{}
	}
	if m.DayRange == nil {
		m.DayRange = []string{}
	}
	return nil
}

func (m *Matching) GetProperty(id string) []string {
	for _, p := range m.Properties {
		if p.ID == id {
			return p.Values
		}
	}
	return []string{}
}

func (m *Matching) GetSingleValueProperty(id string) string {
	for _, p := range m.Properties {
		if p.ID == id && len(p.Values) > 0 {
			return p.Values[0]
		}
	}
	return ""
}

type MatchingProperty struct {
	ID     string   `json:"id"`
	Values []string `json:"values"`
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
	MotionIDs   []string `gorm:"serializer:json;type:json"`
	TopicID     string   `gorm:"index;type:varchar(32)"`
	UserIDs     []string `gorm:"serializer:json;type:json"`

	ConfirmStates  []string `gorm:"serializer:json;type:json"`
	ChatGroupState string   `gorm:"type:varchar(64)"`
	ChatGroupID    string   `gorm:"type:varchar(64)"`

	Closed bool

	MatchingScore int `gorm:"default:100"`

	CreatedBy string `gorm:"type:varchar(64)"`

	// 结束时间，由第一个用户的结束时间来决定
	FinishedAt         *time.Time `gorm:"index"`
	ChatGroupCreatedAt *time.Time `gorm:"index"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (m *MatchingResult) GetMatchingID(userID string) string {
	for i, id := range m.UserIDs {
		if id == userID {
			return m.MatchingIDs[i]
		}
	}
	return ""
}

func (m *MatchingResult) OtherUserIDs(userID string) []string {
	userIDs := []string{}
	for _, id := range m.UserIDs {
		if id != userID {
			userIDs = append(userIDs, id)
		}
	}
	return userIDs
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

	MatchingNum   int
	InvitationNum int

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (MatchingQuota) IsEntity() {}

type MatchingDurationConstraint struct {
	ID int `gorm:"primaryKey"`

	UserID string `gorm:"index;type:varchar(32)"`

	Total  int
	Remain int

	StartDate time.Time `gorm:"index"`
	StopDate  time.Time `gorm:"index"`

	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type MatchingPreview struct {
	UserID    string
	Remark    string
	CreatedAt time.Time
}

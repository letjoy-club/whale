// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package models

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

type Area struct {
	Code string `json:"code"`
}

func (Area) IsEntity() {}

type AvailableMotionOffer struct {
	// 可发起的意向
	Motion *Motion `json:"motion,omitempty"`
	// 下次获得一次配额的时间
	NextQuotaTime *time.Time `json:"nextQuotaTime,omitempty"`
}

type CalendarEvent struct {
	TopicID            string     `json:"topicId"`
	MatchedAt          time.Time  `json:"matchedAt"`
	FinishedAt         time.Time  `json:"finishedAt"`
	ChatGroupCreatedAt *time.Time `json:"chatGroupCreatedAt,omitempty"`
}

type ChatGroup struct {
	ID string `json:"id"`
}

func (ChatGroup) IsEntity() {}

type CitiesTopicsFilter struct {
	CityID *string `json:"cityId,omitempty"`
}

type CityToTopicMatching struct {
	CityID string             `json:"cityId"`
	Topics []*TopicToMatching `json:"topics"`
	City   *Area              `json:"city"`
}

type CreateCityTopicParam struct {
	TopicIds []string `json:"topicIds"`
	CityID   string   `json:"cityId"`
}

type CreateMatchingInvitationParam struct {
	InviteeID string   `json:"inviteeId"`
	Remark    string   `json:"remark"`
	TopicID   string   `json:"topicId"`
	CityID    string   `json:"cityId"`
	AreaIds   []string `json:"areaIds"`
}

type CreateMatchingParam struct {
	TopicID  string     `json:"topicId"`
	AreaIds  []string   `json:"areaIds"`
	CityID   string     `json:"cityId"`
	Gender   Gender     `json:"gender"`
	Remark   *string    `json:"remark,omitempty"`
	Deadline *time.Time `json:"deadline,omitempty"`
}

type CreateMatchingParamV2 struct {
	TopicID string   `json:"topicId"`
	AreaIds []string `json:"areaIds"`
	CityID  string   `json:"cityId"`
	Gender  Gender   `json:"gender"`
	// 特定日期区间，格式 yyyyMMdd，一定要包含两个字符串，字符串区间为闭区间
	DayRange []string `json:"dayRange"`
	// 特定时间区间，如果不限，则长度为0
	PreferredPeriods []DatePeriod `json:"preferredPeriods"`
	// 匹配属性
	Properties []*MatchingPropertyParam `json:"properties"`
	Remark     *string                  `json:"remark,omitempty"`
	Deadline   *time.Time               `json:"deadline,omitempty"`
}

type CreateMotionParam struct {
	TopicID string   `json:"topicId"`
	AreaIds []string `json:"areaIds"`
	CityID  string   `json:"cityId"`
	Gender  Gender   `json:"gender"`
	// 特定日期区间，格式 yyyyMMdd，一定要包含两个字符串，字符串区间为闭区间
	DayRange []string `json:"dayRange"`
	// 特定时间区间，如果不限，则长度为0
	PreferredPeriods []DatePeriod `json:"preferredPeriods"`
	// 匹配属性
	Properties []*MotionPropertyParam `json:"properties"`
	Remark     *string                `json:"remark,omitempty"`
	Deadline   *time.Time             `json:"deadline,omitempty"`
}

type CreateUserJoinTopicParam struct {
	MatchingID string `json:"matchingId"`
}

type DiscoverMotionResult struct {
	Motions   []*Motion `json:"motions"`
	NextToken string    `json:"nextToken"`
}

type DiscoverTopicCategoryMotionFilter struct {
	// 城市 ID，可以不填，不填则为全国
	CityID *string `json:"cityId,omitempty"`
	// 发起人性别，可以不填，不填则为不限
	Gender *Gender `json:"gender,omitempty"`
	// 话题 id，不填则为不限
	TopicIds []string `json:"topicIds,omitempty"`
}

type HotTopicsFilter struct {
	CityID  *string `json:"cityId,omitempty"`
	TopicID *string `json:"topicId,omitempty"`
}

type MatchingFilter struct {
	Before  *time.Time     `json:"before,omitempty"`
	After   *time.Time     `json:"after,omitempty"`
	TopicID *string        `json:"topicId,omitempty"`
	State   *MatchingState `json:"state,omitempty"`
	CityID  *string        `json:"cityId,omitempty"`
	UserID  *string        `json:"userId,omitempty"`
	// 通用关键词, u_ 开头搜用户, t_ 开头搜话题, m_ 开头搜匹配, 6 个数字搜地区
	Keyword *string `json:"keyword,omitempty"`
}

type MatchingInvitationFilter struct {
	UserID *string    `json:"userId,omitempty"`
	Before *time.Time `json:"before,omitempty"`
	After  *time.Time `json:"after,omitempty"`
}

type MatchingPropertyParam struct {
	ID     string   `json:"id"`
	Values []string `json:"values"`
}

type MatchingResultFilter struct {
	UserID *string    `json:"userId,omitempty"`
	Before *time.Time `json:"before,omitempty"`
	After  *time.Time `json:"after,omitempty"`
}

type MotionFilter struct {
	ID     *string `json:"id,omitempty"`
	UserID *string `json:"userId,omitempty"`
	CityID *string `json:"cityId,omitempty"`
	Gender *Gender `json:"gender,omitempty"`
}

type MotionPropertyParam struct {
	ID     string   `json:"id"`
	Values []string `json:"values"`
}

type RecentMatchingFilter struct {
	CityID *string `json:"cityId,omitempty"`
}

type ReviewMatchingParam struct {
	ToUserID string `json:"toUserId"`
	Score    int    `json:"score"`
	Comment  string `json:"comment"`
}

type ReviewMotionParam struct {
	Score   int    `json:"score"`
	Comment string `json:"comment"`
}

type SimpleAvatarUser struct {
	Avatar   string `json:"avatar"`
	Nickname string `json:"nickname"`
}

type Summary struct {
	Count int `json:"count"`
}

type Topic struct {
	ID string `json:"id"`
	// 匹配中用户
	RecentUsers []*SimpleAvatarUser `json:"recentUsers"`
	// 话题下的匹配数量
	MatchingNum int `json:"matchingNum"`
	// 话题下的大致匹配数量，范围是 [9, 999]，显示时建议给 + 表示有更多的数量。9-999 数值展示时，对原数据进行增量处理
	FuzzyMatchingNum int `json:"fuzzyMatchingNum"`
}

func (Topic) IsEntity() {}

type TopicOptionConfig struct {
	TopicID string `json:"topicId"`
}

func (TopicOptionConfig) IsEntity() {}

type TopicToMatching struct {
	TopicID     string   `json:"topicId"`
	MatchingIds []string `json:"matchingIds"`
	Topic       *Topic   `json:"topic"`
}

type UpdateCityTopicParam struct {
	TopicIds []string `json:"topicIds"`
}

type UpdateHotTopicMetricsParam struct {
	TopicID  string `json:"topicId"`
	Heat     int    `json:"heat"`
	Matched  int    `json:"matched"`
	Matching int    `json:"matching"`
}

type UpdateHotTopicParam struct {
	TopicMetrics []*UpdateHotTopicMetricsParam `json:"topicMetrics"`
}

type UpdateMatchingDurationConstraintParam struct {
	Total     *int       `json:"total,omitempty"`
	Remain    *int       `json:"remain,omitempty"`
	StartDate *time.Time `json:"startDate,omitempty"`
	StopDate  *time.Time `json:"stopDate,omitempty"`
}

type UpdateMatchingInvitationParam struct {
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	TopicID   *string    `json:"topicId,omitempty"`
	InviteeID *string    `json:"inviteeId,omitempty"`
	CityID    *string    `json:"cityId,omitempty"`
	Remark    *string    `json:"remark,omitempty"`
}

type UpdateMatchingParam struct {
	TopicID         *string    `json:"topicId,omitempty"`
	AreaIds         []string   `json:"areaIds,omitempty"`
	CityID          *string    `json:"cityId,omitempty"`
	Gender          *Gender    `json:"gender,omitempty"`
	Remark          *string    `json:"remark,omitempty"`
	StartMatchingAt *time.Time `json:"startMatchingAt,omitempty"`
	CreatedAt       *time.Time `json:"createdAt,omitempty"`
	Deadline        *time.Time `json:"deadline,omitempty"`
}

type UpdateMatchingQuotaParam struct {
	Total  *int `json:"total,omitempty"`
	Remain *int `json:"remain,omitempty"`
}

type UpdateMotionParam struct {
	AreaIds []string `json:"areaIds,omitempty"`
	CityID  *string  `json:"cityId,omitempty"`
	Gender  *Gender  `json:"gender,omitempty"`
	// 特定日期区间，格式 yyyyMMdd，一定要包含两个字符串，字符串区间为闭区间
	DayRange []string `json:"dayRange,omitempty"`
	// 特定时间区间，如果不限，则长度为0
	PreferredPeriods []DatePeriod `json:"preferredPeriods,omitempty"`
	// 匹配属性
	Properties []*MotionPropertyParam `json:"properties,omitempty"`
	Remark     *string                `json:"remark,omitempty"`
	Deadline   *time.Time             `json:"deadline,omitempty"`
}

type UpdateRecentMatchingParam struct {
	MatchingIds []string `json:"matchingIds"`
}

type UpdateUserJoinTopicParam struct {
	MatchingID string `json:"matchingId"`
}

type User struct {
	ID            string         `json:"id"`
	MatchingQuota *MatchingQuota `json:"matchingQuota"`
}

func (User) IsEntity() {}

type UserConfirmState struct {
	UserID string                     `json:"userId"`
	State  MatchingResultConfirmState `json:"state"`
}

type UserJoinTopicFilter struct {
	CityID  *string `json:"cityId,omitempty"`
	TopicID *string `json:"topicId,omitempty"`
	UserID  *string `json:"userId,omitempty"`
}

type UserMatchingCalenderParam struct {
	Before      time.Time `json:"before"`
	After       time.Time `json:"after"`
	OtherUserID *string   `json:"otherUserId,omitempty"`
}

type UserMatchingFilter struct {
	State  *MatchingState  `json:"state,omitempty"`
	States []MatchingState `json:"states,omitempty"`
}

type UserMatchingInTheDayParam struct {
	// 日期格式 20060102
	DayStr      string  `json:"dayStr"`
	OtherUserID *string `json:"otherUserId,omitempty"`
}

type UserUpdateMotionParam struct {
	AreaIds          []string               `json:"areaIds,omitempty"`
	Gender           *string                `json:"gender,omitempty"`
	DayRange         []string               `json:"dayRange,omitempty"`
	PreferredPeriods []DatePeriod           `json:"preferredPeriods,omitempty"`
	Properties       []*MotionPropertyParam `json:"properties,omitempty"`
	Remark           *string                `json:"remark,omitempty"`
}

type ChatGroupState string

const (
	// 未创建
	ChatGroupStateUncreated ChatGroupState = "Uncreated"
	// 等待创建
	ChatGroupStateWaitingCreated ChatGroupState = "WaitingCreated"
	// 创建成功
	ChatGroupStateCreated ChatGroupState = "Created"
	// 创建失败
	ChatGroupStateFailed ChatGroupState = "Failed"
	// 已关闭
	ChatGroupStateClosed ChatGroupState = "Closed"
	// 已退出
	ChatGroupStateQuited ChatGroupState = "Quited"
)

var AllChatGroupState = []ChatGroupState{
	ChatGroupStateUncreated,
	ChatGroupStateWaitingCreated,
	ChatGroupStateCreated,
	ChatGroupStateFailed,
	ChatGroupStateClosed,
	ChatGroupStateQuited,
}

func (e ChatGroupState) IsValid() bool {
	switch e {
	case ChatGroupStateUncreated, ChatGroupStateWaitingCreated, ChatGroupStateCreated, ChatGroupStateFailed, ChatGroupStateClosed, ChatGroupStateQuited:
		return true
	}
	return false
}

func (e ChatGroupState) String() string {
	return string(e)
}

func (e *ChatGroupState) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ChatGroupState(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ChatGroupState", str)
	}
	return nil
}

func (e ChatGroupState) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type DatePeriod string

const (
	// 不限
	DatePeriodUnlimited DatePeriod = "Unlimited"
	// 周末
	DatePeriodWeekend DatePeriod = "Weekend"
	// 周末下午
	DatePeriodWeekendAfternoon DatePeriod = "WeekendAfternoon"
	// 周末晚上
	DatePeriodWeekendNight DatePeriod = "WeekendNight"
	// 工作日
	DatePeriodWorkday DatePeriod = "Workday"
	// 工作日下午
	DatePeriodWorkdayAfternoon DatePeriod = "WorkdayAfternoon"
	// 工作日晚上
	DatePeriodWorkdayNight DatePeriod = "WorkdayNight"
)

var AllDatePeriod = []DatePeriod{
	DatePeriodUnlimited,
	DatePeriodWeekend,
	DatePeriodWeekendAfternoon,
	DatePeriodWeekendNight,
	DatePeriodWorkday,
	DatePeriodWorkdayAfternoon,
	DatePeriodWorkdayNight,
}

func (e DatePeriod) IsValid() bool {
	switch e {
	case DatePeriodUnlimited, DatePeriodWeekend, DatePeriodWeekendAfternoon, DatePeriodWeekendNight, DatePeriodWorkday, DatePeriodWorkdayAfternoon, DatePeriodWorkdayNight:
		return true
	}
	return false
}

func (e DatePeriod) String() string {
	return string(e)
}

func (e *DatePeriod) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = DatePeriod(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid DatePeriod", str)
	}
	return nil
}

func (e DatePeriod) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Gender string

const (
	// 女
	GenderF Gender = "F"
	// 男
	GenderM Gender = "M"
	// 不限
	GenderN Gender = "N"
)

var AllGender = []Gender{
	GenderF,
	GenderM,
	GenderN,
}

func (e Gender) IsValid() bool {
	switch e {
	case GenderF, GenderM, GenderN:
		return true
	}
	return false
}

func (e Gender) String() string {
	return string(e)
}

func (e *Gender) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Gender(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Gender", str)
	}
	return nil
}

func (e Gender) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type InvitationConfirmState string

const (
	InvitationConfirmStateConfirmed   InvitationConfirmState = "Confirmed"
	InvitationConfirmStateRejected    InvitationConfirmState = "Rejected"
	InvitationConfirmStateUnconfirmed InvitationConfirmState = "Unconfirmed"
)

var AllInvitationConfirmState = []InvitationConfirmState{
	InvitationConfirmStateConfirmed,
	InvitationConfirmStateRejected,
	InvitationConfirmStateUnconfirmed,
}

func (e InvitationConfirmState) IsValid() bool {
	switch e {
	case InvitationConfirmStateConfirmed, InvitationConfirmStateRejected, InvitationConfirmStateUnconfirmed:
		return true
	}
	return false
}

func (e InvitationConfirmState) String() string {
	return string(e)
}

func (e *InvitationConfirmState) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = InvitationConfirmState(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid InvitationConfirmState", str)
	}
	return nil
}

func (e InvitationConfirmState) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type MatchingResultConfirmState string

const (
	// 未确认
	MatchingResultConfirmStateUnconfirmed MatchingResultConfirmState = "Unconfirmed"
	// 已确认
	MatchingResultConfirmStateConfirmed MatchingResultConfirmState = "Confirmed"
	// 已拒绝
	MatchingResultConfirmStateRejected MatchingResultConfirmState = "Rejected"
)

var AllMatchingResultConfirmState = []MatchingResultConfirmState{
	MatchingResultConfirmStateUnconfirmed,
	MatchingResultConfirmStateConfirmed,
	MatchingResultConfirmStateRejected,
}

func (e MatchingResultConfirmState) IsValid() bool {
	switch e {
	case MatchingResultConfirmStateUnconfirmed, MatchingResultConfirmStateConfirmed, MatchingResultConfirmStateRejected:
		return true
	}
	return false
}

func (e MatchingResultConfirmState) String() string {
	return string(e)
}

func (e *MatchingResultConfirmState) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = MatchingResultConfirmState(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid MatchingResultConfirmState", str)
	}
	return nil
}

func (e MatchingResultConfirmState) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type MatchingState string

const (
	// 匹配中
	MatchingStateMatching MatchingState = "Matching"
	// 匹配成功
	MatchingStateMatched MatchingState = "Matched"
	// 匹配失败
	MatchingStateFailed MatchingState = "Failed"
	// 匹配取消
	MatchingStateCanceled MatchingState = "Canceled"
	// 匹配超时
	MatchingStateTimeout MatchingState = "Timeout"
	// 匹配关闭
	MatchingStateFinished MatchingState = "Finished"
)

var AllMatchingState = []MatchingState{
	MatchingStateMatching,
	MatchingStateMatched,
	MatchingStateFailed,
	MatchingStateCanceled,
	MatchingStateTimeout,
	MatchingStateFinished,
}

func (e MatchingState) IsValid() bool {
	switch e {
	case MatchingStateMatching, MatchingStateMatched, MatchingStateFailed, MatchingStateCanceled, MatchingStateTimeout, MatchingStateFinished:
		return true
	}
	return false
}

func (e MatchingState) String() string {
	return string(e)
}

func (e *MatchingState) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = MatchingState(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid MatchingState", str)
	}
	return nil
}

func (e MatchingState) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type MotionOfferState string

const (
	// 未处理
	MotionOfferStatePending MotionOfferState = "Pending"
	// 被接受
	MotionOfferStateAccepted MotionOfferState = "Accepted"
	// 被拒绝
	MotionOfferStateRejected MotionOfferState = "Rejected"
	// 意向已取消
	MotionOfferStateCanceled MotionOfferState = "Canceled"
	// 处理超时
	MotionOfferStateTimeout MotionOfferState = "Timeout"
	// 结束
	MotionOfferStateFinished MotionOfferState = "Finished"
)

var AllMotionOfferState = []MotionOfferState{
	MotionOfferStatePending,
	MotionOfferStateAccepted,
	MotionOfferStateRejected,
	MotionOfferStateCanceled,
	MotionOfferStateTimeout,
	MotionOfferStateFinished,
}

func (e MotionOfferState) IsValid() bool {
	switch e {
	case MotionOfferStatePending, MotionOfferStateAccepted, MotionOfferStateRejected, MotionOfferStateCanceled, MotionOfferStateTimeout, MotionOfferStateFinished:
		return true
	}
	return false
}

func (e MotionOfferState) String() string {
	return string(e)
}

func (e *MotionOfferState) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = MotionOfferState(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid MotionOfferState", str)
	}
	return nil
}

func (e MotionOfferState) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ResultCreatedBy string

const (
	// 由匹配系统创建的结果
	ResultCreatedByMatching ResultCreatedBy = "Matching"
	// 由邀请创建的结果
	ResultCreatedByInvitation ResultCreatedBy = "Invitation"
	// 由匹配邀约创建的结果
	ResultCreatedByOffer ResultCreatedBy = "Offer"
)

var AllResultCreatedBy = []ResultCreatedBy{
	ResultCreatedByMatching,
	ResultCreatedByInvitation,
	ResultCreatedByOffer,
}

func (e ResultCreatedBy) IsValid() bool {
	switch e {
	case ResultCreatedByMatching, ResultCreatedByInvitation, ResultCreatedByOffer:
		return true
	}
	return false
}

func (e ResultCreatedBy) String() string {
	return string(e)
}

func (e *ResultCreatedBy) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ResultCreatedBy(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ResultCreatedBy", str)
	}
	return nil
}

func (e ResultCreatedBy) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

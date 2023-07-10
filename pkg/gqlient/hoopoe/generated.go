// Code generated by github.com/Khan/genqlient, DO NOT EDIT.

package hoopoe

import (
	"context"
	"encoding/json"

	"github.com/Khan/genqlient/graphql"
)

// CreateMatchingValidResponse is returned by CreateMatchingValid on success.
type CreateMatchingValidResponse struct {
	// 【用户】获取用户信息
	User *CreateMatchingValidUser `json:"user"`
	// 【用户】 用户信息完整性检查
	UserInfoCompletenessCheck *CreateMatchingValidUserInfoCompletenessCheckUserInfoCompleteness `json:"userInfoCompletenessCheck"`
}

// GetUser returns CreateMatchingValidResponse.User, and is useful for accessing the field via an interface.
func (v *CreateMatchingValidResponse) GetUser() *CreateMatchingValidUser { return v.User }

// GetUserInfoCompletenessCheck returns CreateMatchingValidResponse.UserInfoCompletenessCheck, and is useful for accessing the field via an interface.
func (v *CreateMatchingValidResponse) GetUserInfoCompletenessCheck() *CreateMatchingValidUserInfoCompletenessCheckUserInfoCompleteness {
	return v.UserInfoCompletenessCheck
}

// CreateMatchingValidUser includes the requested fields of the GraphQL type User.
type CreateMatchingValidUser struct {
	// 用户封禁信息
	BlockInfo *CreateMatchingValidUserBlockInfo `json:"blockInfo"`
}

// GetBlockInfo returns CreateMatchingValidUser.BlockInfo, and is useful for accessing the field via an interface.
func (v *CreateMatchingValidUser) GetBlockInfo() *CreateMatchingValidUserBlockInfo {
	return v.BlockInfo
}

// CreateMatchingValidUserBlockInfo includes the requested fields of the GraphQL type UserBlockInfo.
type CreateMatchingValidUserBlockInfo struct {
	// 是否用户封禁，封禁则整个APP无法使用
	UserBlocked bool `json:"userBlocked"`
	// 是否匹配封禁，封禁则无法加入新的匹配
	MatchingBlocked bool `json:"matchingBlocked"`
}

// GetUserBlocked returns CreateMatchingValidUserBlockInfo.UserBlocked, and is useful for accessing the field via an interface.
func (v *CreateMatchingValidUserBlockInfo) GetUserBlocked() bool { return v.UserBlocked }

// GetMatchingBlocked returns CreateMatchingValidUserBlockInfo.MatchingBlocked, and is useful for accessing the field via an interface.
func (v *CreateMatchingValidUserBlockInfo) GetMatchingBlocked() bool { return v.MatchingBlocked }

// CreateMatchingValidUserInfoCompletenessCheckUserInfoCompleteness includes the requested fields of the GraphQL type UserInfoCompleteness.
type CreateMatchingValidUserInfoCompletenessCheckUserInfoCompleteness struct {
	// 信息是否全部填写
	Filled bool `json:"filled"`
}

// GetFilled returns CreateMatchingValidUserInfoCompletenessCheckUserInfoCompleteness.Filled, and is useful for accessing the field via an interface.
func (v *CreateMatchingValidUserInfoCompletenessCheckUserInfoCompleteness) GetFilled() bool {
	return v.Filled
}

type ExtraOptionKey string

const (
	ExtraOptionKeyNone    ExtraOptionKey = "None"
	ExtraOptionKeySinger  ExtraOptionKey = "Singer"
	ExtraOptionKeyConcert ExtraOptionKey = "Concert"
	ExtraOptionKeyPoi     ExtraOptionKey = "POI"
	ExtraOptionKeyCity    ExtraOptionKey = "City"
)

// GetAreaArea includes the requested fields of the GraphQL type Area.
type GetAreaArea struct {
	// 地区代码
	Code string `json:"code"`
	// 名称
	Name string `json:"name"`
	// 层级
	Depth int `json:"depth"`
}

// GetCode returns GetAreaArea.Code, and is useful for accessing the field via an interface.
func (v *GetAreaArea) GetCode() string { return v.Code }

// GetName returns GetAreaArea.Name, and is useful for accessing the field via an interface.
func (v *GetAreaArea) GetName() string { return v.Name }

// GetDepth returns GetAreaArea.Depth, and is useful for accessing the field via an interface.
func (v *GetAreaArea) GetDepth() int { return v.Depth }

// GetAreaResponse is returned by GetArea on success.
type GetAreaResponse struct {
	// 【地区】根据地区代码查询地区信息
	Area *GetAreaArea `json:"area"`
}

// GetArea returns GetAreaResponse.Area, and is useful for accessing the field via an interface.
func (v *GetAreaResponse) GetArea() *GetAreaArea { return v.Area }

// GetAvatarByIDsGetUserByIdsUser includes the requested fields of the GraphQL type User.
type GetAvatarByIDsGetUserByIdsUser struct {
	// 用户ID
	Id string `json:"id"`
	// 昵称
	Nickname string `json:"nickname"`
	// 头像URL
	Avatar string `json:"avatar"`
}

// GetId returns GetAvatarByIDsGetUserByIdsUser.Id, and is useful for accessing the field via an interface.
func (v *GetAvatarByIDsGetUserByIdsUser) GetId() string { return v.Id }

// GetNickname returns GetAvatarByIDsGetUserByIdsUser.Nickname, and is useful for accessing the field via an interface.
func (v *GetAvatarByIDsGetUserByIdsUser) GetNickname() string { return v.Nickname }

// GetAvatar returns GetAvatarByIDsGetUserByIdsUser.Avatar, and is useful for accessing the field via an interface.
func (v *GetAvatarByIDsGetUserByIdsUser) GetAvatar() string { return v.Avatar }

// GetAvatarByIDsResponse is returned by GetAvatarByIDs on success.
type GetAvatarByIDsResponse struct {
	// 【用户】根据ID批量获取
	GetUserByIds []*GetAvatarByIDsGetUserByIdsUser `json:"getUserByIds"`
}

// GetGetUserByIds returns GetAvatarByIDsResponse.GetUserByIds, and is useful for accessing the field via an interface.
func (v *GetAvatarByIDsResponse) GetGetUserByIds() []*GetAvatarByIDsGetUserByIdsUser {
	return v.GetUserByIds
}

// GetBlacklistRelationshipBlacklistRelationshipUserPair includes the requested fields of the GraphQL type UserPair.
type GetBlacklistRelationshipBlacklistRelationshipUserPair struct {
	A string `json:"a"`
	B string `json:"b"`
}

// GetA returns GetBlacklistRelationshipBlacklistRelationshipUserPair.A, and is useful for accessing the field via an interface.
func (v *GetBlacklistRelationshipBlacklistRelationshipUserPair) GetA() string { return v.A }

// GetB returns GetBlacklistRelationshipBlacklistRelationshipUserPair.B, and is useful for accessing the field via an interface.
func (v *GetBlacklistRelationshipBlacklistRelationshipUserPair) GetB() string { return v.B }

// GetBlacklistRelationshipResponse is returned by GetBlacklistRelationship on success.
type GetBlacklistRelationshipResponse struct {
	// 【黑名单】获取指定 userId 之间的黑名单关系
	BlacklistRelationship []*GetBlacklistRelationshipBlacklistRelationshipUserPair `json:"blacklistRelationship"`
}

// GetBlacklistRelationship returns GetBlacklistRelationshipResponse.BlacklistRelationship, and is useful for accessing the field via an interface.
func (v *GetBlacklistRelationshipResponse) GetBlacklistRelationship() []*GetBlacklistRelationshipBlacklistRelationshipUserPair {
	return v.BlacklistRelationship
}

// GetTopicConfigOptionResponse is returned by GetTopicConfigOption on success.
type GetTopicConfigOptionResponse struct {
	TopicOptionConfig *GetTopicConfigOptionTopicOptionConfig `json:"topicOptionConfig"`
}

// GetTopicOptionConfig returns GetTopicConfigOptionResponse.TopicOptionConfig, and is useful for accessing the field via an interface.
func (v *GetTopicConfigOptionResponse) GetTopicOptionConfig() *GetTopicConfigOptionTopicOptionConfig {
	return v.TopicOptionConfig
}

// GetTopicConfigOptionTopicOptionConfig includes the requested fields of the GraphQL type TopicOptionConfig.
type GetTopicConfigOptionTopicOptionConfig struct {
	TopicOptionConfigFields `json:"-"`
}

// GetTopicId returns GetTopicConfigOptionTopicOptionConfig.TopicId, and is useful for accessing the field via an interface.
func (v *GetTopicConfigOptionTopicOptionConfig) GetTopicId() string {
	return v.TopicOptionConfigFields.TopicId
}

// GetTimeWeight returns GetTopicConfigOptionTopicOptionConfig.TimeWeight, and is useful for accessing the field via an interface.
func (v *GetTopicConfigOptionTopicOptionConfig) GetTimeWeight() int {
	return v.TopicOptionConfigFields.TimeWeight
}

// GetThreshold returns GetTopicConfigOptionTopicOptionConfig.Threshold, and is useful for accessing the field via an interface.
func (v *GetTopicConfigOptionTopicOptionConfig) GetThreshold() int {
	return v.TopicOptionConfigFields.Threshold
}

// GetFuzzyMatchingTopic returns GetTopicConfigOptionTopicOptionConfig.FuzzyMatchingTopic, and is useful for accessing the field via an interface.
func (v *GetTopicConfigOptionTopicOptionConfig) GetFuzzyMatchingTopic() bool {
	return v.TopicOptionConfigFields.FuzzyMatchingTopic
}

// GetAllowFuzzyMatching returns GetTopicConfigOptionTopicOptionConfig.AllowFuzzyMatching, and is useful for accessing the field via an interface.
func (v *GetTopicConfigOptionTopicOptionConfig) GetAllowFuzzyMatching() bool {
	return v.TopicOptionConfigFields.AllowFuzzyMatching
}

// GetDelayMinuteToPairWithFuzzyTopic returns GetTopicConfigOptionTopicOptionConfig.DelayMinuteToPairWithFuzzyTopic, and is useful for accessing the field via an interface.
func (v *GetTopicConfigOptionTopicOptionConfig) GetDelayMinuteToPairWithFuzzyTopic() int {
	return v.TopicOptionConfigFields.DelayMinuteToPairWithFuzzyTopic
}

// GetProperties returns GetTopicConfigOptionTopicOptionConfig.Properties, and is useful for accessing the field via an interface.
func (v *GetTopicConfigOptionTopicOptionConfig) GetProperties() []*TopicOptionConfigFieldsPropertiesTopicOptionProperty {
	return v.TopicOptionConfigFields.Properties
}

func (v *GetTopicConfigOptionTopicOptionConfig) UnmarshalJSON(b []byte) error {

	if string(b) == "null" {
		return nil
	}

	var firstPass struct {
		*GetTopicConfigOptionTopicOptionConfig
		graphql.NoUnmarshalJSON
	}
	firstPass.GetTopicConfigOptionTopicOptionConfig = v

	err := json.Unmarshal(b, &firstPass)
	if err != nil {
		return err
	}

	err = json.Unmarshal(
		b, &v.TopicOptionConfigFields)
	if err != nil {
		return err
	}
	return nil
}

type __premarshalGetTopicConfigOptionTopicOptionConfig struct {
	TopicId string `json:"topicId"`

	TimeWeight int `json:"timeWeight"`

	Threshold int `json:"threshold"`

	FuzzyMatchingTopic bool `json:"fuzzyMatchingTopic"`

	AllowFuzzyMatching bool `json:"allowFuzzyMatching"`

	DelayMinuteToPairWithFuzzyTopic int `json:"delayMinuteToPairWithFuzzyTopic"`

	Properties []*TopicOptionConfigFieldsPropertiesTopicOptionProperty `json:"properties"`
}

func (v *GetTopicConfigOptionTopicOptionConfig) MarshalJSON() ([]byte, error) {
	premarshaled, err := v.__premarshalJSON()
	if err != nil {
		return nil, err
	}
	return json.Marshal(premarshaled)
}

func (v *GetTopicConfigOptionTopicOptionConfig) __premarshalJSON() (*__premarshalGetTopicConfigOptionTopicOptionConfig, error) {
	var retval __premarshalGetTopicConfigOptionTopicOptionConfig

	retval.TopicId = v.TopicOptionConfigFields.TopicId
	retval.TimeWeight = v.TopicOptionConfigFields.TimeWeight
	retval.Threshold = v.TopicOptionConfigFields.Threshold
	retval.FuzzyMatchingTopic = v.TopicOptionConfigFields.FuzzyMatchingTopic
	retval.AllowFuzzyMatching = v.TopicOptionConfigFields.AllowFuzzyMatching
	retval.DelayMinuteToPairWithFuzzyTopic = v.TopicOptionConfigFields.DelayMinuteToPairWithFuzzyTopic
	retval.Properties = v.TopicOptionConfigFields.Properties
	return &retval, nil
}

// GetTopicConfigOptionsResponse is returned by GetTopicConfigOptions on success.
type GetTopicConfigOptionsResponse struct {
	TopicOptionConfigs []*GetTopicConfigOptionsTopicOptionConfigsTopicOptionConfig `json:"topicOptionConfigs"`
}

// GetTopicOptionConfigs returns GetTopicConfigOptionsResponse.TopicOptionConfigs, and is useful for accessing the field via an interface.
func (v *GetTopicConfigOptionsResponse) GetTopicOptionConfigs() []*GetTopicConfigOptionsTopicOptionConfigsTopicOptionConfig {
	return v.TopicOptionConfigs
}

// GetTopicConfigOptionsTopicOptionConfigsTopicOptionConfig includes the requested fields of the GraphQL type TopicOptionConfig.
type GetTopicConfigOptionsTopicOptionConfigsTopicOptionConfig struct {
	TopicOptionConfigFields `json:"-"`
}

// GetTopicId returns GetTopicConfigOptionsTopicOptionConfigsTopicOptionConfig.TopicId, and is useful for accessing the field via an interface.
func (v *GetTopicConfigOptionsTopicOptionConfigsTopicOptionConfig) GetTopicId() string {
	return v.TopicOptionConfigFields.TopicId
}

// GetTimeWeight returns GetTopicConfigOptionsTopicOptionConfigsTopicOptionConfig.TimeWeight, and is useful for accessing the field via an interface.
func (v *GetTopicConfigOptionsTopicOptionConfigsTopicOptionConfig) GetTimeWeight() int {
	return v.TopicOptionConfigFields.TimeWeight
}

// GetThreshold returns GetTopicConfigOptionsTopicOptionConfigsTopicOptionConfig.Threshold, and is useful for accessing the field via an interface.
func (v *GetTopicConfigOptionsTopicOptionConfigsTopicOptionConfig) GetThreshold() int {
	return v.TopicOptionConfigFields.Threshold
}

// GetFuzzyMatchingTopic returns GetTopicConfigOptionsTopicOptionConfigsTopicOptionConfig.FuzzyMatchingTopic, and is useful for accessing the field via an interface.
func (v *GetTopicConfigOptionsTopicOptionConfigsTopicOptionConfig) GetFuzzyMatchingTopic() bool {
	return v.TopicOptionConfigFields.FuzzyMatchingTopic
}

// GetAllowFuzzyMatching returns GetTopicConfigOptionsTopicOptionConfigsTopicOptionConfig.AllowFuzzyMatching, and is useful for accessing the field via an interface.
func (v *GetTopicConfigOptionsTopicOptionConfigsTopicOptionConfig) GetAllowFuzzyMatching() bool {
	return v.TopicOptionConfigFields.AllowFuzzyMatching
}

// GetDelayMinuteToPairWithFuzzyTopic returns GetTopicConfigOptionsTopicOptionConfigsTopicOptionConfig.DelayMinuteToPairWithFuzzyTopic, and is useful for accessing the field via an interface.
func (v *GetTopicConfigOptionsTopicOptionConfigsTopicOptionConfig) GetDelayMinuteToPairWithFuzzyTopic() int {
	return v.TopicOptionConfigFields.DelayMinuteToPairWithFuzzyTopic
}

// GetProperties returns GetTopicConfigOptionsTopicOptionConfigsTopicOptionConfig.Properties, and is useful for accessing the field via an interface.
func (v *GetTopicConfigOptionsTopicOptionConfigsTopicOptionConfig) GetProperties() []*TopicOptionConfigFieldsPropertiesTopicOptionProperty {
	return v.TopicOptionConfigFields.Properties
}

func (v *GetTopicConfigOptionsTopicOptionConfigsTopicOptionConfig) UnmarshalJSON(b []byte) error {

	if string(b) == "null" {
		return nil
	}

	var firstPass struct {
		*GetTopicConfigOptionsTopicOptionConfigsTopicOptionConfig
		graphql.NoUnmarshalJSON
	}
	firstPass.GetTopicConfigOptionsTopicOptionConfigsTopicOptionConfig = v

	err := json.Unmarshal(b, &firstPass)
	if err != nil {
		return err
	}

	err = json.Unmarshal(
		b, &v.TopicOptionConfigFields)
	if err != nil {
		return err
	}
	return nil
}

type __premarshalGetTopicConfigOptionsTopicOptionConfigsTopicOptionConfig struct {
	TopicId string `json:"topicId"`

	TimeWeight int `json:"timeWeight"`

	Threshold int `json:"threshold"`

	FuzzyMatchingTopic bool `json:"fuzzyMatchingTopic"`

	AllowFuzzyMatching bool `json:"allowFuzzyMatching"`

	DelayMinuteToPairWithFuzzyTopic int `json:"delayMinuteToPairWithFuzzyTopic"`

	Properties []*TopicOptionConfigFieldsPropertiesTopicOptionProperty `json:"properties"`
}

func (v *GetTopicConfigOptionsTopicOptionConfigsTopicOptionConfig) MarshalJSON() ([]byte, error) {
	premarshaled, err := v.__premarshalJSON()
	if err != nil {
		return nil, err
	}
	return json.Marshal(premarshaled)
}

func (v *GetTopicConfigOptionsTopicOptionConfigsTopicOptionConfig) __premarshalJSON() (*__premarshalGetTopicConfigOptionsTopicOptionConfigsTopicOptionConfig, error) {
	var retval __premarshalGetTopicConfigOptionsTopicOptionConfigsTopicOptionConfig

	retval.TopicId = v.TopicOptionConfigFields.TopicId
	retval.TimeWeight = v.TopicOptionConfigFields.TimeWeight
	retval.Threshold = v.TopicOptionConfigFields.Threshold
	retval.FuzzyMatchingTopic = v.TopicOptionConfigFields.FuzzyMatchingTopic
	retval.AllowFuzzyMatching = v.TopicOptionConfigFields.AllowFuzzyMatching
	retval.DelayMinuteToPairWithFuzzyTopic = v.TopicOptionConfigFields.DelayMinuteToPairWithFuzzyTopic
	retval.Properties = v.TopicOptionConfigFields.Properties
	return &retval, nil
}

// GetTopicNameResponse is returned by GetTopicName on success.
type GetTopicNameResponse struct {
	// 【话题】话题查询
	Topic *GetTopicNameTopic `json:"topic"`
}

// GetTopic returns GetTopicNameResponse.Topic, and is useful for accessing the field via an interface.
func (v *GetTopicNameResponse) GetTopic() *GetTopicNameTopic { return v.Topic }

// GetTopicNameTopic includes the requested fields of the GraphQL type Topic.
type GetTopicNameTopic struct {
	// 话题ID
	Id string `json:"id"`
	// 名称
	Name string `json:"name"`
}

// GetId returns GetTopicNameTopic.Id, and is useful for accessing the field via an interface.
func (v *GetTopicNameTopic) GetId() string { return v.Id }

// GetName returns GetTopicNameTopic.Name, and is useful for accessing the field via an interface.
func (v *GetTopicNameTopic) GetName() string { return v.Name }

// GetTopicResponse is returned by GetTopic on success.
type GetTopicResponse struct {
	// 【话题】话题查询
	Topic *GetTopicTopic `json:"topic"`
}

// GetTopic returns GetTopicResponse.Topic, and is useful for accessing the field via an interface.
func (v *GetTopicResponse) GetTopic() *GetTopicTopic { return v.Topic }

// GetTopicTopic includes the requested fields of the GraphQL type Topic.
type GetTopicTopic struct {
	// 话题ID
	Id string `json:"id"`
}

// GetId returns GetTopicTopic.Id, and is useful for accessing the field via an interface.
func (v *GetTopicTopic) GetId() string { return v.Id }

// GetTopicsResponse is returned by GetTopics on success.
type GetTopicsResponse struct {
	// 【话题】列表查询
	Topics []*GetTopicsTopicsTopic `json:"topics"`
}

// GetTopics returns GetTopicsResponse.Topics, and is useful for accessing the field via an interface.
func (v *GetTopicsResponse) GetTopics() []*GetTopicsTopicsTopic { return v.Topics }

// GetTopicsTopicsTopic includes the requested fields of the GraphQL type Topic.
type GetTopicsTopicsTopic struct {
	// 话题ID
	Id string `json:"id"`
	// 名称
	Name string `json:"name"`
}

// GetId returns GetTopicsTopicsTopic.Id, and is useful for accessing the field via an interface.
func (v *GetTopicsTopicsTopic) GetId() string { return v.Id }

// GetName returns GetTopicsTopicsTopic.Name, and is useful for accessing the field via an interface.
func (v *GetTopicsTopicsTopic) GetName() string { return v.Name }

// GetUserByIDsGetUserByIdsUser includes the requested fields of the GraphQL type User.
type GetUserByIDsGetUserByIdsUser struct {
	// 用户ID
	Id string `json:"id"`
	// 性别：男——M；女——F
	Gender string `json:"gender"`
	// 用户等级
	Level int `json:"level"`
}

// GetId returns GetUserByIDsGetUserByIdsUser.Id, and is useful for accessing the field via an interface.
func (v *GetUserByIDsGetUserByIdsUser) GetId() string { return v.Id }

// GetGender returns GetUserByIDsGetUserByIdsUser.Gender, and is useful for accessing the field via an interface.
func (v *GetUserByIDsGetUserByIdsUser) GetGender() string { return v.Gender }

// GetLevel returns GetUserByIDsGetUserByIdsUser.Level, and is useful for accessing the field via an interface.
func (v *GetUserByIDsGetUserByIdsUser) GetLevel() int { return v.Level }

// GetUserByIDsResponse is returned by GetUserByIDs on success.
type GetUserByIDsResponse struct {
	// 【用户】根据ID批量获取
	GetUserByIds []*GetUserByIDsGetUserByIdsUser `json:"getUserByIds"`
}

// GetGetUserByIds returns GetUserByIDsResponse.GetUserByIds, and is useful for accessing the field via an interface.
func (v *GetUserByIDsResponse) GetGetUserByIds() []*GetUserByIDsGetUserByIdsUser {
	return v.GetUserByIds
}

type GraphQLPaginator struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

// GetPage returns GraphQLPaginator.Page, and is useful for accessing the field via an interface.
func (v *GraphQLPaginator) GetPage() int { return v.Page }

// GetSize returns GraphQLPaginator.Size, and is useful for accessing the field via an interface.
func (v *GraphQLPaginator) GetSize() int { return v.Size }

// TopicOptionConfigFields includes the GraphQL fields of TopicOptionConfig requested by the fragment TopicOptionConfigFields.
type TopicOptionConfigFields struct {
	TopicId    string `json:"topicId"`
	TimeWeight int    `json:"timeWeight"`
	Threshold  int    `json:"threshold"`
	// 是否是模糊匹配话题
	FuzzyMatchingTopic bool `json:"fuzzyMatchingTopic"`
	// 是否允许被模糊匹配
	AllowFuzzyMatching bool `json:"allowFuzzyMatching"`
	// 等待多少分钟后可进行模糊匹配
	DelayMinuteToPairWithFuzzyTopic int                                                     `json:"delayMinuteToPairWithFuzzyTopic"`
	Properties                      []*TopicOptionConfigFieldsPropertiesTopicOptionProperty `json:"properties"`
}

// GetTopicId returns TopicOptionConfigFields.TopicId, and is useful for accessing the field via an interface.
func (v *TopicOptionConfigFields) GetTopicId() string { return v.TopicId }

// GetTimeWeight returns TopicOptionConfigFields.TimeWeight, and is useful for accessing the field via an interface.
func (v *TopicOptionConfigFields) GetTimeWeight() int { return v.TimeWeight }

// GetThreshold returns TopicOptionConfigFields.Threshold, and is useful for accessing the field via an interface.
func (v *TopicOptionConfigFields) GetThreshold() int { return v.Threshold }

// GetFuzzyMatchingTopic returns TopicOptionConfigFields.FuzzyMatchingTopic, and is useful for accessing the field via an interface.
func (v *TopicOptionConfigFields) GetFuzzyMatchingTopic() bool { return v.FuzzyMatchingTopic }

// GetAllowFuzzyMatching returns TopicOptionConfigFields.AllowFuzzyMatching, and is useful for accessing the field via an interface.
func (v *TopicOptionConfigFields) GetAllowFuzzyMatching() bool { return v.AllowFuzzyMatching }

// GetDelayMinuteToPairWithFuzzyTopic returns TopicOptionConfigFields.DelayMinuteToPairWithFuzzyTopic, and is useful for accessing the field via an interface.
func (v *TopicOptionConfigFields) GetDelayMinuteToPairWithFuzzyTopic() int {
	return v.DelayMinuteToPairWithFuzzyTopic
}

// GetProperties returns TopicOptionConfigFields.Properties, and is useful for accessing the field via an interface.
func (v *TopicOptionConfigFields) GetProperties() []*TopicOptionConfigFieldsPropertiesTopicOptionProperty {
	return v.Properties
}

// TopicOptionConfigFieldsPropertiesTopicOptionProperty includes the requested fields of the GraphQL type TopicOptionProperty.
type TopicOptionConfigFieldsPropertiesTopicOptionProperty struct {
	Id               string                                                                    `json:"id"`
	Required         bool                                                                      `json:"required"`
	Name             string                                                                    `json:"name"`
	Weight           int                                                                       `json:"weight"`
	Comparable       bool                                                                      `json:"comparable"`
	Enabled          bool                                                                      `json:"enabled"`
	MaxSelection     int                                                                       `json:"maxSelection"`
	DefaultSelectAll bool                                                                      `json:"defaultSelectAll"`
	Options          []*TopicOptionConfigFieldsPropertiesTopicOptionPropertyOptionsTopicOption `json:"options"`
	ExtraOptionKey   ExtraOptionKey                                                            `json:"extraOptionKey"`
}

// GetId returns TopicOptionConfigFieldsPropertiesTopicOptionProperty.Id, and is useful for accessing the field via an interface.
func (v *TopicOptionConfigFieldsPropertiesTopicOptionProperty) GetId() string { return v.Id }

// GetRequired returns TopicOptionConfigFieldsPropertiesTopicOptionProperty.Required, and is useful for accessing the field via an interface.
func (v *TopicOptionConfigFieldsPropertiesTopicOptionProperty) GetRequired() bool { return v.Required }

// GetName returns TopicOptionConfigFieldsPropertiesTopicOptionProperty.Name, and is useful for accessing the field via an interface.
func (v *TopicOptionConfigFieldsPropertiesTopicOptionProperty) GetName() string { return v.Name }

// GetWeight returns TopicOptionConfigFieldsPropertiesTopicOptionProperty.Weight, and is useful for accessing the field via an interface.
func (v *TopicOptionConfigFieldsPropertiesTopicOptionProperty) GetWeight() int { return v.Weight }

// GetComparable returns TopicOptionConfigFieldsPropertiesTopicOptionProperty.Comparable, and is useful for accessing the field via an interface.
func (v *TopicOptionConfigFieldsPropertiesTopicOptionProperty) GetComparable() bool {
	return v.Comparable
}

// GetEnabled returns TopicOptionConfigFieldsPropertiesTopicOptionProperty.Enabled, and is useful for accessing the field via an interface.
func (v *TopicOptionConfigFieldsPropertiesTopicOptionProperty) GetEnabled() bool { return v.Enabled }

// GetMaxSelection returns TopicOptionConfigFieldsPropertiesTopicOptionProperty.MaxSelection, and is useful for accessing the field via an interface.
func (v *TopicOptionConfigFieldsPropertiesTopicOptionProperty) GetMaxSelection() int {
	return v.MaxSelection
}

// GetDefaultSelectAll returns TopicOptionConfigFieldsPropertiesTopicOptionProperty.DefaultSelectAll, and is useful for accessing the field via an interface.
func (v *TopicOptionConfigFieldsPropertiesTopicOptionProperty) GetDefaultSelectAll() bool {
	return v.DefaultSelectAll
}

// GetOptions returns TopicOptionConfigFieldsPropertiesTopicOptionProperty.Options, and is useful for accessing the field via an interface.
func (v *TopicOptionConfigFieldsPropertiesTopicOptionProperty) GetOptions() []*TopicOptionConfigFieldsPropertiesTopicOptionPropertyOptionsTopicOption {
	return v.Options
}

// GetExtraOptionKey returns TopicOptionConfigFieldsPropertiesTopicOptionProperty.ExtraOptionKey, and is useful for accessing the field via an interface.
func (v *TopicOptionConfigFieldsPropertiesTopicOptionProperty) GetExtraOptionKey() ExtraOptionKey {
	return v.ExtraOptionKey
}

// TopicOptionConfigFieldsPropertiesTopicOptionPropertyOptionsTopicOption includes the requested fields of the GraphQL type TopicOption.
type TopicOptionConfigFieldsPropertiesTopicOptionPropertyOptionsTopicOption struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

// GetName returns TopicOptionConfigFieldsPropertiesTopicOptionPropertyOptionsTopicOption.Name, and is useful for accessing the field via an interface.
func (v *TopicOptionConfigFieldsPropertiesTopicOptionPropertyOptionsTopicOption) GetName() string {
	return v.Name
}

// GetValue returns TopicOptionConfigFieldsPropertiesTopicOptionPropertyOptionsTopicOption.Value, and is useful for accessing the field via an interface.
func (v *TopicOptionConfigFieldsPropertiesTopicOptionPropertyOptionsTopicOption) GetValue() int {
	return v.Value
}

// __CreateMatchingValidInput is used internally by genqlient
type __CreateMatchingValidInput struct {
	UserId string `json:"userId"`
}

// GetUserId returns __CreateMatchingValidInput.UserId, and is useful for accessing the field via an interface.
func (v *__CreateMatchingValidInput) GetUserId() string { return v.UserId }

// __GetAreaInput is used internally by genqlient
type __GetAreaInput struct {
	Id string `json:"id"`
}

// GetId returns __GetAreaInput.Id, and is useful for accessing the field via an interface.
func (v *__GetAreaInput) GetId() string { return v.Id }

// __GetAvatarByIDsInput is used internally by genqlient
type __GetAvatarByIDsInput struct {
	Ids []string `json:"ids"`
}

// GetIds returns __GetAvatarByIDsInput.Ids, and is useful for accessing the field via an interface.
func (v *__GetAvatarByIDsInput) GetIds() []string { return v.Ids }

// __GetBlacklistRelationshipInput is used internally by genqlient
type __GetBlacklistRelationshipInput struct {
	Ids []string `json:"ids"`
}

// GetIds returns __GetBlacklistRelationshipInput.Ids, and is useful for accessing the field via an interface.
func (v *__GetBlacklistRelationshipInput) GetIds() []string { return v.Ids }

// __GetTopicConfigOptionInput is used internally by genqlient
type __GetTopicConfigOptionInput struct {
	TopicId string `json:"topicId"`
}

// GetTopicId returns __GetTopicConfigOptionInput.TopicId, and is useful for accessing the field via an interface.
func (v *__GetTopicConfigOptionInput) GetTopicId() string { return v.TopicId }

// __GetTopicConfigOptionsInput is used internally by genqlient
type __GetTopicConfigOptionsInput struct {
	Paginator *GraphQLPaginator `json:"paginator,omitempty"`
}

// GetPaginator returns __GetTopicConfigOptionsInput.Paginator, and is useful for accessing the field via an interface.
func (v *__GetTopicConfigOptionsInput) GetPaginator() *GraphQLPaginator { return v.Paginator }

// __GetTopicInput is used internally by genqlient
type __GetTopicInput struct {
	Id string `json:"id"`
}

// GetId returns __GetTopicInput.Id, and is useful for accessing the field via an interface.
func (v *__GetTopicInput) GetId() string { return v.Id }

// __GetTopicNameInput is used internally by genqlient
type __GetTopicNameInput struct {
	Id string `json:"id"`
}

// GetId returns __GetTopicNameInput.Id, and is useful for accessing the field via an interface.
func (v *__GetTopicNameInput) GetId() string { return v.Id }

// __GetUserByIDsInput is used internally by genqlient
type __GetUserByIDsInput struct {
	Ids []string `json:"ids"`
}

// GetIds returns __GetUserByIDsInput.Ids, and is useful for accessing the field via an interface.
func (v *__GetUserByIDsInput) GetIds() []string { return v.Ids }

// The query or mutation executed by CreateMatchingValid.
const CreateMatchingValid_Operation = `
query CreateMatchingValid ($userId: String!) {
	user(id: $userId) {
		blockInfo {
			userBlocked
			matchingBlocked
		}
	}
	userInfoCompletenessCheck(userId: $userId) {
		filled
	}
}
`

func CreateMatchingValid(
	ctx context.Context,
	client graphql.Client,
	userId string,
) (*CreateMatchingValidResponse, error) {
	req := &graphql.Request{
		OpName: "CreateMatchingValid",
		Query:  CreateMatchingValid_Operation,
		Variables: &__CreateMatchingValidInput{
			UserId: userId,
		},
	}
	var err error

	var data CreateMatchingValidResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}

// The query or mutation executed by GetArea.
const GetArea_Operation = `
query GetArea ($id: AreaCode!) {
	area(code: $id) {
		code
		name
		depth
	}
}
`

func GetArea(
	ctx context.Context,
	client graphql.Client,
	id string,
) (*GetAreaResponse, error) {
	req := &graphql.Request{
		OpName: "GetArea",
		Query:  GetArea_Operation,
		Variables: &__GetAreaInput{
			Id: id,
		},
	}
	var err error

	var data GetAreaResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}

// The query or mutation executed by GetAvatarByIDs.
const GetAvatarByIDs_Operation = `
query GetAvatarByIDs ($ids: [String!]!) {
	getUserByIds(ids: $ids) {
		id
		nickname
		avatar
	}
}
`

func GetAvatarByIDs(
	ctx context.Context,
	client graphql.Client,
	ids []string,
) (*GetAvatarByIDsResponse, error) {
	req := &graphql.Request{
		OpName: "GetAvatarByIDs",
		Query:  GetAvatarByIDs_Operation,
		Variables: &__GetAvatarByIDsInput{
			Ids: ids,
		},
	}
	var err error

	var data GetAvatarByIDsResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}

// The query or mutation executed by GetBlacklistRelationship.
const GetBlacklistRelationship_Operation = `
query GetBlacklistRelationship ($ids: [String!]!) {
	blacklistRelationship(ids: $ids) {
		a
		b
	}
}
`

func GetBlacklistRelationship(
	ctx context.Context,
	client graphql.Client,
	ids []string,
) (*GetBlacklistRelationshipResponse, error) {
	req := &graphql.Request{
		OpName: "GetBlacklistRelationship",
		Query:  GetBlacklistRelationship_Operation,
		Variables: &__GetBlacklistRelationshipInput{
			Ids: ids,
		},
	}
	var err error

	var data GetBlacklistRelationshipResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}

// The query or mutation executed by GetTopic.
const GetTopic_Operation = `
query GetTopic ($id: String!) {
	topic(id: $id) {
		id
	}
}
`

func GetTopic(
	ctx context.Context,
	client graphql.Client,
	id string,
) (*GetTopicResponse, error) {
	req := &graphql.Request{
		OpName: "GetTopic",
		Query:  GetTopic_Operation,
		Variables: &__GetTopicInput{
			Id: id,
		},
	}
	var err error

	var data GetTopicResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}

// The query or mutation executed by GetTopicConfigOption.
const GetTopicConfigOption_Operation = `
query GetTopicConfigOption ($topicId: String!) {
	topicOptionConfig(topicId: $topicId) {
		... TopicOptionConfigFields
	}
}
fragment TopicOptionConfigFields on TopicOptionConfig {
	topicId
	timeWeight
	threshold
	fuzzyMatchingTopic
	allowFuzzyMatching
	delayMinuteToPairWithFuzzyTopic
	properties {
		id
		required
		name
		weight
		comparable
		enabled
		maxSelection
		defaultSelectAll
		options {
			name
			value
		}
		extraOptionKey
	}
}
`

func GetTopicConfigOption(
	ctx context.Context,
	client graphql.Client,
	topicId string,
) (*GetTopicConfigOptionResponse, error) {
	req := &graphql.Request{
		OpName: "GetTopicConfigOption",
		Query:  GetTopicConfigOption_Operation,
		Variables: &__GetTopicConfigOptionInput{
			TopicId: topicId,
		},
	}
	var err error

	var data GetTopicConfigOptionResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}

// The query or mutation executed by GetTopicConfigOptions.
const GetTopicConfigOptions_Operation = `
query GetTopicConfigOptions ($paginator: GraphQLPaginator!) {
	topicOptionConfigs(paginator: $paginator) {
		... TopicOptionConfigFields
	}
}
fragment TopicOptionConfigFields on TopicOptionConfig {
	topicId
	timeWeight
	threshold
	fuzzyMatchingTopic
	allowFuzzyMatching
	delayMinuteToPairWithFuzzyTopic
	properties {
		id
		required
		name
		weight
		comparable
		enabled
		maxSelection
		defaultSelectAll
		options {
			name
			value
		}
		extraOptionKey
	}
}
`

func GetTopicConfigOptions(
	ctx context.Context,
	client graphql.Client,
	paginator *GraphQLPaginator,
) (*GetTopicConfigOptionsResponse, error) {
	req := &graphql.Request{
		OpName: "GetTopicConfigOptions",
		Query:  GetTopicConfigOptions_Operation,
		Variables: &__GetTopicConfigOptionsInput{
			Paginator: paginator,
		},
	}
	var err error

	var data GetTopicConfigOptionsResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}

// The query or mutation executed by GetTopicName.
const GetTopicName_Operation = `
query GetTopicName ($id: String!) {
	topic(id: $id) {
		id
		name
	}
}
`

func GetTopicName(
	ctx context.Context,
	client graphql.Client,
	id string,
) (*GetTopicNameResponse, error) {
	req := &graphql.Request{
		OpName: "GetTopicName",
		Query:  GetTopicName_Operation,
		Variables: &__GetTopicNameInput{
			Id: id,
		},
	}
	var err error

	var data GetTopicNameResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}

// The query or mutation executed by GetTopics.
const GetTopics_Operation = `
query GetTopics {
	topics {
		id
		name
	}
}
`

func GetTopics(
	ctx context.Context,
	client graphql.Client,
) (*GetTopicsResponse, error) {
	req := &graphql.Request{
		OpName: "GetTopics",
		Query:  GetTopics_Operation,
	}
	var err error

	var data GetTopicsResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}

// The query or mutation executed by GetUserByIDs.
const GetUserByIDs_Operation = `
query GetUserByIDs ($ids: [String!]!) {
	getUserByIds(ids: $ids) {
		id
		gender
		level
	}
}
`

func GetUserByIDs(
	ctx context.Context,
	client graphql.Client,
	ids []string,
) (*GetUserByIDsResponse, error) {
	req := &graphql.Request{
		OpName: "GetUserByIDs",
		Query:  GetUserByIDs_Operation,
		Variables: &__GetUserByIDsInput{
			Ids: ids,
		},
	}
	var err error

	var data GetUserByIDsResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}

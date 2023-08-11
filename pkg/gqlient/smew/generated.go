// Code generated by github.com/Khan/genqlient, DO NOT EDIT.

package smew

import (
	"context"

	"github.com/Khan/genqlient/graphql"
)

// CreateChatGroupResponse is returned by CreateChatGroup on success.
type CreateChatGroupResponse struct {
	// 【群组】创建匹配群组
	CreateGroup string `json:"createGroup"`
}

// GetCreateGroup returns CreateChatGroupResponse.CreateGroup, and is useful for accessing the field via an interface.
func (v *CreateChatGroupResponse) GetCreateGroup() string { return v.CreateGroup }

type CreateMotionGroupParam struct {
	// 匹配结果ID
	ResultId int `json:"resultId"`
	// 话题ID
	TopicId string `json:"topicId"`
	// 连接发起方ID
	FromUserId string `json:"fromUserId"`
	// 连接发起方Motion ID
	FromMotionId string `json:"fromMotionId"`
	// 连接发起方ID
	ToUserId string `json:"toUserId"`
	// 连接发起方Motion ID
	ToMotionId string `json:"toMotionId"`
}

// GetResultId returns CreateMotionGroupParam.ResultId, and is useful for accessing the field via an interface.
func (v *CreateMotionGroupParam) GetResultId() int { return v.ResultId }

// GetTopicId returns CreateMotionGroupParam.TopicId, and is useful for accessing the field via an interface.
func (v *CreateMotionGroupParam) GetTopicId() string { return v.TopicId }

// GetFromUserId returns CreateMotionGroupParam.FromUserId, and is useful for accessing the field via an interface.
func (v *CreateMotionGroupParam) GetFromUserId() string { return v.FromUserId }

// GetFromMotionId returns CreateMotionGroupParam.FromMotionId, and is useful for accessing the field via an interface.
func (v *CreateMotionGroupParam) GetFromMotionId() string { return v.FromMotionId }

// GetToUserId returns CreateMotionGroupParam.ToUserId, and is useful for accessing the field via an interface.
func (v *CreateMotionGroupParam) GetToUserId() string { return v.ToUserId }

// GetToMotionId returns CreateMotionGroupParam.ToMotionId, and is useful for accessing the field via an interface.
func (v *CreateMotionGroupParam) GetToMotionId() string { return v.ToMotionId }

// CreateMotionGroupResponse is returned by CreateMotionGroup on success.
type CreateMotionGroupResponse struct {
	// 【群组】创建Motion群组
	CreateMotionGroup string `json:"createMotionGroup"`
}

// GetCreateMotionGroup returns CreateMotionGroupResponse.CreateMotionGroup, and is useful for accessing the field via an interface.
func (v *CreateMotionGroupResponse) GetCreateMotionGroup() string { return v.CreateMotionGroup }

// CreateTimGroupResponse is returned by CreateTimGroup on success.
type CreateTimGroupResponse struct {
	// 【群组】TIM群组创建——对内部群组适用
	CreateTimGroup string `json:"createTimGroup"`
}

// GetCreateTimGroup returns CreateTimGroupResponse.CreateTimGroup, and is useful for accessing the field via an interface.
func (v *CreateTimGroupResponse) GetCreateTimGroup() string { return v.CreateTimGroup }

// DestroyGroupResponse is returned by DestroyGroup on success.
type DestroyGroupResponse struct {
	// 【群组】解散群组
	DestroyGroup string `json:"destroyGroup"`
}

// GetDestroyGroup returns DestroyGroupResponse.DestroyGroup, and is useful for accessing the field via an interface.
func (v *DestroyGroupResponse) GetDestroyGroup() string { return v.DestroyGroup }

// GroupMemberLeaveResponse is returned by GroupMemberLeave on success.
type GroupMemberLeaveResponse struct {
	// 【群组】群组成员离开
	GroupMemberLeave string `json:"groupMemberLeave"`
}

// GetGroupMemberLeave returns GroupMemberLeaveResponse.GroupMemberLeave, and is useful for accessing the field via an interface.
func (v *GroupMemberLeaveResponse) GetGroupMemberLeave() string { return v.GroupMemberLeave }

// SendTextMessageResponse is returned by SendTextMessage on success.
type SendTextMessageResponse struct {
	// 【消息】发送文字消息
	SendTextMessage string `json:"sendTextMessage"`
}

// GetSendTextMessage returns SendTextMessageResponse.SendTextMessage, and is useful for accessing the field via an interface.
func (v *SendTextMessageResponse) GetSendTextMessage() string { return v.SendTextMessage }

// __CreateChatGroupInput is used internally by genqlient
type __CreateChatGroupInput struct {
	ResultId  int      `json:"resultId"`
	TopicId   string   `json:"topicId"`
	MemberIds []string `json:"memberIds"`
}

// GetResultId returns __CreateChatGroupInput.ResultId, and is useful for accessing the field via an interface.
func (v *__CreateChatGroupInput) GetResultId() int { return v.ResultId }

// GetTopicId returns __CreateChatGroupInput.TopicId, and is useful for accessing the field via an interface.
func (v *__CreateChatGroupInput) GetTopicId() string { return v.TopicId }

// GetMemberIds returns __CreateChatGroupInput.MemberIds, and is useful for accessing the field via an interface.
func (v *__CreateChatGroupInput) GetMemberIds() []string { return v.MemberIds }

// __CreateMotionGroupInput is used internally by genqlient
type __CreateMotionGroupInput struct {
	Param CreateMotionGroupParam `json:"param"`
}

// GetParam returns __CreateMotionGroupInput.Param, and is useful for accessing the field via an interface.
func (v *__CreateMotionGroupInput) GetParam() CreateMotionGroupParam { return v.Param }

// __CreateTimGroupInput is used internally by genqlient
type __CreateTimGroupInput struct {
	ChatGroupId string `json:"chatGroupId"`
}

// GetChatGroupId returns __CreateTimGroupInput.ChatGroupId, and is useful for accessing the field via an interface.
func (v *__CreateTimGroupInput) GetChatGroupId() string { return v.ChatGroupId }

// __DestroyGroupInput is used internally by genqlient
type __DestroyGroupInput struct {
	GroupId string `json:"groupId"`
}

// GetGroupId returns __DestroyGroupInput.GroupId, and is useful for accessing the field via an interface.
func (v *__DestroyGroupInput) GetGroupId() string { return v.GroupId }

// __GroupMemberLeaveInput is used internally by genqlient
type __GroupMemberLeaveInput struct {
	GroupId string `json:"groupId"`
	UserId  string `json:"userId"`
}

// GetGroupId returns __GroupMemberLeaveInput.GroupId, and is useful for accessing the field via an interface.
func (v *__GroupMemberLeaveInput) GetGroupId() string { return v.GroupId }

// GetUserId returns __GroupMemberLeaveInput.UserId, and is useful for accessing the field via an interface.
func (v *__GroupMemberLeaveInput) GetUserId() string { return v.UserId }

// __SendTextMessageInput is used internally by genqlient
type __SendTextMessageInput struct {
	ChatGroupId string `json:"chatGroupId"`
	Sender      string `json:"sender"`
	Text        string `json:"text"`
}

// GetChatGroupId returns __SendTextMessageInput.ChatGroupId, and is useful for accessing the field via an interface.
func (v *__SendTextMessageInput) GetChatGroupId() string { return v.ChatGroupId }

// GetSender returns __SendTextMessageInput.Sender, and is useful for accessing the field via an interface.
func (v *__SendTextMessageInput) GetSender() string { return v.Sender }

// GetText returns __SendTextMessageInput.Text, and is useful for accessing the field via an interface.
func (v *__SendTextMessageInput) GetText() string { return v.Text }

// The query or mutation executed by CreateChatGroup.
const CreateChatGroup_Operation = `
mutation CreateChatGroup ($resultId: Int!, $topicId: String!, $memberIds: [String!]!) {
	createGroup(resultId: $resultId, topicId: $topicId, memberIds: $memberIds)
}
`

func CreateChatGroup(
	ctx context.Context,
	client graphql.Client,
	resultId int,
	topicId string,
	memberIds []string,
) (*CreateChatGroupResponse, error) {
	req := &graphql.Request{
		OpName: "CreateChatGroup",
		Query:  CreateChatGroup_Operation,
		Variables: &__CreateChatGroupInput{
			ResultId:  resultId,
			TopicId:   topicId,
			MemberIds: memberIds,
		},
	}
	var err error

	var data CreateChatGroupResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}

// The query or mutation executed by CreateMotionGroup.
const CreateMotionGroup_Operation = `
mutation CreateMotionGroup ($param: CreateMotionGroupParam!) {
	createMotionGroup(param: $param)
}
`

func CreateMotionGroup(
	ctx context.Context,
	client graphql.Client,
	param CreateMotionGroupParam,
) (*CreateMotionGroupResponse, error) {
	req := &graphql.Request{
		OpName: "CreateMotionGroup",
		Query:  CreateMotionGroup_Operation,
		Variables: &__CreateMotionGroupInput{
			Param: param,
		},
	}
	var err error

	var data CreateMotionGroupResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}

// The query or mutation executed by CreateTimGroup.
const CreateTimGroup_Operation = `
mutation CreateTimGroup ($chatGroupId: String!) {
	createTimGroup(chatGroupId: $chatGroupId)
}
`

func CreateTimGroup(
	ctx context.Context,
	client graphql.Client,
	chatGroupId string,
) (*CreateTimGroupResponse, error) {
	req := &graphql.Request{
		OpName: "CreateTimGroup",
		Query:  CreateTimGroup_Operation,
		Variables: &__CreateTimGroupInput{
			ChatGroupId: chatGroupId,
		},
	}
	var err error

	var data CreateTimGroupResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}

// The query or mutation executed by DestroyGroup.
const DestroyGroup_Operation = `
mutation DestroyGroup ($groupId: String!) {
	destroyGroup(groupId: $groupId)
}
`

func DestroyGroup(
	ctx context.Context,
	client graphql.Client,
	groupId string,
) (*DestroyGroupResponse, error) {
	req := &graphql.Request{
		OpName: "DestroyGroup",
		Query:  DestroyGroup_Operation,
		Variables: &__DestroyGroupInput{
			GroupId: groupId,
		},
	}
	var err error

	var data DestroyGroupResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}

// The query or mutation executed by GroupMemberLeave.
const GroupMemberLeave_Operation = `
mutation GroupMemberLeave ($groupId: String!, $userId: String!) {
	groupMemberLeave(groupId: $groupId, userId: $userId)
}
`

func GroupMemberLeave(
	ctx context.Context,
	client graphql.Client,
	groupId string,
	userId string,
) (*GroupMemberLeaveResponse, error) {
	req := &graphql.Request{
		OpName: "GroupMemberLeave",
		Query:  GroupMemberLeave_Operation,
		Variables: &__GroupMemberLeaveInput{
			GroupId: groupId,
			UserId:  userId,
		},
	}
	var err error

	var data GroupMemberLeaveResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}

// The query or mutation executed by SendTextMessage.
const SendTextMessage_Operation = `
mutation SendTextMessage ($chatGroupId: String!, $sender: String!, $text: String!) {
	sendTextMessage(chatGroupId: $chatGroupId, sender: $sender, text: $text)
}
`

func SendTextMessage(
	ctx context.Context,
	client graphql.Client,
	chatGroupId string,
	sender string,
	text string,
) (*SendTextMessageResponse, error) {
	req := &graphql.Request{
		OpName: "SendTextMessage",
		Query:  SendTextMessage_Operation,
		Variables: &__SendTextMessageInput{
			ChatGroupId: chatGroupId,
			Sender:      sender,
			Text:        text,
		},
	}
	var err error

	var data SendTextMessageResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}

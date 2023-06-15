package whalecode

import "github.com/letjoy-club/mida-tool/midacode"

var (
	ErrCannotModifyOtherMatched      = midacode.NewError("CANNOT_MODIFY_OTHER_MATCHED", "不能修改别人的匹配", midacode.LogLevelWarn)
	ErrMatchingAlreadyCanceled       = midacode.NewError("MATCHING_ALREADY_CANCELED", "匹配已经取消", midacode.LogLevelWarn)
	ErrMatchingStateShouldBeMatched  = midacode.NewError("MATCHING_STATE_SHOULD_BE_MATCHED", "匹配状态应该为已匹配", midacode.LogLevelWarn)
	ErrMatchingStateShouldBeMatching = midacode.NewError("MATCHING_STATE_SHOULD_BE_MATCHING", "匹配状态应该为匹配中", midacode.LogLevelWarn)
	ErrResourceBusy                  = midacode.NewError("RESOURCE_BUSY", "资源繁忙", midacode.LogLevelWarn)

	ErrCannotInviteSelf         = midacode.NewError("CANNOT_INVITE_SELF", "不能邀请自己", midacode.LogLevelWarn)
	ErrMatchingQuotaNotEnough   = midacode.NewError("MATCHING_QUOTA_NOT_ENOUGH", "匹配次数不足，请结束现有聊天后再继续", midacode.LogLevelWarn)
	ErrQueryDurationTooLong     = midacode.NewError("QUERY_DURATION_TOO_LONG", "查询时间段过长", midacode.LogLevelWarn)
	ErrTopicIsAlreadyInMatching = midacode.NewError("TOPIC_IS_ALREADY_IN_MATCHING", "该话题已经在匹配中", midacode.LogLevelWarn)
	ErrUserGenderIsNotSet       = midacode.NewError("USER_GENDER_IS_NOT_SET", "账号性别信息为空，请先补充性别", midacode.LogLevelWarn)
	ErrUserIDCannotBeEmpty      = midacode.NewError("USER_ID_CANNOT_BE_EMPTY", "用户 id 不能为空", midacode.LogLevelInfo)

	ErrCannotInviteOtherTwiceForTheSameTopic          = midacode.NewError("CANNOT_INVITE_OTHER_TWICE_FOR_THE_SAME_TOPIC", "不能连续邀请同一个人两次同个话题", midacode.LogLevelWarn)
	ErrCannotPerformActionWhenChatGroupAlreadyCreated = midacode.NewError("CANNOT_PERFORM_ACTION_WHEN_CHAT_GROUP_ALREADY_CREATED", "聊天室已经创建，不能执行该操作", midacode.LogLevelWarn)
	ErrMatchingDurationQuotaNotEnough                 = midacode.NewError("MATCHING_DURATION_QUOTA_NOT_ENOUGH", "当前时间段匹配次数用尽", midacode.LogLevelWarn)
	ErrMatchingNotMatchWithTopic                      = midacode.NewError("MATCHING_NOT_MATCH_WITH_TOPIC", "匹配信息与话题地区/话题不匹配", midacode.LogLevelWarn)
)

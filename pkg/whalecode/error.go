package whalecode

import "github.com/letjoy-club/mida-tool/midacode"

var (
	ErrCannotModifyOtherMatched      = midacode.NewError("CANNOT_MODIFY_OTHER_MATCHED", "不能修改别人的匹配", midacode.LogLevelWarn)
	ErrMatchingAlreadyCanceled       = midacode.NewError("MATCHING_ALREADY_CANCELED", "匹配已经取消", midacode.LogLevelWarn)
	ErrMatchingStateShouldBeMatched  = midacode.NewError("MATCHING_STATE_SHOULD_BE_MATCHED", "匹配状态应该为已匹配", midacode.LogLevelWarn)
	ErrMatchingStateShouldBeMatching = midacode.NewError("MATCHING_STATE_SHOULD_BE_MATCHING", "匹配状态应该为匹配中", midacode.LogLevelWarn)
	ErrDayRangeNumInvalid            = midacode.NewError("DAY_RANGE_NUM_INVALID", "日期区间不成对", midacode.LogLevelWarn)
	ErrDayRangeInvalid               = midacode.NewError("DAY_RANGE_INVALID", "日期区间不正确", midacode.LogLevelWarn)
	ErrDayRangeDateFormatInvalid     = midacode.NewError("DAY_RANGE_DATE_FORMAT_INVALID", "日期区间日期格式不正确", midacode.LogLevelWarn)
	ErrResourceBusy                  = midacode.NewError("RESOURCE_BUSY", "资源繁忙", midacode.LogLevelWarn)

	ErrAreaNotSupport           = midacode.NewError("AREA_NOT_SUPPORT", "该地区非法或未上线", midacode.LogLevelWarn)
	ErrCannotInviteSelf         = midacode.NewError("CANNOT_INVITE_SELF", "不能邀请自己", midacode.LogLevelWarn)
	ErrCannotSelectGender       = midacode.NewError("CANNOT_SELECT_GENDER", "您暂时无法选择性别", midacode.LogLevelWarn)
	ErrMatchingQuotaNotEnough   = midacode.NewError("MATCHING_QUOTA_NOT_ENOUGH", "匹配次数不足，请结束现有聊天后再继续", midacode.LogLevelWarn)
	ErrQueryDurationTooLong     = midacode.NewError("QUERY_DURATION_TOO_LONG", "查询时间段过长", midacode.LogLevelWarn)
	ErrTopicIsAlreadyInMatching = midacode.NewError("TOPIC_IS_ALREADY_IN_MATCHING", "该话题已经在匹配中", midacode.LogLevelWarn)
	ErrTopicNotExisted          = midacode.NewError("TOPIC_NOT_EXISTED", "话题不存在", midacode.LogLevelWarn)
	ErrUserIDCannotBeEmpty      = midacode.NewError("USER_ID_CANNOT_BE_EMPTY", "用户 id 不能为空", midacode.LogLevelWarn)
	ErrUserInfoNotComplete      = midacode.NewError("USER_INFO_NOT_COMPLETE", "用户信息未完善，请先补充", midacode.LogLevelWarn)
	ErrUserBlocked              = midacode.NewError("USER_BLOCKED", "用户已被封禁，无法发起匹配", midacode.LogLevelWarn)
	ErrUserInBlacklist          = midacode.NewError("USER_IN_BLACKLIST", "对方已将您加入黑名单，无法邀请", midacode.LogLevelWarn)
	ErrInviteeInBlacklist       = midacode.NewError("INVITEE_IN_BLACKLIST", "您已将对方加入黑名单，无法邀请", midacode.LogLevelWarn)
	ErrInviteeNotExist          = midacode.NewError("INVITEE_NOT_EXISTED", "您邀请的人不存在", midacode.LogLevelWarn)
	ErrInviterBlocked           = midacode.NewError("INVITER_BLOCKED", "您已被封禁，无法发起邀请", midacode.LogLevelWarn)
	ErrInviteeBlocked           = midacode.NewError("INVITEE_BLOCKED", "对方已被封禁，无法发起邀请", midacode.LogLevelWarn)

	ErrIsAlreadyHasActiveMotionOfTopic              = midacode.NewError("IS_ALREADY_HAS_ACTIVE_MOTION_OF_TOPIC", "已经有该话题的活动卡片", midacode.LogLevelWarn)
	ErrMotionOfferQuotaNotEnough                    = midacode.NewError("MOTION_OFFER_QUOTA_NOT_ENOUGH", "邀约次数不足", midacode.LogLevelWarn)
	MotionOfferIsNotPending                         = midacode.NewError("MOTION_OFFER_IS_NOT_PENDING", "邀约不是待处理状态", midacode.LogLevelWarn)
	ErrMatchingOfferIsClosed                        = midacode.NewError("MATCHING_OFFER_IS_CLOSED", "邀约已经关闭", midacode.LogLevelWarn)
	ErrChatChanceNotEnough                          = midacode.NewError("CHAT_CHANCE_NOT_ENOUGH", "当前邀请聊天机会不足", midacode.LogLevelWarn)
	ErrOnlyChatWhenNotClosed                        = midacode.NewError("ONLY_CHAT_WHEN_NOT_CLOSED", "聊天已结束，无法发送消息", midacode.LogLevelWarn)
	ErrCannotSendMatchingOfferToSelf                = midacode.NewError("CANNOT_SEND_MATCHING_OFFER_TO_SELF", "不能给自己发邀约", midacode.LogLevelWarn)
	ErrCannotSendMatchingOfferToDifferentTopic      = midacode.NewError("CANNOT_SEND_MATCHING_OFFER_TO_DIFFERENT_TOPIC", "邀约话题不一致", midacode.LogLevelWarn)
	ErrYourMatchingNotInMatchingState               = midacode.NewError("YOUR_MATCHING_NOT_IN_MATCHING_STATE", "你的匹配不在匹配中", midacode.LogLevelWarn)
	ErrTheMatchingIsNotInMatchingState              = midacode.NewError("THE_MATCHING_IS_NOT_IN_MATCHING_STATE", "对方匹配不在匹配中", midacode.LogLevelWarn)
	ErrAlreadySentOutMatchingOffer                  = midacode.NewError("ALREADY_SENT_OUT_MATCHING_OFFER", "已经发出过邀约", midacode.LogLevelWarn)
	ErrMatchingOfferNotFound                        = midacode.NewError("MATCHING_OFFER_NOT_FOUND", "邀约不存在", midacode.LogLevelWarn)
	ErrMatchingOfferIsCancelableOnlyWhenUnprocessed = midacode.NewError("MATCHING_OFFER_IS_CANCELABLE_ONLY_WHEN_UNPROCESSED", "邀约只有在未处理状态下才能取消", midacode.LogLevelWarn)
	ErrMatchingOfferIsAcceptableOnlyWhenUnprocessed = midacode.NewError("MATCHING_OFFER_IS_ACCEPTABLE_ONLY_WHEN_UNPROCESSED", "邀约只有在未处理状态下才能接受", midacode.LogLevelWarn)
	ErrMatchingOfferIsNotActive                     = midacode.NewError("MATCHING_OFFER_IS_NOT_ACTIVE", "邀约已经失效", midacode.LogLevelWarn)
	ErrYourMotionIsNotActive                        = midacode.NewError("YOUR_MOTION_IS_NOT_ACTIVE", "你的活动已关闭", midacode.LogLevelWarn)
	ErrTheMotionIsNotActive                         = midacode.NewError("THE_MOTION_IS_NOT_ACTIVE", "对方的活动不在已关闭", midacode.LogLevelWarn)
	ErrFinishChatBeforeCloseMotion                  = midacode.NewError("FINISH_CHAT_BEFORE_CLOSE_MOTION", "请先结束聊天再关闭活动", midacode.LogLevelWarn)

	ErrCannotInviteOtherTwiceForTheSameTopic          = midacode.NewError("CANNOT_INVITE_OTHER_TWICE_FOR_THE_SAME_TOPIC", "不能连续邀请同一个人两次同个话题", midacode.LogLevelWarn)
	ErrCannotPerformActionWhenChatGroupAlreadyCreated = midacode.NewError("CANNOT_PERFORM_ACTION_WHEN_CHAT_GROUP_ALREADY_CREATED", "聊天室已经创建，不能执行该操作", midacode.LogLevelWarn)
	ErrMatchingDurationQuotaNotEnough                 = midacode.NewError("MATCHING_DURATION_QUOTA_NOT_ENOUGH", "当前时间段匹配次数用尽", midacode.LogLevelWarn)
	ErrStartDateShouldLessThanNow                     = midacode.NewError("START_DATE_SHOULD_LESS_THAN_NOW", "开始时间不能大于当前时间", midacode.LogLevelWarn)
	ErrStopDateShoulGreateThanNow                     = midacode.NewError("STOP_DATE_SHOULD_GREATE_THAN_NOW", "结束时间不能小于当前时间", midacode.LogLevelWarn)
	ErrMatchingNotMatchWithTopic                      = midacode.NewError("MATCHING_NOT_MATCH_WITH_TOPIC", "匹配信息与话题地区/话题不匹配", midacode.LogLevelWarn)
)

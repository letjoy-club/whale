package whalecode

import "github.com/letjoy-club/mida-tool/midacode"

var (
	ErrMatchingStateShouldBeMatched  = midacode.NewError("MATCHING_STATE_SHOULD_BE_MATCHED", "匹配状态应该为已匹配", midacode.LogLevelWarn)
	ErrMatchingStateShouldBeMatching = midacode.NewError("MATCHING_STATE_SHOULD_BE_MATCHING", "匹配状态应该为匹配中", midacode.LogLevelWarn)
	ErrCannotModifyOtherMatched      = midacode.NewError("CANNOT_MODIFY_OTHER_MATCHED", "不能修改别人的匹配", midacode.LogLevelWarn)

	ErrUserIDCannotBeEmpty      = midacode.NewError("USER_ID_CANNOT_BE_EMPTY", "用户 id 不能为空", midacode.LogLevelInfo)
	ErrMatchingQuotaNotEnough   = midacode.NewError("MATCHING_QUOTA_NOT_ENOUGH", "匹配次数不足", midacode.LogLevelWarn)
	ErrTopicIsAlreadyInMatching = midacode.NewError("TOPIC_IS_ALREADY_IN_MATCHING", "该话题已经在匹配中", midacode.LogLevelWarn)
)

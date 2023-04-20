package whalecode

import "github.com/letjoy-club/mida-tool/midacode"

var (
	ErrMatchingStateShouldBeMatched  = midacode.NewError("MATCHING_STATE_SHOULD_BE_MATCHED", "匹配状态应该为已匹配", true)
	ErrMatchingStateShouldBeMatching = midacode.NewError("MATCHING_STATE_SHOULD_BE_MATCHING", "匹配状态应该为匹配中", true)
	ErrCannotModifyOtherMatched      = midacode.NewError("CANNOT_MODIFY_OTHER_MATCHED", "不能修改别人的匹配", true)

	ErrUserIDCannotBeEmpty = midacode.NewError("USER_ID_CANNOT_BE_EMPTY", "用户 id 不能为空", true)
)

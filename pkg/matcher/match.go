package matcher

import (
	"context"
	"whale/pkg/models"
)

type MatchingReason string

const (
	MatchingReasonOk                MatchingReason = "Ok"
	MatchingReasonAreaNotOverlapse  MatchingReason = "AreaNotOverlapse"
	MatchingReasonGenderNotMatched  MatchingReason = "GenderNotMatched"
	MatchingReasonCannotMatchItSelf MatchingReason = "CannotMatchItSelf"
)

func Matched(ctx context.Context, m1 *models.Matching, m2 *models.Matching) (MatchingReason, bool) {
	if m1.UserID == m2.UserID {
		return MatchingReasonCannotMatchItSelf, false
	}
	mc := GetMatchingContext(ctx)
	userIGender := mc.userProfiles[m1.UserID].Gender.String()
	userIWantGender := m1.Gender
	userJGender := mc.userProfiles[m2.UserID].Gender.String()
	userJWantGender := m2.Gender
	// 如果 I 用户 希望的性别和 J 用户一致
	if userIWantGender != models.GenderN.String() && userIWantGender != userJGender {
		return MatchingReasonGenderNotMatched, false
	}
	// 如果 J 用户 希望的性别和 I 用户一致
	if userJWantGender != models.GenderN.String() && userJWantGender != userIGender {
		return MatchingReasonGenderNotMatched, false
	}
	_, hasOverlapse := AreaOverlapse(m1.AreaIDs, m2.AreaIDs)
	if hasOverlapse {
		return MatchingReasonOk, true
	}
	return MatchingReasonAreaNotOverlapse, false
}
package modelutil

import (
	"context"
	"fmt"
	"time"
	"whale/pkg/dbquery"
	"whale/pkg/gqlient/hoopoe"
	"whale/pkg/gqlient/scream"
	"whale/pkg/gqlient/smew"
	"whale/pkg/loader"
	"whale/pkg/models"
	"whale/pkg/whalecode"

	"github.com/letjoy-club/mida-tool/dbutil"
	"github.com/letjoy-club/mida-tool/midacode"
	"github.com/letjoy-club/mida-tool/midacontext"
	"github.com/letjoy-club/mida-tool/shortid"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

func GetMatchingAndCheckUser(ctx context.Context, matchingID, uid string) (*models.Matching, error) {
	thunk := midacontext.GetLoader[loader.Loader](ctx).Matching.Load(ctx, matchingID)
	matching, err := thunk()
	if err != nil {
		return nil, err
	}
	if matching.State != models.MatchingStateMatching.String() {
		return nil, whalecode.ErrMatchingAlreadyCanceled
	}

	if uid != "" {
		if matching.UserID != uid {
			return nil, whalecode.ErrCannotModifyOtherMatched
		}
	}
	return matching, nil
}

func ConfirmMatching(ctx context.Context, db *gorm.DB, matchingResult *models.MatchingResult, matchingID, uid string, confirmed bool) error {
	action := &models.MatchingResultConfirmAction{
		UserID:           uid,
		Confirmed:        confirmed,
		MatchingResultID: matchingResult.ID,
	}
	err := dbquery.Use(db).MatchingResultConfirmAction.WithContext(ctx).Create(action)
	if err != nil {
		return err
	}
	defer func() {
		midacontext.GetLoader[loader.Loader](ctx).MatchingResult.Clear(ctx, matchingResult.ID)
		for _, id := range matchingResult.MatchingIDs {
			midacontext.GetLoader[loader.Loader](ctx).Matching.Clear(ctx, id)
		}
	}()

	if matchingResult.ChatGroupState != models.ChatGroupStateUncreated.String() {
		// 已经创建了聊天室
		return nil
	}
	index := lo.IndexOf(matchingResult.MatchingIDs, matchingID)
	if index == -1 {
		return nil
	}
	if confirmed {
		matchingResult.ConfirmStates[index] = models.MatchingResultConfirmStateConfirmed.String()
	} else {
		matchingResult.ConfirmStates[index] = models.MatchingResultConfirmStateRejected.String()
		matchingResult.Closed = true
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		MatchingResult := dbquery.Use(tx).MatchingResult
		_, err = MatchingResult.WithContext(ctx).
			Where(MatchingResult.ID.Eq(matchingResult.ID)).
			Select(MatchingResult.ConfirmStates, MatchingResult.Closed).
			Updates(matchingResult)
		if err != nil {
			return err
		}
		if !matchingResult.Closed {
			return nil
		}
		// 如果有人拒绝匹配
		Matching := dbquery.Use(tx).Matching
		_, err = Matching.WithContext(ctx).
			Where(Matching.ID.In(matchingResult.MatchingIDs...)).
			UpdateSimple(
				// 回到匹配状态
				Matching.State.Value(string(models.MatchingStateMatching)),
				Matching.MatchedAt.Null(),
				Matching.ResultID.Value(0),
			)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func CheckMatchingResultAndCreateChatGroup(ctx context.Context, m *models.MatchingResult) error {
	// 所有匹配都是确认了
	if lo.Count(m.ConfirmStates, models.InvitationConfirmStateConfirmed.String()) != len(m.ConfirmStates) {
		return nil
	}

	resp, err := smew.CreateChatGroup(ctx, midacontext.GetServices(ctx).Smew, m.ID, m.TopicID, m.UserIDs)
	if err != nil {
		return err
	}
	groupID := resp.CreateGroup
	db := dbutil.GetDB(ctx)
	err = db.Transaction(func(db *gorm.DB) error {
		Matching := dbquery.Use(db).Matching
		MatchingResult := dbquery.Use(db).MatchingResult
		_, err := MatchingResult.WithContext(ctx).
			Where(MatchingResult.ID.Eq(m.ID)).
			UpdateSimple(
				MatchingResult.ChatGroupState.Value(models.ChatGroupStateCreated.String()),
				MatchingResult.ChatGroupID.Value(groupID),
				MatchingResult.ChatGroupCreatedAt.Value(time.Now()),
			)
		if err != nil {
			return err
		}

		_, err = Matching.WithContext(ctx).
			Where(Matching.ID.In(m.MatchingIDs...)).
			UpdateSimple(
				Matching.InChatGroup.Value(true),
				Matching.ChatGroupState.Value(models.ChatGroupStateCreated.String()),
			)
		return err
	})
	return err
}

func CreateMatching(ctx context.Context, uid string, param models.CreateMatchingParam) (*models.Matching, error) {
	err := checkMatchingParam(ctx, uid, param.TopicID, param.CityID, param.Gender)
	if err != nil {
		return nil, err
	}

	db := dbutil.GetDB(ctx)
	Matching := dbquery.Use(db).Matching

	_, err = Matching.WithContext(ctx).Where(
		Matching.TopicID.Eq(param.TopicID),
		Matching.UserID.Eq(uid),
		Matching.State.In(
			string(models.MatchingStateMatching),
			string(models.MatchingStateMatched),
		),
	).Take()
	// 如果有找到，或者其他数据库错误
	notFoundErr := midacode.ItemCustomNotFound(err, midacode.ErrItemNotFound)
	if notFoundErr != midacode.ErrItemNotFound {
		return nil, whalecode.ErrTopicIsAlreadyInMatching
	}

	thunk := midacontext.GetLoader[loader.Loader](ctx).UserProfile.Load(ctx, uid)
	profile, err := thunk()

	startMatchingAt := time.Now().Add(time.Minute * 30)
	matching := &models.Matching{
		ID:              shortid.NewWithTime("m_", 6),
		TopicID:         param.TopicID,
		UserID:          uid,
		State:           models.MatchingStateMatching.String(),
		Gender:          param.Gender.String(),
		Remark:          *param.Remark,
		CityID:          param.CityID,
		ChatGroupState:  string(models.ChatGroupStateUncreated),
		Deadline:        time.Now().Add(time.Hour * 24 * 7),
		AreaIDs:         param.AreaIds,
		MyGender:        profile.Gender.String(),
		Discoverable:    true,
		StartMatchingAt: &startMatchingAt,
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		err = dbquery.Use(tx).Matching.WithContext(ctx).Create(matching)
		if err != nil {
			return err
		}

		MatchingQuota := dbquery.Use(tx).MatchingQuota
		_, err = MatchingQuota.WithContext(ctx).
			Where(MatchingQuota.UserID.Eq(uid)).
			UpdateSimple(
				MatchingQuota.Remain.Add(-1),
				MatchingQuota.MatchingNum.Add(1),
			)
		if err != nil {
			return err
		}

		MatchingDurationConstraint := dbquery.Use(tx).MatchingDurationConstraint
		_, err = MatchingDurationConstraint.WithContext(ctx).
			Where(MatchingDurationConstraint.UserID.Eq(uid)).
			UpdateSimple(MatchingDurationConstraint.Remain.Add(-1))
		return err
	})
	if err == nil {
		midacontext.GetLoader[loader.Loader](ctx).MatchingQuota.Clear(ctx, uid)
		midacontext.GetLoader[loader.Loader](ctx).MatchingDurationConstraint.Clear(ctx, uid)
		if err := PublishMatchingCreatedEvent(ctx, matching); err != nil {
			fmt.Println("publishMatchingCreatedEvent err:", err)
		}
	}
	RecordUserJoinTopic(ctx, matching.TopicID, matching.CityID, matching.UserID, matching.ID)
	return matching, err
}

func checkMatchingParam(ctx context.Context, uid, topicID, cityID string, gender models.Gender) error {
	// 基础检查
	res, err := hoopoe.CreateMatchingCheck(ctx, midacontext.GetServices(ctx).Hoopoe, topicID, cityID, uid)
	if err != nil {
		return err
	}
	if res.Topic == nil || !res.Topic.Enable {
		return whalecode.ErrTopicNotExisted
	}
	if res.Area == nil || !res.Area.Enabled {
		return whalecode.ErrAreaNotSupport
	}
	if !res.GetUserInfoCompletenessCheck().Filled {
		return whalecode.ErrUserInfoNotComplete
	}
	if res.User.BlockInfo.UserBlocked || res.User.BlockInfo.MatchingBlocked {
		return whalecode.ErrUserBlocked
	}
	//if gender != models.GenderN {
	//	if !res.LevelDetail.Rights.GenderSelection {
	//		return whalecode.ErrCannotSelectGender
	//	}
	//}
	// 额度检查
	thunk := midacontext.GetLoader[loader.Loader](ctx).MatchingQuota.Load(ctx, uid)
	quota, err := thunk()
	if err != nil {
		return err
	}
	if quota.Remain <= 0 {
		return whalecode.ErrMatchingQuotaNotEnough
	}

	constraintThunk := midacontext.GetLoader[loader.Loader](ctx).MatchingDurationConstraint.Load(ctx, uid)
	constraint, err := constraintThunk()
	if err != nil {
		return err
	}
	if constraint.Remain <= 0 {
		return whalecode.ErrMatchingDurationQuotaNotEnough
	}
	return nil
}

func SimplifyPreferredPeriods(periods []models.DatePeriod) []string {
	if lo.Contains(periods, models.DatePeriodUnlimited) {
		return []string{}
	}
	periodSet := map[models.DatePeriod]struct{}{}
	for _, period := range periods {
		periodSet[period] = struct{}{}
	}

	// 如果周末晚上和下午都可以，合并为周末
	if _, ok := periodSet[models.DatePeriodWeekendNight]; ok {
		if _, ok := periodSet[models.DatePeriodWeekendAfternoon]; ok {
			periodSet[models.DatePeriodWeekend] = struct{}{}
		}
	}
	// 如果工作日晚上和下午都可以，合并为工作日
	if _, ok := periodSet[models.DatePeriodWorkdayNight]; ok {
		if _, ok := periodSet[models.DatePeriodWorkdayAfternoon]; ok {
			periodSet[models.DatePeriodWorkday] = struct{}{}
		}
	}
	// 如果周末和工作日都可以，合并为不限
	if _, ok := periodSet[models.DatePeriodWeekend]; ok {
		if _, ok := periodSet[models.DatePeriodWorkday]; ok {
			return []string{}
		}
	}

	if _, ok := periodSet[models.DatePeriodWeekend]; ok {
		delete(periodSet, models.DatePeriodWeekendNight)
		delete(periodSet, models.DatePeriodWeekendAfternoon)
	}
	if _, ok := periodSet[models.DatePeriodWorkday]; ok {
		delete(periodSet, models.DatePeriodWorkdayNight)
		delete(periodSet, models.DatePeriodWorkdayAfternoon)
	}

	return lo.MapToSlice(periodSet, func(period models.DatePeriod, _ struct{}) string {
		return period.String()
	})
}

func CreateMatchingV2(ctx context.Context, uid string, param models.CreateMatchingParamV2) (*models.Matching, error) {
	err := checkMatchingParam(ctx, uid, param.TopicID, param.CityID, param.Gender)
	if err != nil {
		return nil, err
	}

	db := dbutil.GetDB(ctx)
	Matching := dbquery.Use(db).Matching

	_, err = Matching.WithContext(ctx).Where(
		Matching.TopicID.Eq(param.TopicID),
		Matching.UserID.Eq(uid),
		Matching.State.In(
			string(models.MatchingStateMatching),
			string(models.MatchingStateMatched),
		),
	).Take()
	// 如果有找到，或者其他数据库错误
	notFoundErr := midacode.ItemCustomNotFound(err, midacode.ErrItemNotFound)
	if notFoundErr != midacode.ErrItemNotFound {
		return nil, whalecode.ErrTopicIsAlreadyInMatching
	}

	matching := &models.Matching{
		ID:               shortid.NewWithTime("m_", 6),
		TopicID:          param.TopicID,
		UserID:           uid,
		State:            models.MatchingStateMatching.String(),
		Gender:           param.Gender.String(),
		Remark:           *param.Remark,
		DayRange:         param.DayRange,
		PreferredPeriods: SimplifyPreferredPeriods(param.PreferredPeriods),
		Properties: lo.Map(param.Properties, func(p *models.MatchingPropertyParam, i int) models.MatchingProperty {
			return models.MatchingProperty{ID: p.ID, Values: p.Values}
		}),
		CityID:         param.CityID,
		ChatGroupState: string(models.ChatGroupStateUncreated),
		Deadline:       time.Now().Add(time.Hour * 24 * 7),
		AreaIDs:        param.AreaIds,
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		err = dbquery.Use(tx).Matching.WithContext(ctx).Create(matching)
		if err != nil {
			return err
		}

		MatchingQuota := dbquery.Use(tx).MatchingQuota
		_, err = MatchingQuota.WithContext(ctx).
			Where(MatchingQuota.UserID.Eq(uid)).
			UpdateSimple(
				MatchingQuota.Remain.Add(-1),
				MatchingQuota.MatchingNum.Add(1),
			)
		if err != nil {
			return err
		}

		MatchingDurationConstraint := dbquery.Use(tx).MatchingDurationConstraint
		_, err = MatchingDurationConstraint.WithContext(ctx).
			Where(MatchingDurationConstraint.UserID.Eq(uid)).
			UpdateSimple(MatchingDurationConstraint.Remain.Add(-1))
		return err
	})
	if err == nil {
		midacontext.GetLoader[loader.Loader](ctx).MatchingQuota.Clear(ctx, uid)
		midacontext.GetLoader[loader.Loader](ctx).MatchingDurationConstraint.Clear(ctx, uid)
		if err := PublishMatchingCreatedEvent(ctx, matching); err != nil {
			fmt.Println("publishMatchingCreatedEvent err:", err)
		}
	}
	RecordUserJoinTopic(ctx, matching.TopicID, matching.CityID, matching.UserID, matching.ID)
	return matching, err
}

func CreateMatchingInvitation(ctx context.Context, uid string, param models.CreateMatchingInvitationParam) (*models.MatchingInvitation, error) {
	// 额度检查
	thunk := midacontext.GetLoader[loader.Loader](ctx).MatchingQuota.Load(ctx, uid)
	quota, err := thunk()
	if err != nil {
		return nil, err
	}
	if quota.Remain <= 0 {
		return nil, whalecode.ErrMatchingQuotaNotEnough
	}

	constraintThunk := midacontext.GetLoader[loader.Loader](ctx).MatchingDurationConstraint.Load(ctx, uid)
	constraint, err := constraintThunk()
	if err != nil {
		return nil, err
	}
	if constraint.Remain <= 0 {
		return nil, whalecode.ErrMatchingDurationQuotaNotEnough
	}
	// 基础检查
	ids := []string{uid, param.InviteeID}
	res, err := hoopoe.CreateMatchingInvitationCheck(ctx, midacontext.GetServices(ctx).Hoopoe, param.TopicID, param.CityID, ids)
	if res.Topic == nil || !res.Topic.Enable {
		return nil, whalecode.ErrTopicNotExisted
	}
	if res.Area == nil || !res.Area.Enabled {
		return nil, whalecode.ErrAreaNotSupport
	}
	if res.BlacklistRelationship != nil {
		for _, pair := range res.BlacklistRelationship {
			if pair.A == uid {
				return nil, whalecode.ErrInviteeInBlacklist
			} else {
				return nil, whalecode.ErrUserInBlacklist
			}
		}
	}
	if len(res.GetGetUserByIdsV2()) != 2 {
		return nil, whalecode.ErrInviteeNotExist
	}
	for _, user := range res.GetGetUserByIdsV2() {
		if user.BlockInfo.UserBlocked || user.BlockInfo.MatchingBlocked {
			if user.Id == uid {
				return nil, whalecode.ErrInviterBlocked
			} else {
				return nil, whalecode.ErrInviteeBlocked
			}
		}
	}

	db := dbutil.GetDB(ctx)
	invitation := models.MatchingInvitation{
		ID:               shortid.NewWithTime("mi_", 6),
		UserID:           uid,
		InviteeID:        param.InviteeID,
		Remark:           param.Remark,
		TopicID:          param.TopicID,
		CityID:           param.CityID,
		AreaIDs:          param.AreaIds,
		ConfirmState:     models.InvitationConfirmStateUnconfirmed.String(),
		MatchingIds:      []string{},
		MatchingResultId: 0,
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		MatchingInvitation := dbquery.Use(tx).MatchingInvitation
		err = MatchingInvitation.WithContext(ctx).Create(&invitation)
		if err != nil {
			return err
		}
		MatchingQuota := dbquery.Use(tx).MatchingQuota
		_, err = MatchingQuota.WithContext(ctx).Where(MatchingQuota.UserID.Eq(uid)).
			UpdateSimple(
				MatchingQuota.Remain.Add(-1),
				MatchingQuota.InvitationNum.Add(1),
			)
		if err != nil {
			return err
		}
		MatchingDurationConstraint := dbquery.Use(tx).MatchingDurationConstraint
		_, err = MatchingDurationConstraint.
			WithContext(ctx).
			Where(MatchingDurationConstraint.ID.Eq(constraint.ID)).
			UpdateSimple(MatchingDurationConstraint.Remain.Add(-1))
		return err
	})

	midacontext.GetLoader[loader.Loader](ctx).MatchingQuota.Clear(ctx, uid)
	midacontext.GetLoader[loader.Loader](ctx).MatchingInvitation.Clear(ctx, uid)

	if err != nil {
		return nil, err
	}

	_, err = scream.InvitationCreated(ctx, midacontext.GetServices(ctx).Scream, scream.InvitationCreatedParam{
		InvitationId: invitation.ID,
		InviterId:    uid,
		InviteeId:    param.InviteeID,
		TopicId:      param.TopicID,
		AreaIds:      []string{param.CityID},
	})
	if err != nil {
		fmt.Println("failed to create invitation notification:", err)
	}

	userThunk := midacontext.GetLoader[loader.Loader](ctx).UserAvatarNickname.Load(ctx, uid)
	profile, err := userThunk()
	if err != nil {
		fmt.Println("failed to get user nickname", err)
		return nil, nil
	}

	topicName := res.Topic.Name
	if topicName != "" {
		_, err = scream.SendUserNotification(ctx, midacontext.GetServices(ctx).Scream, scream.UserNotificationKindInvitationrecieved, invitation.InviteeID, map[string]interface{}{
			"topicName":    topicName,
			"userId":       uid,
			"userName":     profile.Nickname,
			"invitationId": invitation.ID,
		})
		if err != nil {
			fmt.Println("failed to send invitation notification:", err)
		}
	}
	return &invitation, nil
}

func FinishMatching(ctx context.Context, matchingID string, uid string) error {
	matchingThunk := midacontext.GetLoader[loader.Loader](ctx).Matching.Load(ctx, matchingID)
	matching, err := matchingThunk()
	if err != nil {
		return err
	}
	if uid != "" {
		if matching.UserID != uid {
			return midacode.ErrNotPermitted
		}
	}

	if matching.State != models.MatchingStateMatched.String() {
		return midacode.ErrStateMayHaveChanged
	}

	matchingResultThunk := midacontext.GetLoader[loader.Loader](ctx).MatchingResult.Load(ctx, matching.ResultID)
	matchingResult, err := matchingResultThunk()
	if err != nil {
		return err
	}

	_, err = smew.GroupMemberLeave(ctx, midacontext.GetServices(ctx).Smew, matchingResult.ChatGroupID, matching.UserID)
	if err != nil {
		return err
	}

	db := dbutil.GetDB(ctx)

	err = db.Transaction(func(tx *gorm.DB) error {
		Matching := dbquery.Use(tx).Matching
		MatchingQuota := dbquery.Use(tx).MatchingQuota
		MatchingResult := dbquery.Use(tx).MatchingResult

		ret, err := Matching.WithContext(ctx).
			Where(Matching.ID.Eq(matchingID)).
			Where(Matching.State.Eq(models.MatchingStateMatched.String())).
			UpdateSimple(
				Matching.State.Value(models.MatchingStateFinished.String()),
				Matching.InChatGroup.Value(false),
				Matching.FinishedAt.Value(time.Now()),
			)

		if err != nil {
			return err
		}

		if ret.RowsAffected != 1 {
			return midacode.ErrStateMayHaveChanged
		}

		_, err = MatchingQuota.
			WithContext(ctx).
			Where(MatchingQuota.UserID.Eq(matching.UserID)).
			UpdateSimple(MatchingQuota.Remain.Add(1))
		if err != nil {
			return err
		}

		_, err = MatchingResult.WithContext(ctx).
			Where(MatchingResult.ID.Eq(matching.ResultID)).
			Where(MatchingResult.ChatGroupState.Eq(models.ChatGroupStateCreated.String())).
			UpdateSimple(
				MatchingResult.ChatGroupState.Value(models.ChatGroupStateClosed.String()),
				MatchingResult.FinishedAt.Value(time.Now()),
			)
		return err
	})
	for _, matching := range matchingResult.MatchingIDs {
		midacontext.GetLoader[loader.Loader](ctx).Matching.Clear(ctx, matching)
	}
	midacontext.GetLoader[loader.Loader](ctx).MatchingResult.Clear(ctx, matching.ResultID)
	midacontext.GetLoader[loader.Loader](ctx).MatchingQuota.Clear(ctx, matching.UserID)
	PublishMatchingFinishedEvent(ctx, matching, matchingResult.CreatedBy)
	return err
}

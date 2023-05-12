package modelutil

import (
	"context"
	"fmt"
	"time"
	"whale/pkg/dbquery"
	"whale/pkg/gqlient/hoopoe"
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
	index := lo.IndexOf(matchingResult.MatchingIDs, matchingID)
	if index != -1 {
		if confirmed {
			matchingResult.ConfirmStates[index] = models.MatchingResultConfirmStateConfirmed.String()
		} else {
			matchingResult.ConfirmStates[index] = models.MatchingResultConfirmStateRejected.String()
			matchingResult.Closed = true
		}
		MatchingResult := dbquery.Use(db).MatchingResult
		_, err = MatchingResult.WithContext(ctx).
			Where(MatchingResult.ID.Eq(matchingResult.ID)).
			Select(MatchingResult.ConfirmStates, MatchingResult.Closed).
			Updates(matchingResult)
		if err != nil {
			return err
		}
		if matchingResult.Closed {
			// 如果有人拒绝匹配
			Matching := dbquery.Use(db).Matching
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
		}
	}
	return nil
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
	thunk := midacontext.GetLoader[loader.Loader](ctx).MatchingQuota.Load(ctx, uid)
	quota, err := thunk()
	if err != nil {
		return nil, err
	}
	if quota.Remain <= 0 {
		return nil, whalecode.ErrMatchingQuotaNotEnough
	}

	users, err := hoopoe.GetUserByIDs(ctx, midacontext.GetServices(ctx).Hoopoe, []string{uid})
	if err != nil {
		return nil, err
	}

	if users.GetUserByIds[0].Gender == "" {
		return nil, whalecode.ErrMatchingQuotaNotEnough
	}

	_, err = hoopoe.GetTopic(ctx, midacontext.GetServices(ctx).Hoopoe, param.TopicID)
	if err != nil {
		return nil, err
	}

	_, err = hoopoe.GetArea(ctx, midacontext.GetServices(ctx).Hoopoe, param.CityID)
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
		ID:             shortid.New("m_", 8),
		TopicID:        param.TopicID,
		UserID:         uid,
		State:          models.MatchingStateMatching.String(),
		Gender:         param.Gender.String(),
		Remark:         *param.Remark,
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
		_, err = MatchingQuota.WithContext(ctx).Where(MatchingQuota.UserID.Eq(uid)).UpdateSimple(MatchingQuota.Remain.Add(-1))
		return err
	})
	if err == nil {
		midacontext.GetLoader[loader.Loader](ctx).MatchingQuota.Clear(ctx, uid)
	}
	return matching, err
}

func CreateMatchingInvitation(ctx context.Context, uid string, param models.CreateMatchingInvitationParam) (*models.MatchingInvitation, error) {
	thunk := midacontext.GetLoader[loader.Loader](ctx).MatchingQuota.Load(ctx, uid)
	quota, err := thunk()
	if err != nil {
		return nil, err
	}
	if quota.Remain <= 0 {
		return nil, whalecode.ErrMatchingQuotaNotEnough
	}

	_, err = hoopoe.GetTopic(ctx, midacontext.GetServices(ctx).Hoopoe, param.TopicID)
	if err != nil {
		return nil, err
	}

	_, err = hoopoe.GetArea(ctx, midacontext.GetServices(ctx).Hoopoe, param.CityID)
	if err != nil {
		return nil, err
	}

	db := dbutil.GetDB(ctx)
	MatchingInvitation := dbquery.Use(db).MatchingInvitation

	invitation := models.MatchingInvitation{
		ID:               shortid.New("mi_", 8),
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
	err = MatchingInvitation.WithContext(ctx).Create(&invitation)
	if err != nil {
		return nil, err
	}
	MatchingQuota := dbquery.Use(db).MatchingQuota
	_, err = MatchingQuota.WithContext(ctx).Where(MatchingQuota.UserID.Eq(uid)).UpdateSimple(MatchingQuota.Remain.Add(-1))
	if err != nil {
		return nil, err
	}
	midacontext.GetLoader[loader.Loader](ctx).MatchingQuota.Clear(ctx, uid)
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
	Matching := dbquery.Use(db).Matching
	MatchingQuota := dbquery.Use(db).MatchingQuota
	MatchingResult := dbquery.Use(db).MatchingResult

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
	if err != nil {
		fmt.Println(err)
	}

	for _, matching := range matchingResult.MatchingIDs {
		midacontext.GetLoader[loader.Loader](ctx).Matching.Clear(ctx, matching)
	}
	midacontext.GetLoader[loader.Loader](ctx).MatchingResult.Clear(ctx, matching.ResultID)
	midacontext.GetLoader[loader.Loader](ctx).MatchingQuota.Clear(ctx, matching.UserID)
	return nil
}

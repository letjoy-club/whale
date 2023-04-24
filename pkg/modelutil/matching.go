package modelutil

import (
	"context"
	"whale/pkg/dbquery"
	"whale/pkg/loader"
	"whale/pkg/models"
	"whale/pkg/whalecode"

	"github.com/letjoy-club/mida-tool/midacontext"
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
		return nil, err
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
			Matching := dbquery.Use(db).Matching
			_, err = Matching.WithContext(ctx).
				Where(Matching.ResultID.Eq(matchingResult.ID)).
				UpdateSimple(
					Matching.State.Value(string(models.MatchingStateMatching)),
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

	var mut struct {
		CreateGroup string `graphql:"createGroup(resultId: $resultId, topicId: $topicId, memberIds: $memberIds)"`
	}
	if err := midacontext.GetServices(ctx).Smew.Mutate(ctx, &mut, map[string]interface{}{
		"resultId":  m.ID,
		"topicId":   m.TopicID,
		"memberIds": m.UserIDs,
	}); err != nil {
		return err
	}

	db := midacontext.GetDB(ctx)
	err := db.Transaction(func(db *gorm.DB) error {
		Matching := dbquery.Use(db).Matching
		MatchingResult := dbquery.Use(db).MatchingResult
		_, err := MatchingResult.WithContext(ctx).
			Where(MatchingResult.ID).
			UpdateSimple(
				MatchingResult.ChatGroupState.Value(models.ChatGroupStateCreated.String()),
				MatchingResult.ChatGroupID.Value(mut.CreateGroup),
			)
		if err != nil {
			return err
		}

		_, err = Matching.WithContext(ctx).
			Where(Matching.ResultID.Eq(m.ID)).
			UpdateSimple(
				Matching.InChatGroup.Value(true),
			)
		return err
	})
	return err
}

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
	index := lo.IndexOf(matchingResult.MatchingIDs, matchingID)
	if index != -1 {
		if confirmed {
			matchingResult.ConfirmStates[index] = models.MatchingResultConfirmStateConfirmed.String()
		} else {
			matchingResult.ConfirmStates[index] = models.MatchingResultConfirmStateRejected.String()
		}
		MatchingResult := dbquery.Use(db).MatchingResult
		_, err = MatchingResult.WithContext(ctx).Where(MatchingResult.ID.Eq(matchingResult.ID)).Select(MatchingResult.ConfirmStates).Updates(matchingResult)
		if err != nil {
			return err
		}
	}
	return nil
}

func PublishMatching() error {
	return nil
}

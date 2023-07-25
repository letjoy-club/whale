package modelutil

import (
	"context"
	"time"
	"whale/pkg/dbquery"
	"whale/pkg/keyer"
	"whale/pkg/loader"
	"whale/pkg/models"
	"whale/pkg/whalecode"

	"github.com/letjoy-club/mida-tool/dbutil"
	"github.com/letjoy-club/mida-tool/midacode"
	"github.com/letjoy-club/mida-tool/midacontext"
	"github.com/letjoy-club/mida-tool/redisutil"
	"go.uber.org/multierr"
	"gorm.io/gorm"
)

func SendOutOffer(ctx context.Context, myUserID, myMatchingID, targetMatchingID string) error {
	matchingThunk := midacontext.GetLoader[loader.Loader](ctx).Matching.LoadMany(ctx, []string{myMatchingID, targetMatchingID})
	matchings, errors := matchingThunk()
	if errors != nil {
		return multierr.Combine(errors...)
	}
	myMatching := matchings[0]
	targetMatching := matchings[1]

	if myMatching.UserID == targetMatching.UserID {
		return whalecode.ErrCannotSendMatchingOfferToSelf
	}

	if myMatching.TopicID != targetMatching.TopicID {
		return whalecode.ErrCannotSendMatchingOfferToDifferentTopic
	}

	if myUserID != "" {
		if myMatching.UserID != myUserID {
			return midacode.ErrNotPermitted
		}
	}

	if myMatching.State != models.MatchingStateMatching.String() {
		return whalecode.ErrYourMatchingNotInMatchingState
	}

	if targetMatching.State != models.MatchingStateMatching.String() {
		// 可能已经匹配到或者已经取消
		return whalecode.ErrTheMatchingIsNotInMatchingState
	}

	release, err := redisutil.LockAll(ctx, keyer.UserMatching(myMatching.UserID), keyer.UserMatching(targetMatching.UserID))
	if err != nil {
		return err
	}
	defer release(ctx)

	mySummary, err := GetMatchingOfferSummary(ctx, myMatching)
	if err != nil {
		return err
	}

	targetSummary, err := GetMatchingOfferSummary(ctx, targetMatching)
	if err != nil {
		return err
	}

	if mySummary.Active == false || targetSummary.Active == false {
		return whalecode.ErrMatchingOfferIsClosed
	}

	db := dbutil.GetDB(ctx)
	err = dbquery.Use(db).Transaction(func(tx *dbquery.Query) error {
		MatchingOfferSummary := tx.MatchingOfferSummary
		MatchingOfferRecord := tx.MatchingOfferRecord

		err := MatchingOfferRecord.WithContext(ctx).Create(&models.MatchingOfferRecord{
			MatchingID:   myMatchingID,
			UserID:       myMatching.UserID,
			ToMatchingID: targetMatchingID,
			State:        models.MatchingOfferStateUnprocessed.String(),
		})
		if err != nil {
			return whalecode.ErrAlreadySentOutMatchingOffer
		}

		mySummary.RemainQuota -= 1
		mySummary.OutOfferNum++
		targetSummary.InOfferNum++

		if err := MatchingOfferSummary.WithContext(ctx).Save(mySummary, targetSummary); err != nil {
			return err
		}
		return nil
	})

	summaryLoader := midacontext.GetLoader[loader.Loader](ctx).MatchingOfferSummary
	if err == nil {
		summaryLoader.Prime(ctx, myMatchingID, mySummary)
		summaryLoader.Prime(ctx, targetMatchingID, targetSummary)
	} else {
		summaryLoader.Clear(ctx, myMatchingID)
		summaryLoader.Clear(ctx, targetMatchingID)
	}
	return err
}

func CancelOutOffer(ctx context.Context, userID, myMatchingID, targetMatchingID string) error {
	matchingThunk := midacontext.GetLoader[loader.Loader](ctx).Matching.LoadMany(ctx, []string{myMatchingID, targetMatchingID})
	matchings, errors := matchingThunk()
	if errors != nil {
		return multierr.Combine(errors...)
	}
	myMatching := matchings[0]
	targetMatching := matchings[1]

	if userID != "" {
		if myMatching.UserID != userID {
			return midacode.ErrNotPermitted
		}
	}

	release, err := redisutil.LockAll(ctx, keyer.UserMatching(myMatching.UserID), keyer.UserMatching(targetMatching.UserID))
	if err != nil {
		return err
	}
	defer release(ctx)

	myMatchingOfferSummary, err := GetMatchingOfferSummary(ctx, myMatching)
	if err != nil {
		return err
	}

	targetMatchingOfferSummary, err := GetMatchingOfferSummary(ctx, targetMatching)
	if err != nil {
		return err
	}

	db := dbutil.GetDB(ctx)
	err = db.Transaction(func(tx *gorm.DB) error {
		MatchingOfferSummary := dbquery.Use(tx).MatchingOfferSummary
		MatchingOfferRecord := dbquery.Use(tx).MatchingOfferRecord

		record, err := MatchingOfferRecord.WithContext(ctx).Where(MatchingOfferRecord.MatchingID.Eq(myMatchingID), MatchingOfferRecord.ToMatchingID.Eq(targetMatchingID)).Take()
		if err != nil {
			return whalecode.ErrMatchingOfferNotFound
		}
		if record.State != models.MatchingOfferStateUnprocessed.String() {
			return whalecode.ErrMatchingOfferIsCancelableOnlyWhenUnprocessed
		}

		myMatchingOfferSummary.OutOfferNum--
		targetMatchingOfferSummary.InOfferNum--

		if err := MatchingOfferSummary.WithContext(ctx).Save(myMatchingOfferSummary, targetMatchingOfferSummary); err != nil {
			return err
		}

		if _, err := MatchingOfferRecord.WithContext(ctx).Where(MatchingOfferRecord.MatchingID.Eq(myMatchingID), MatchingOfferRecord.ToMatchingID.Eq(targetMatchingID)).UpdateSimple(MatchingOfferRecord.State.Value(models.MatchingOfferStateCanceled.String())); err != nil {
			return err
		}

		// 回复我的匹配配额
		_, err = MatchingOfferSummary.WithContext(ctx).Where(MatchingOfferSummary.MatchingID.Eq(myMatchingID)).UpdateSimple(MatchingOfferSummary.RemainQuota.Add(1))
		if err != nil {
			return err
		}

		return nil
	})

	summaryLoader := midacontext.GetLoader[loader.Loader](ctx).MatchingOfferSummary
	if err == nil {
		summaryLoader.Prime(ctx, myMatchingID, myMatchingOfferSummary)
		summaryLoader.Prime(ctx, targetMatchingID, targetMatchingOfferSummary)

		inOfferLoader := midacontext.GetLoader[loader.Loader](ctx).InMatchingOfferRecord
		inOfferLoader.Clear(ctx, targetMatchingID)
		outOfferLoader := midacontext.GetLoader[loader.Loader](ctx).OutMatchingOfferRecord
		outOfferLoader.Clear(ctx, myMatchingID)
	} else {
		summaryLoader.Clear(ctx, myMatchingID)
		summaryLoader.Clear(ctx, targetMatchingID)
	}
	return nil
}

func AcceptInOffer(ctx context.Context, userID, myMatchingID, targetMatchingID string) error {
	matchingThunk := midacontext.GetLoader[loader.Loader](ctx).Matching.LoadMany(ctx, []string{myMatchingID, targetMatchingID})
	matchings, errors := matchingThunk()
	if errors != nil {
		return multierr.Combine(errors...)
	}

	myMatching := matchings[0]
	targetMatching := matchings[1]

	if userID != "" {
		if myMatching.UserID != userID {
			return midacode.ErrNotPermitted
		}
	}

	if myMatching.State != models.MatchingStateMatching.String() {
		return whalecode.ErrYourMatchingNotInMatchingState
	}

	if targetMatching.State != models.MatchingStateMatching.String() {
		// 可能已经匹配到或者已经取消
		return whalecode.ErrTheMatchingIsNotInMatchingState
	}

	release, err := redisutil.LockAll(ctx, keyer.UserMatching(myMatching.UserID), keyer.UserMatching(targetMatching.UserID))
	if err != nil {
		return err
	}
	defer release(ctx)

	myMatchingOfferSummary, err := GetMatchingOfferSummary(ctx, myMatching)
	if err != nil {
		return err
	}

	targetMatchingOfferSummary, err := GetMatchingOfferSummary(ctx, targetMatching)
	if err != nil {
		return err
	}

	var myInOrderIDs []string
	var myOutOrderIDs []string

	db := dbutil.GetDB(ctx)
	err = dbquery.Use(db).Transaction(func(tx *dbquery.Query) error {
		MatchingOfferSummary := tx.MatchingOfferSummary
		MatchingOfferRecord := tx.MatchingOfferRecord

		myMatchingOfferSummary.Active = false
		targetMatchingOfferSummary.Active = false

		now := time.Now()

		// 更新邀约状态
		rx, err := MatchingOfferRecord.WithContext(ctx).Where(
			MatchingOfferRecord.MatchingID.Eq(targetMatchingID),
			MatchingOfferRecord.ToMatchingID.Eq(myMatchingID),
			MatchingOfferRecord.State.Eq(models.MatchingOfferStateUnprocessed.String()),
		).UpdateSimple(
			MatchingOfferRecord.State.Value(models.MatchingOfferStateAccepted.String()),
			MatchingOfferRecord.ReactedAt.Value(now),
		)
		if err != nil {
			return err
		}
		if rx.RowsAffected == 0 {
			return whalecode.ErrMatchingOfferNotFound
		}

		if err := CloseMatchingOfferRelatedToMatchingID(ctx, tx, myMatchingID); err != nil {
			return err
		}
		rx, err = MatchingOfferSummary.WithContext(ctx).Where(
			MatchingOfferSummary.MatchingID.In(myMatchingID, targetMatchingID),
		).UpdateSimple(MatchingOfferSummary.Active.Value(false))
		if err != nil {
			return err
		}
		if rx.RowsAffected != 2 {
			return midacode.ErrStateMayHaveChanged
		}

		return nil
	})

	summaryLoader := midacontext.GetLoader[loader.Loader](ctx).MatchingOfferSummary
	outOfferLoader := midacontext.GetLoader[loader.Loader](ctx).OutMatchingOfferRecord
	inOfferLoader := midacontext.GetLoader[loader.Loader](ctx).OutMatchingOfferRecord
	if err == nil {
		for _, id := range myInOrderIDs {
			inOfferLoader.Clear(ctx, id)
			summaryLoader.Clear(ctx, id)
		}
		for _, id := range myOutOrderIDs {
			outOfferLoader.Clear(ctx, id)
		}
		summaryLoader.Prime(ctx, myMatchingID, myMatchingOfferSummary)
		summaryLoader.Prime(ctx, targetMatchingID, targetMatchingOfferSummary)
	} else {
		summaryLoader.Clear(ctx, myMatchingID)
		summaryLoader.Clear(ctx, targetMatchingID)
	}
	return err
}

func RejectInOffer(ctx context.Context, userID, myMatchingID, targetMatchingID string) error {
	db := dbutil.GetDB(ctx)

	matchingThunk := midacontext.GetLoader[loader.Loader](ctx).Matching.LoadMany(ctx, []string{myMatchingID, targetMatchingID})
	matchings, errors := matchingThunk()
	if errors != nil {
		return multierr.Combine(errors...)
	}

	myMatching := matchings[0]
	targetMatching := matchings[1]

	release, err := redisutil.LockAll(ctx, keyer.UserMatching(myMatching.UserID), keyer.UserMatching(targetMatching.UserID))
	if err != nil {
		return err
	}
	defer release(ctx)

	myMatchingOfferSummary, err := GetMatchingOfferSummary(ctx, myMatching)
	if err != nil {
		return err
	}

	if !myMatchingOfferSummary.Active {
		return whalecode.ErrMatchingOfferIsNotActive
	}

	err = dbquery.Use(db).Transaction(func(tx *dbquery.Query) error {
		// 更新当前的邀约状态
		MatchingOfferRecord := tx.MatchingOfferRecord
		_, err = MatchingOfferRecord.WithContext(ctx).
			Where(
				MatchingOfferRecord.ToMatchingID.Eq(myMatchingID),
				MatchingOfferRecord.MatchingID.Eq(targetMatchingID),
			).
			UpdateSimple(MatchingOfferRecord.State.Value(models.MatchingOfferStateRejected.String()))
		if err != nil {
			return err
		}

		// 给被拒绝的邀约的配额回复配额
		MatchingOfferSummary := tx.MatchingOfferSummary
		_, err = MatchingOfferSummary.WithContext(ctx).Where(
			MatchingOfferSummary.MatchingID.Eq(targetMatchingID),
		).UpdateSimple(
			MatchingOfferSummary.RemainQuota.Add(1),
		)
		return err
	})
	if err != nil {
		return err
	}

	inOfferLoader := midacontext.GetLoader[loader.Loader](ctx).InMatchingOfferRecord
	inOfferLoader.Clear(ctx, myMatchingID)
	outOfferLoader := midacontext.GetLoader[loader.Loader](ctx).OutMatchingOfferRecord
	outOfferLoader.Clear(ctx, targetMatchingID)
	return nil
}

func GetMatchingOfferSummary(ctx context.Context, matching *models.Matching) (*models.MatchingOfferSummary, error) {
	db := dbutil.GetDB(ctx)
	MatchingOfferSummary := dbquery.Use(db).MatchingOfferSummary
	ret, err := MatchingOfferSummary.WithContext(ctx).Where(MatchingOfferSummary.MatchingID.Eq(matching.ID)).Find()
	if err != nil {
		return nil, err
	}
	if len(ret) == 1 {
		return ret[0], nil
	}

	summary := &models.MatchingOfferSummary{
		MatchingID: matching.ID,
		UserID:     matching.UserID,
		Active:     matching.State == models.MatchingStateMatching.String(),
	}

	// 不存在则创建
	if err := MatchingOfferSummary.WithContext(ctx).Create(summary); err != nil {
		return nil, err
	}
	return summary, nil
}

func DeactiveMatchingSummaryAndUpdateMatchingOffer(ctx context.Context, tx *dbquery.Query, matchingID string) error {
	MatchingOfferSummary := tx.MatchingOfferSummary

	// 关闭邀约
	rx, err := MatchingOfferSummary.WithContext(ctx).Where(
		MatchingOfferSummary.MatchingID.Eq(matchingID),
		MatchingOfferSummary.Active.Is(true),
	).UpdateSimple(MatchingOfferSummary.Active.Value(false))
	if err != nil {
		return err
	}
	if rx.RowsAffected != 1 {
		return whalecode.ErrMatchingOfferIsNotActive
	}
	return CloseMatchingOfferRelatedToMatchingID(ctx, tx, matchingID)
}

// CloseMatchingOfferRelatedToMatchingID 关闭与 matchingID 相关的所有邀约，并恢复配额
func CloseMatchingOfferRelatedToMatchingID(ctx context.Context, tx *dbquery.Query, matchingID string) error {
	MatchingOfferSummary := tx.MatchingOfferSummary
	MatchingOfferRecord := tx.MatchingOfferRecord

	// 关闭 matching id 中已发起的所有邀约
	_, err := MatchingOfferRecord.WithContext(ctx).Where(
		MatchingOfferRecord.MatchingID.Eq(matchingID),
		MatchingOfferRecord.State.Eq(models.MatchingOfferStateUnprocessed.String()),
	).UpdateSimple(MatchingOfferRecord.State.Value(models.MatchingOfferStateRejected.String()))
	if err != nil {
		return err
	}

	// 关闭 matching id 中已收到的所有邀约，并回复给相应匹配配额
	recievedOfferIDs := []string{}
	MatchingOfferRecord.WithContext(ctx).Where(
		MatchingOfferRecord.ToMatchingID.Eq(matchingID),
		MatchingOfferRecord.State.Eq(models.MatchingOfferStateUnprocessed.String()),
	).Pluck(MatchingOfferRecord.MatchingID, &recievedOfferIDs)

	// 回复配额
	if len(recievedOfferIDs) > 0 {
		rx, err := MatchingOfferSummary.WithContext(ctx).Where(
			MatchingOfferSummary.MatchingID.In(recievedOfferIDs...),
		).
			UpdateSimple(MatchingOfferSummary.RemainQuota.Add(1))
		if err != nil {
			return err
		}
		if rx.RowsAffected != int64(len(recievedOfferIDs)) {
			return midacode.ErrStateMayHaveChanged
		}
	}

	// 关闭 matching id 中已收到的所有邀约
	_, err = MatchingOfferRecord.WithContext(ctx).Where(
		MatchingOfferRecord.ToMatchingID.Eq(matchingID),
		MatchingOfferRecord.State.Eq(models.MatchingOfferStateUnprocessed.String()),
	).UpdateSimple(
		MatchingOfferRecord.State.Value(models.MatchingOfferStateRejected.String()),
		MatchingOfferRecord.ReactedAt.Value(time.Now()),
	)
	if err != nil {
		return err
	}
	return nil
}

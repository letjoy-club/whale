package modelutil

import (
	"context"
	"time"
	"whale/pkg/dbquery"
	"whale/pkg/gqlient/smew"
	"whale/pkg/loader"
	"whale/pkg/models"
	"whale/pkg/utils"
	"whale/pkg/whalecode"

	"github.com/letjoy-club/mida-tool/dbutil"
	"github.com/letjoy-club/mida-tool/keyer"
	"github.com/letjoy-club/mida-tool/midacode"
	"github.com/letjoy-club/mida-tool/midacontext"
	"github.com/letjoy-club/mida-tool/redisutil"
	"gorm.io/gen/field"
)

func CreateMotionOffer(ctx context.Context, myUserID, myMotionID, targetMotionID string) error {
	motionsThunk := midacontext.GetLoader[loader.Loader](ctx).Motion.LoadMany(ctx, []string{myMotionID, targetMotionID})
	motions, err := utils.ReturnThunk(motionsThunk)
	if err != nil {
		return err
	}
	myMotion := motions[0]
	targetMotion := motions[1]

	if myMotion.UserID == myUserID {
		return whalecode.ErrCannotSendMatchingOfferToSelf
	}

	if myMotion.TopicID != targetMotion.TopicID {
		return whalecode.ErrCannotSendMatchingOfferToDifferentTopic
	}

	if myUserID != "" {
		if myMotion.UserID != myUserID {
			return midacode.ErrNotPermitted
		}
	}

	if !myMotion.Active || !myMotion.Discoverable {
		return whalecode.ErrYourMotionIsNotActive
	}

	if !targetMotion.Active || !targetMotion.Discoverable {
		return whalecode.ErrTheMotionIsNotActive
	}

	release, err := redisutil.LockAll(ctx, keyer.UserMotion(myMotion.UserID), keyer.UserMotion(targetMotion.UserID))
	if err != nil {
		return err
	}
	defer release(ctx)

	db := dbutil.GetDB(ctx).Debug()
	err = dbquery.Use(db).Transaction(func(tx *dbquery.Query) error {
		MatchingResult := tx.MatchingResult
		matchingResult := &models.MatchingResult{
			UserIDs:        []string{myMotion.UserID, targetMotion.UserID},
			TopicID:        myMotion.TopicID,
			MotionIDs:      []string{myMotionID, targetMotionID},
			CreatedBy:      string(models.ResultCreatedByOffer),
			ConfirmStates:  []string{string(models.MatchingResultConfirmStateUnconfirmed), string(models.MatchingResultConfirmStateUnconfirmed)},
			ChatGroupState: models.ChatGroupStateUncreated.String(),
		}
		if err := MatchingResult.WithContext(ctx).Create(matchingResult); err != nil {
			return err
		}
		MotionOfferRecord := tx.MotionOfferRecord
		record := &models.MotionOfferRecord{
			UserID:     myMotion.UserID,
			MotionID:   myMotion.ID,
			ToUserID:   targetMotion.UserID,
			ToMotionID: targetMotion.ID,
			ExpiredAt:  time.Now().Add(time.Hour * 24),
			State:      string(models.MotionOfferStatePending),
		}
		err := MotionOfferRecord.WithContext(ctx).Create(record)
		if err != nil {
			return err
		}

		Motion := tx.Motion
		rx, err := Motion.WithContext(ctx).Where(Motion.ID.Eq(myMotion.ID)).UpdateSimple(
			Motion.OutOfferNum.Add(1),
			Motion.PendingOutNum.Add(1),
		)
		if err != nil {
			return err
		}
		if rx.RowsAffected != 1 {
			return midacode.ErrUnknownError
		}

		rx, err = Motion.WithContext(ctx).Where(Motion.ID.Eq(targetMotion.ID)).UpdateSimple(Motion.InOfferNum.Add(1), Motion.PendingInNum.Add(1))
		if err != nil {
			return err
		}
		if rx.RowsAffected != 1 {
			return midacode.ErrUnknownError
		}

		resp, err := smew.CreateMotionGroup(ctx, midacontext.GetServices(ctx).Smew, smew.CreateMotionGroupParam{
			ToUserId:     targetMotion.UserID,
			FromUserId:   myMotion.UserID,
			ToMotionId:   targetMotion.ID,
			FromMotionId: myMotion.ID,
			TopicId:      targetMotion.TopicID,
			ResultId:     matchingResult.ID,
		})
		if err != nil {
			return err
		}

		rx, err = MotionOfferRecord.WithContext(ctx).Where(MotionOfferRecord.ID.Eq(record.ID)).UpdateSimple(MotionOfferRecord.ChatGroupID.Value(resp.CreateMotionGroup))
		if err != nil {
			return err
		}
		if rx.RowsAffected != 1 {
			return midacode.ErrUnknownError
		}

		return nil
	})

	if err == nil {
		midacontext.GetLoader[loader.Loader](ctx).Motion.Clear(ctx, myMotionID)
		midacontext.GetLoader[loader.Loader](ctx).Motion.Clear(ctx, targetMotionID)
		midacontext.GetLoader[loader.Loader](ctx).InMotionOfferRecord.Clear(ctx, myMotionID)
		midacontext.GetLoader[loader.Loader](ctx).OutMotionOfferRecord.Clear(ctx, targetMotionID)
	}
	return err
}

func AcceptMotionOffer(ctx context.Context, myUserID, myMotionID, targetMotionID string) error {
	motionsThunk := midacontext.GetLoader[loader.Loader](ctx).Motion.LoadMany(ctx, []string{myMotionID, targetMotionID})
	motions, err := utils.ReturnThunk(motionsThunk)
	if err != nil {
		return err
	}

	myMotion := motions[0]
	targetMotion := motions[1]

	if myUserID != "" {
		if myMotion.UserID != myUserID {
			return midacode.ErrNotPermitted
		}
	}

	if !myMotion.Active || !myMotion.Discoverable {
		return whalecode.ErrYourMotionIsNotActive
	}

	if !targetMotion.Active || !targetMotion.Discoverable {
		return whalecode.ErrTheMotionIsNotActive
	}

	db := dbutil.GetDB(ctx)
	now := time.Now()
	err = dbquery.Use(db).Transaction(func(tx *dbquery.Query) error {
		MotionOfferRecord := tx.MotionOfferRecord
		record, err := MotionOfferRecord.WithContext(ctx).Where(
			MotionOfferRecord.MotionID.Eq(myMotion.ID),
			MotionOfferRecord.ToMotionID.Eq(targetMotion.ID),
		).Take()
		if err != nil {
			return err
		}

		if _, err := smew.CreateTimGroup(ctx, midacontext.GetServices(ctx).Smew, record.ChatGroupID); err != nil {
			return err
		}

		rx, err := MotionOfferRecord.WithContext(ctx).Where(
			MotionOfferRecord.MotionID.Eq(myMotion.ID),
			MotionOfferRecord.ToMotionID.Eq(targetMotion.ID),
			MotionOfferRecord.State.Eq(string(models.MotionOfferStatePending)),
		).UpdateSimple(
			MotionOfferRecord.State.Value(string(models.MotionOfferStateAccepted)),
			MotionOfferRecord.ReactAt.Value(now),
		)
		if err != nil {
			return err
		}
		if rx.RowsAffected != 1 {
			return midacode.ErrStateMayHaveChanged
		}
		// 修改我的 motion
		Motion := tx.Motion
		fields := []field.AssignExpr{Motion.PendingInNum.Add(-1), Motion.ActiveNum.Add(1)}
		if !myMotion.Discoverable {
			// 没关闭，但是不可见
			if myMotion.PendingInNum+myMotion.PendingOutNum == 1 {
				// 设为关闭
				fields = append(fields, Motion.Active.Value(false))
			}
		}

		rx, err = Motion.WithContext(ctx).Where(Motion.ID.Eq(myMotion.ID)).UpdateSimple(fields...)
		if err != nil {
			return err
		}
		if rx.RowsAffected != 1 {
			return midacode.ErrStateMayHaveChanged
		}
		// 修改对方的 motion
		fields = []field.AssignExpr{Motion.PendingOutNum.Add(-1), Motion.ActiveNum.Add(1)}
		if !targetMotion.Discoverable {
			// 没关闭，但是不可见
			if targetMotion.PendingInNum+targetMotion.PendingOutNum == 1 {
				// 设为关闭
				fields = append(fields, Motion.Active.Value(false))
			}
		}
		rx, err = Motion.WithContext(ctx).Where(Motion.ID.Eq(targetMotion.ID)).UpdateSimple(fields...)
		if err != nil {
			return err
		}
		if rx.RowsAffected != 1 {
			return midacode.ErrStateMayHaveChanged
		}
		return nil
	})

	if err != nil {
		return err
	}
	midacontext.GetLoader[loader.Loader](ctx).Motion.Clear(ctx, myMotionID)
	midacontext.GetLoader[loader.Loader](ctx).Motion.Clear(ctx, targetMotionID)
	midacontext.GetLoader[loader.Loader](ctx).InMotionOfferRecord.Clear(ctx, myMotionID)
	midacontext.GetLoader[loader.Loader](ctx).OutMotionOfferRecord.Clear(ctx, targetMotionID)
	return nil
}

func RejectMotionOffer(ctx context.Context, userID, myMotionID, targetMotionID string) error {
	motionThunk := midacontext.GetLoader[loader.Loader](ctx).Motion.LoadMany(ctx, []string{myMotionID, targetMotionID})
	motions, err := utils.ReturnThunk(motionThunk)
	if err != nil {
		return err
	}

	myMotion := motions[0]
	targetMotion := motions[1]

	if userID != "" {
		if myMotion.UserID != userID {
			return midacode.ErrNotPermitted
		}
	}

	motionOfferThunk := midacontext.GetLoader[loader.Loader](ctx).InMotionOfferRecord.Load(ctx, myMotionID)
	motionOffer, err := motionOfferThunk()
	if err != nil {
		return err
	}
	found := false
	for _, offer := range motionOffer.Offers {
		if offer.MotionID != targetMotionID {
			if offer.State != string(models.MotionOfferStatePending) {
				return whalecode.MotionOfferIsNotPending
			}
			found = true
			break
		}
	}
	if !found {
		return whalecode.ErrMatchingOfferNotFound
	}

	release, err := redisutil.LockAll(ctx, keyer.UserMotion(myMotion.UserID), keyer.UserMotion(targetMotion.UserID))
	if err != nil {
		return err
	}
	defer release(ctx)

	db := dbutil.GetDB(ctx)
	err = dbquery.Use(db).Transaction(func(tx *dbquery.Query) error {
		MotionOfferRecord := tx.MotionOfferRecord
		motionOffer, err := MotionOfferRecord.WithContext(ctx).Where(MotionOfferRecord.MotionID.Eq(targetMotionID), MotionOfferRecord.ToMotionID.Eq(myMotionID)).Take()
		if err != nil {
			return err
		}

		if motionOffer.ChatGroupID != "" {
			_, err := smew.DestroyGroup(ctx, midacontext.GetServices(ctx).Smew, motionOffer.ChatGroupID)
			if err != nil {
				return err
			}
		}

		rx, err := MotionOfferRecord.WithContext(ctx).Where(
			MotionOfferRecord.MotionID.Eq(targetMotionID),
			MotionOfferRecord.ToMotionID.Eq(myMotionID),
			MotionOfferRecord.State.Eq(string(models.MotionOfferStatePending)),
		).UpdateSimple(
			MotionOfferRecord.State.Value(string(models.MotionOfferStateRejected)),
			MotionOfferRecord.ReactAt.Value(time.Now()),
		)
		if err != nil {
			return err
		}
		if rx.RowsAffected != 1 {
			return midacode.ErrStateMayHaveChanged
		}

		Motion := tx.Motion
		fields := []field.AssignExpr{Motion.PendingInNum.Add(-1)}
		if !targetMotion.Discoverable {
			// 没关闭，但是不可见
			if targetMotion.PendingInNum+targetMotion.PendingOutNum == 1 {
				// 设为关闭
				fields = append(fields, Motion.Active.Value(false))
			}
		}
		// 修改对方的 motion
		rx, err = Motion.WithContext(ctx).Where(Motion.ID.Eq(targetMotionID)).UpdateSimple(fields...)
		if err != nil {
			return err
		}
		if rx.RowsAffected != 1 {
			return midacode.ErrStateMayHaveChanged
		}

		// 修改我的 motion
		fields = []field.AssignExpr{Motion.PendingOutNum.Add(-1)}
		if !myMotion.Discoverable {
			// 没关闭，但是不可见
			if myMotion.PendingInNum+myMotion.PendingOutNum == 1 {
				// 设为关闭
				fields = append(fields, Motion.Active.Value(false))
			}
		}
		rx, err = Motion.WithContext(ctx).Where(Motion.ID.Eq(myMotionID)).UpdateSimple(fields...)
		if err != nil {
			return err
		}
		if rx.RowsAffected != 1 {
			return midacode.ErrStateMayHaveChanged
		}

		return nil
	})

	if err != nil {
		return err
	}

	loader := midacontext.GetLoader[loader.Loader](ctx)
	loader.Motion.Clear(ctx, myMotionID)
	loader.Motion.Clear(ctx, targetMotionID)
	loader.InMotionOfferRecord.Clear(ctx, myMotionID)
	loader.OutMotionOfferRecord.Clear(ctx, targetMotionID)
	return nil
}

func CancelMotionOffer(ctx context.Context, myUserID, myMotionID, targetMotionID string) error {
	motionsThunk := midacontext.GetLoader[loader.Loader](ctx).Motion.LoadMany(ctx, []string{myMotionID, targetMotionID})
	motions, err := utils.ReturnThunk(motionsThunk)
	if err != nil {
		return err
	}

	myMotion := motions[0]
	targetMotion := motions[1]

	if myUserID != "" {
		if myMotion.UserID != myUserID {
			return midacode.ErrNotPermitted
		}
	}

	motionOfferThunk := midacontext.GetLoader[loader.Loader](ctx).OutMotionOfferRecord.Load(ctx, myMotionID)
	motionOffer, err := motionOfferThunk()
	if err != nil {
		return err
	}
	found := false
	for _, offer := range motionOffer.Offers {
		if offer.ToMotionID != targetMotionID {
			continue
		}
		if offer.State != string(models.MotionOfferStatePending) {
			return whalecode.MotionOfferIsNotPending
		}
		found = true
	}
	if !found {
		return whalecode.ErrMatchingOfferNotFound
	}

	release, err := redisutil.LockAll(ctx, keyer.UserMotion(myMotion.UserID), keyer.UserMotion(targetMotion.UserID))
	if err != nil {
		return err
	}
	defer release(ctx)

	db := dbutil.GetDB(ctx)
	err = dbquery.Use(db).Transaction(func(tx *dbquery.Query) error {
		MotionOfferRecord := tx.MotionOfferRecord
		rx, err := MotionOfferRecord.WithContext(ctx).Where(
			MotionOfferRecord.MotionID.Eq(myMotionID),
			MotionOfferRecord.ToMotionID.Eq(targetMotionID),
			MotionOfferRecord.State.Eq(string(models.MotionOfferStatePending)),
		).UpdateSimple(
			MotionOfferRecord.State.Value(string(models.MotionOfferStateCanceled)),
		)
		if err != nil {
			return err
		}
		if rx.RowsAffected != 1 {
			return midacode.ErrStateMayHaveChanged
		}

		Motion := tx.Motion
		// 修改对方的 motion
		rx, err = Motion.WithContext(ctx).Where(Motion.ID.Eq(targetMotionID)).UpdateSimple(
			Motion.PendingInNum.Add(-1),
		)
		if err != nil {
			return err
		}
		if rx.RowsAffected != 1 {
			return midacode.ErrStateMayHaveChanged
		}

		// 修改我的 motion
		fields := []field.AssignExpr{}
		fields = append(fields, Motion.PendingOutNum.Add(-1))

		if myMotion.Active && !myMotion.Discoverable {
			// 如果广场不可见，但是没被关闭，需要关闭
			if myMotion.PendingInNum+myMotion.PendingOutNum == 1 {
				fields = append(fields, Motion.Active.Value(false))
			}
		}

		rx, err = Motion.WithContext(ctx).Where(Motion.ID.Eq(myMotionID)).UpdateSimple(fields...)
		if err != nil {
			return err
		}
		if rx.RowsAffected != 1 {
			return midacode.ErrStateMayHaveChanged
		}
		return nil
	})

	if err != nil {
		return err
	}

	loader := midacontext.GetLoader[loader.Loader](ctx)
	loader.Motion.Clear(ctx, myMotionID)
	loader.Motion.Clear(ctx, targetMotionID)
	loader.InMotionOfferRecord.Clear(ctx, myMotionID)
	loader.OutMotionOfferRecord.Clear(ctx, targetMotionID)
	return nil
}

func MarkMotionDiscoverable(ctx context.Context, userID, motionID string, discoverable bool) error {
	motionThunk := midacontext.GetLoader[loader.Loader](ctx).Motion.Load(ctx, motionID)
	motion, err := motionThunk()
	if err != nil {
		return err
	}
	if userID != "" {
		if motion.UserID != userID {
			return midacode.ErrNotPermitted
		}
	}
	if !motion.Active {
		return whalecode.ErrMatchingOfferIsNotActive
	}
	if motion.Discoverable == discoverable {
		return nil
	}
	db := dbutil.GetDB(ctx)
	Motion := dbquery.Use(db).Motion
	_, err = Motion.WithContext(ctx).Where(Motion.ID.Eq(motionID)).UpdateSimple(
		Motion.Discoverable.Value(discoverable),
	)
	if err != nil {
		return err
	}
	midacontext.GetLoader[loader.Loader](ctx).Motion.Clear(ctx, motionID)
	return nil
}

func ClearOutdateMotionOffer(ctx context.Context) error {
	db := dbutil.GetDB(ctx)
	MotionOfferRecord := dbquery.Use(db).MotionOfferRecord
	records, err := MotionOfferRecord.WithContext(ctx).Where(
		MotionOfferRecord.ExpiredAt.Gt(time.Now()),
		MotionOfferRecord.State.Eq(string(models.MotionOfferStatePending)),
	).Select(MotionOfferRecord.MotionID, MotionOfferRecord.ToMotionID).Find()
	if err != nil {
		return err
	}

	queryer := dbquery.Use(db)
	loader := midacontext.GetLoader[loader.Loader](ctx)

	for _, record := range records {
		err := queryer.Transaction(func(tx *dbquery.Query) error {
			MotionOfferRecord := tx.MotionOfferRecord
			Motion := tx.Motion
			_, err := MotionOfferRecord.WithContext(ctx).Where(
				MotionOfferRecord.ID.Eq(record.ID),
			).UpdateSimple(
				MotionOfferRecord.State.Value(string(models.MatchingStateTimeout)),
			)
			if err != nil {
				return err
			}
			_, err = Motion.WithContext(ctx).Where(Motion.ID.Eq(record.MotionID)).UpdateSimple(
				Motion.PendingOutNum.Add(-1),
			)
			if err != nil {
				return err
			}
			_, err = Motion.WithContext(ctx).Where(Motion.ID.Eq(record.ToMotionID)).UpdateSimple(
				Motion.PendingInNum.Add(-1),
			)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}

		loader.Motion.Clear(ctx, record.MotionID)
		loader.Motion.Clear(ctx, record.ToMotionID)
		loader.OutMotionOfferRecord.Clear(ctx, record.MotionID)
		loader.InMotionOfferRecord.Clear(ctx, record.ToMotionID)
	}

	motionIDs := map[string]struct{}{}
	for _, record := range records {
		motionIDs[record.MotionID] = struct{}{}
		motionIDs[record.ToMotionID] = struct{}{}
	}

	Motion := queryer.Motion
	_, err = queryer.Motion.WithContext(ctx).Where(
		Motion.Active.Is(true),
		Motion.Discoverable.Is(false),
		Motion.PendingOutNum.Eq(0),
		Motion.PendingInNum.Eq(0),
	).UpdateSimple(Motion.Active.Value(false))
	return err
}

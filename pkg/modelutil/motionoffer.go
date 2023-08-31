package modelutil

import (
	"context"
	"errors"
	"github.com/letjoy-club/mida-tool/logger"
	"go.uber.org/zap"
	"time"
	"whale/pkg/dbquery"
	"whale/pkg/gqlient/hoopoe"
	"whale/pkg/gqlient/smew"
	"whale/pkg/loader"
	"whale/pkg/models"
	"whale/pkg/utils"
	"whale/pkg/whalecode"

	"gorm.io/gorm"

	"github.com/letjoy-club/mida-tool/dbutil"
	"github.com/letjoy-club/mida-tool/keyer"
	"github.com/letjoy-club/mida-tool/midacode"
	"github.com/letjoy-club/mida-tool/midacontext"
	"github.com/letjoy-club/mida-tool/redisutil"
	"gorm.io/gen/field"
)

func CreateMotionOffer(ctx context.Context, myUserID, myMotionID, targetMotionID string) (string, error) {
	motionsThunk := midacontext.GetLoader[loader.Loader](ctx).Motion.LoadMany(ctx, []string{myMotionID, targetMotionID})
	motions, err := utils.ReturnThunk(motionsThunk)
	if err != nil {
		return "", err
	}
	myMotion := motions[0]
	targetMotion := motions[1]

	if myUserID != "" {
		if myMotion.UserID != myUserID {
			return "", midacode.ErrNotPermitted
		}
	}
	if targetMotion.UserID == myUserID {
		return "", whalecode.ErrCannotSendMatchingOfferToSelf
	}
	if myMotion.TopicID != targetMotion.TopicID {
		return "", whalecode.ErrCannotSendMatchingOfferToDifferentTopic
	}
	if !myMotion.Active {
		return "", whalecode.ErrYourMotionIsNotActive
	}
	if !targetMotion.Active {
		return "", whalecode.ErrTheMotionIsNotActive
	}

	// 拉黑检查
	resp, err := hoopoe.GetBlacklistRelationship(ctx, midacontext.GetServices(ctx).Hoopoe, []string{myMotion.UserID, targetMotion.UserID})
	if err != nil {
		return "", err
	}
	if resp.BlacklistRelationship != nil && len(resp.BlacklistRelationship) > 0 {
		for _, pair := range resp.BlacklistRelationship {
			if pair.A == myMotion.UserID {
				return "", whalecode.ErrUserInYourBlacklist
			} else {
				return "", whalecode.ErrYouAreInUserBlacklist
			}
		}
	}

	db := dbutil.GetDB(ctx)
	MotionOfferRecord := dbquery.Use(db).MotionOfferRecord
	record, err := MotionOfferRecord.WithContext(ctx).Where(
		MotionOfferRecord.MotionID.Eq(myMotionID),
		MotionOfferRecord.ToMotionID.Eq(targetMotionID),
	).Take()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}
	if record != nil {
		return "", whalecode.ErrAlreadySentOutMatchingOffer
	}

	release, err := redisutil.LockAll(ctx, keyer.UserMotion(myMotion.UserID), keyer.UserMotion(targetMotion.UserID))
	if err != nil {
		return "", err
	}
	defer release(ctx)

	var groupId string
	err = dbquery.Use(db).Transaction(func(tx *dbquery.Query) error {
		MatchingResult := tx.MatchingResult
		matchingResult := &models.MatchingResult{
			UserIDs:        []string{myMotion.UserID, targetMotion.UserID},
			TopicID:        myMotion.TopicID,
			MotionIDs:      []string{myMotionID, targetMotionID},
			CreatedBy:      string(models.ResultCreatedByOffer),
			ConfirmStates:  []string{string(models.MatchingResultConfirmStateConfirmed), string(models.MatchingResultConfirmStateUnconfirmed)},
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
			ExpiredAt:  time.Now().Add(time.Hour * 24 * 3),
			ChatChance: 1,
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

		rx, err = Motion.WithContext(ctx).Where(Motion.ID.Eq(targetMotion.ID)).UpdateSimple(
			Motion.InOfferNum.Add(1),
			Motion.PendingInNum.Add(1),
		)
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
		groupId = resp.CreateMotionGroup

		rx, err = MotionOfferRecord.WithContext(ctx).Where(MotionOfferRecord.ID.Eq(record.ID)).UpdateSimple(MotionOfferRecord.ChatGroupID.Value(groupId))
		if err != nil {
			return err
		}
		if rx.RowsAffected != 1 {
			return midacode.ErrStateMayHaveChanged
		}
		if rx, err := MatchingResult.WithContext(ctx).Where(MatchingResult.ID.Eq(matchingResult.ID)).UpdateSimple(
			MatchingResult.ChatGroupState.Value(models.ChatGroupStateCreated.String()),
			MatchingResult.ChatGroupID.Value(resp.CreateMotionGroup),
			MatchingResult.ChatGroupCreatedAt.Value(time.Now()),
		); err != nil {
			return err
		} else if rx.RowsAffected != 1 {
			return midacode.ErrStateMayHaveChanged
		}

		return nil
	})

	if err == nil {
		midacontext.GetLoader[loader.Loader](ctx).Motion.Clear(ctx, myMotionID)
		midacontext.GetLoader[loader.Loader](ctx).Motion.Clear(ctx, targetMotionID)
		midacontext.GetLoader[loader.Loader](ctx).InMotionOfferRecord.Clear(ctx, myMotionID)
		midacontext.GetLoader[loader.Loader](ctx).OutMotionOfferRecord.Clear(ctx, targetMotionID)
	}
	return groupId, err
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

	release, err := redisutil.LockAll(ctx, keyer.UserMotion(myMotion.UserID), keyer.UserMotion(targetMotion.UserID))
	if err != nil {
		return err
	}
	defer release(ctx)

	db := dbutil.GetDB(ctx)
	MotionOfferRecord := dbquery.Use(db).MotionOfferRecord
	record, err := MotionOfferRecord.WithContext(ctx).Where(
		MotionOfferRecord.MotionID.Eq(targetMotion.ID),
		MotionOfferRecord.ToMotionID.Eq(myMotion.ID),
	).Take()
	if err != nil {
		return err
	}
	if record.State != string(models.MotionOfferStatePending) {
		return whalecode.ErrMotionOfferIsNotPending
	}

	now := time.Now()
	matchingResultID := 0
	err = dbquery.Use(db).Transaction(func(tx *dbquery.Query) error {
		if _, err := smew.CreateTimGroup(ctx, midacontext.GetServices(ctx).Smew, record.ChatGroupID); err != nil {
			return err
		}

		MotionOfferRecord := tx.MotionOfferRecord
		rx, err := MotionOfferRecord.WithContext(ctx).Where(
			MotionOfferRecord.MotionID.Eq(targetMotion.ID),
			MotionOfferRecord.ToMotionID.Eq(myMotion.ID),
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
		// 修改matchingResult
		MatchingResult := tx.MatchingResult
		matchingResult, err := MatchingResult.WithContext(ctx).Where(MatchingResult.ChatGroupID.Eq(record.ChatGroupID)).Take()
		if err != nil {
			return err
		}
		matchingResultID = matchingResult.ID
		matchingResult.ConfirmStates = []string{models.MatchingResultConfirmStateConfirmed.String(), models.MatchingResultConfirmStateConfirmed.String()}
		if rx, err := MatchingResult.WithContext(ctx).Where(MatchingResult.ID.Eq(matchingResultID)).
			Select(MatchingResult.ConfirmStates).Updates(matchingResult); err != nil {
			return err
		} else if rx.RowsAffected != 1 {
			return midacode.ErrStateMayHaveChanged
		}
		// 修改我的 motion
		Motion := tx.Motion
		fields := []field.AssignExpr{Motion.PendingInNum.Add(-1), Motion.ActiveNum.Add(1)}
		if !myMotion.Active {
			if myMotion.PendingInNum+myMotion.PendingOutNum == 1 { // 设为不可见
				fields = append(fields, Motion.Discoverable.Value(false))
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
		if !targetMotion.Active {
			if targetMotion.PendingInNum+targetMotion.PendingOutNum == 1 { // 设为不可见
				fields = append(fields, Motion.Discoverable.Value(false))
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
	midacontext.GetLoader[loader.Loader](ctx).MatchingResult.Clear(ctx, matchingResultID)
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

	release, err := redisutil.LockAll(ctx, keyer.UserMotion(myMotion.UserID), keyer.UserMotion(targetMotion.UserID))
	if err != nil {
		return err
	}
	defer release(ctx)

	db := dbutil.GetDB(ctx)
	MotionOfferRecord := dbquery.Use(db).MotionOfferRecord
	motionOffer, err := MotionOfferRecord.WithContext(ctx).Where(
		MotionOfferRecord.MotionID.Eq(targetMotion.ID),
		MotionOfferRecord.ToMotionID.Eq(myMotion.ID),
	).Take()
	if err != nil {
		return err
	}
	if motionOffer.State != string(models.MotionOfferStatePending) {
		return whalecode.ErrMotionOfferIsNotPending
	}

	matchingResultID := 0
	now := time.Now()
	err = dbquery.Use(db).Transaction(func(tx *dbquery.Query) error {
		MotionOfferRecord := tx.MotionOfferRecord
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

		// 修改matchingResult
		MatchingResult := tx.MatchingResult
		matchingResult, err := MatchingResult.WithContext(ctx).Where(MatchingResult.ChatGroupID.Eq(motionOffer.ChatGroupID)).Take()
		if err != nil {
			return err
		}
		updates := &models.MatchingResult{
			ConfirmStates:  []string{models.MatchingResultConfirmStateConfirmed.String(), models.MatchingResultConfirmStateRejected.String()},
			ChatGroupState: models.ChatGroupStateClosed.String(),
			Closed:         true,
			FinishedAt:     &now,
		}
		matchingResultID = matchingResult.ID
		if rx, err := MatchingResult.WithContext(ctx).Where(MatchingResult.ID.Eq(matchingResultID)).
			Updates(updates); err != nil {
			return err
		} else if rx.RowsAffected != 1 {
			return midacode.ErrStateMayHaveChanged
		}

		Motion := tx.Motion
		// 修改对方的 motion
		fields := []field.AssignExpr{Motion.PendingOutNum.Add(-1)}
		if !targetMotion.Active {
			if targetMotion.PendingInNum+targetMotion.PendingOutNum == 1 { // 设为不可见
				fields = append(fields, Motion.Discoverable.Value(false))
			}
		}
		rx, err = Motion.WithContext(ctx).Where(Motion.ID.Eq(targetMotionID)).UpdateSimple(fields...)
		if err != nil {
			return err
		}
		if rx.RowsAffected != 1 {
			return midacode.ErrStateMayHaveChanged
		}

		// 修改我的 motion
		fields = []field.AssignExpr{Motion.PendingInNum.Add(-1)}
		if !myMotion.Active {
			if myMotion.PendingInNum+myMotion.PendingOutNum == 1 { // 设为不可见
				fields = append(fields, Motion.Discoverable.Value(false))
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
	loader.MatchingResult.Clear(ctx, matchingResultID)
	return nil
}

// CancelMotionOffer 取消MotionOffer，暂时不可用
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
			return whalecode.ErrMotionOfferIsNotPending
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

// ClearOutDateMotionOffer 清理过期的MotionOffer
func ClearOutDateMotionOffer(ctx context.Context) error {
	db := dbutil.GetDB(ctx)
	MotionOfferRecord := dbquery.Use(db).MotionOfferRecord
	records, err := MotionOfferRecord.WithContext(ctx).Where(
		MotionOfferRecord.ExpiredAt.Lt(time.Now()),
		MotionOfferRecord.State.Eq(string(models.MotionOfferStatePending)),
	).Select(MotionOfferRecord.MotionID, MotionOfferRecord.ToMotionID).Find()
	if err != nil {
		return err
	}

	queryer := dbquery.Use(db)
	loader := midacontext.GetLoader[loader.Loader](ctx)
	for _, record := range records {
		motionThunk := loader.Motion.LoadMany(ctx, []string{record.MotionID, record.ToMotionID})
		motions, err := utils.ReturnThunk(motionThunk)
		if err != nil {
			return err
		}

		fromMotion := motions[0]
		toMotion := motions[1]

		matchingResultID := 0
		if err := queryer.Transaction(func(tx *dbquery.Query) error {
			MotionOfferRecord := tx.MotionOfferRecord
			Motion := tx.Motion
			MatchingResult := tx.MatchingResult

			if _, err := MotionOfferRecord.WithContext(ctx).Where(MotionOfferRecord.ID.Eq(record.ID)).UpdateSimple(
				MotionOfferRecord.State.Value(string(models.MatchingStateTimeout)),
			); err != nil {
				return err
			}

			// fromMotion
			fields := []field.AssignExpr{Motion.PendingOutNum.Add(-1)}
			if !fromMotion.Active {
				if fromMotion.PendingInNum+fromMotion.PendingOutNum == 1 { // 设为不可见
					fields = append(fields, Motion.Discoverable.Value(false))
				}
			}
			if rx, err := Motion.WithContext(ctx).Where(Motion.ID.Eq(fromMotion.ID)).UpdateSimple(fields...); err != nil {
				return err
			} else if rx.RowsAffected != 1 {
				return midacode.ErrStateMayHaveChanged
			}
			// toMotion
			fields = []field.AssignExpr{Motion.PendingOutNum.Add(-1)}
			if !toMotion.Active {
				if toMotion.PendingInNum+toMotion.PendingOutNum == 1 { // 设为不可见
					fields = append(fields, Motion.Discoverable.Value(false))
				}
			}
			if rx, err := Motion.WithContext(ctx).Where(Motion.ID.Eq(toMotion.ID)).UpdateSimple(fields...); err != nil {
				return err
			} else if rx.RowsAffected != 1 {
				return midacode.ErrStateMayHaveChanged
			}

			if record.ChatGroupID != "" {
				matchingResult, err := MatchingResult.WithContext(ctx).Where(MatchingResult.ChatGroupID.Eq(record.ChatGroupID)).Take()
				if err != nil {
					return err
				}
				matchingResultID = matchingResult.ID
				if rx, err := MatchingResult.WithContext(ctx).Where(MatchingResult.ID.Eq(matchingResultID)).UpdateSimple(
					MatchingResult.Closed.Value(true),
					MatchingResult.ChatGroupState.Value(models.ChatGroupStateClosed.String()),
					MatchingResult.FinishedAt.Value(time.Now()),
				); err != nil {
					return err
				} else if rx.RowsAffected != 1 {
					return midacode.ErrStateMayHaveChanged
				}

				if _, err := smew.DestroyGroup(ctx, midacontext.GetServices(ctx).Smew, record.ChatGroupID); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return err
		}

		loader.Motion.Clear(ctx, record.MotionID)
		loader.Motion.Clear(ctx, record.ToMotionID)
		loader.OutMotionOfferRecord.Clear(ctx, record.MotionID)
		loader.InMotionOfferRecord.Clear(ctx, record.ToMotionID)
		loader.MatchingResult.Clear(ctx, matchingResultID)
	}
	return nil
}

func SendChatInOffer(ctx context.Context, myUserID, myMotionID, targetMotionID, sentence string) error {
	db := dbutil.GetDB(ctx)
	MotionOfferRecord := dbquery.Use(db).MotionOfferRecord
	record, err := MotionOfferRecord.WithContext(ctx).Where(MotionOfferRecord.MotionID.Eq(myMotionID)).Where(MotionOfferRecord.ToMotionID.Eq(targetMotionID)).Take()
	if err != nil {
		return midacode.ItemMayNotFound(err)
	}
	if record.ChatChance <= 0 {
		return whalecode.ErrChatChanceNotEnough
	}
	if record.State != string(models.MotionOfferStatePending) {
		return whalecode.ErrOnlyChatWhenNotAccepted
	}
	if myUserID != "" {
		if record.UserID != myUserID {
			return midacode.ErrNotPermitted
		}
	}

	_, err = smew.SendTextMessage(ctx, midacontext.GetServices(ctx).Smew, record.ChatGroupID, record.UserID, sentence)
	if err != nil {
		return err
	}
	_, err = MotionOfferRecord.WithContext(ctx).Where(MotionOfferRecord.ID.Eq(record.ID)).UpdateSimple(MotionOfferRecord.ChatChance.Add(-1))
	return err
}

// FinishMotionOffer 结束邀约 fromMatchingID 表示邀约的一方，toMatchingID 表示被邀约的一方
func FinishMotionOffer(ctx context.Context, myUserID, fromMatchingID, toMatchingID string) error {
	db := dbutil.GetDB(ctx)
	MotionOfferRecord := dbquery.Use(db).MotionOfferRecord
	record, err := MotionOfferRecord.WithContext(ctx).Where(
		MotionOfferRecord.MotionID.Eq(fromMatchingID), MotionOfferRecord.ToMotionID.Eq(toMatchingID),
	).Take()
	if err != nil {
		return midacode.ItemMayNotFound(err)
	}
	if record.State != string(models.MotionOfferStateAccepted) && record.State != string(models.MotionOfferStateFinished) {
		return whalecode.ErrMotionCanOnlyFinishedWhenAccepted
	}
	if myUserID != "" {
		// 可以是被邀请方或者邀请方来结束邀约
		if record.UserID != myUserID && record.ToUserID != myUserID {
			return midacode.ErrNotPermitted
		}
	}

	matchingResultID := 0
	if err = dbquery.Use(db).Transaction(func(tx *dbquery.Query) error {
		if record.ChatGroupID != "" {
			smewClient := midacontext.GetServices(ctx).Smew
			if myUserID == "" { // 管理端关闭
				if _, err := smew.DestroyGroup(ctx, smewClient, record.ChatGroupID); err != nil {
					return err
				}
			} else { // 用户关闭
				if _, err := smew.GroupMemberLeave(ctx, smewClient, record.ChatGroupID, myUserID); err != nil {
					return err
				}
			}
		}
		if record.State == string(models.MotionOfferStateAccepted) {
			MotionOfferRecord := tx.MotionOfferRecord
			if _, err = MotionOfferRecord.WithContext(ctx).Where(MotionOfferRecord.ID.Eq(record.ID)).
				UpdateSimple(MotionOfferRecord.State.Value(string(models.MotionOfferStateFinished))); err != nil {
				return err
			}

			MatchingResult := tx.MatchingResult
			matchingResult, err := MatchingResult.WithContext(ctx).Where(MatchingResult.ChatGroupID.Eq(record.ChatGroupID)).Take()
			if err != nil {
				return err
			}
			matchingResultID = matchingResult.ID
			if rx, err := MatchingResult.WithContext(ctx).Where(MatchingResult.ID.Eq(matchingResultID)).UpdateSimple(
				MatchingResult.Closed.Value(true),
				MatchingResult.ChatGroupState.Value(models.ChatGroupStateClosed.String()),
				MatchingResult.FinishedAt.Value(time.Now()),
			); err != nil {
				return err
			} else if rx.RowsAffected != 1 {
				return midacode.ErrStateMayHaveChanged
			}
		}

		return nil
	}); err != nil {
		return err
	}

	midacontext.GetLoader[loader.Loader](ctx).InMotionOfferRecord.Clear(ctx, toMatchingID)
	midacontext.GetLoader[loader.Loader](ctx).OutMotionOfferRecord.Clear(ctx, fromMatchingID)
	midacontext.GetLoader[loader.Loader](ctx).MatchingResult.Clear(ctx, matchingResultID)
	return nil
}

func RefreshMotionState(ctx context.Context) error {
	db := dbutil.GetDB(ctx)
	Motion := dbquery.Use(db).Motion
	MotionOfferRecord := dbquery.Use(db).MotionOfferRecord

	offset := 0
	for {
		motions, err := Motion.WithContext(ctx).Offset(offset).Limit(10).Find()
		if err != nil {
			return err
		}
		if motions == nil || len(motions) == 0 {
			break
		}

		for _, motion := range motions {
			// Motion是否Close
			active := true
			if !motion.Active || !motion.Discoverable {
				active = false
			}

			inOffers, err := MotionOfferRecord.WithContext(ctx).Where(MotionOfferRecord.ToMotionID.Eq(motion.ID)).Find()
			if err != nil {
				logger.L.Error("query inOffers error", zap.Error(err), zap.String("motionId", motion.ID))
				return err
			}
			outOffers, err := MotionOfferRecord.WithContext(ctx).Where(MotionOfferRecord.MotionID.Eq(motion.ID)).Find()
			if err != nil {
				logger.L.Error("query outOffers error", zap.Error(err), zap.String("motionId", motion.ID))
				return err
			}
			inOfferNum := len(inOffers)
			outOfferNum := len(outOffers)
			pendingInNum := 0
			pendingOutNum := 0
			activeNum := 0

			for _, offer := range inOffers {
				if offer.State == string(models.MotionOfferStatePending) {
					pendingInNum++
				}
				if offer.State == string(models.MotionOfferStateAccepted) ||
					offer.State == string(models.MotionOfferStateFinished) {
					activeNum++
				}
			}
			for _, offer := range outOffers {
				if offer.State == string(models.MotionOfferStatePending) {
					pendingOutNum++
				}
				if offer.State == string(models.MotionOfferStateAccepted) ||
					offer.State == string(models.MotionOfferStateFinished) {
					activeNum++
				}
			}

			discoverable := true
			if !active && pendingInNum+pendingOutNum == 0 {
				discoverable = false
			}

			if _, err := Motion.WithContext(ctx).Where(Motion.ID.Eq(motion.ID)).UpdateSimple(
				Motion.Active.Value(active),
				Motion.Discoverable.Value(discoverable),
				Motion.InOfferNum.Value(inOfferNum),
				Motion.OutOfferNum.Value(outOfferNum),
				Motion.PendingInNum.Value(pendingInNum),
				Motion.PendingOutNum.Value(pendingOutNum),
				Motion.ActiveNum.Value(activeNum),
			); err != nil {
				logger.L.Error("update motion error", zap.Error(err), zap.String("motionId", motion.ID))
				return err
			}

		}
		offset += 10
	}

	loader := midacontext.GetLoader[loader.Loader](ctx)
	loader.Motion.ClearAll()
	loader.InMotionOfferRecord.ClearAll()
	loader.OutMotionOfferRecord.ClearAll()

	return nil
}

package modelutil

import (
	"context"
	"time"
	"whale/pkg/dbquery"
	"whale/pkg/loader"
	"whale/pkg/models"
	"whale/pkg/whalecode"

	"github.com/letjoy-club/mida-tool/dbutil"
	"github.com/letjoy-club/mida-tool/graphqlutil"
	"github.com/letjoy-club/mida-tool/keyer"
	"github.com/letjoy-club/mida-tool/midacode"
	"github.com/letjoy-club/mida-tool/midacontext"
	"github.com/letjoy-club/mida-tool/redisutil"
	"github.com/letjoy-club/mida-tool/shortid"
	"github.com/samber/lo"
	"gorm.io/gen/field"
)

func CreateMotion(ctx context.Context, userID string, param *models.CreateMotionParam) (*models.Motion, error) {
	profileThunk := midacontext.GetLoader[loader.Loader](ctx).UserProfile.Load(ctx, userID)
	profile, err := profileThunk()
	if err != nil {
		return nil, err
	}

	db := dbutil.GetDB(ctx)
	Motion := dbquery.Use(db).Motion
	basicQuota := 2
	motion := &models.Motion{
		ID:       shortid.NewWithTime("mo_", 4),
		UserID:   userID,
		Gender:   param.Gender.String(),
		CityID:   param.CityID,
		Active:   true,
		Remark:   *param.Remark,
		Deadline: time.Now().Add(time.Hour * 24 * 7),
		MyGender: string(profile.Gender),
		TopicID:  param.TopicID,
		Properties: lo.Map(param.Properties, func(p *models.MotionPropertyParam, i int) models.MotionProperty {
			return models.MotionProperty{ID: p.ID, Values: p.Values}
		}),
		BasicQuota:       basicQuota,
		Discoverable:     true,
		AreaIDs:          param.AreaIds,
		DayRange:         param.DayRange,
		PreferredPeriods: SimplifyPreferredPeriods(param.PreferredPeriods),
		RemainQuota:      basicQuota,
	}
	if err := Motion.WithContext(ctx).Create(motion); err != nil {
		return nil, err
	}
	return motion, nil
}

func UpdateMotion(ctx context.Context, motionID string, param *models.UpdateMotionParam) error {
	db := dbutil.GetDB(ctx)
	Motion := dbquery.Use(db).Motion

	fields := []field.AssignExpr{}
	if param.Gender != nil {
		fields = append(fields, Motion.Gender.Value(param.Gender.String()))
	}
	if param.CityID != nil {
		fields = append(fields, Motion.CityID.Value(*param.CityID))
	}
	if param.Remark != nil {
		fields = append(fields, Motion.Remark.Value(*param.Remark))
	}
	if param.AreaIds != nil {
		fields = append(fields, Motion.AreaIDs.Value(graphqlutil.ElementList[string](param.AreaIds)))
	}
	if param.Deadline != nil {
		fields = append(fields, Motion.Deadline.Value(*param.Deadline))
	}
	if param.DayRange != nil {
		fields = append(fields, Motion.DayRange.Value(graphqlutil.ElementList[string](param.DayRange)))
	}
	if param.PreferredPeriods != nil {
		fields = append(fields, Motion.PreferredPeriods.Value(graphqlutil.ElementList[string](SimplifyPreferredPeriods(param.PreferredPeriods))))
	}
	if param.Properties != nil {
		fields = append(fields, Motion.Properties.Value(
			graphqlutil.ElementList[*models.MotionProperty](
				lo.Map(param.Properties, func(p *models.MotionPropertyParam, i int) *models.MotionProperty {
					return &models.MotionProperty{ID: p.ID, Values: p.Values}
				})),
		),
		)
	}

	_, err := Motion.WithContext(ctx).Where(Motion.ID.Eq(motionID)).UpdateSimple(fields...)
	if err != nil {
		return err
	}
	return nil
}

func CloseMotion(ctx context.Context, userID, motionID string) error {
	motionThunk := midacontext.GetLoader[loader.Loader](ctx).Motion.Load(ctx, motionID)
	motion, err := motionThunk()
	if err != nil {
		return err
	}

	if motion.UserID != userID {
		return midacode.ErrNotPermitted
	}

	if !motion.Active {
		return whalecode.ErrMatchingOfferIsNotActive
	}

	inMotionOfferThunk := midacontext.GetLoader[loader.Loader](ctx).InMotionOfferRecord.Load(ctx, motionID)
	inMotionOffer, err := inMotionOfferThunk()
	if err != nil {
		return err
	}
	outMotionOfferThunk := midacontext.GetLoader[loader.Loader](ctx).OutMotionOfferRecord.Load(ctx, motionID)
	outMotionOffer, err := outMotionOfferThunk()
	if err != nil {
		return err
	}

	release, err := redisutil.LockAll(ctx, keyer.UserMotion(motion.UserID))
	if err != nil {
		return err
	}
	defer release(ctx)

	myOutOfferIDs := []string{}
	myInOfferIDs := []string{}

	db := dbutil.GetDB(ctx)
	err = dbquery.Use(db).Transaction(func(tx *dbquery.Query) error {
		Motion := tx.Motion
		MotionOfferRecord := tx.MotionOfferRecord
		latestMotion, err := Motion.WithContext(ctx).Where(Motion.ID.Eq(motionID)).Take()
		if err != nil {
			return err
		}
		if latestMotion.UpdatedAt != motion.UpdatedAt {
			return midacode.ErrStateMayHaveChanged
		}

		// 关闭我收到的所有邀约
		for _, offer := range inMotionOffer.Offers {
			if offer.State == string(models.MotionOfferStatePending) {
				myInOfferIDs = append(myInOfferIDs, offer.MotionID)
			}
		}
		if len(myInOfferIDs) > 0 {
			// 关闭对方的邀约，回复对方的配额
			rx, err := Motion.WithContext(ctx).Where(Motion.ID.In(myInOfferIDs...)).UpdateSimple(
				Motion.RemainQuota.Add(1),
				Motion.PendingOutNum.Add(-1),
			)
			if err != nil {
				return err
			}
			if rx.RowsAffected != int64(len(myInOfferIDs)) {
				return midacode.ErrStateMayHaveChanged
			}
			// 对方的邀约状态改为 rejected
			rx, err = MotionOfferRecord.WithContext(ctx).Where(MotionOfferRecord.MotionID.In(myInOfferIDs...)).UpdateSimple(
				MotionOfferRecord.State.Value(string(models.MatchingOfferStateRejected)),
				MotionOfferRecord.ReactAt.Value(time.Now()),
			)
			if err != nil {
				return err
			}
			if rx.RowsAffected != int64(len(myInOfferIDs)) {
				return midacode.ErrStateMayHaveChanged
			}
		}

		// 关闭我发给对方的所有邀约
		for _, offer := range outMotionOffer.Offers {
			if offer.State == string(models.MotionOfferStatePending) {
				myOutOfferIDs = append(myOutOfferIDs, offer.ToMotionID)
			}
		}
		if len(myOutOfferIDs) > 0 {
			// 关闭我发给对方的邀约
			rx, err := Motion.WithContext(ctx).Where(Motion.ID.In(myOutOfferIDs...)).UpdateSimple(
				Motion.PendingInNum.Add(-1),
			)
			if err != nil {
				return err
			}
			if rx.RowsAffected != int64(len(myOutOfferIDs)) {
				return midacode.ErrStateMayHaveChanged
			}
		}

		// 关闭我的 motion
		rx, err := Motion.WithContext(ctx).Where(Motion.ID.Eq(motionID)).UpdateSimple(
			Motion.Active.Value(false),
			Motion.RemainQuota.Add(0),
			Motion.PendingOutNum.Add(0),
			Motion.PendingInNum.Add(0),
		)
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
	// 清理缓存
	loader := midacontext.GetLoader[loader.Loader](ctx)
	if len(myInOfferIDs) > 0 {
		for _, id := range myInOfferIDs {
			loader.Motion.Clear(ctx, id)
			loader.OutMotionOfferRecord.Clear(ctx, id)
		}
	}
	if len(myOutOfferIDs) > 0 {
		for _, id := range myOutOfferIDs {
			loader.Motion.Clear(ctx, id)
			loader.InMotionOfferRecord.Clear(ctx, id)
		}
	}
	loader.InMotionOfferRecord.Clear(ctx, motionID)
	loader.OutMotionOfferRecord.Clear(ctx, motionID)
	loader.Motion.Clear(ctx, motionID)
	return nil
}

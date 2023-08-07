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

	if err := checkMatchingParam(ctx, userID, param.TopicID, param.CityID, param.Gender); err != nil {
		return nil, err
	}

	db := dbutil.GetDB(ctx)
	Motion := dbquery.Use(db).Motion

	count, err := Motion.WithContext(ctx).Where(
		Motion.Active.Is(true),
		Motion.TopicID.Eq(param.TopicID),
		Motion.UserID.Eq(userID),
	).Count()
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, whalecode.ErrIsAlreadyHasActiveMotionOfTopic
	}

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
		Discoverable:     true,
		AreaIDs:          param.AreaIds,
		DayRange:         param.DayRange,
		PreferredPeriods: SimplifyPreferredPeriods(param.PreferredPeriods),
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

	if userID != "" {
		if motion.UserID != userID {
			return midacode.ErrNotPermitted
		}
	}

	if !motion.Active {
		return nil
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

	for _, offer := range inMotionOffer.Offers {
		if offer.State == models.MotionOfferStatePending.String() {
			myInOfferIDs = append(myInOfferIDs, offer.MotionID)
		}
	}

	for _, offer := range outMotionOffer.Offers {
		if offer.State == models.MotionOfferStatePending.String() {
			myOutOfferIDs = append(myOutOfferIDs, offer.ToMotionID)
		}
	}

	db := dbutil.GetDB(ctx)
	err = dbquery.Use(db).Transaction(func(tx *dbquery.Query) error {
		Motion := tx.Motion
		latestMotion, err := Motion.WithContext(ctx).Where(Motion.ID.Eq(motionID)).Take()
		if err != nil {
			return err
		}
		if latestMotion.UpdatedAt != motion.UpdatedAt {
			return midacode.ErrStateMayHaveChanged
		}

		fields := []field.AssignExpr{}
		fields = append(fields, Motion.Discoverable.Value(false))

		if len(myInOfferIDs) > 0 || len(myOutOfferIDs) > 0 {
			// 不关闭，但是不可见
			fields = append(fields, Motion.Active.Value(false))
		}
		rx, err := Motion.WithContext(ctx).Where(Motion.ID.Eq(motionID)).UpdateSimple(fields...)
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

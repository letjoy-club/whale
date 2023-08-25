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
	"github.com/letjoy-club/mida-tool/logger"
	"github.com/letjoy-club/mida-tool/midacode"
	"github.com/letjoy-club/mida-tool/midacontext"
	"github.com/letjoy-club/mida-tool/redisutil"
	"github.com/letjoy-club/mida-tool/shortid"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"gorm.io/gen/field"
)

func CreateMotion(ctx context.Context, userID string, param *models.CreateMotionParam) (*models.Motion, error) {
	profileThunk := midacontext.GetLoader[loader.Loader](ctx).UserProfile.Load(ctx, userID)
	profile, err := profileThunk()
	if err != nil {
		return nil, err
	}

	// allAreas := areaTable
	// if len(allAreas) == 0 {
	// 	return nil, whalecode.ErrAreaNotSupport
	// }
	//
	// if len(param.AreaIds) > 0 {
	// 	param.AreaIds = lo.Uniq(param.AreaIds)
	// 	for _, selectedAreaID := range param.AreaIds {
	// 		index := lo.IndexOf(allAreas, func(areaID string, i int) bool {
	// 			return areaID == selectedAreaID
	// 		})
	// 		if index == -1 {
	// 			return nil, whalecode.ErrAreaNotSupport
	// 		}
	// 	}
	// }
	//
	// todo: 临时下掉限制，后续重新测试上线
	//durationConstraintThunk := midacontext.GetLoader[loader.Loader](ctx).DurationConstraint.Load(ctx, userID)
	//durationConstraint, err := durationConstraintThunk()
	//if err != nil {
	//	return nil, err
	//}
	//if durationConstraint.RemainMotionQuota <= 0 {
	//	return nil, whalecode.ErrMotionQuotaNotEnough
	//}

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

	matchingStartAt := time.Now().Add(time.Hour * 24 * 7)

	matching := &models.Matching{
		ID:               shortid.NewWithTime("m_", 4),
		UserID:           userID,
		Gender:           param.Gender.String(),
		CityID:           param.CityID,
		AreaIDs:          param.AreaIds,
		DayRange:         param.DayRange,
		PreferredPeriods: SimplifyPreferredPeriods(param.PreferredPeriods),
		TopicID:          param.TopicID,
		Properties: lo.Map(param.Properties, func(p *models.MotionPropertyParam, i int) models.MatchingProperty {
			return models.MatchingProperty{ID: p.ID, Values: p.Values}
		}),
		MyGender:        string(profile.Gender),
		State:           string(models.MatchingStateMatching),
		Deadline:        time.Now().Add(time.Hour * 24 * 7),
		StartMatchingAt: &matchingStartAt,
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

	matching.RelatedMotionID = motion.ID
	motion.RelatedMatchingID = matching.ID

	err = dbquery.Use(db).Transaction(func(tx *dbquery.Query) error {
		Matching := tx.Matching
		if existMatching, err := Matching.WithContext(ctx).Where(
			Matching.TopicID.Eq(motion.TopicID), Matching.UserID.Eq(motion.UserID), Matching.State.Eq(string(models.MatchingStateMatching)),
		).Take(); err != nil {
			// 如果没找到对应话题正在进行中的匹配，说明可以为用户创建
			if err := Matching.WithContext(ctx).Create(matching); err != nil {
				// 创建失败，忽略
				motion.RelatedMatchingID = ""
				logger.L.Error("failed to create corresponding matching", zap.Error(err))
			}
		} else if existMatching.RelatedMotionID == "" {
			// 这个匹配是由 motion 创建的，需要主动关闭，并创建新匹配
			_, err := Matching.WithContext(ctx).Where(Matching.ID.Eq(existMatching.ID)).UpdateSimple(Matching.State.Value(string(models.MatchingStateCanceled)))
			if err != nil {
				// 关闭匹配失败，不创建新匹配
				motion.RelatedMatchingID = ""
				logger.L.Error("failed to close previous matching created by motion", zap.Error(err))
			} else {
				midacontext.GetLoader[loader.Loader](ctx).Matching.Clear(ctx, existMatching.ID)
				// 已经关闭原有匹配，创建新匹配
				if err := Matching.WithContext(ctx).Create(matching); err != nil {
					// 创建失败，忽略
					motion.RelatedMatchingID = ""
					logger.L.Error("failed to create corresponding matching", zap.Error(err))
				}
			}
		}

		if err := Motion.WithContext(ctx).Create(motion); err != nil {
			return err
		}

		// todo: 临时下掉限制，后续重新测试上线
		//tx.DurationConstraint.WithContext(ctx).Where(tx.DurationConstraint.ID.Eq(durationConstraint.ID)).UpdateSimple(tx.DurationConstraint.RemainMotionQuota.Value(durationConstraint.RemainMotionQuota - 1))
		//midacontext.GetLoader[loader.Loader](ctx).DurationConstraint.Clear(ctx, userID)
		return nil
	})

	return motion, err
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
		))
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
		if motion.RelatedMatchingID != "" {
			// 关闭没有结束的匹配
			Matching := tx.Matching
			_, err := Matching.WithContext(ctx).Where(
				Matching.ID.Eq(motion.RelatedMatchingID),
				Matching.State.Eq(models.MatchingStateMatched.String()),
			).UpdateSimple(Matching.State.Value(string(models.MatchingStateCanceled)))
			if err != nil {
				return err
			}
			midacontext.GetLoader[loader.Loader](ctx).Matching.Clear(ctx, motion.RelatedMatchingID)
		}
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

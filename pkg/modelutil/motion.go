package modelutil

import (
	"context"
	"time"
	"unicode/utf8"
	"whale/pkg/dbquery"
	"whale/pkg/gqlient/hoopoe"
	"whale/pkg/loader"
	"whale/pkg/models"
	"whale/pkg/whalecode"

	"github.com/golang-module/carbon"

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
	// 创建前检查
	if err := checkCreateMotionParam(ctx, userID, param); err != nil {
		return nil, err
	}
	// 额度检查
	durationConstraintThunk := midacontext.GetLoader[loader.Loader](ctx).DurationConstraint.Load(ctx, userID)
	durationConstraint, err := durationConstraintThunk()
	if err != nil {
		return nil, err
	}
	if durationConstraint.RemainMotionQuota <= 0 { // 限制每周可发起的 motion 次数
		return nil, whalecode.ErrMotionQuotaNotEnough
	}

	profileThunk := midacontext.GetLoader[loader.Loader](ctx).UserProfile.Load(ctx, userID)
	profile, err := profileThunk()
	if err != nil {
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

	matchingStartAt := time.Now().Add(time.Minute * 30)

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
		Deadline:        *param.Deadline,
		StartMatchingAt: &matchingStartAt,
	}

	motion := &models.Motion{
		ID:       shortid.NewWithTime("mo_", 4),
		UserID:   userID,
		Gender:   param.Gender.String(),
		CityID:   param.CityID,
		Active:   true,
		Remark:   *param.Remark,
		Deadline: *param.Deadline,
		MyGender: string(profile.Gender),
		TopicID:  param.TopicID,
		Properties: lo.Map(param.Properties, func(p *models.MotionPropertyParam, i int) models.MotionProperty {
			return models.MotionProperty{ID: p.ID, Values: p.Values}
		}),
		Quick:            *param.Quick,
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

		midacontext.GetLoader[loader.Loader](ctx).AllMotion.AppendNewMotion(ctx, motion)
		// 更新用户的剩余匹配次数
		tx.DurationConstraint.WithContext(ctx).Where(tx.DurationConstraint.ID.Eq(durationConstraint.ID)).
			UpdateSimple(tx.DurationConstraint.RemainMotionQuota.Add(-1))
		midacontext.GetLoader[loader.Loader](ctx).DurationConstraint.Clear(ctx, userID)
		return nil
	})

	return motion, err
}

// checkCreateMotionParam 创建动议前检查
func checkCreateMotionParam(ctx context.Context, userID string, param *models.CreateMotionParam) error {
	// 字段检查
	if userID == "" {
		return whalecode.ErrUserIDCannotBeEmpty
	}
	if param.TopicID == "" {
		return whalecode.ErrTopicIdShouldNotBeEmpty
	}
	if param.CityID == "" {
		return whalecode.ErrCityIdShouldNotBeEmpty
	}
	if param.Remark == nil || utf8.RuneCountInString(*param.Remark) < 5 {
		return whalecode.ErrRemarkTooShort
	}
	if utf8.RuneCountInString(*param.Remark) > 250 {
		return whalecode.ErrRemarkTooLong
	}

	defaultDeadline := time.Now().Add(time.Hour * 24 * 7)

	if param.Deadline != nil {
		if param.Deadline.Before(time.Now()) {
			return whalecode.ErrDeadlineShouldNotBeBeforeNow
		}
	} else {
		param.Deadline = &defaultDeadline
	}

	if param.Quick != nil {
		if *param.Quick {
			// 极速搭过期时间为当天结束
			defaultDeadline = carbon.Now().EndOfDay().ToStdTime()
			param.Deadline = &defaultDeadline
		}
	} else {
		defaultQuick := false
		param.Quick = &defaultQuick
	}

	// 基础检查
	res, err := hoopoe.CreateMotionCheck(ctx, midacontext.GetServices(ctx).Hoopoe, param.TopicID, param.CityID, userID)
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

	// 内容检查，长于一定长度才进行
	if len(*param.Remark) > 1 {
		contentCheckRes, err := hoopoe.TextCheck(ctx, midacontext.GetServices(ctx).Hoopoe, userID, *param.Remark)
		if err != nil {
			logger.L.Error("failed to check motion content", zap.Error(err))
		} else {
			if contentCheckRes != nil && contentCheckRes.TextCheck == hoopoe.TextCheckResultRisky {
				return whalecode.ErrMotionContentRisky
			}
		}
	}
	return nil
}

func UpdateMotion(ctx context.Context, motionID string, param *models.UpdateMotionParam) error {
	thunk := midacontext.GetLoader[loader.Loader](ctx).Motion.Load(ctx, motionID)
	motion, err := thunk()
	if err != nil {
		return err
	}
	if motion == nil {
		return midacode.ErrItemNotFound
	}
	if !motion.Active {
		return whalecode.ErrYourMotionIsNotActive
	}

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
	if param.Quick != nil {
		fields = append(fields, Motion.Quick.Value(*param.Quick))
	}

	if _, err := Motion.WithContext(ctx).Where(Motion.ID.Eq(motionID)).UpdateSimple(fields...); err != nil {
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

	release, err := redisutil.LockAll(ctx, keyer.UserMotion(motion.UserID))
	if err != nil {
		return err
	}
	defer release(ctx)

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
		if !latestMotion.Active {
			return nil
		}

		fields := []field.AssignExpr{}
		fields = append(fields, Motion.Active.Value(false))

		if motion.PendingInNum+motion.PendingOutNum == 0 { // 关联的Offer全部处理后，最终设置为不可见
			fields = append(fields, Motion.Discoverable.Value(false))
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
	loader.InMotionOfferRecord.Clear(ctx, motionID)
	loader.OutMotionOfferRecord.Clear(ctx, motionID)
	loader.Motion.Clear(ctx, motionID)
	return nil
}

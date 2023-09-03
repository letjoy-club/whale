package modelutil

import (
	"context"
	"errors"
	"github.com/golang-module/carbon"
	"github.com/letjoy-club/mida-tool/dbutil"
	"github.com/letjoy-club/mida-tool/logger"
	"github.com/letjoy-club/mida-tool/midacontext"
	"gorm.io/gorm"
	"time"
	"whale/pkg/dbquery"
	"whale/pkg/loader"
	"whale/pkg/models"
)

// UpdateUserRights 更新用户等级权益
func UpdateUserRights(ctx context.Context, userId string, level int) error {
	if userId == "" || level == 0 {
		logger.L.Error("UpdateUserRights - event payload illegal")
		return nil
	}

	config, err := GetLevelRightsConfig(ctx, level)
	if err != nil {
		logger.L.Error("UpdateUserRights - can not find user level rights config")
		return err
	}

	db := dbutil.GetDB(ctx)
	DurationConstraint := dbquery.Use(db).DurationConstraint
	durationConstraint, err := DurationConstraint.WithContext(ctx).Where(
		DurationConstraint.UserID.Eq(userId),
		DurationConstraint.StopDate.Gt(time.Now()),
	).Take()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	weekStart := carbon.Now().StartOfWeek().ToStdTime()
	weekEnd := carbon.Now().EndOfWeek().ToStdTime()
	if durationConstraint == nil {
		if err := DurationConstraint.WithContext(ctx).Create(&models.DurationConstraint{
			StartDate:         weekStart,
			StopDate:          weekEnd,
			TotalMotionQuota:  config.MotionQuota,
			RemainMotionQuota: config.MotionQuota,
			TotalOfferQuota:   config.OfferQuota,
			RemainOfferQuota:  config.OfferQuota,
		}); err != nil {
			return err
		}
	} else {
		motionQuotaAdd := config.MotionQuota - durationConstraint.TotalMotionQuota
		if motionQuotaAdd < 0 { // 额度只增不减
			motionQuotaAdd = 0
		}
		offerQuotaAdd := config.OfferQuota - durationConstraint.TotalOfferQuota
		if offerQuotaAdd < 0 { // 额度只增不减
			offerQuotaAdd = 0
		}
		if _, err := DurationConstraint.WithContext(ctx).Where(DurationConstraint.ID.Eq(durationConstraint.ID)).UpdateSimple(
			DurationConstraint.TotalMotionQuota.Add(motionQuotaAdd),
			DurationConstraint.RemainMotionQuota.Add(motionQuotaAdd),
			DurationConstraint.TotalOfferQuota.Add(offerQuotaAdd),
			DurationConstraint.RemainOfferQuota.Add(offerQuotaAdd),
		); err != nil {
			return err
		}
	}
	midacontext.GetLoader[loader.Loader](ctx).DurationConstraint.Clear(ctx, userId)

	return nil
}

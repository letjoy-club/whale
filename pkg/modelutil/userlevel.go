package modelutil

import (
	"context"
	"errors"
	"github.com/letjoy-club/mida-tool/dbutil"
	"github.com/letjoy-club/mida-tool/logger"
	"github.com/letjoy-club/mida-tool/midacontext"
	"gorm.io/gorm"
	"time"
	"whale/pkg/dbquery"
	"whale/pkg/loader"
	"whale/pkg/models"
	"whale/pkg/utils"
)

// UpdateUserRights 更新用户等级权益
func UpdateUserRights(ctx context.Context, userId string, level int) error {
	if userId == "" || level == 0 {
		logger.L.Error("event payload illegal")
		return nil
	}

	config, err := GetLevelRightsConfig(ctx, level)
	if err != nil {
		logger.L.Error("can not find user level rights config")
		return err
	}

	db := dbutil.GetDB(ctx)
	MatchingQuota := dbquery.Use(db).MatchingQuota
	MatchingDurationConstraint := dbquery.Use(db).MatchingDurationConstraint
	currQuota, err := MatchingQuota.WithContext(ctx).Where(MatchingQuota.UserID.Eq(userId)).Take()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	currConstraint, err := MatchingDurationConstraint.WithContext(ctx).
		Where(MatchingDurationConstraint.UserID.Eq(userId),
			MatchingDurationConstraint.StartDate.Lte(time.Now()),
			MatchingDurationConstraint.StopDate.Gt(time.Now()),
		).Order(MatchingDurationConstraint.ID.Desc()).Take()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		// 更新同时匹配额度
		MatchingQuota := dbquery.Use(tx).MatchingQuota
		if currQuota == nil {
			if err := MatchingQuota.WithContext(ctx).Create(&models.MatchingQuota{
				UserID: userId,
				Remain: config.MatchingQuota,
				Total:  config.MatchingQuota,
			}); err != nil {
				return err
			}
		} else {
			quotaAdd := config.MatchingQuota - currQuota.Total
			if quotaAdd > 0 { // 额度只增不减
				if _, err := MatchingQuota.WithContext(ctx).Where(MatchingQuota.UserID.Eq(userId)).UpdateSimple(
					MatchingQuota.Total.Add(quotaAdd),
					MatchingQuota.Remain.Add(quotaAdd),
				); err != nil {
					return err
				}
			}
		}
		// 更新每周匹配次数额度
		MatchingDurationConstraint := dbquery.Use(tx).MatchingDurationConstraint
		if currConstraint == nil {
			startDate := utils.StartTimeOfWeek(time.Now())
			if err := MatchingDurationConstraint.WithContext(ctx).Create(&models.MatchingDurationConstraint{
				UserID:    userId,
				Total:     config.MatchingDurationConstraint,
				Remain:    config.MatchingDurationConstraint,
				StartDate: startDate,
				StopDate:  startDate.AddDate(0, 0, 7),
			}); err != nil {
				return err
			}
		} else {
			constraintAdd := config.MatchingDurationConstraint - currConstraint.Total
			if constraintAdd > 0 { // 额度只增不减
				if _, err := MatchingDurationConstraint.WithContext(ctx).
					Where(MatchingDurationConstraint.ID.Eq(currConstraint.ID)).
					UpdateSimple(
						MatchingDurationConstraint.Total.Add(constraintAdd),
						MatchingDurationConstraint.Remain.Add(constraintAdd),
					); err != nil {
					return err
				}
			}
		}
		return nil
	}); err != nil {
		return err
	}
	midacontext.GetLoader[loader.Loader](ctx).MatchingQuota.Clear(ctx, userId)
	midacontext.GetLoader[loader.Loader](ctx).MatchingDurationConstraint.Clear(ctx, userId)

	return nil
}

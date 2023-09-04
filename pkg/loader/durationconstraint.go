package loader

import (
	"context"
	"github.com/letjoy-club/mida-tool/midacontext"
	"time"
	"whale/pkg/dbquery"
	"whale/pkg/gqlient/hoopoe"
	"whale/pkg/models"

	"github.com/golang-module/carbon"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/letjoy-club/mida-tool/loaderutil"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

func NewDurationConstraintLoader(db *gorm.DB) *dataloader.Loader[string, *models.DurationConstraint] {
	DurationConstraint := dbquery.Use(db).DurationConstraint
	return loaderutil.NewItemLoader(db, func(ctx context.Context, userIDs []string) ([]*models.DurationConstraint, error) {
		return DurationConstraint.WithContext(ctx).Where(
			DurationConstraint.UserID.In(userIDs...),
			DurationConstraint.StopDate.Gt(time.Now()),
		).Find()
	}, func(m map[string]*models.DurationConstraint, v *models.DurationConstraint) {
		m[v.UserID] = v
	}, time.Minute, loaderutil.CreateIfNotFound(func(ctx context.Context, userIDs []string) ([]*models.DurationConstraint, []error) {
		// 查询用户信息
		services := midacontext.GetServices(ctx)
		resp, _ := hoopoe.GetUserByIDs(ctx, services.Hoopoe, userIDs)
		userMap := make(map[string]*hoopoe.GetUserByIDsGetUserByIdsUser)
		if resp != nil && resp.GetUserByIds != nil {
			for _, user := range resp.GetUserByIds {
				userMap[user.Id] = user
			}
		}
		// 查询用户等级权益配置
		thunk := midacontext.GetLoader[Loader](ctx).WhaleConfig.Load(ctx, models.ConfigLevelRights)
		conf, _ := thunk()
		levelConfig := &models.UserLevelConfig{}
		if conf != nil {
			levelConfig.Parse(conf)
		}
		rightsMap := make(map[int]*models.LevelRights)
		if levelConfig != nil {
			for _, levelRight := range levelConfig.Rights {
				rightsMap[levelRight.Level] = levelRight
			}
		}
		// 创建记录
		weekStart := carbon.Now().StartOfWeek().ToStdTime()
		weekEnd := carbon.Now().EndOfWeek().ToStdTime()
		constraints := lo.Map(userIDs, func(userID string, i int) *models.DurationConstraint {
			level := userMap[userID].Level
			levelRights, _ := lo.Find(levelConfig.Rights, func(item *models.LevelRights) bool {
				return item.Level == level
			})
			// todo: 默认值待确认
			motionQuota := 10
			offerQuota := 10
			if levelRights != nil && levelRights.MotionQuota != 0 {
				motionQuota = levelRights.MotionQuota
			}
			if levelRights != nil && levelRights.OfferQuota != 0 {
				offerQuota = levelRights.OfferQuota
			}
			return &models.DurationConstraint{
				UserID:            userID,
				StartDate:         weekStart,
				StopDate:          weekEnd,
				TotalMotionQuota:  motionQuota,
				RemainMotionQuota: motionQuota,
				TotalOfferQuota:   offerQuota,
				RemainOfferQuota:  offerQuota,
			}
		})
		err := DurationConstraint.WithContext(ctx).Create(constraints...)
		if err != nil {
			return nil, lo.Map(userIDs, func(userID string, i int) error { return err })
		}
		return constraints, nil
	}))
}

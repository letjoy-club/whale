package loader

import (
	"context"
	"github.com/letjoy-club/mida-tool/midacontext"
	"time"
	"whale/pkg/dbquery"
	"whale/pkg/gqlient/hoopoe"
	"whale/pkg/models"
	"whale/pkg/utils"

	"github.com/graph-gophers/dataloader/v7"
	"gorm.io/gorm"
)

func NewMatchingDurationConstraintLoader(db *gorm.DB) *dataloader.Loader[string, *models.MatchingDurationConstraint] {
	MatchingDurationConstraint := dbquery.Use(db).MatchingDurationConstraint
	return NewSingleLoader(db, func(ctx context.Context, keys []string) ([]*models.MatchingDurationConstraint, error) {
		constraints, err := MatchingDurationConstraint.WithContext(ctx).
			Where(MatchingDurationConstraint.UserID.In(keys...)).
			Where(MatchingDurationConstraint.StartDate.Lte(time.Now())).
			Where(MatchingDurationConstraint.StopDate.Gt(time.Now())).
			Find()
		if err != nil {
			return nil, err
		}

		constraintMap := map[string]struct{}{}
		for _, constraint := range constraints {
			constraintMap[constraint.UserID] = struct{}{}
		}

		notFoundIds := []string{}
		for _, key := range keys {
			if _, ok := constraintMap[key]; !ok {
				notFoundIds = append(notFoundIds, key)
			}
		}

		toBeAdded := []*models.MatchingDurationConstraint{}
		if len(notFoundIds) > 0 {
			// 查询用户信息
			services := midacontext.GetServices(ctx)
			resp, _ := hoopoe.GetUserByIDs(ctx, services.Hoopoe, notFoundIds)
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
			// 初始化记录
			startDate := utils.StartTimeOfWeek(time.Now())
			for _, id := range notFoundIds {
				user := userMap[id]
				level := 1
				constraint := 10
				if user != nil && user.Level > 0 { // 用户信息未找到，兜底使用
					level = user.Level
				}
				config := rightsMap[level]
				if config != nil && config.MatchingQuota > 0 {
					constraint = config.MatchingDurationConstraint
				}
				toBeAdded = append(toBeAdded, &models.MatchingDurationConstraint{
					UserID:    id,
					Remain:    constraint,
					Total:     constraint,
					StartDate: startDate,
					StopDate:  startDate.AddDate(0, 0, 7),
				})
			}
			if err = MatchingDurationConstraint.WithContext(ctx).Create(toBeAdded...); err != nil {
				return nil, err
			}
		}
		return append(constraints, toBeAdded...), nil
	},
		func(k map[string]*models.MatchingDurationConstraint, v *models.MatchingDurationConstraint) {
			k[v.UserID] = v
		}, time.Minute,
	)
}

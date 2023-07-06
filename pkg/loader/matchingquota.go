package loader

import (
	"context"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/letjoy-club/mida-tool/midacontext"
	"gorm.io/gorm"
	"time"
	"whale/pkg/dbquery"
	"whale/pkg/gqlient/hoopoe"
	"whale/pkg/models"
)

func NewMatchingQuotaLoader(db *gorm.DB) *dataloader.Loader[string, *models.MatchingQuota] {
	MatchingQuota := dbquery.Use(db).MatchingQuota
	return NewSingleLoader(db, func(ctx context.Context, keys []string) ([]*models.MatchingQuota, error) {
		matchingQuotas, err := MatchingQuota.WithContext(ctx).Where(MatchingQuota.UserID.In(keys...)).Find()
		if err != nil {
			return nil, err
		}

		matchingQuotaMap := map[string]struct{}{}
		for _, matchingQuota := range matchingQuotas {
			matchingQuotaMap[matchingQuota.UserID] = struct{}{}
		}

		notFoundIds := []string{}
		for _, key := range keys {
			if _, ok := matchingQuotaMap[key]; !ok {
				notFoundIds = append(notFoundIds, key)
			}
		}

		toBeAdded := []*models.MatchingQuota{}
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
			for _, id := range notFoundIds {
				user := userMap[id]
				level := 1
				quota := 3
				if user != nil && user.Level > 0 { // 用户信息未找到，兜底使用
					level = user.Level
				}
				config := rightsMap[level]
				if config != nil && config.MatchingQuota > 0 {
					quota = config.MatchingQuota
				}
				toBeAdded = append(toBeAdded, &models.MatchingQuota{
					UserID: id,
					Remain: quota,
					Total:  quota,
				})
			}
			if err := MatchingQuota.WithContext(ctx).Create(toBeAdded...); err != nil {
				return nil, err
			}
		}
		return append(matchingQuotas, toBeAdded...), nil
	}, func(k map[string]*models.MatchingQuota, v *models.MatchingQuota) { k[v.UserID] = v }, time.Minute)
}

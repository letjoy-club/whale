package modelutil

import (
	"context"
	"github.com/letjoy-club/mida-tool/midacode"
	"github.com/letjoy-club/mida-tool/midacontext"
	"whale/pkg/loader"
	"whale/pkg/models"
)

// GetUserLevelConfig 获取用户等级权益配置
func GetUserLevelConfig(ctx context.Context) (*models.UserLevelConfig, error) {
	thunk := midacontext.GetLoader[loader.Loader](ctx).WhaleConfig.Load(ctx, models.ConfigLevelRights)
	conf, err := thunk()
	if err != nil {
		return nil, err
	}

	config := &models.UserLevelConfig{}
	if err := config.Parse(conf); err != nil {
		return nil, err
	}

	return config, nil
}

// GetLevelRightsConfig 获取用户等级权益配置
func GetLevelRightsConfig(ctx context.Context, level int) (*models.LevelRights, error) {
	config, err := GetUserLevelConfig(ctx)
	if err != nil {
		return nil, err
	}

	for _, levelRights := range config.Rights {
		if levelRights.Level == level {
			return levelRights, nil
		}
	}

	return nil, midacode.ErrItemNotFound
}

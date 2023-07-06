package loader

import (
	"context"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/letjoy-club/mida-tool/midacode"
	"gorm.io/gorm"
	"time"
	"whale/pkg/dbquery"
	"whale/pkg/models"
)

func NewWhaleConfigLoader(db *gorm.DB) *dataloader.Loader[string, *models.WhaleConfig] {
	WhaleConfig := dbquery.Use(db).WhaleConfig
	return NewSingleLoader(db, func(ctx context.Context, keys []string) ([]*models.WhaleConfig, error) {
		var results []*models.WhaleConfig
		now := time.Now()
		for _, key := range keys {
			config, err := WhaleConfig.WithContext(ctx).Where(
				WhaleConfig.Name.Eq(key),
				WhaleConfig.StartAt.Lte(now),
				WhaleConfig.Enable.Is(true),
			).Order(WhaleConfig.StartAt.Desc()).Take()
			if err != nil {
				return nil, err
			}
			if config != nil && config.EndAt != nil && config.EndAt.Before(now) {
				return nil, midacode.ErrItemNotFound
			}
			results = append(results, config)
		}
		return results, nil
	}, func(k map[string]*models.WhaleConfig, v *models.WhaleConfig) { k[v.Name] = v }, time.Second*10)
}

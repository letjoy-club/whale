package loader

import (
	"context"
	"time"
	"whale/pkg/dbquery"
	"whale/pkg/models"

	"github.com/graph-gophers/dataloader/v7"
	"gorm.io/gorm"
)

func NewMatchingQuotaLoader(db *gorm.DB) *dataloader.Loader[string, *models.MatchingQuota] {
	MatchingQuota := dbquery.Use(db).MatchingQuota
	return NewSingleLoader(db, func(ctx context.Context, keys []string) ([]*models.MatchingQuota, error) {
		matchingQuotas, err := MatchingQuota.WithContext(ctx).Where(MatchingQuota.UserID.In(keys...)).Find()
		if err != nil {
			return nil, err
		}
		notFound := map[string]struct{}{}
		for _, id := range keys {
			notFound[id] = struct{}{}
		}
		for _, matchingQuota := range matchingQuotas {
			delete(notFound, matchingQuota.UserID)
		}
		toBeAdded := []*models.MatchingQuota{}
		if len(notFound) > 0 {
			for id := range notFound {
				toBeAdded = append(toBeAdded, &models.MatchingQuota{
					UserID: id,
					Remain: 3,
					Total:  3,
				})
			}
			if err := MatchingQuota.WithContext(ctx).Create(toBeAdded...); err != nil {
				return nil, err
			}
		}
		return append(matchingQuotas, toBeAdded...), nil
	}, func(k map[string]*models.MatchingQuota, v *models.MatchingQuota) { k[v.UserID] = v }, time.Minute)
}

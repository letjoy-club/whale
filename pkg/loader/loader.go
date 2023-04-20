package loader

import (
	"context"
	"time"
	"whale/pkg/dbquery"
	"whale/pkg/models"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/letjoy-club/mida-tool/midacode"
	"github.com/letjoy-club/mida-tool/ttlcache"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type Loader struct {
	Matching       *dataloader.Loader[string, *models.Matching]
	MatchingQuota  *dataloader.Loader[string, *models.MatchingQuota]
	MatchingResult *dataloader.Loader[int, *models.MatchingResult]
}

func NewLoader(db *gorm.DB) *Loader {
	return &Loader{
		Matching: NewSingleLoader(db, func(ctx context.Context, keys []string) ([]*models.Matching, error) {
			Matching := dbquery.Use(db).Matching
			return Matching.WithContext(ctx).Where(Matching.ID.In(keys...)).Find()
		}, func(k map[string]*models.Matching, v *models.Matching) {
			k[v.ID] = v
		}, time.Second*10),
		MatchingQuota: NewSingleLoader(db, func(ctx context.Context, keys []string) ([]*models.MatchingQuota, error) {
			MatchingQuota := dbquery.Use(db).MatchingQuota
			return MatchingQuota.WithContext(ctx).Where(MatchingQuota.UserID.In(keys...)).Find()
		}, func(k map[string]*models.MatchingQuota, v *models.MatchingQuota) {
			k[v.UserID] = v
		}, time.Second*60),
		MatchingResult: NewSingleLoader(db, func(ctx context.Context, keys []int) ([]*models.MatchingResult, error) {
			MatchingResult := dbquery.Use(db).MatchingResult
			return MatchingResult.WithContext(ctx).Where(MatchingResult.ID.In(keys...)).Find()
		}, func(k map[int]*models.MatchingResult, v *models.MatchingResult) {
			k[v.ID] = v
		}, time.Second*100),
	}
}

func NewSingleLoader[K comparable, V any](
	db *gorm.DB,
	loader func(ctx context.Context, keys []K) ([]V, error),
	dataMaper func(k map[K]V, v V),
	duration time.Duration,
) *dataloader.Loader[K, V] {
	c := ttlcache.New[K, V](duration)
	return dataloader.NewBatchedLoader(func(ctx context.Context, keys []K) []*dataloader.Result[V] {
		data, err := loader(ctx, keys)

		if err != nil {
			return lo.Map(data, func(m V, i int) *dataloader.Result[V] {
				return &dataloader.Result[V]{Error: err}
			})
		}

		dataMap := map[K]V{}
		for _, m := range data {
			dataMaper(dataMap, m)
		}

		return lo.Map(keys, func(key K, i int) *dataloader.Result[V] {
			ret, itemNotFound := dataMap[key]
			if itemNotFound {
				return &dataloader.Result[V]{Data: ret}
			} else {
				return &dataloader.Result[V]{Error: midacode.ErrItemNotFound}
			}
		})
	}, dataloader.WithCache[K, V](&c),
	)
}

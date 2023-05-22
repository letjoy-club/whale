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
	Matching           *dataloader.Loader[string, *models.Matching]
	MatchingInvitation *dataloader.Loader[string, *models.MatchingInvitation]
	MatchingQuota      *dataloader.Loader[string, *models.MatchingQuota]
	MatchingResult     *dataloader.Loader[int, *models.MatchingResult]
	MatchingReviewed   *dataloader.Loader[string, MatchingReviewed]

	CityTopicMatchings  *dataloader.Loader[CityTopicKey, CityTopicMatchings]
	CityTopicRequestNum *dataloader.Loader[string, CityTopicRequestNum]

	UserProfile        *dataloader.Loader[string, UserProfile]
	UserAvatarNickname *dataloader.Loader[string, UserAvatarNickname]
	HotTopics          *dataloader.Loader[string, *models.HotTopicsInArea]
}

func NewLoader(db *gorm.DB) *Loader {
	return &Loader{
		CityTopicMatchings:  NewCityTopicMatchingLoader(db),
		CityTopicRequestNum: NewCityTopicRequestNumLoader(db),

		Matching:           NewMatchingLoader(db),
		MatchingInvitation: NewMatchingInvitationLoader(db),
		MatchingQuota:      NewMatchingQuotaLoader(db),

		MatchingResult:     NEwMatchingResultLoader(db),
		MatchingReviewed:   NewMatchingReviewedLoader(db),
		UserProfile:        NewUserProfileLoader(db),
		UserAvatarNickname: NewUserAvatarNicknameLoader(db),
		HotTopics: NewSingleLoader(db, func(ctx context.Context, keys []string) ([]*models.HotTopicsInArea, error) {
			HotTopicsInArea := dbquery.Use(db).HotTopicsInArea
			topics, err := HotTopicsInArea.WithContext(ctx).Where(HotTopicsInArea.CityID.In(keys...)).Find()
			return topics, err
		}, func(k map[string]*models.HotTopicsInArea, v *models.HotTopicsInArea) { k[v.CityID] = v }, time.Second*60),
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

package loader

import (
	"context"
	"time"
	"whale/pkg/models"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/letjoy-club/mida-tool/midacode"
	"github.com/letjoy-club/mida-tool/ttlcache"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type Loader struct {
	Matching                   *dataloader.Loader[string, *models.Matching]
	MatchingInvitation         *dataloader.Loader[string, *models.MatchingInvitation]
	MatchingQuota              *dataloader.Loader[string, *models.MatchingQuota]
	MatchingResult             *dataloader.Loader[int, *models.MatchingResult]
	MatchingReviewed           *dataloader.Loader[string, MatchingReviewed]
	MatchingDurationConstraint *dataloader.Loader[string, *models.MatchingDurationConstraint]

	UserJoinTopic  *dataloader.Loader[int, *models.UserJoinTopic]
	RecentMatching *dataloader.Loader[string, *models.RecentMatching]

	// 从 recentMatching 中查询最近的 city, topic 对应的 matching id 信息
	CityTopicMatchings *dataloader.Loader[CityTopicKey, CityTopicMatchings]
	// 从 matching 表中获取最近的 topic 匹配中/已匹配数量
	CityTopicRequestNum *dataloader.Loader[string, CityTopicRequestNum]
	// 首屏的话题推荐
	CityTopics *dataloader.Loader[string, *models.CityTopics]

	UserProfile        *dataloader.Loader[string, UserProfile]
	UserAvatarNickname *dataloader.Loader[string, UserAvatarNickname]
	// 查询城市的热门话题
	HotTopics *dataloader.Loader[string, *models.HotTopicsInArea]
}

func NewLoader(db *gorm.DB) *Loader {
	return &Loader{
		CityTopicMatchings:  NewCityTopicMatchingLoader(db),
		CityTopicRequestNum: NewCityTopicRequestNumLoader(db),
		CityTopics:          NewCityTopicLoader(db),
		HotTopics:           NewHotTopicLoader(db),

		Matching:                   NewMatchingLoader(db),
		MatchingInvitation:         NewMatchingInvitationLoader(db),
		MatchingQuota:              NewMatchingQuotaLoader(db),
		MatchingResult:             NewMatchingResultLoader(db),
		MatchingReviewed:           NewMatchingReviewedLoader(db),
		MatchingDurationConstraint: NewMatchingDurationConstraintLoader(db),

		UserProfile:        NewUserProfileLoader(db),
		UserAvatarNickname: NewUserAvatarNicknameLoader(db),
		UserJoinTopic:      NewUserJoinTopicLoader(db),

		RecentMatching: NewRecentMatchingLoader(db),
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

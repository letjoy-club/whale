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
	AllMatching *AllMatchingLoader
	AllMotion   *AllMotionLoader

	Matching                   *dataloader.Loader[string, *models.Matching]
	MatchingInvitation         *dataloader.Loader[string, *models.MatchingInvitation]
	MatchingQuota              *dataloader.Loader[string, *models.MatchingQuota]
	MatchingResult             *dataloader.Loader[int, *models.MatchingResult]
	MatchingReviewed           *dataloader.Loader[string, MatchingReviewed]
	MatchingDurationConstraint *dataloader.Loader[string, *models.MatchingDurationConstraint]
	MatchingOfferSummary       *dataloader.Loader[string, *models.MatchingOfferSummary]

	UserJoinTopic  *dataloader.Loader[int, *models.UserJoinTopic]
	RecentMatching *dataloader.Loader[string, *models.RecentMatching]

	Motion               *dataloader.Loader[string, *models.Motion]
	UserLikeMotion       *dataloader.Loader[string, *UserLikeMotion]
	UserThumbsUpMotion   *dataloader.Loader[string, *UserThumbsUpMotions]
	UserSubmitMotion     *dataloader.Loader[string, *UserSubmitMotion]
	InMotionOfferRecord  *dataloader.Loader[string, *MotionOffers]
	OutMotionOfferRecord *dataloader.Loader[string, *MotionOffers]
	MotionReviewed       *dataloader.Loader[int, *MotionReviewed]

	DurationConstraint *dataloader.Loader[string, *models.DurationConstraint]

	// 从 recentMatching 中查询最近的 city, topic 对应的 matching id 信息
	CityTopicMatchings *dataloader.Loader[CityTopicKey, CityTopicMatchings]
	// 从 matching 表中获取最近的 topic 匹配中/已匹配数量
	CityTopicRequestNum *dataloader.Loader[string, CityTopicRequestNum]
	// 首屏的话题推荐
	CityTopics *dataloader.Loader[string, *models.CityTopics]

	UserProfile        *dataloader.Loader[string, UserProfile]
	UserAvatarNickname *dataloader.Loader[string, UserAvatarNickname]

	// 查询城市的热门话题
	HotTopics   *dataloader.Loader[string, *models.HotTopicsInArea]
	HotTopicsV2 *HotTopicV2Loader

	TopicOptionConfig *TopicOptionConfigLoader
	TopicCategory     *TopicCategoryLoader
	// 配置
	WhaleConfig *dataloader.Loader[string, *models.WhaleConfig]
}

func NewLoader(db *gorm.DB) *Loader {
	return &Loader{
		AllMatching: NewAllMatchingLoader(db),
		AllMotion:   NewAllMotionLoader(db),

		CityTopicMatchings:  NewCityTopicMatchingLoader(db),
		CityTopicRequestNum: NewCityTopicRequestNumLoader(db),
		CityTopics:          NewCityTopicLoader(db),
		HotTopics:           NewHotTopicLoader(db),
		HotTopicsV2:         NewHotTopicV2Loader(db),

		Matching:                   NewMatchingLoader(db),
		MatchingInvitation:         NewMatchingInvitationLoader(db),
		MatchingQuota:              NewMatchingQuotaLoader(db),
		MatchingResult:             NewMatchingResultLoader(db),
		MatchingReviewed:           NewMatchingReviewedLoader(db),
		MatchingDurationConstraint: NewMatchingDurationConstraintLoader(db),

		DurationConstraint: NewDurationConstraintLoader(db),

		InMotionOfferRecord:  NewInMotionOfferLoader(db),
		OutMotionOfferRecord: NewOutMotionOfferLoader(db),
		MotionReviewed:       NewMotionReviewedLoader(db),
		Motion:               NewMotionLoader(db),
		UserLikeMotion:       NewUserLikeMotionLoader(db),
		UserSubmitMotion:     NewUserSubmitMotionLoader(db),

		UserProfile:        NewUserProfileLoader(db),
		UserAvatarNickname: NewUserAvatarNicknameLoader(db),
		UserJoinTopic:      NewUserJoinTopicLoader(db),
		UserThumbsUpMotion: NewUserThumbsUpMotionLoader(db),

		RecentMatching: NewRecentMatchingLoader(db),

		TopicOptionConfig: NewTopicOptionConfigLoader(),
		TopicCategory:     NewTopicCategoryLoader(),

		WhaleConfig: NewWhaleConfigLoader(db),
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

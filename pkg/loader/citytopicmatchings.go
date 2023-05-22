package loader

import (
	"context"
	"time"
	"whale/pkg/dbquery"
	"whale/pkg/models"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/letjoy-club/mida-tool/dbutil"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CityTopicKey struct {
	CityID  string
	TopicID string
}

func (c CityTopicKey) Key() string {
	return c.CityID + "-" + c.TopicID
}

type CityTopicMatchings struct {
	CityTopicID string

	CityID  string
	TopicID string

	MatchingIDs []string
}

func (c CityTopicMatchings) CityTopicKey() CityTopicKey {
	return CityTopicKey{
		CityID:  c.CityID,
		TopicID: c.TopicID,
	}
}

func FetchCityTopicMatchingsFromDB(ctx context.Context, cityTopic CityTopicKey) (*models.RecentMatching, error) {
	db := dbutil.GetDB(ctx)

	UserJoinTopic := dbquery.Use(db).UserJoinTopic

	matchingIDs := []string{}

	err := UserJoinTopic.WithContext(ctx).
		Where(UserJoinTopic.CityID.Eq(cityTopic.CityID), UserJoinTopic.TopicID.Eq(cityTopic.TopicID)).
		Order(UserJoinTopic.CreatedAt.Desc()).
		Limit(16).
		Pluck(UserJoinTopic.LatestMatchingID, &matchingIDs)
	if err != nil {
		return nil, err
	}

	return &models.RecentMatching{
		ID:          cityTopic.Key(),
		CityID:      cityTopic.CityID,
		TopicID:     cityTopic.TopicID,
		MatchingIDs: matchingIDs,
	}, nil

}

func CreateCityTopicMatchings(ctx context.Context, recentMatchings []*models.RecentMatching) error {
	db := dbutil.GetDB(ctx)
	RecentMatching := dbquery.Use(db).RecentMatching
	return RecentMatching.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Create(recentMatchings...)
}

func NewCityTopicMatchingLoader(db *gorm.DB) *dataloader.Loader[CityTopicKey, CityTopicMatchings] {
	RecentMatching := dbquery.Use(db).RecentMatching
	return NewSingleLoader(db, func(ctx context.Context, keys []CityTopicKey) ([]CityTopicMatchings, error) {
		ids := lo.Map(keys, func(k CityTopicKey, i int) string {
			return k.Key()
		})
		recentMatchings, err := RecentMatching.WithContext(ctx).Where(RecentMatching.ID.In(ids...)).Find()
		if err != nil {
			return nil, err
		}
		existKeySet := map[string]struct{}{}
		for _, rm := range recentMatchings {
			existKeySet[rm.ID] = struct{}{}
		}
		notExistKeys := []CityTopicKey{}
		for _, key := range keys {
			if _, ok := existKeySet[key.Key()]; !ok {
				notExistKeys = append(notExistKeys, key)
			}
		}
		newRecentMatchings := []*models.RecentMatching{}
		for _, key := range notExistKeys {
			cityTopicMatching, err := FetchCityTopicMatchingsFromDB(ctx, key)
			if err != nil {
				return nil, err
			}
			newRecentMatchings = append(newRecentMatchings, &models.RecentMatching{
				ID:          key.Key(),
				CityID:      key.CityID,
				TopicID:     key.TopicID,
				MatchingIDs: cityTopicMatching.MatchingIDs,
			})
		}
		if len(newRecentMatchings) > 0 {
			CreateCityTopicMatchings(ctx, newRecentMatchings)
		}
		return lo.Map(append(recentMatchings, newRecentMatchings...), func(rm *models.RecentMatching, i int) CityTopicMatchings {
			return CityTopicMatchings{
				CityTopicID: rm.ID,
				CityID:      rm.CityID,
				TopicID:     rm.TopicID,
				MatchingIDs: rm.MatchingIDs,
			}
		}), nil
	}, func(k map[CityTopicKey]CityTopicMatchings, v CityTopicMatchings) { k[v.CityTopicKey()] = v }, time.Minute)
}

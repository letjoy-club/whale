package modelutil

import (
	"context"
	"fmt"
	"sort"
	"time"
	"whale/pkg/dbquery"
	"whale/pkg/gqlient/hoopoe"
	"whale/pkg/loader"
	"whale/pkg/models"

	"github.com/letjoy-club/mida-tool/dbutil"
	"github.com/letjoy-club/mida-tool/midacontext"
	"github.com/samber/lo"
	"gorm.io/gorm/clause"
)

func RefreshHotTopic(ctx context.Context, before time.Time) error {
	db := dbutil.GetDB(ctx)

	Matching := dbquery.Use(db).Matching

	matchings, err := Matching.WithContext(ctx).
		Where(Matching.CreatedAt.Gt(before)).
		Find()
	if err != nil {
		return err
	}

	matchingsOfCity := map[string][]*models.Matching{}
	lo.ForEach(matchings, func(m *models.Matching, i int) {
		matchingsOfCity[m.CityID] = append(matchingsOfCity[m.CityID], m)
	})

	hotTopics := []*models.HotTopicsInArea{}
	for cityId, matchings := range matchingsOfCity {
		tms := topicMetrics(matchings)
		hotTopics = append(hotTopics, &models.HotTopicsInArea{
			CityID:       cityId,
			TopicMetrics: tms,
		})
	}

	HotTopicsInArea := dbquery.Use(db).HotTopicsInArea
	err = HotTopicsInArea.
		WithContext(ctx).
		Clauses(clause.OnConflict{
			DoUpdates: clause.AssignmentColumns([]string{HotTopicsInArea.TopicMetrics.ColumnName().String()}),
		}).
		Create(hotTopics...)
		// Heat is the resolver for the heat field.
	fmt.Println("update/create hot topics:", len(hotTopics))
	return err
}

func topicMetrics(matchings []*models.Matching) []models.TopicMetrics {
	topicMetrics := map[string]models.TopicMetrics{}
	lo.ForEach(matchings, func(m *models.Matching, index int) {
		tm := topicMetrics[m.TopicID]
		tm.ID = m.TopicID
		if m.State == string(models.MatchingStateMatching) {
			tm.Matching++
		} else if m.State == string(models.MatchingStateMatched) || m.State == string(models.MatchingStateFinished) {
			tm.Matched++
		}
		topicMetrics[m.TopicID] = tm
	})

	lo.ForEach(matchings, func(m *models.Matching, i int) {
		tm := topicMetrics[m.TopicID]
		tm.Heat = 5*tm.Matched + 13*tm.Matching + 100
		topicMetrics[m.TopicID] = tm
	})

	tms := lo.Values(topicMetrics)
	sort.Slice(tms, func(i, j int) bool {
		return tms[i].Heat > tms[j].Heat
	})
	return tms
}

func HotTopicsInArea(ctx context.Context, cityID string) (*models.HotTopicsInArea, error) {
	thunk := midacontext.GetLoader[loader.Loader](ctx).HotTopics.Load(ctx, cityID)
	topics, _ := thunk()
	if topics != nil {
		return &models.HotTopicsInArea{
			CityID:       cityID,
			TopicMetrics: []models.TopicMetrics{},
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}, nil
	}
	return topics, nil
}

func RecordUserJoinTopic(ctx context.Context, topicID, cityID, userID, matchingID string) {
	db := dbutil.GetDB(ctx)
	UserJoinTopic := dbquery.Use(db).UserJoinTopic
	userjoined, err := UserJoinTopic.WithContext(ctx).Where(UserJoinTopic.TopicID.Eq(topicID)).Where(UserJoinTopic.CityID.Eq(cityID)).Where(UserJoinTopic.UserID.Eq(userID)).Take()
	if err != nil {
		// 没有找到记录
		UserJoinTopic.WithContext(ctx).Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(&models.UserJoinTopic{
			TopicID:          topicID,
			CityID:           cityID,
			UserID:           userID,
			LatestMatchingID: matchingID,
			Times:            1,
		})
	} else {
		// 已经有记录了
		userjoined.Times++
		userjoined.LatestMatchingID = matchingID
		UserJoinTopic.WithContext(ctx).Save(userjoined)
	}
}

func GetTopicName(ctx context.Context, topicID string) string {
	topic, err := hoopoe.GetTopicName(ctx, midacontext.GetServices(ctx).Hoopoe, topicID)
	if err != nil {
		fmt.Println("GetTopicName err:", err)
		return ""
	}
	return topic.Topic.Name
}

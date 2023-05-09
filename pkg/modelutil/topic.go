package modelutil

import (
	"context"
	"fmt"
	"sort"
	"time"
	"whale/pkg/dbquery"
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
	tms := lo.Values(topicMetrics)
	sort.Slice(tms, func(i, j int) bool {
		return tms[i].Total() > tms[j].Total()
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

package loader

import (
	"context"
	"time"
	"whale/pkg/dbquery"
	"whale/pkg/models"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type CityTopicRequestNum struct {
	CityID     string
	RequestNum map[string]int
}

type cityTopicRequest struct {
	CityID  string
	TopicID string
	Total   int
}

func NewCityTopicRequestNumLoader(db *gorm.DB) *dataloader.Loader[string, CityTopicRequestNum] {
	Matching := dbquery.Use(db).Matching
	return NewSingleLoader(db, func(ctx context.Context, cityIDs []string) ([]CityTopicRequestNum, error) {
		reqs := []cityTopicRequest{}
		err := Matching.WithContext(ctx).
			Select(Matching.CityID, Matching.TopicID, Matching.TopicID.Count().As("total")).
			Where(Matching.CityID.In(cityIDs...)).
			Where(Matching.State.In(models.MatchingStateMatching.String(), models.MatchingStateMatched.String())).
			Group(Matching.CityID, Matching.TopicID).
			Scan(&reqs)
		if err != nil {
			return nil, err
		}
		m := map[string]*CityTopicRequestNum{}
		for _, cityID := range cityIDs {
			m[cityID] = &CityTopicRequestNum{CityID: cityID, RequestNum: map[string]int{}}
		}
		for _, cityTopic := range reqs {
			m[cityTopic.CityID].RequestNum[cityTopic.TopicID] = cityTopic.Total
		}
		return lo.Map(cityIDs, func(cityID string, i int) CityTopicRequestNum {
			return *m[cityID]
		}), nil
	}, func(k map[string]CityTopicRequestNum, v CityTopicRequestNum) {
		k[v.CityID] = v
	}, time.Minute)
}

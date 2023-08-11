package loader

import (
	"context"
	"sync"
	"time"
	"whale/pkg/dbquery"
	"whale/pkg/gqlient/hoopoe"
	"whale/pkg/models"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/letjoy-club/mida-tool/midacontext"
	"gorm.io/gorm"
)

func NewHotTopicLoader(db *gorm.DB) *dataloader.Loader[string, *models.HotTopicsInArea] {
	return NewSingleLoader(db, func(ctx context.Context, keys []string) ([]*models.HotTopicsInArea, error) {
		HotTopicsInArea := dbquery.Use(db).HotTopicsInArea
		topics, err := HotTopicsInArea.WithContext(ctx).Where(HotTopicsInArea.CityID.In(keys...)).Find()
		return topics, err
	}, func(k map[string]*models.HotTopicsInArea, v *models.HotTopicsInArea) { k[v.CityID] = v }, time.Second*60)
}

func NewUserJoinTopicLoader(db *gorm.DB) *dataloader.Loader[int, *models.UserJoinTopic] {
	UserJoinTopic := dbquery.Use(db).UserJoinTopic
	return NewSingleLoader(db, func(ctx context.Context, keys []int) ([]*models.UserJoinTopic, error) {
		return UserJoinTopic.WithContext(ctx).Where(UserJoinTopic.ID.In(keys...)).Find()
	}, func(k map[int]*models.UserJoinTopic, v *models.UserJoinTopic) { k[v.ID] = v }, time.Second*10)
}

func NewCityTopicLoader(db *gorm.DB) *dataloader.Loader[string, *models.CityTopics] {
	return NewSingleLoader(db, func(ctx context.Context, keys []string) ([]*models.CityTopics, error) {
		CityTopics := dbquery.Use(db).CityTopics
		return CityTopics.WithContext(ctx).Where(CityTopics.CityID.In(keys...)).Find()
	}, func(k map[string]*models.CityTopics, v *models.CityTopics) { k[v.CityID] = v }, time.Second*60)
}

type TopicCategoryLoader struct {
	categoryHasTopicIDs map[string][]string
	topic2category      map[string]string

	lastUpdate time.Time
	mu         sync.Mutex
}

func (t *TopicCategoryLoader) Load(ctx context.Context) {
	if time.Since(t.lastUpdate) < time.Minute*5 {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	// double check
	if time.Since(t.lastUpdate) < time.Minute*5 {
		return
	}

	resp, err := hoopoe.GetTopics(ctx, midacontext.GetServices(ctx).Hoopoe, &hoopoe.GraphQLPaginator{
		Page: 1,
		Size: 9999,
	})
	if err != nil {
		return
	}

	categoryHasTopicIDs := make(map[string][]string)
	topic2category := make(map[string]string)
	for _, topic := range resp.Topics {
		categoryHasTopicIDs[topic.Category] = append(categoryHasTopicIDs[topic.Category], topic.Id)
		topic2category[topic.Id] = topic.Category
	}

	t.categoryHasTopicIDs = categoryHasTopicIDs
	t.topic2category = topic2category
	t.lastUpdate = time.Now()
}

func NewTopicCategoryLoader() *TopicCategoryLoader {
	return &TopicCategoryLoader{
		categoryHasTopicIDs: make(map[string][]string),
		topic2category:      make(map[string]string),
	}
}

func (t *TopicCategoryLoader) Topics(categoryID string) []string {
	return t.categoryHasTopicIDs[categoryID]
}

func (t *TopicCategoryLoader) Category(topicID string) string {
	return t.topic2category[topicID]
}

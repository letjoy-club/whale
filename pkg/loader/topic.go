package loader

import (
	"context"
	"math/rand"
	"sync"
	"time"
	"whale/pkg/dbquery"
	"whale/pkg/gqlient/hoopoe"
	"whale/pkg/models"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/letjoy-club/mida-tool/midacontext"
	"golang.org/x/exp/slices"
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

type HotTopicV2Loader struct {
	metrics    []models.TopicMetrics
	db         *gorm.DB
	lastUpdate time.Time
	mu         sync.Mutex
}

type TopicFreq struct {
	TopicID string
	Freq    int
}

// Load 热度数据逻辑, topic 卡片数量 * topic 总数 + random(1000)
func (h *HotTopicV2Loader) Load(ctx context.Context) error {
	if time.Since(h.lastUpdate) < time.Minute*15 {
		return nil
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	// double check
	if time.Since(h.lastUpdate) < time.Minute*15 {
		return nil
	}

	Motion := dbquery.Use(h.db).Motion
	topicFreqs := []TopicFreq{}
	err := Motion.WithContext(ctx).Select(Motion.TopicID, Motion.TopicID.Count().As("freq")).Group(Motion.TopicID).Scan(&topicFreqs)
	if err != nil {
		return err
	}

	slices.SortFunc(topicFreqs, func(a, b TopicFreq) int {
		return b.Freq - a.Freq
	})

	allNum := 0
	for _, topicFreq := range topicFreqs {
		allNum += topicFreq.Freq
	}

	metrics := make([]models.TopicMetrics, 0, 3)
	for ix, topicFreq := range topicFreqs {
		if ix >= 3 {
			break
		}

		metrics = append(metrics, models.TopicMetrics{
			ID:   topicFreq.TopicID,
			Heat: topicFreq.Freq*allNum + rand.Intn(1000),
		})
	}
	h.metrics = metrics
	return nil
}

func (h *HotTopicV2Loader) Metrics() []models.TopicMetrics {
	return h.metrics
}

func NewHotTopicV2Loader(db *gorm.DB) *HotTopicV2Loader {
	return &HotTopicV2Loader{
		db: db,
	}
}

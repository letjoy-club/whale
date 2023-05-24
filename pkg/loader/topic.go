package loader

import (
	"context"
	"time"
	"whale/pkg/dbquery"
	"whale/pkg/models"

	"github.com/graph-gophers/dataloader/v7"
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

package loader

import (
	"context"
	"time"
	"whale/pkg/dbquery"
	"whale/pkg/models"

	"github.com/graph-gophers/dataloader/v7"
	"gorm.io/gorm"
)

func NewMatchingInvitationLoader(db *gorm.DB) *dataloader.Loader[string, *models.MatchingInvitation] {
	MatchingInvitation := dbquery.Use(db).MatchingInvitation
	return NewSingleLoader(db, func(ctx context.Context, keys []string) ([]*models.MatchingInvitation, error) {
		return MatchingInvitation.WithContext(ctx).Where(MatchingInvitation.ID.In(keys...)).Find()
	}, func(k map[string]*models.MatchingInvitation, v *models.MatchingInvitation) { k[v.ID] = v }, time.Second*10)
}

func NewMatchingLoader(db *gorm.DB) *dataloader.Loader[string, *models.Matching] {
	Matching := dbquery.Use(db).Matching
	return NewSingleLoader(db, func(ctx context.Context, keys []string) ([]*models.Matching, error) {
		return Matching.WithContext(ctx).Where(Matching.ID.In(keys...)).Find()
	}, func(k map[string]*models.Matching, v *models.Matching) { k[v.ID] = v }, time.Second*10)
}

func NEwMatchingResultLoader(db *gorm.DB) *dataloader.Loader[int, *models.MatchingResult] {
	MatchingResult := dbquery.Use(db).MatchingResult
	return NewSingleLoader(db, func(ctx context.Context, keys []int) ([]*models.MatchingResult, error) {
		return MatchingResult.WithContext(ctx).Where(MatchingResult.ID.In(keys...)).Find()
	}, func(k map[int]*models.MatchingResult, v *models.MatchingResult) { k[v.ID] = v }, time.Second*100)
}

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

type MatchingReviewed struct {
	MatchingReviewed string
	Reviewed         bool
}

func NewMatchingReviewedLoader(db *gorm.DB) *dataloader.Loader[string, MatchingReviewed] {
	MatchingReview := dbquery.Use(db).MatchingReview
	return NewSingleLoader(db, func(ctx context.Context, keys []string) ([]MatchingReviewed, error) {
		reviews, err := MatchingReview.WithContext(ctx).Where(MatchingReview.MatchingID.In(keys...)).Find()
		if err != nil {
			return nil, err
		}
		return lo.Map(reviews, func(r *models.MatchingReview, i int) MatchingReviewed {
			return MatchingReviewed{MatchingReviewed: r.MatchingID, Reviewed: true}
		}), nil
	}, func(k map[string]MatchingReviewed, v MatchingReviewed) { k[v.MatchingReviewed] = v }, time.Second*30)
}

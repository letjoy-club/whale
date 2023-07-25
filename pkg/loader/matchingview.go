package loader

import (
	"context"
	"time"
	"whale/pkg/dbquery"
	"whale/pkg/models"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/letjoy-club/mida-tool/loaderutil"
	"gorm.io/gorm"
)

func NewMatchingViewLoader(db *gorm.DB) *dataloader.Loader[string, *models.MatchingView] {
	MatchingView := dbquery.Use(db).MatchingView
	return loaderutil.NewItemLoader(db, func(ctx context.Context, ids []string) ([]*models.MatchingView, error) {
		matchingViews, err := MatchingView.WithContext(ctx).Where(MatchingView.MatchingID.In(ids...)).Select(MatchingView.ViewCount, MatchingView.MatchingID).Find()
		if err != nil {
			return nil, err
		}
		return matchingViews, nil
	}, func(m map[string]*models.MatchingView, v *models.MatchingView) {
		m[v.MatchingID] = v
	}, time.Minute*5, loaderutil.Placeholder(func(ctx context.Context, id string) (*models.MatchingView, error) {
		return &models.MatchingView{MatchingID: id}, nil
	}),
	)
}

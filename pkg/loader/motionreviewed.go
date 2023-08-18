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

// 该对象实现的是以用户视角，查询自己的 motion 是否发出过评价（而不是用来查询是否被评价），已 review 的 motion offer 以 motion id 为维度进行聚合
type MotionReviewed struct {
	MotionID          string
	ReviewedMotionIDs []string
}

func (m *MotionReviewed) IsReviewed(motionID string) bool {
	for _, id := range m.ReviewedMotionIDs {
		if id == motionID {
			return true
		}
	}
	return false
}

func NewMotionReviewedLoader(db *gorm.DB) *dataloader.Loader[string, *MotionReviewed] {
	MotionReview := dbquery.Use(db).MotionReview
	return loaderutil.NewAggregatorLoader(db, func(ctx context.Context, keys []string) ([]*models.MotionReview, error) {
		reviews, err := MotionReview.WithContext(ctx).Where(MotionReview.MotionID.In(keys...)).Find()
		return reviews, err
	}, func(k map[string]*MotionReviewed, v *models.MotionReview) {
		if _, ok := k[v.MotionID]; !ok {
			k[v.MotionID] = &MotionReviewed{MotionID: v.MotionID, ReviewedMotionIDs: []string{}}
		} else {
			k[v.MotionID].ReviewedMotionIDs = append(k[v.MotionID].ReviewedMotionIDs, v.MotionID)
		}
	}, time.Minute*5, loaderutil.Placeholder(func(ctx context.Context, id string) (*MotionReviewed, error) {
		return &MotionReviewed{MotionID: id, ReviewedMotionIDs: []string{}}, nil
	}))
}

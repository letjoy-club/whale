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

// 该对象实现的是以用户视角，查询自己在 motion offer 是否发出过评价（而不是用来查询是否被评价）
type MotionReviewed struct {
	MotionOfferID  int
	ReviewedUserID []string
}

func (m *MotionReviewed) IsReviewed(userID string) bool {
	for _, id := range m.ReviewedUserID {
		if id == userID {
			return true
		}
	}
	return false
}

func NewMotionReviewedLoader(db *gorm.DB) *dataloader.Loader[int, *MotionReviewed] {
	MotionReview := dbquery.Use(db).MotionReview
	return loaderutil.NewAggregatorLoader(db, func(ctx context.Context, keys []int) ([]*models.MotionReview, error) {
		reviews, err := MotionReview.WithContext(ctx).Where(MotionReview.MotionOfferID.In(keys...)).Find()
		return reviews, err
	}, func(k map[int]*MotionReviewed, v *models.MotionReview) {
		if _, ok := k[v.MotionOfferID]; !ok {
			k[v.MotionOfferID] = &MotionReviewed{MotionOfferID: v.MotionOfferID, ReviewedUserID: []string{}}
		} else {
			k[v.MotionOfferID].ReviewedUserID = append(k[v.MotionOfferID].ReviewedUserID, v.ReviewerID)
		}
	}, time.Minute*5, loaderutil.Placeholder(func(ctx context.Context, id int) (*MotionReviewed, error) {
		return &MotionReviewed{MotionOfferID: id, ReviewedUserID: []string{}}, nil
	}))
}

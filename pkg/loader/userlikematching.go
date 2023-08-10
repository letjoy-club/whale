package loader

import (
	"context"
	"time"
	"whale/pkg/models"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/letjoy-club/mida-tool/loaderutil"
	"gorm.io/gorm"
)

type UserThumbsUpMotions struct {
	userID    string
	MotionIDs []string
}

func (u *UserThumbsUpMotions) ThumbsUpped(motionID string) bool {
	for _, id := range u.MotionIDs {
		if id == motionID {
			return true
		}
	}
	return false
}

func NewUserThumbsUpMotionLoader(db *gorm.DB) *dataloader.Loader[string, *UserThumbsUpMotions] {
	return loaderutil.NewAggregatorLoader(db, func(ctx context.Context, keys []string) ([]*models.UserThumbsUpMotion, error) {
		return nil, nil
	}, func(m map[string]*UserThumbsUpMotions, v *models.UserThumbsUpMotion) {}, time.Minute, loaderutil.Placeholder(func(ctx context.Context, id string) (*UserThumbsUpMotions, error) {
		return nil, nil
	}))
}

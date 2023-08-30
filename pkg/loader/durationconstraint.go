package loader

import (
	"context"
	"time"
	"whale/pkg/dbquery"
	"whale/pkg/models"

	"github.com/golang-module/carbon"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/letjoy-club/mida-tool/loaderutil"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

func NewDurationConstraintLoader(db *gorm.DB) *dataloader.Loader[string, *models.DurationConstraint] {
	DurationConstraint := dbquery.Use(db).DurationConstraint
	return loaderutil.NewItemLoader(db, func(ctx context.Context, userIDs []string) ([]*models.DurationConstraint, error) {
		return DurationConstraint.WithContext(ctx).Where(DurationConstraint.UserID.In(userIDs...), DurationConstraint.StopDate.Gt(time.Now())).Find()
	}, func(m map[string]*models.DurationConstraint, v *models.DurationConstraint) {
		m[v.UserID] = v
	}, time.Minute, loaderutil.CreateIfNotFound(func(ctx context.Context, userIDs []string) ([]*models.DurationConstraint, []error) {
		weekStart := carbon.Now().StartOfWeek().ToStdTime()
		weekEnd := carbon.Now().EndOfWeek().ToStdTime()
		constraints := lo.Map(userIDs, func(userID string, i int) *models.DurationConstraint {
			return &models.DurationConstraint{
				UserID:            userID,
				StartDate:         weekStart,
				StopDate:          weekEnd,
				TotalMotionQuota:  20,
				RemainMotionQuota: 20,
			}
		})
		err := DurationConstraint.WithContext(ctx).Create(constraints...)
		if err != nil {
			return nil, lo.Map(userIDs, func(userID string, i int) error { return err })
		}
		return constraints, nil
	}))
}

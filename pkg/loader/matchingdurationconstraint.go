package loader

import (
	"context"
	"time"
	"whale/pkg/dbquery"
	"whale/pkg/models"

	"github.com/graph-gophers/dataloader/v7"
	"gorm.io/gorm"
)

func NewMatchingDurationConstraintLoader(db *gorm.DB) *dataloader.Loader[string, *models.MatchingDurationConstraint] {
	MatchingDurationConstraint := dbquery.Use(db).MatchingDurationConstraint
	return NewSingleLoader(db, func(ctx context.Context, keys []string) ([]*models.MatchingDurationConstraint, error) {
		constraints, err := MatchingDurationConstraint.WithContext(ctx).
			Where(MatchingDurationConstraint.UserID.In(keys...)).
			Where(MatchingDurationConstraint.StartDate.Lte(time.Now())).
			Where(MatchingDurationConstraint.StopDate.Gt(time.Now())).
			Find()
		if err != nil {
			return nil, err
		}
		existed := map[string]struct{}{}
		for _, constraint := range constraints {
			existed[constraint.UserID] = struct{}{}
		}
		toBeAdded := []*models.MatchingDurationConstraint{}
		for _, id := range keys {
			if _, ok := existed[id]; !ok {
				toBeAdded = append(toBeAdded, &models.MatchingDurationConstraint{
					UserID:    id,
					Total:     10,
					Remain:    10,
					StartDate: time.Now(),
					StopDate:  time.Now().Add(time.Hour * 24 * 7),
				})
			}
		}
		err = MatchingDurationConstraint.WithContext(ctx).Create(toBeAdded...)
		if err != nil {
			return nil, err
		}
		constraints = append(constraints, toBeAdded...)
		return constraints, err
	},
		func(k map[string]*models.MatchingDurationConstraint, v *models.MatchingDurationConstraint) {
			k[v.UserID] = v
		}, time.Minute,
	)
}

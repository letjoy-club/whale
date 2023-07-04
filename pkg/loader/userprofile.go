package loader

import (
	"context"
	"time"
	"whale/pkg/gqlient/hoopoe"
	"whale/pkg/models"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/letjoy-club/mida-tool/midacontext"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type UserProfile struct {
	ID     string
	Gender models.Gender
}

func NewUserProfileLoader(db *gorm.DB) *dataloader.Loader[string, UserProfile] {
	return NewSingleLoader(db, func(ctx context.Context, keys []string) ([]UserProfile, error) {
		ret, err := hoopoe.GetUserByIDs(ctx, midacontext.GetServices(ctx).Hoopoe, keys)
		if err != nil {
			return nil, err
		}
		return lo.Map(ret.GetUserByIds, func(u *hoopoe.GetUserByIDsGetUserByIdsUser, i int) UserProfile {
			return UserProfile{ID: u.Id, Gender: models.Gender(u.Gender)}
		}), nil
	}, func(k map[string]UserProfile, v UserProfile) { k[v.ID] = v }, time.Minute)
}

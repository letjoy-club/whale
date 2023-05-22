package loader

import (
	"context"
	"time"
	"whale/pkg/gqlient/hoopoe"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/letjoy-club/mida-tool/midacontext"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type UserAvatarNickname struct {
	ID       string
	Avatar   string
	Nickname string
}

func NewUserAvatarNicknameLoader(db *gorm.DB) *dataloader.Loader[string, UserAvatarNickname] {
	return NewSingleLoader(db, func(ctx context.Context, keys []string) ([]UserAvatarNickname, error) {
		ret, err := hoopoe.GetAvatarByIDs(ctx, midacontext.GetServices(ctx).Hoopoe, keys)
		if err != nil {
			return nil, err
		}
		return lo.Map(ret.GetUserByIds, func(u hoopoe.GetAvatarByIDsGetUserByIdsUser, i int) UserAvatarNickname {
			return UserAvatarNickname{ID: u.Id, Avatar: u.Avatar, Nickname: u.Nickname}
		}), nil
	}, func(k map[string]UserAvatarNickname, v UserAvatarNickname) { k[v.ID] = v }, time.Minute)
}

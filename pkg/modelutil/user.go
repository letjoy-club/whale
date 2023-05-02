package modelutil

import (
	"context"
	"whale/pkg/gqlient/hoopoe"

	"github.com/letjoy-club/mida-tool/midacontext"
)

func QueryUserByIDs(ctx context.Context, ids []string) (map[string]hoopoe.GetUserByIDsGetUserByIdsUser, error) {

	services := midacontext.GetServices(ctx)

	resp, err := hoopoe.GetUserByIDs(ctx, services.Hoopoe, ids)
	if err != nil {
		return nil, err
	}
	m := map[string]hoopoe.GetUserByIDsGetUserByIdsUser{}
	for _, user := range resp.GetGetUserByIds() {
		m[user.Id] = hoopoe.GetUserByIDsGetUserByIdsUser{
			Id:     user.Id,
			Gender: user.Gender,
		}
	}
	return m, nil
}

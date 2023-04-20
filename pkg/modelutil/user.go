package modelutil

import (
	"context"

	"github.com/letjoy-club/mida-tool/midacontext"
)

type UserAndGender struct {
	ID     string
	Gender string
}

func QueryUserByIDs(ctx context.Context, ids []string) (map[string]UserAndGender, error) {
	var q struct {
		User []struct {
			ID     string
			Gender string
		} `graphql:"userByIDs(ids: $ids)"`
	}

	services := midacontext.GetServices(ctx)

	err := services.Hoopoe.Query(context.Background(), &q, map[string]interface{}{
		"ids": ids,
	})
	if err != nil {
		return nil, err
	}
	m := map[string]UserAndGender{}
	for _, user := range q.User {
		m[user.ID] = UserAndGender{
			ID:     user.ID,
			Gender: user.Gender,
		}
	}
	return m, nil
}

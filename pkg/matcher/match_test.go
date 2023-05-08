package matcher_test

import (
	"context"
	"fmt"
	"testing"
	"whale/pkg/loader"
	"whale/pkg/matcher"
	"whale/pkg/models"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/letjoy-club/mida-tool/midacontext"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestMatched(t *testing.T) {
	ctx := context.Background()
	u1Female := loader.UserProfile{ID: "u_1", Gender: "F"}
	u2Unknown := loader.UserProfile{ID: "u_2", Gender: ""}
	u3Male := loader.UserProfile{ID: "u_3", Gender: "M"}
	u4Male := loader.UserProfile{ID: "u_4", Gender: "M"}
	u5Female := loader.UserProfile{ID: "u_5", Gender: "F"}

	userMaps := map[string]loader.UserProfile{
		u1Female.ID:  u1Female,
		u2Unknown.ID: u2Unknown,
		u3Male.ID:    u3Male,
		u4Male.ID:    u4Male,
		u5Female.ID:  u5Female,
	}
	m1 := models.Matching{
		ID:      "m_1",
		UserID:  u1Female.ID,
		AreaIDs: []string{"310101"},
		Gender:  models.GenderN.String(),
	}
	m2 := models.Matching{
		ID:      "m_2",
		UserID:  u2Unknown.ID,
		AreaIDs: []string{"310101"},
		Gender:  models.GenderM.String(),
	}
	m3 := models.Matching{
		ID:      "m_3",
		UserID:  u3Male.ID,
		AreaIDs: []string{"310100"},
		Gender:  models.GenderF.String(),
	}
	m4 := models.Matching{
		ID:      "m_4",
		UserID:  u4Male.ID,
		AreaIDs: []string{"310100"},
		Gender:  models.GenderM.String(),
	}
	m5 := models.Matching{
		ID:      "m_5",
		UserID:  u5Female.ID,
		AreaIDs: []string{"310100"},
		Gender:  models.GenderM.String(),
	}
	ctx = midacontext.WithLoader(ctx, &loader.Loader{
		UserProfile: dataloader.NewBatchedLoader(func(ctx context.Context, k []string) []*dataloader.Result[loader.UserProfile] {
			return lo.Map(k, func(uid string, i int) *dataloader.Result[loader.UserProfile] {
				return &dataloader.Result[loader.UserProfile]{Data: userMaps[uid]}
			})
		}),
	})

	ctx = matcher.WithMatchingContext(ctx, []*models.Matching{&m1, &m2, &m3, &m4, &m5})
	{
		reason, matched := matcher.Matched(ctx, &m1, &m2)
		assert.False(t, matched)
		fmt.Println(reason)
	}
	{
		reason, matched := matcher.Matched(ctx, &m1, &m3)
		assert.False(t, matched)
		fmt.Println(reason)
	}
	{
		reason, matched := matcher.Matched(ctx, &m3, &m4)
		assert.False(t, matched)
		fmt.Println(reason)
	}
	{
		reason, matched := matcher.Matched(ctx, &m3, &m5)
		assert.True(t, matched)
		fmt.Println(reason)
	}
}

package matcher

import (
	"context"
	"whale/pkg/loader"
	"whale/pkg/models"

	"github.com/letjoy-club/mida-tool/midacontext"
	"github.com/samber/lo"
)

type MatchingContext struct {
	topic2matchings map[string][]*models.Matching
	used            map[string]bool
	topics          []string
	userProfiles    map[string]loader.UserProfile
}

func (mc *MatchingContext) TopicMatchings(topicID string) []*models.Matching {
	return mc.topic2matchings[topicID]
}

func (mc *MatchingContext) Topics() []string {
	return mc.topics
}

func (mc *MatchingContext) Use(matchingID string) {
	mc.used[matchingID] = true
}

func (mc *MatchingContext) Used(matchingID string) bool {
	return mc.used[matchingID]
}

type mcKey struct{}

func GetMatchingContext(ctx context.Context) *MatchingContext {
	return ctx.Value(mcKey{}).(*MatchingContext)
}

func WithMatchingContext(ctx context.Context, matchings []*models.Matching) context.Context {
	topicMap := make(map[string][]*models.Matching)
	userMap := map[string]struct{}{}

	for _, matching := range matchings {
		topicMap[matching.TopicID] = append(topicMap[matching.TopicID], matching)
		userMap[matching.UserID] = struct{}{}
	}

	// 预加载用户信息
	userIDs := lo.Keys(userMap)
	thunk := midacontext.GetLoader[loader.Loader](ctx).UserProfile.LoadMany(ctx, userIDs)
	users, _ := thunk()

	userProfiles := map[string]loader.UserProfile{}
	for _, u := range users {
		userProfiles[u.ID] = u
	}

	mc := &MatchingContext{
		topic2matchings: topicMap,
		used:            make(map[string]bool),
		topics:          lo.Keys(topicMap),
		userProfiles:    userProfiles,
	}

	return context.WithValue(ctx, mcKey{}, mc)
}

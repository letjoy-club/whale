package matcher

import (
	"context"
	"whale/pkg/models"

	"github.com/samber/lo"
)

type MatchingContext struct {
	topic2matchings map[int][]*models.Matching
	used            map[string]bool
	topics          []int
}

func (mc *MatchingContext) TopicMatchings(topicID int) []*models.Matching {
	return mc.topic2matchings[topicID]
}

func (mc *MatchingContext) Topics() []int {
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
	topicMap := make(map[int][]*models.Matching)

	for _, matching := range matchings {
		topicMap[matching.TopicID] = append(topicMap[matching.TopicID], matching)
	}

	mc := &MatchingContext{
		topic2matchings: topicMap,
		used:            make(map[string]bool),
		topics:          lo.Keys(topicMap),
	}
	return context.WithValue(ctx, mcKey{}, mc)
}

package modelutil

import (
	"context"
	"fmt"
	"time"
	"whale/pkg/models"
	"whale/pkg/whaleconf"

	"github.com/letjoy-club/mida-tool/pulsarutil"
)

type matchingCanceled struct {
	MatchingID string    `json:"matchingId"`
	UserID     string    `json:"userId"`
	TopicID    string    `json:"topicId"`
	CityID     string    `json:"cityId"`
	CreatedAt  time.Time `json:"createdAt"`
	CanceledAt time.Time `json:"canceledAt"`
}

func PublishMatchingCanceledEvent(ctx context.Context, matching *models.Matching) error {
	pub := pulsarutil.GetMQ[*whaleconf.Publisher](ctx)
	return pub.Pub(ctx, matching.ID+":canceled", matchingCanceled{
		MatchingID: matching.ID,
		UserID:     matching.UserID,
		TopicID:    matching.TopicID,
		CityID:     matching.CityID,
		CreatedAt:  matching.CreatedAt,
		CanceledAt: time.Now(),
	})
}

type matchingCreated struct {
	MatchingID string    `json:"matchingId"`
	UserID     string    `json:"userId"`
	TopicID    string    `json:"topicId"`
	CityID     string    `json:"cityId"`
	CreatedAt  time.Time `json:"createdAt"`
}

func PublishMatchingCreatedEvent(ctx context.Context, matching *models.Matching) error {
	pub := pulsarutil.GetMQ[*whaleconf.Publisher](ctx)
	return pub.Pub(ctx, matching.ID+":created", matchingCreated{
		MatchingID: matching.ID,
		UserID:     matching.UserID,
		TopicID:    matching.TopicID,
		CityID:     matching.CityID,
		CreatedAt:  matching.CreatedAt,
	})
}

type matchingMatched struct {
	MatchingID string    `json:"matchingId"`
	TopicID    string    `json:"topicId"`
	UserID     string    `json:"userId"`
	CityID     string    `json:"cityId"`
	MatchedAt  time.Time `json:"matchedAt"`
	CreatedAt  time.Time `json:"createdAt"`
}

func PublishMatchedEvent(ctx context.Context, matching *models.Matching) error {
	pub := pulsarutil.GetMQ[*whaleconf.Publisher](ctx)
	return pub.Pub(ctx, matching.ID+":matched", matchingMatched{
		MatchingID: matching.ID,
		UserID:     matching.UserID,
		TopicID:    matching.TopicID,
		CityID:     matching.CityID,
		CreatedAt:  matching.CreatedAt,
	})
}

type matchingTimemout struct {
	MatchingID string    `json:"matchingId"`
	UserID     string    `json:"userId"`
	TopicID    string    `json:"topicId"`
	CityID     string    `json:"cityId"`
	CreatedAt  time.Time `json:"createdAt"`
	TimeoutAt  time.Time `json:"timeoutAt"`
}

func PublishMatchingTimeoutEvent(ctx context.Context, matching *models.Matching) error {
	mq := pulsarutil.GetMQ[*whaleconf.Publisher](ctx)
	return mq.Pub(ctx, matching.ID+":timeout", matchingTimemout{
		MatchingID: matching.ID,
		UserID:     matching.UserID,
		TopicID:    matching.TopicID,
		CityID:     matching.CityID,
		CreatedAt:  matching.CreatedAt,
		TimeoutAt:  matching.Deadline,
	})
}

type matchingFinished struct {
	MatchingID string    `json:"matchingId"`
	UserID     string    `json:"userId"`
	TopicID    string    `json:"topicId"`
	CityID     string    `json:"cityId"`
	CreatedAt  time.Time `json:"createdAt"`
	FinishedAt time.Time `json:"finishedAt"`
	CreatedBy  string    `json:"createdBy"`
}

func PublishMatchingFinishedEvent(ctx context.Context, matching *models.Matching, createdBy string) error {
	mq := pulsarutil.GetMQ[*whaleconf.Publisher](ctx)
	if mq == nil {
		fmt.Println("skip due to no mq")
		return nil
	}
	if createdBy == "" {
		createdBy = models.ResultCreatedByMatching.String()
	}
	return mq.Pub(ctx, matching.ID+":finished", matchingFinished{
		MatchingID: matching.ID,
		UserID:     matching.UserID,
		TopicID:    matching.TopicID,
		CreatedAt:  matching.CreatedAt,
		FinishedAt: *matching.FinishedAt,
		CityID:     matching.CityID,
		CreatedBy:  createdBy,
	})
}

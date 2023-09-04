package modelutil

import (
	"context"
	"time"
	"whale/pkg/models"
	"whale/pkg/whaleconf"

	"github.com/letjoy-club/mida-tool/pulsarutil"
)

type motionCreated struct {
	MotionID  string    `json:"motionId"`
	UserID    string    `json:"userId"`
	TopicID   string    `json:"topicId"`
	CityID    string    `json:"cityId"`
	CreatedAt time.Time `json:"createdAt"`
}

func PublishMotionCreatedEvent(ctx context.Context, motion *models.Motion) error {
	pub := pulsarutil.GetMQ[*whaleconf.Publisher](ctx)
	return pub.Pub(ctx, "MotionCreated", motionCreated{
		MotionID:  motion.ID,
		UserID:    motion.UserID,
		TopicID:   motion.TopicID,
		CityID:    motion.CityID,
		CreatedAt: motion.CreatedAt,
	})
}

type motionClosed struct {
	MotionID  string    `json:"motionId"`
	UserID    string    `json:"userId"`
	TopicID   string    `json:"topicId"`
	CityID    string    `json:"cityId"`
	CreatedAt time.Time `json:"createdAt"`
	ClosedAt  time.Time `json:"closedAt"`
}

func PublishMotionClosedEvent(ctx context.Context, motion *models.Motion) error {
	pub := pulsarutil.GetMQ[*whaleconf.Publisher](ctx)
	return pub.Pub(ctx, "MotionClosed", motionClosed{
		MotionID:  motion.ID,
		UserID:    motion.UserID,
		TopicID:   motion.TopicID,
		CityID:    motion.CityID,
		CreatedAt: motion.CreatedAt,
		ClosedAt:  time.Now(),
	})
}

type motionOfferCreated struct {
	MotionOfferID int       `json:"motionOfferId"`
	FromMotionID  string    `json:"fromMotionId"`
	ToMotionID    string    `json:"toMotionId"`
	FromUserID    string    `json:"fromUserId"`
	ToUserID      string    `json:"toUserId"`
	CreatedAt     time.Time `json:"createdAt"`
}

func PublishMotionOfferCreatedEvent(ctx context.Context, record *models.MotionOfferRecord) error {
	pub := pulsarutil.GetMQ[*whaleconf.Publisher](ctx)
	return pub.Pub(ctx, "MotionOfferCreated", motionOfferCreated{
		MotionOfferID: record.ID,
		FromMotionID:  record.MotionID,
		ToMotionID:    record.ToMotionID,
		FromUserID:    record.UserID,
		ToUserID:      record.ToUserID,
		CreatedAt:     record.CreatedAt,
	})
}

type motionOfferAccepted struct {
	MotionOfferID int       `json:"motionOfferId"`
	FromMotionID  string    `json:"fromMotionId"`
	ToMotionID    string    `json:"toMotionId"`
	FromUserID    string    `json:"fromUserId"`
	ToUserID      string    `json:"toUserId"`
	CreatedAt     time.Time `json:"createdAt"`
	AcceptedAt    time.Time `json:"acceptedAt"`
}

func PublishMotionOfferAcceptedEvent(ctx context.Context, record *models.MotionOfferRecord) error {
	pub := pulsarutil.GetMQ[*whaleconf.Publisher](ctx)
	return pub.Pub(ctx, "MotionOfferAccepted", motionOfferAccepted{
		MotionOfferID: record.ID,
		FromMotionID:  record.MotionID,
		ToMotionID:    record.ToMotionID,
		FromUserID:    record.UserID,
		ToUserID:      record.ToUserID,
		CreatedAt:     record.CreatedAt,
		AcceptedAt:    time.Now(),
	})
}

type motionOfferRejected struct {
	MotionOfferID int       `json:"motionOfferId"`
	FromMotionID  string    `json:"fromMotionId"`
	ToMotionID    string    `json:"toMotionId"`
	FromUserID    string    `json:"fromUserId"`
	ToUserID      string    `json:"toUserId"`
	CreatedAt     time.Time `json:"createdAt"`
	RejectedAt    time.Time `json:"rejectedAt"`
}

func PublishMotionOfferRejectedEvent(ctx context.Context, record *models.MotionOfferRecord) error {
	pub := pulsarutil.GetMQ[*whaleconf.Publisher](ctx)
	return pub.Pub(ctx, "MotionOfferRejected", motionOfferRejected{
		MotionOfferID: record.ID,
		FromMotionID:  record.MotionID,
		ToMotionID:    record.ToMotionID,
		FromUserID:    record.UserID,
		ToUserID:      record.ToUserID,
		CreatedAt:     record.CreatedAt,
		RejectedAt:    time.Now(),
	})
}

type motionOfferCanceled struct {
	MotionOfferID int       `json:"motionOfferId"`
	FromMotionID  string    `json:"fromMotionId"`
	ToMotionID    string    `json:"toMotionId"`
	FromUserID    string    `json:"fromUserId"`
	ToUserID      string    `json:"toUserId"`
	CreatedAt     time.Time `json:"createdAt"`
	CanceledAt    time.Time `json:"canceledAt"`
}

func PublishMotionOfferCanceledEvent(ctx context.Context, record *models.MotionOfferRecord) error {
	pub := pulsarutil.GetMQ[*whaleconf.Publisher](ctx)
	return pub.Pub(ctx, "MotionOfferCanceled", motionOfferCanceled{
		MotionOfferID: record.ID,
		FromMotionID:  record.MotionID,
		ToMotionID:    record.ToMotionID,
		FromUserID:    record.UserID,
		ToUserID:      record.ToUserID,
		CreatedAt:     record.CreatedAt,
		CanceledAt:    time.Now(),
	})
}

type motionOfferTimeout struct {
	MotionOfferID int       `json:"motionOfferId"`
	FromMotionID  string    `json:"fromMotionId"`
	ToMotionID    string    `json:"toMotionId"`
	FromUserID    string    `json:"fromUserId"`
	ToUserID      string    `json:"toUserId"`
	CreatedAt     time.Time `json:"createdAt"`
	TimeoutAt     time.Time `json:"timeoutAt"`
}

func PublishMotionOfferTimeoutEvent(ctx context.Context, record *models.MotionOfferRecord) error {
	pub := pulsarutil.GetMQ[*whaleconf.Publisher](ctx)
	return pub.Pub(ctx, "MotionOfferTimeout", motionOfferTimeout{
		MotionOfferID: record.ID,
		FromMotionID:  record.MotionID,
		ToMotionID:    record.ToMotionID,
		FromUserID:    record.UserID,
		ToUserID:      record.ToUserID,
		CreatedAt:     record.CreatedAt,
		TimeoutAt:     time.Now(),
	})
}

type motionOfferFinished struct {
	MotionOfferID int       `json:"motionOfferId"`
	FromMotionID  string    `json:"fromMotionId"`
	ToMotionID    string    `json:"toMotionId"`
	FromUserID    string    `json:"fromUserId"`
	ToUserID      string    `json:"toUserId"`
	CreatedAt     time.Time `json:"createdAt"`
	FinishedBy    string    `json:"finishedBy"`
	FinishedAt    time.Time `json:"finishedAt"`
}

func PublishMotionOfferFinishedEvent(ctx context.Context, record *models.MotionOfferRecord, finishedBy string) error {
	pub := pulsarutil.GetMQ[*whaleconf.Publisher](ctx)
	return pub.Pub(ctx, "MotionOfferFinished", motionOfferFinished{
		MotionOfferID: record.ID,
		FromMotionID:  record.MotionID,
		ToMotionID:    record.ToMotionID,
		FromUserID:    record.UserID,
		ToUserID:      record.ToUserID,
		CreatedAt:     record.CreatedAt,
		FinishedBy:    finishedBy,
		FinishedAt:    time.Now(),
	})
}

type userReview struct {
	UserId       string    `json:"userId"`
	ToUserId     string    `json:"toUserId"`
	ActivityType string    `json:"activityType"`
	ActivityId   string    `json:"activityId"`
	ReviewAt     time.Time `json:"reviewAt"`
}

func PublishUserReviewEvent(ctx context.Context, userId, toUserId, activityType, activityId string) error {
	pub := pulsarutil.GetMQ[*whaleconf.Publisher](ctx)
	return pub.Pub(ctx, "UserReview", userReview{
		UserId:       userId,
		ToUserId:     toUserId,
		ActivityType: activityType,
		ActivityId:   activityId,
		ReviewAt:     time.Now(),
	})
}

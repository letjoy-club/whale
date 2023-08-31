package modelutil

import (
	"context"
	"time"
	"whale/pkg/dbquery"
	"whale/pkg/gqlient/scream"
	"whale/pkg/loader"
	"whale/pkg/models"

	"github.com/letjoy-club/mida-tool/dbutil"
	"github.com/letjoy-club/mida-tool/logger"
	"github.com/letjoy-club/mida-tool/midacontext"
	"go.uber.org/zap"
)

// ReminderMotionOffer is a function that sends a reminder to the user that they have a motion offer
func NotifyNewMotionOffer(ctx context.Context, beginTime, endTime time.Time) error {
	db := dbutil.GetDB(ctx)

	MotionOfferRecord := dbquery.Use(db).MotionOfferRecord
	motionOfferRecords, err := MotionOfferRecord.WithContext(ctx).
		Where(MotionOfferRecord.CreatedAt.Between(beginTime, endTime)).
		Where(MotionOfferRecord.State.Eq(string(models.MotionOfferStatePending))).
		Find()
	if err != nil {
		return err
	}

	userRecieveOffer := map[string][]*models.MotionOfferRecord{}

	for _, motionOfferRecord := range motionOfferRecords {
		userRecieveOffer[motionOfferRecord.UserID] = append(userRecieveOffer[motionOfferRecord.UserID], motionOfferRecord)
	}

	for userID, motionOfferRecords := range userRecieveOffer {
		err := SendMotionOfferRecievedMessage(ctx, userID, motionOfferRecords)
		if err != nil {
			logger.L.Error("send motion offer recieved message error", zap.Error(err))
		}
	}

	return nil
}

func SendMotionOfferRecievedMessage(ctx context.Context, userID string, motionOfferRecords []*models.MotionOfferRecord) error {
	motionThunk := midacontext.GetLoader[loader.Loader](ctx).Motion.Load(ctx, motionOfferRecords[0].MotionID)
	motion, err := motionThunk()
	if err != nil {
		return err
	}

	_, err = scream.SendMotionOfferRecieved(ctx, midacontext.GetServices(ctx).Scream, scream.MotionOfferRecievedParam{
		UserId:      userID,
		TopicIds:    []string{motion.TopicID},
		RecievedNum: len(motionOfferRecords),
	})
	return err
}

func SendMotionOfferAcceptedMessage(ctx context.Context, topicID string, motionOfferRecord *models.MotionOfferRecord) error {
	_, err := scream.SendMotionOfferAccepted(ctx, midacontext.GetServices(ctx).Scream, scream.MotionOfferAcceptedParam{
		UserId:      motionOfferRecord.UserID,
		TopicId:     topicID,
		PartnerId:   motionOfferRecord.ToUserID,
		ChatGroupId: motionOfferRecord.ChatGroupID,
	})
	return err
}

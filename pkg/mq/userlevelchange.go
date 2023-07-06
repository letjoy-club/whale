package mq

import (
	"context"
	"encoding/json"
	"github.com/letjoy-club/mida-tool/logger"
	"github.com/letjoy-club/mida-tool/pulsarutil"
	"go.uber.org/zap"
	"time"
	"whale/pkg/modelutil"
	"whale/pkg/whaleconf"
)

type UserLevelChangeEvent struct {
	UserId  string    `json:"userId"`
	Level   int       `json:"level"`
	MsgTime time.Time `json:"msgTime"`
}

func UserLevelChangeListener(ctx context.Context) {
	reader := pulsarutil.GetMQ[*whaleconf.Subscriber](ctx).UserLevelChange
	logger := logger.L.With(zap.String("subscription", reader.Name()))
	for {
		msg, err := reader.Receive(ctx)
		if err != nil {
			logger.Error("failed to pull data from pulsar", zap.Error(err))
			if ctx.Err() != nil {
				reader.Close()
				return
			}
			continue
		}
		logger.Info("receive message", zap.Any("payload", string(msg.Payload())))

		var ev UserLevelChangeEvent
		err = json.Unmarshal(msg.Payload(), &ev)
		if err != nil {
			logger.Error("failed to parse payload into UserLevelChangeListener", zap.Error(err))
			continue
		}
		if err := modelutil.UpdateUserRights(ctx, ev.UserId, ev.Level); err == nil {
			reader.Ack(msg)
		} else {
			reader.Nack(msg)
		}
	}
}

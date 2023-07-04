package matcher

import (
	"context"
	"fmt"
	"time"
	"whale/pkg/dbquery"
	"whale/pkg/gqlient/scream"
	"whale/pkg/loader"
	"whale/pkg/models"
	"whale/pkg/modelutil"

	"github.com/letjoy-club/mida-tool/dbutil"
	"github.com/letjoy-club/mida-tool/logger"
	"github.com/letjoy-club/mida-tool/midacode"
	"github.com/letjoy-club/mida-tool/midacontext"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Matcher struct {
}

func (m *Matcher) Match(ctx context.Context) error {
	db := dbutil.GetDB(ctx)
	Matching := dbquery.Use(db).Matching
	{
		matchings, err := Matching.WithContext(ctx).
			Where(Matching.State.Eq(string(models.MatchingStateMatching)), Matching.Deadline.Lt(time.Now())).Find()
		if err != nil {
			return err
		}
		for _, matching := range matchings {
			err := modelutil.PublishMatchingTimeoutEvent(ctx, matching)
			if err != nil {
				fmt.Println("failed to publish matching timeout event", err)
			}
		}
		if len(matchings) > 0 {
			rx, err := Matching.WithContext(ctx).
				Where(Matching.ID.In(lo.Map(matchings, func(m *models.Matching, i int) string { return m.ID })...)).
				UpdateSimple(Matching.State.Value(string(models.MatchingStateTimeout)))
			if err != nil {
				return err
			}
			if rx.RowsAffected > 0 {
				logger.L.Info("mark matching as timeout", zap.Int64("count", rx.RowsAffected))
			}
		}
	}

	matchings, err := Matching.WithContext(ctx).Where(Matching.State.Eq(string(models.MatchingStateMatching))).Find()
	if err != nil {
		return err
	}

	ctx = WithMatchingContext(ctx, matchings)
	return m.MatchTopics(ctx)
}

func (m *Matcher) MatchTopics(ctx context.Context) error {
	mc := GetMatchingContext(ctx)
	for _, topicID := range mc.Topics() {
		if err := m.MatchTopic(ctx, topicID); err != nil {
			logger.L.Error("failed to match topic", zap.Error(err), zap.String("topic-id", topicID))
		}
	}
	return nil
}

func (m *Matcher) MatchTopic(ctx context.Context, topicID string) error {
	return m.MatchPair(ctx, topicID)
}

func (m *Matcher) MatchPair(ctx context.Context, topicID string) error {
	mc := GetMatchingContext(ctx)
	matchings := mc.TopicMatchings(topicID)
	city2matchings := lo.GroupBy(matchings, func(m *models.Matching) string {
		return m.CityID
	})
	for _, cityMatchings := range city2matchings {
		if len(cityMatchings) < 2 {
			continue
		}
		err := MatchingInArea(ctx, cityMatchings)
		if err != nil {
			return err
		}
	}
	return nil
}

func MatchingInArea(ctx context.Context, matchings []*models.Matching) error {
	mc := GetMatchingContext(ctx)
	if len(matchings) < 2 {
		return nil
	}
	for i := 0; i < len(matchings); i += 1 {
		if mc.Used(matchings[i].ID) {
			continue
		}
		for j := i + 1; j < len(matchings); j += 1 {
			if mc.Used(matchings[j].ID) {
				continue
			}
			if _, matched := Matched(ctx, matchings[i], matchings[j]); matched {
				topicOption := mc.TopicOption(matchings[i].TopicID)
				if topicOption != nil {
					score := Evaluate(topicOption, matchings[i], matchings[j])
					if score.Score >= topicOption.Threshold {
						// 分数大于阈值，创建匹配结果
						if _, err := NewMatchingResult(ctx, []*models.Matching{matchings[i], matchings[j]}, score.Score); err != nil {
							logger.L.Error("failed to create matching result", zap.Error(err), zap.String("matching-1", matchings[i].ID), zap.String("matching-2", matchings[j].ID))
							return err
						}
					} else {
						// 分数小于阈值，不创建匹配结果
						continue
					}
				} else {
					// 没有选项，默认直接生成匹配结果
					if _, err := NewMatchingResult(ctx, []*models.Matching{matchings[i], matchings[j]}, 100); err != nil {
						logger.L.Error("failed to create matching result", zap.Error(err), zap.String("matching-1", matchings[i].ID), zap.String("matching-2", matchings[j].ID))
						return err
					}
				}

				mc.Use(matchings[i].ID)
				mc.Use(matchings[j].ID)
			}
		}
	}
	return nil
}

func NewMatchingResult(ctx context.Context, matchings []*models.Matching, score int) (*models.MatchingResult, error) {
	db := dbutil.GetDB(ctx)
	MatchingResult := dbquery.Use(db).MatchingResult

	userIDs := make([]string, len(matchings))
	states := make([]string, len(matchings))
	matchingIDs := make([]string, len(matchings))
	for i := range matchings {
		matchingIDs[i] = matchings[i].ID
		states[i] = models.MatchingResultConfirmStateConfirmed.String()
		userIDs[i] = matchings[i].UserID
	}

	matchingResult := &models.MatchingResult{
		MatchingIDs:    matchingIDs,
		UserIDs:        userIDs,
		ChatGroupState: models.ChatGroupStateUncreated.String(),
		ConfirmStates:  states,
		MatchingScore:  score,
		CreatedBy:      string(models.ResultCreatedByMatching),
		TopicID:        matchings[0].TopicID,
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		err := MatchingResult.WithContext(ctx).Create(matchingResult)
		if err != nil {
			return err
		}
		Matching := dbquery.Use(tx).Matching
		matched := time.Now()
		rx, err := Matching.
			WithContext(ctx).
			Where(Matching.ID.In(matchingIDs...), Matching.State.Eq(models.MatchingStateMatching.String())).
			Select(Matching.ResultID, Matching.State, Matching.MatchedAt).
			Updates(&models.Matching{
				ResultID:  matchingResult.ID,
				State:     models.MatchingStateMatched.String(),
				MatchedAt: &matched,
			})
		if err != nil {
			return err
		}
		if rx.RowsAffected != int64(len(matchingIDs)) {
			return midacode.ErrStateMayHaveChanged
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	for _, matching := range matchings {
		err := modelutil.PublishMatchedEvent(ctx, matching)
		if err != nil {
			fmt.Println("failed to publish matched event", err)
		}
	}

	err = modelutil.CheckMatchingResultAndCreateChatGroup(ctx, matchingResult)
	if err != nil {
		fmt.Println("failed to create chat group", err)
	}

	// 通知
	areaIDs := []string{matchings[0].CityID}
	_, err = scream.MatchingGroupCreated(ctx, midacontext.GetServices(ctx).Scream, scream.MatchingGroupCreatedParam{
		MatchingId: matchingResult.MatchingIDs[0],
		UserId:     userIDs[0],
		PartnerId:  userIDs[1],
		TopicId:    matchingResult.TopicID,
		AreaIds:    areaIDs,
	})
	if err != nil {
		fmt.Println("failed send notification: create matching group err", err)
	}
	_, err = scream.MatchingGroupCreated(ctx, midacontext.GetServices(ctx).Scream, scream.MatchingGroupCreatedParam{
		MatchingId: matchingResult.MatchingIDs[1],
		UserId:     userIDs[1],
		PartnerId:  userIDs[0],
		TopicId:    matchingResult.TopicID,
		AreaIds:    areaIDs,
	})
	if err != nil {
		fmt.Println("err", err)
	}
	topicName := modelutil.GetTopicName(ctx, matchingResult.TopicID)
	if topicName == "" {
		return nil, nil
	}
	thunk := midacontext.GetLoader[loader.Loader](ctx).UserAvatarNickname.LoadMany(ctx, userIDs)
	users, _ := thunk()
	if len(users) > 0 {
		_, err = scream.SendUserNotification(ctx,
			midacontext.GetServices(ctx).Scream,
			scream.UserNotificationKindMatched,
			userIDs[1],
			map[string]interface{}{
				"userName":   users[0].Nickname,
				"userId":     users[0].ID,
				"topicName":  topicName,
				"matchingId": matchingResult.MatchingIDs[0],
			},
		)
		if err != nil {
			fmt.Println("err", err)
		}
		_, err = scream.SendUserNotification(ctx,
			midacontext.GetServices(ctx).Scream,
			scream.UserNotificationKindMatched,
			userIDs[0],
			map[string]interface{}{
				"userName":   users[1].Nickname,
				"userId":     users[1].ID,
				"topicName":  topicName,
				"matchingId": matchingResult.MatchingIDs[1],
			},
		)
		if err != nil {
			fmt.Println("err", err)
		}
	}

	return matchingResult, nil
}

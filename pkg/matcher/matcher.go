package matcher

import (
	"context"
	"fmt"
	"time"
	"whale/pkg/dbquery"
	"whale/pkg/gqlient/scream"
	"whale/pkg/keyer"
	"whale/pkg/loader"
	"whale/pkg/models"
	"whale/pkg/modelutil"

	"github.com/letjoy-club/mida-tool/dbutil"
	"github.com/letjoy-club/mida-tool/logger"
	"github.com/letjoy-club/mida-tool/midacode"
	"github.com/letjoy-club/mida-tool/midacontext"
	"github.com/letjoy-club/mida-tool/redisutil"
	"github.com/samber/lo"
	"go.uber.org/zap"
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
				Where(
					Matching.ID.In(lo.Map(matchings, func(m *models.Matching, i int) string { return m.ID })...),
				).
				UpdateSimple(Matching.State.Value(string(models.MatchingStateTimeout)))
			if err != nil {
				return err
			}
			if rx.RowsAffected > 0 {
				logger.L.Info("mark matching as timeout", zap.Int64("count", rx.RowsAffected))
			}
		}
	}

	matchings, err := Matching.WithContext(ctx).Where(
		Matching.State.Eq(string(models.MatchingStateMatching)),
		Matching.WithContext(ctx).Or(
			Matching.StartMatchingAt.Lt(time.Now()),
			Matching.StartMatchingAt.IsNull(),
		),
	).Find()
	if err != nil {
		return err
	}

	ctx = WithMatchingContext(ctx, matchings)
	return m.MatchTopics(ctx)
}

func (m *Matcher) MatchTopics(ctx context.Context) error {
	mc := GetMatchingContext(ctx)
	// 普通匹配
	for _, topicID := range mc.Topics() {
		config := mc.TopicOption(topicID)
		if config == nil {
			if err := m.MatchTopic(ctx, topicID); err != nil {
				logger.L.Error("failed to match topic", zap.Error(err), zap.String("topic-id", topicID))
			}
		} else {
			if !config.FuzzyMatchingTopic {
				if err := m.MatchTopic(ctx, topicID); err != nil {
					logger.L.Error("failed to match topic", zap.Error(err), zap.String("topic-id", topicID))
				}
			}
		}
	}

	// 模糊匹配
	pairableMatchings := m.PairableMatchings(ctx)
	fuzzyMatchings := m.FuzzyMatchings(ctx)
	if err := m.MatchingFuzzyTopic(ctx, fuzzyMatchings, pairableMatchings); err != nil {
		logger.L.Error("failed to fuzzy match topic", zap.Error(err))
	}
	return nil
}

func (m *Matcher) MatchingFuzzyTopic(ctx context.Context, fuzzyMatchings, matchings []*models.Matching) error {
	mc := GetMatchingContext(ctx)
	city2matchings := lo.GroupBy(matchings, func(m *models.Matching) string {
		return m.CityID
	})
	for _, fuzzyMatching := range fuzzyMatchings {
		if mc.Used(fuzzyMatching.ID) {
			continue
		}
		cityMatchings := city2matchings[fuzzyMatching.CityID]
		if len(cityMatchings) == 0 {
			continue
		}
		for _, matching := range cityMatchings {
			if mc.Used(matching.ID) {
				continue
			}
			reason, matched := Matched(ctx, matching, fuzzyMatching)
			if !matched {
				continue
			}
			topicOption := mc.TopicOption(fuzzyMatching.TopicID)
			if topicOption == nil {
				continue
			}
			score := EvaluateWithFuzzyMatching(topicOption, matching, fuzzyMatching)
			score.FailedReason = reason
			// 只考虑时间分数
			if score.TimeScore >= topicOption.Threshold {
				// 分数大于阈值，创建匹配结果
				if _, err := NewMatchingResult(ctx, []*models.Matching{matching, fuzzyMatching}, score.Score); err != nil {
					logger.L.Error("failed to create matching result", zap.Error(err), zap.String("matching", matching.ID), zap.String("fuzzy-matching", fuzzyMatching.ID))
					return err
				}
			} else {
				// 分数小于阈值，不创建匹配结果
				continue
			}
		}
	}
	return nil
}

func (m *Matcher) MatchTopic(ctx context.Context, topicID string) error {
	return m.MatchPair(ctx, topicID)
}

func (m *Matcher) FuzzyMatchings(ctx context.Context) []*models.Matching {
	mc := GetMatchingContext(ctx)
	matchings := []*models.Matching{}
	for _, topic := range mc.Topics() {
		config := mc.TopicOption(topic)
		if config == nil {
			continue
		}
		if !config.FuzzyMatchingTopic {
			continue
		}
		ms := mc.TopicMatchings(topic)
		if ms == nil {
			continue
		}
		for _, m := range ms {
			if mc.Used(m.ID) {
				continue
			}
			matchings = append(matchings, m)
		}
	}
	return matchings
}

func (m *Matcher) PairableMatchings(ctx context.Context) []*models.Matching {
	mc := GetMatchingContext(ctx)

	matchings := []*models.Matching{}
	for _, topic := range mc.Topics() {
		config := mc.TopicOption(topic)
		if config == nil {
			continue
		}
		if !config.AllowFuzzyMatching || config.FuzzyMatchingTopic {
			continue
		}
		if config.DelayMinuteToPairWithFuzzyTopic < 0 {
			continue
		}
		ms := mc.TopicMatchings(topic)
		if ms == nil {
			continue
		}
		collectBefore := time.Now().Add(-time.Duration(config.DelayMinuteToPairWithFuzzyTopic) * time.Minute)
		for _, m := range ms {
			if mc.Used(m.ID) {
				continue
			}
			if m.CreatedAt.Before(collectBefore) {
				matchings = append(matchings, m)
			}
		}
	}
	return matchings
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
	// 给两个用户加锁
	release, err := redisutil.LockAll(ctx, keyer.UserMatching(matchings[0].UserID), keyer.UserMatching(matchings[1].UserID))
	if err != nil {
		return nil, err
	}
	defer release(ctx)

	db := dbutil.GetDB(ctx)

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

	err = dbquery.Use(db).Transaction(func(tx *dbquery.Query) error {
		MatchingResult := tx.MatchingResult
		Matching := tx.Matching
		err := MatchingResult.WithContext(ctx).Create(matchingResult)
		if err != nil {
			return err
		}
		matched := time.Now()
		// 更新 matching 的状态
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

		// 确保 matching offer summary 存在
		_, err = modelutil.GetMatchingOfferSummary(ctx, matchings[0])
		if err != nil {
			return err
		}
		_, err = modelutil.GetMatchingOfferSummary(ctx, matchings[1])
		if err != nil {
			return err
		}

		// 更新 matchingOfferSummary 状态
		err = modelutil.DeactiveMatchingSummaryAndUpdateMatchingOffer(ctx, tx, matchings[0].ID)
		if err != nil {
			return err
		}
		err = modelutil.DeactiveMatchingSummaryAndUpdateMatchingOffer(ctx, tx, matchings[1].ID)
		if err != nil {
			return err
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

package matcher

import (
	"context"
	"time"
	"whale/pkg/dbquery"
	"whale/pkg/models"

	"github.com/letjoy-club/mida-tool/dbutil"
	"github.com/letjoy-club/mida-tool/logger"
	"github.com/letjoy-club/mida-tool/midacode"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Matcher struct {
}

func (m *Matcher) Match(ctx context.Context) error {
	db := dbutil.GetDB(ctx)
	Matching := dbquery.Use(db).Matching
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
				if _, err := NewMatchingResult(ctx, []*models.Matching{matchings[i], matchings[j]}); err != nil {
					logger.L.Error("failed to create matching result", zap.Error(err), zap.String("matching-1", matchings[i].ID), zap.String("matching-2", matchings[j].ID))
					return err
				}
				mc.Use(matchings[i].ID)
				mc.Use(matchings[j].ID)
			}
		}
	}
	return nil
}

func NewMatchingResult(ctx context.Context, matchings []*models.Matching) (*models.MatchingResult, error) {
	db := dbutil.GetDB(ctx)
	MatchingResult := dbquery.Use(db).MatchingResult

	userIDs := make([]string, len(matchings))
	states := make([]string, len(matchings))
	matchingIDs := make([]string, len(matchings))
	for i := range matchings {
		matchingIDs[i] = matchings[i].ID
		states[i] = models.MatchingResultConfirmStateUnconfirmed.String()
		userIDs[i] = matchings[i].UserID
	}

	matchingResult := &models.MatchingResult{
		MatchingIDs:    matchingIDs,
		UserIDs:        userIDs,
		ChatGroupState: models.ChatGroupStateUncreated.String(),
		ConfirmStates:  states,
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
	return matchingResult, err
}

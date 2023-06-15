package modelutil

import (
	"context"
	"whale/pkg/dbquery"
	"whale/pkg/loader"
	"whale/pkg/models"

	"github.com/letjoy-club/mida-tool/dbutil"
	"github.com/letjoy-club/mida-tool/midacontext"
	"github.com/samber/lo"
	"go.uber.org/multierr"
	"gorm.io/gorm/clause"
)

func AddMatchingToRecent(ctx context.Context, matchingID string) (*models.RecentMatching, error) {
	db := dbutil.GetDB(ctx)
	RecentMatching := dbquery.Use(db).RecentMatching

	matchingThunk := midacontext.GetLoader[loader.Loader](ctx).Matching.Load(ctx, matchingID)
	matching, err := matchingThunk()
	if err != nil {
		return nil, err
	}

	recentThunk := midacontext.GetLoader[loader.Loader](ctx).RecentMatching.Load(ctx, matching.CityID+"-"+matching.TopicID)
	recentMatching, err := recentThunk()
	if err != nil {
		recentMatching = &models.RecentMatching{
			ID:          matching.CityID + "-" + matching.TopicID,
			CityID:      matching.CityID,
			TopicID:     matching.TopicID,
			MatchingIDs: []string{},
		}
	}
	if lo.Contains(recentMatching.MatchingIDs, matchingID) {
		return recentMatching, nil
	}
	matchingIDs := append([]string{matchingID}, recentMatching.MatchingIDs...)
	matchingsThunk := midacontext.GetLoader[loader.Loader](ctx).Matching.LoadMany(ctx, matchingIDs)
	matchings, errors := matchingsThunk()
	if errors != nil {
		return nil, multierr.Combine(errors...)
	}
	users := map[string]struct{}{}
	newMatchingIDs := []string{}
	for _, matching := range matchings {
		if _, ok := users[matching.UserID]; ok {
			continue
		}
		users[matching.UserID] = struct{}{}
		newMatchingIDs = append(newMatchingIDs, matching.ID)
	}

	recentMatching.MatchingIDs = newMatchingIDs
	if len(recentMatching.MatchingIDs) > 8 {
		recentMatching.MatchingIDs = recentMatching.MatchingIDs[:8]
	}

	err = RecentMatching.WithContext(ctx).
		Clauses(clause.OnConflict{DoUpdates: clause.AssignmentColumns([]string{
			RecentMatching.MatchingIDs.ColumnName().String(),
		})}).
		Save(recentMatching)

	midacontext.GetLoader[loader.Loader](ctx).RecentMatching.Prime(ctx, recentMatching.ID, recentMatching)
	midacontext.GetLoader[loader.Loader](ctx).CityTopicMatchings.Clear(ctx, loader.CityTopicKey{
		CityID:  matching.CityID,
		TopicID: matching.TopicID,
	})
	return recentMatching, err
}

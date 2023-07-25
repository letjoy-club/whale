package loader

import (
	"context"
	"time"
	"whale/pkg/dbquery"
	"whale/pkg/models"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/letjoy-club/mida-tool/loaderutil"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type UserLikeMatchingSummary struct {
	userID string

	matchingMap map[string]struct{}
	matchingIds []*models.UserLikeMatching
}

func (u *UserLikeMatchingSummary) LikeNum() int {
	return len(u.matchingMap)
}

func (u *UserLikeMatchingSummary) Slice(offset, limit int) []*models.UserLikeMatching {
	return lo.Slice(u.matchingIds, offset, offset+limit)
}

func (u *UserLikeMatchingSummary) IsLike(matchingID string) bool {
	_, ok := u.matchingMap[matchingID]
	return ok
}

func (u *UserLikeMatchingSummary) CancelLike(matchingID string) {
	newMap := map[string]struct{}{}
	for k, v := range u.matchingMap {
		if k != matchingID {
			newMap[k] = v
		}
	}
	u.matchingMap = newMap
}

func NewUserLikeMatchingLoader(db *gorm.DB) *dataloader.Loader[string, *UserLikeMatchingSummary] {
	UserLikeMatching := dbquery.Use(db).UserLikeMatching
	return loaderutil.NewAggregatorLoader(db, func(ctx context.Context, userIDs []string) ([]*models.UserLikeMatching, error) {
		matchings, err := UserLikeMatching.WithContext(ctx).
			Where(UserLikeMatching.UserID.In(userIDs...)).
			Select(UserLikeMatching.UserID, UserLikeMatching.ToMatchingID, UserLikeMatching.CreatedAt).
			Find()
		if err != nil {
			return nil, err
		}
		return matchings, nil
	}, func(m map[string]*UserLikeMatchingSummary, v *models.UserLikeMatching) {
		if _, ok := m[v.UserID]; !ok {
			m[v.UserID] = &UserLikeMatchingSummary{
				userID:      v.UserID,
				matchingMap: map[string]struct{}{v.ToMatchingID: {}},
				matchingIds: []*models.UserLikeMatching{v},
			}
		}
		s := m[v.UserID]
		s.matchingMap[v.ToMatchingID] = struct{}{}
		s.matchingIds = append(s.matchingIds, v)
	}, time.Minute, loaderutil.Placeholder(func(ctx context.Context, id string) (*UserLikeMatchingSummary, error) {
		return &UserLikeMatchingSummary{
			userID:      id,
			matchingMap: map[string]struct{}{},
			matchingIds: []*models.UserLikeMatching{},
		}, nil
	}),
	)
}

func NewMatchingReceiveLikeLoader(db *gorm.DB) *dataloader.Loader[string, *models.MatchingReceiveLike] {
	MatchingReceiveLike := dbquery.Use(db).MatchingReceiveLike
	return loaderutil.NewItemLoader(db, func(ctx context.Context, ids []string) ([]*models.MatchingReceiveLike, error) {
		matchingReceiveLikes, err := MatchingReceiveLike.WithContext(ctx).Where(MatchingReceiveLike.MatchingID.In(ids...)).Select(MatchingReceiveLike.MatchingID, MatchingReceiveLike.LikeNum).Find()
		if err != nil {
			return nil, err
		}
		return matchingReceiveLikes, nil
	}, func(m map[string]*models.MatchingReceiveLike, v *models.MatchingReceiveLike) {}, time.Minute*5, loaderutil.Placeholder(func(ctx context.Context, id string) (*models.MatchingReceiveLike, error) {
		return &models.MatchingReceiveLike{MatchingID: id}, nil
	}))
}

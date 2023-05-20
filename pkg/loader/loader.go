package loader

import (
	"context"
	"time"
	"whale/pkg/dbquery"
	"whale/pkg/gqlient/hoopoe"
	"whale/pkg/models"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/letjoy-club/mida-tool/midacode"
	"github.com/letjoy-club/mida-tool/midacontext"
	"github.com/letjoy-club/mida-tool/ttlcache"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type Loader struct {
	Matching           *dataloader.Loader[string, *models.Matching]
	MatchingInvitation *dataloader.Loader[string, *models.MatchingInvitation]
	MatchingQuota      *dataloader.Loader[string, *models.MatchingQuota]
	MatchingResult     *dataloader.Loader[int, *models.MatchingResult]
	MatchingReviewed   *dataloader.Loader[string, MatchingReviewed]

	UserProfile        *dataloader.Loader[string, UserProfile]
	UserAvatarNickname *dataloader.Loader[string, UserAvatarNickname]
	HotTopics          *dataloader.Loader[string, *models.HotTopicsInArea]
}

func NewLoader(db *gorm.DB) *Loader {
	return &Loader{
		Matching: NewSingleLoader(db, func(ctx context.Context, keys []string) ([]*models.Matching, error) {
			Matching := dbquery.Use(db).Matching
			return Matching.WithContext(ctx).Where(Matching.ID.In(keys...)).Find()
		}, func(k map[string]*models.Matching, v *models.Matching) {
			k[v.ID] = v
		}, time.Second*10),
		MatchingInvitation: NewSingleLoader(db, func(ctx context.Context, keys []string) ([]*models.MatchingInvitation, error) {
			MatchingInvitation := dbquery.Use(db).MatchingInvitation
			return MatchingInvitation.WithContext(ctx).Where(MatchingInvitation.ID.In(keys...)).Find()
		}, func(k map[string]*models.MatchingInvitation, v *models.MatchingInvitation) {
			k[v.ID] = v
		}, time.Second*10),
		MatchingQuota: NewSingleLoader(db, func(ctx context.Context, keys []string) ([]*models.MatchingQuota, error) {
			MatchingQuota := dbquery.Use(db).MatchingQuota
			matchingQuotas, err := MatchingQuota.WithContext(ctx).Where(MatchingQuota.UserID.In(keys...)).Find()
			if err != nil {
				return nil, err
			}
			notFound := map[string]struct{}{}
			for _, id := range keys {
				notFound[id] = struct{}{}
			}
			for _, matchingQuota := range matchingQuotas {
				delete(notFound, matchingQuota.UserID)
			}
			toBeAdded := []*models.MatchingQuota{}
			if len(notFound) > 0 {
				for id := range notFound {
					toBeAdded = append(toBeAdded, &models.MatchingQuota{
						UserID: id,
						Remain: 3,
						Total:  3,
					})
				}
				if err := MatchingQuota.WithContext(ctx).Create(toBeAdded...); err != nil {
					return nil, err
				}
			}
			return append(matchingQuotas, toBeAdded...), nil
		}, func(k map[string]*models.MatchingQuota, v *models.MatchingQuota) { k[v.UserID] = v }, time.Second*60),
		MatchingResult: NewSingleLoader(db, func(ctx context.Context, keys []int) ([]*models.MatchingResult, error) {
			MatchingResult := dbquery.Use(db).MatchingResult
			return MatchingResult.WithContext(ctx).Where(MatchingResult.ID.In(keys...)).Find()
		}, func(k map[int]*models.MatchingResult, v *models.MatchingResult) { k[v.ID] = v }, time.Second*100),
		MatchingReviewed: NewSingleLoader(db, func(ctx context.Context, keys []string) ([]MatchingReviewed, error) {
			MatchingReview := dbquery.Use(db).MatchingReview
			reviews, err := MatchingReview.WithContext(ctx).Where(MatchingReview.MatchingID.In(keys...)).Find()
			if err != nil {
				return nil, err
			}
			return lo.Map(reviews, func(r *models.MatchingReview, i int) MatchingReviewed {
				return MatchingReviewed{
					MatchingReviewed: r.MatchingID,
					Reviewed:         true,
				}
			}), nil
		}, func(k map[string]MatchingReviewed, v MatchingReviewed) { k[v.MatchingReviewed] = v }, time.Second*30),
		UserProfile: NewSingleLoader(db, func(ctx context.Context, keys []string) ([]UserProfile, error) {
			ret, err := hoopoe.GetUserByIDs(ctx, midacontext.GetServices(ctx).Hoopoe, keys)
			if err != nil {
				return nil, err
			}
			return lo.Map(ret.GetUserByIds, func(u hoopoe.GetUserByIDsGetUserByIdsUser, i int) UserProfile {
				return UserProfile{ID: u.Id, Gender: models.Gender(u.Gender)}
			}), nil
		}, func(k map[string]UserProfile, v UserProfile) { k[v.ID] = v }, time.Minute),
		UserAvatarNickname: NewSingleLoader(db, func(ctx context.Context, keys []string) ([]UserAvatarNickname, error) {
			ret, err := hoopoe.GetAvatarByIDs(ctx, midacontext.GetServices(ctx).Hoopoe, keys)
			if err != nil {
				return nil, err
			}
			return lo.Map(ret.GetUserByIds, func(u hoopoe.GetAvatarByIDsGetUserByIdsUser, i int) UserAvatarNickname {
				return UserAvatarNickname{ID: u.Id, Avatar: u.Avatar, Nickname: u.Nickname}
			}), nil
		}, func(k map[string]UserAvatarNickname, u UserAvatarNickname) { k[u.ID] = u }, time.Minute),
		HotTopics: NewSingleLoader(db, func(ctx context.Context, keys []string) ([]*models.HotTopicsInArea, error) {
			HotTopicsInArea := dbquery.Use(db).HotTopicsInArea
			topics, err := HotTopicsInArea.WithContext(ctx).Where(HotTopicsInArea.CityID.In(keys...)).Find()
			return topics, err
		}, func(k map[string]*models.HotTopicsInArea, v *models.HotTopicsInArea) { k[v.CityID] = v }, time.Second*60),
	}
}

func NewSingleLoader[K comparable, V any](
	db *gorm.DB,
	loader func(ctx context.Context, keys []K) ([]V, error),
	dataMaper func(k map[K]V, v V),
	duration time.Duration,
) *dataloader.Loader[K, V] {
	c := ttlcache.New[K, V](duration)
	return dataloader.NewBatchedLoader(func(ctx context.Context, keys []K) []*dataloader.Result[V] {
		data, err := loader(ctx, keys)

		if err != nil {
			return lo.Map(data, func(m V, i int) *dataloader.Result[V] {
				return &dataloader.Result[V]{Error: err}
			})
		}

		dataMap := map[K]V{}
		for _, m := range data {
			dataMaper(dataMap, m)
		}

		return lo.Map(keys, func(key K, i int) *dataloader.Result[V] {
			ret, itemNotFound := dataMap[key]
			if itemNotFound {
				return &dataloader.Result[V]{Data: ret}
			} else {
				return &dataloader.Result[V]{Error: midacode.ErrItemNotFound}
			}
		})
	}, dataloader.WithCache[K, V](&c),
	)
}

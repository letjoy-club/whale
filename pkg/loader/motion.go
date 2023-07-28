package loader

import (
	"context"
	"sort"
	"time"
	"whale/pkg/dbquery"
	"whale/pkg/models"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/letjoy-club/mida-tool/loaderutil"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

func NewMotionLoader(db *gorm.DB) *dataloader.Loader[string, *models.Motion] {
	Motion := dbquery.Use(db).Motion
	return loaderutil.NewItemLoader(db, func(ctx context.Context, keys []string) (items []*models.Motion, err error) {
		return Motion.WithContext(ctx).Where(Motion.ID.In(keys...)).Find()
	}, func(m map[string]*models.Motion, v *models.Motion) {
		m[v.ID] = v
	}, time.Minute)
}

type MotionOffers struct {
	motionID string
	Offers   []*models.MotionOfferRecord
}

func NewInMotionOfferLoader(db *gorm.DB) *dataloader.Loader[string, *MotionOffers] {
	MotionOfferRecord := dbquery.Use(db).MotionOfferRecord
	return loaderutil.NewAggregatorLoader(db, func(ctx context.Context, keys []string) (items []*models.MotionOfferRecord, err error) {
		return MotionOfferRecord.WithContext(ctx).Where(MotionOfferRecord.ToMotionID.In(keys...)).Find()
	}, func(m map[string]*MotionOffers, v *models.MotionOfferRecord) {
		if _, ok := m[v.ToMotionID]; !ok {
			m[v.ToMotionID] = &MotionOffers{
				motionID: v.ToMotionID,
			}
		} else {
			m[v.ToMotionID].Offers = append(m[v.ToMotionID].Offers, v)
		}
	}, time.Minute, loaderutil.Placeholder(func(ctx context.Context, keys string) (item *MotionOffers, err error) {
		return &MotionOffers{Offers: []*models.MotionOfferRecord{}}, nil
	}))
}

func NewOutMotionOfferLoader(db *gorm.DB) *dataloader.Loader[string, *MotionOffers] {
	MotionOfferRecord := dbquery.Use(db).MotionOfferRecord
	return loaderutil.NewAggregatorLoader(db, func(ctx context.Context, keys []string) (items []*models.MotionOfferRecord, err error) {
		return MotionOfferRecord.WithContext(ctx).Where(MotionOfferRecord.MotionID.In(keys...)).Find()
	}, func(m map[string]*MotionOffers, v *models.MotionOfferRecord) {
		if _, ok := m[v.MotionID]; !ok {
			m[v.MotionID] = &MotionOffers{
				motionID: v.MotionID,
			}
		} else {
			m[v.MotionID].Offers = append(m[v.MotionID].Offers, v)
		}
	}, time.Minute, loaderutil.Placeholder(func(ctx context.Context, keys string) (item *MotionOffers, err error) {
		return &MotionOffers{Offers: []*models.MotionOfferRecord{}}, nil
	}))
}

type UserSubmitMotion struct {
	UserID    string
	MotionIDs []string
}

func (u *UserSubmitMotion) Submit(motionID string) {
	if u.IsSubmitted(motionID) {
		return
	}
	// 有并发问题，但是不影响稳定性
	u.MotionIDs = remove(u.MotionIDs, motionID)
}

func NewUserSubmitMotionLoader(db *gorm.DB) *dataloader.Loader[string, *UserSubmitMotion] {
	MatchingOfferRecord := dbquery.Use(db).MatchingOfferRecord
	return loaderutil.NewAggregatorLoader(db, func(ctx context.Context, keys []string) (items []*models.MatchingOfferRecord, err error) {
		motions, err := MatchingOfferRecord.WithContext(ctx).Where(MatchingOfferRecord.UserID.In(keys...)).Select(MatchingOfferRecord.ToMatchingID, MatchingOfferRecord.UserID).Find()
		if err != nil {
			return nil, err
		}
		// 按 motion id 排序，后面会用二分查找
		sort.Slice(motions, func(i, j int) bool {
			return motions[i].ToMatchingID < motions[j].ToMatchingID
		})
		return motions, nil
	}, func(m map[string]*UserSubmitMotion, v *models.MatchingOfferRecord) {
		if _, ok := m[v.UserID]; !ok {
			m[v.UserID] = &UserSubmitMotion{
				UserID:    v.UserID,
				MotionIDs: []string{v.ToMatchingID},
			}
		} else {
			m[v.UserID].MotionIDs = append(m[v.UserID].MotionIDs, v.ToMatchingID)
		}
	}, time.Minute, loaderutil.Placeholder(func(ctx context.Context, keys string) (item *UserSubmitMotion, err error) {
		return &UserSubmitMotion{MotionIDs: []string{}}, nil
	}))
}

func (u UserSubmitMotion) IsSubmitted(motionID string) bool {
	return searchString(motionID, u.MotionIDs) != -1
}

type UserLikeMotion struct {
	UserID    string
	MotionIDs []string
}

func (u *UserLikeMotion) Like(motionID string) {
	if u.IsLike(motionID) {
		return
	}
	// 有并发问题，但是不影响稳定性
	u.MotionIDs = insert(u.MotionIDs, motionID)
}

func (u *UserLikeMotion) Unlike(motionID string) {
	if !u.IsLike(motionID) {
		return
	}
	// 有并发问题，但是不影响稳定性
	u.MotionIDs = remove(u.MotionIDs, motionID)
}

func (u UserLikeMotion) IsLike(motionID string) bool {
	return searchString(motionID, u.MotionIDs) != -1
}

func NewUserLikeMotionLoader(db *gorm.DB) *dataloader.Loader[string, *UserLikeMotion] {
	LikedMotion := dbquery.Use(db).LikeMotion
	return loaderutil.NewAggregatorLoader(db, func(ctx context.Context, keys []string) (items []*models.LikeMotion, err error) {
		motions, err := LikedMotion.WithContext(ctx).Where(LikedMotion.UserID.In(keys...)).Find()
		if err != nil {
			return nil, err
		}
		// 按 motion id 排序，后面会用二分查找
		sort.Slice(motions, func(i, j int) bool {
			return motions[i].MotionID < motions[j].MotionID
		})
		return motions, nil
	}, func(m map[string]*UserLikeMotion, v *models.LikeMotion) {
		if _, ok := m[v.UserID]; !ok {
			m[v.UserID] = &UserLikeMotion{
				UserID:    v.UserID,
				MotionIDs: []string{v.MotionID},
			}
		} else {
			m[v.UserID].MotionIDs = append(m[v.UserID].MotionIDs, v.MotionID)
		}
	}, time.Minute*5, loaderutil.Placeholder(func(ctx context.Context, keys string) (item *UserLikeMotion, err error) {
		return &UserLikeMotion{MotionIDs: []string{}}, nil
	}))
}

func searchString(id string, ids []string) int {
	l := 0
	r := len(ids) - 1
	for l <= r {
		m := (l + r) / 2
		if ids[m] == id {
			return m
		} else if ids[m] < id {
			l = m + 1
		} else {
			r = m - 1
		}
	}
	return -1
}

func insert[T constraints.Ordered](ts []T, t T) []T {
	var dummy T
	ts = append(ts, dummy) // extend the slice

	i := slices.BinarySearch(ts, t) // find slot
	copy(ts[i+1:], ts[i:])          // make room
	ts[i] = t
	return ts
}

func remove[T constraints.Ordered](ts []T, t T) []T {
	i := slices.BinarySearch(ts, t)
	if i < 0 {
		return ts
	}
	copy(ts[i:], ts[i+1:])
	return ts[:len(ts)-1]
}

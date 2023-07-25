package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.35

import (
	"context"
	"whale/pkg/loader"
	"whale/pkg/models"
	"whale/pkg/modelutil"

	"github.com/letjoy-club/mida-tool/graphqlutil"
	"github.com/letjoy-club/mida-tool/midacode"
	"github.com/letjoy-club/mida-tool/midacontext"
	"github.com/samber/lo"
	"go.uber.org/multierr"
)

// PreferredPeriods is the resolver for the preferredPeriods field.
func (r *discoverMatchingResolver) PreferredPeriods(ctx context.Context, obj *models.Matching) ([]models.DatePeriod, error) {
	return lo.Map(obj.PreferredPeriods, func(m string, i int) models.DatePeriod {
		return models.DatePeriod(m)
	}), nil
}

// TopicOptionConfig is the resolver for the topicOptionConfig field.
func (r *discoverMatchingResolver) TopicOptionConfig(ctx context.Context, obj *models.Matching) (*models.TopicOptionConfig, error) {
	return &models.TopicOptionConfig{TopicID: obj.TopicID}, nil
}

// Liked is the resolver for the liked field.
func (r *discoverMatchingResolver) Liked(ctx context.Context, obj *models.Matching) (bool, error) {
	thunk := midacontext.GetLoader[loader.Loader](ctx).UserLikeMatching.Load(ctx, obj.UserID)
	like, err := thunk()
	if err != nil {
		return false, err
	}
	return like.IsLike(obj.ID), nil
}

// ViewCount is the resolver for the viewCount field.
func (r *discoverMatchingResolver) ViewCount(ctx context.Context, obj *models.Matching) (int, error) {
	thunk := midacontext.GetLoader[loader.Loader](ctx).MatchingView.Load(ctx, obj.ID)
	view, err := thunk()
	if err != nil {
		return 0, err
	}
	return view.ViewCount, nil
}

// LikeCount is the resolver for the likeCount field.
func (r *discoverMatchingResolver) LikeCount(ctx context.Context, obj *models.Matching) (int, error) {
	thunk := midacontext.GetLoader[loader.Loader](ctx).MatchingReceiveLike.Load(ctx, obj.ID)
	like, err := thunk()
	if err != nil {
		return 0, err
	}
	return like.LikeNum, nil
}

// User is the resolver for the user field.
func (r *discoverMatchingResolver) User(ctx context.Context, obj *models.Matching) (*models.User, error) {
	return &models.User{ID: obj.UserID}, nil
}

// City is the resolver for the city field.
func (r *discoverMatchingResolver) City(ctx context.Context, obj *models.Matching) (*models.Area, error) {
	return &models.Area{Code: obj.CityID}, nil
}

// Areas is the resolver for the areas field.
func (r *discoverMatchingResolver) Areas(ctx context.Context, obj *models.Matching) ([]*models.Area, error) {
	return lo.Map(obj.AreaIDs, func(m string, i int) *models.Area {
		return &models.Area{Code: m}
	}), nil
}

// State is the resolver for the state field.
func (r *matchingOfferRecordResolver) State(ctx context.Context, obj *models.MatchingOfferRecord) (models.MatchingOfferState, error) {
	return models.MatchingOfferState(obj.State), nil
}

// ToMatching is the resolver for the toMatching field.
func (r *matchingOfferRecordResolver) ToMatching(ctx context.Context, obj *models.MatchingOfferRecord) (*models.Matching, error) {
	thunk := midacontext.GetLoader[loader.Loader](ctx).Matching.Load(ctx, obj.ToMatchingID)
	matching, err := thunk()
	if err != nil {
		return nil, err
	}
	return matching, nil
}

// Matching is the resolver for the matching field.
func (r *matchingOfferRecordResolver) Matching(ctx context.Context, obj *models.MatchingOfferRecord) (*models.Matching, error) {
	thunk := midacontext.GetLoader[loader.Loader](ctx).Matching.Load(ctx, obj.MatchingID)
	matching, err := thunk()
	if err != nil {
		return nil, err
	}
	return matching, nil
}

// UnprocessedInOfferNum is the resolver for the unprocessedInOfferNum field.
func (r *matchingOfferSummaryResolver) UnprocessedInOfferNum(ctx context.Context, obj *models.MatchingOfferSummary) (int, error) {
	thunk := midacontext.GetLoader[loader.Loader](ctx).InMatchingOfferRecord.Load(ctx, obj.MatchingID)
	record, err := thunk()
	if err != nil {
		return 0, err
	}
	return record.UnprocessCount(), nil
}

// UnprocessedOutOfferNum is the resolver for the unprocessedOutOfferNum field.
func (r *matchingOfferSummaryResolver) UnprocessedOutOfferNum(ctx context.Context, obj *models.MatchingOfferSummary) (int, error) {
	thunk := midacontext.GetLoader[loader.Loader](ctx).OutMatchingOfferRecord.Load(ctx, obj.MatchingID)
	record, err := thunk()
	if err != nil {
		return 0, err
	}
	return record.UnprocessCount(), nil
}

// InMatchingOffers is the resolver for the inMatchingOffers field.
func (r *matchingOfferSummaryResolver) InMatchingOffers(ctx context.Context, obj *models.MatchingOfferSummary) ([]*models.MatchingOfferRecord, error) {
	thunk := midacontext.GetLoader[loader.Loader](ctx).InMatchingOfferRecord.Load(ctx, obj.MatchingID)
	record, err := thunk()
	if err != nil {
		return nil, err
	}
	return record.OfferRecords(), nil
}

// OutMatchingOffers is the resolver for the outMatchingOffers field.
func (r *matchingOfferSummaryResolver) OutMatchingOffers(ctx context.Context, obj *models.MatchingOfferSummary) ([]*models.MatchingOfferRecord, error) {
	thunk := midacontext.GetLoader[loader.Loader](ctx).OutMatchingOfferRecord.Load(ctx, obj.MatchingID)
	record, err := thunk()
	if err != nil {
		return nil, err
	}
	return record.OfferRecords(), nil
}

// SendMatchingOffer is the resolver for the sendMatchingOffer field.
func (r *mutationResolver) SendMatchingOffer(ctx context.Context, myMatchingID string, targetMatchingID string) (*string, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() && !token.IsUser() {
		return nil, midacode.ErrNotPermitted
	}
	uid := token.UserID()
	err := modelutil.SendOutOffer(ctx, uid, myMatchingID, targetMatchingID)
	return nil, err
}

// CancelMatchingOffer is the resolver for the cancelMatchingOffer field.
func (r *mutationResolver) CancelMatchingOffer(ctx context.Context, myMatchingID string, targetMatchingID string) (*string, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() && !token.IsUser() {
		return nil, midacode.ErrNotPermitted
	}
	uid := token.UserID()
	err := modelutil.CancelOutOffer(ctx, uid, myMatchingID, targetMatchingID)
	return nil, err
}

// AcceptMatchingOffer is the resolver for the acceptMatchingOffer field.
func (r *mutationResolver) AcceptMatchingOffer(ctx context.Context, myMatchingID string, targetMatchingID string) (*string, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() && !token.IsUser() {
		return nil, midacode.ErrNotPermitted
	}
	uid := token.UserID()
	err := modelutil.AcceptInOffer(ctx, uid, myMatchingID, targetMatchingID)
	return nil, err
}

// RejectMatchingOffer is the resolver for the rejectMatchingOffer field.
func (r *mutationResolver) RejectMatchingOffer(ctx context.Context, myMatchingID string, targetMatchingID string) (*string, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() && !token.IsUser() {
		return nil, midacode.ErrNotPermitted
	}
	uid := token.UserID()
	err := modelutil.RejectInOffer(ctx, uid, myMatchingID, targetMatchingID)
	return nil, err
}

// DiscoverMatchingOfTopic is the resolver for the discoverMatchingOfTopic field.
func (r *queryResolver) DiscoverMatchingOfTopic(ctx context.Context, userID *string, topicID string, filter *models.DiscoverMatchingFilter, nextToken *string) (*models.DiscoverMatchingResult, error) {
	token := midacontext.GetClientToken(ctx)
	uid := graphqlutil.GetID(token, userID)
	allMatchingLoader := midacontext.GetLoader[loader.Loader](ctx).AllMatching
	allMatchingLoader.Load(ctx)
	opt := loader.UserDiscoverOpt{N: 4}

	if filter != nil {
		if filter.CityID != nil {
			opt.CityID = *filter.CityID
		}
		if filter.Gender != nil {
			opt.Gender = models.Gender(*filter.Gender)
		} else {
			opt.Gender = models.GenderN
		}
	}

	var ids []string
	var next string
	if uid == "" {
		ids = allMatchingLoader.LoadForAnoumynous(ctx, topicID, opt)
	} else {
		ids, next = allMatchingLoader.LoadForUser(ctx, uid, topicID, opt)
	}
	if len(ids) == 0 {
		return &models.DiscoverMatchingResult{Matchings: []*models.Matching{}, NextToken: next}, nil
	}
	thunk := midacontext.GetLoader[loader.Loader](ctx).Matching.LoadMany(ctx, ids)
	matchings, errs := thunk()
	if len(errs) > 0 {
		return nil, multierr.Combine(errs...)
	}
	return &models.DiscoverMatchingResult{Matchings: matchings, NextToken: next}, nil
}

// InMatchingOffer is the resolver for the inMatchingOffer field.
func (r *queryResolver) InMatchingOffer(ctx context.Context, matchingID string) ([]*models.MatchingOfferRecord, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() && !token.IsUser() {
		return nil, midacode.ErrNotPermitted
	}

	thunk := midacontext.GetLoader[loader.Loader](ctx).MatchingOfferSummary.Load(ctx, matchingID)
	summary, err := thunk()
	if err != nil {
		return nil, err
	}
	if token.IsUser() {
		if summary.UserID != token.UserID() {
			return nil, midacode.ErrNotPermitted
		}
	}

	recordsThunk := midacontext.GetLoader[loader.Loader](ctx).InMatchingOfferRecord.Load(ctx, matchingID)
	records, err := recordsThunk()
	if err != nil {
		return nil, err
	}
	return records.OfferRecords(), nil
}

// OutMatchingOffer is the resolver for the outMatchingOffer field.
func (r *queryResolver) OutMatchingOffer(ctx context.Context, matchingID string) ([]*models.MatchingOfferRecord, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() && !token.IsUser() {
		return nil, midacode.ErrNotPermitted
	}
	thunk := midacontext.GetLoader[loader.Loader](ctx).MatchingOfferSummary.Load(ctx, matchingID)
	summary, err := thunk()
	if err != nil {
		return nil, err
	}
	if token.IsUser() {
		if summary.UserID != token.UserID() {
			return nil, midacode.ErrNotPermitted
		}
	}

	recordsThunk := midacontext.GetLoader[loader.Loader](ctx).OutMatchingOfferRecord.Load(ctx, matchingID)
	records, err := recordsThunk()
	if err != nil {
		return nil, err
	}
	return records.OfferRecords(), nil
}

// DiscoverMatching returns DiscoverMatchingResolver implementation.
func (r *Resolver) DiscoverMatching() DiscoverMatchingResolver { return &discoverMatchingResolver{r} }

// MatchingOfferRecord returns MatchingOfferRecordResolver implementation.
func (r *Resolver) MatchingOfferRecord() MatchingOfferRecordResolver {
	return &matchingOfferRecordResolver{r}
}

// MatchingOfferSummary returns MatchingOfferSummaryResolver implementation.
func (r *Resolver) MatchingOfferSummary() MatchingOfferSummaryResolver {
	return &matchingOfferSummaryResolver{r}
}

type discoverMatchingResolver struct{ *Resolver }
type matchingOfferRecordResolver struct{ *Resolver }
type matchingOfferSummaryResolver struct{ *Resolver }

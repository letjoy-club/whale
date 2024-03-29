package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.35

import (
	"context"
	"time"
	"whale/pkg/dbquery"
	"whale/pkg/loader"
	"whale/pkg/models"
	"whale/pkg/modelutil"
	"whale/pkg/utils"
	"whale/pkg/whalecode"

	"github.com/golang-module/carbon"
	"github.com/letjoy-club/mida-tool/dbutil"
	"github.com/letjoy-club/mida-tool/graphqlutil"
	"github.com/letjoy-club/mida-tool/midacode"
	"github.com/letjoy-club/mida-tool/midacontext"
	"github.com/samber/lo"
)

// Gender is the resolver for the gender field.
func (r *discoverMotionResolver) Gender(ctx context.Context, obj *models.Motion) (models.Gender, error) {
	return models.Gender(obj.Gender), nil
}

// PreferredPeriods is the resolver for the preferredPeriods field.
func (r *discoverMotionResolver) PreferredPeriods(ctx context.Context, obj *models.Motion) ([]models.DatePeriod, error) {
	return lo.Map(obj.PreferredPeriods, func(v string, i int) models.DatePeriod {
		return models.DatePeriod(v)
	}), nil
}

// Liked is the resolver for the liked field.
func (r *discoverMotionResolver) Liked(ctx context.Context, obj *models.Motion, userID *string) (bool, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsUser() && !token.IsAdmin() {
		return false, nil
	}
	uid := graphqlutil.GetID(token, userID)
	if uid == "" {
		return false, nil
	}
	thunk := midacontext.GetLoader[loader.Loader](ctx).UserLikeMotion.Load(ctx, uid)
	u, err := thunk()
	if err != nil {
		return false, err
	}
	return u.IsLike(obj.ID), nil
}

// Submitted is the resolver for the submitted field.
func (r *discoverMotionResolver) Submitted(ctx context.Context, obj *models.Motion, userID *string) (bool, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsUser() && !token.IsAdmin() {
		return false, nil
	}
	uid := graphqlutil.GetID(token, userID)
	if uid == "" {
		return false, nil
	}
	thunk := midacontext.GetLoader[loader.Loader](ctx).UserSubmitMotion.Load(ctx, uid)
	u, err := thunk()
	if err != nil {
		return false, err
	}
	return u.IsSubmitted(obj.ID), nil
}

// ThumbsUp is the resolver for the thumbsUp field.
func (r *discoverMotionResolver) ThumbsUp(ctx context.Context, obj *models.Motion, userID *string) (bool, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsUser() && !token.IsAdmin() {
		return false, nil
	}
	uid := graphqlutil.GetID(token, userID)
	if uid == "" {
		return false, nil
	}
	thunk := midacontext.GetLoader[loader.Loader](ctx).UserThumbsUpMotion.Load(ctx, uid)
	u, err := thunk()
	if err != nil {
		return false, err
	}
	return u.ThumbsUp(obj.ID), nil
}

// Topic is the resolver for the topic field.
func (r *discoverMotionResolver) Topic(ctx context.Context, obj *models.Motion) (*models.Topic, error) {
	return &models.Topic{ID: obj.TopicID}, nil
}

// TopicOptionConfig is the resolver for the topicOptionConfig field.
func (r *discoverMotionResolver) TopicOptionConfig(ctx context.Context, obj *models.Motion) (*models.TopicOptionConfig, error) {
	return &models.TopicOptionConfig{TopicID: obj.TopicID}, nil
}

// User is the resolver for the user field.
func (r *discoverMotionResolver) User(ctx context.Context, obj *models.Motion) (*models.User, error) {
	return &models.User{ID: obj.UserID}, nil
}

// City is the resolver for the city field.
func (r *discoverMotionResolver) City(ctx context.Context, obj *models.Motion) (*models.Area, error) {
	return &models.Area{Code: obj.CityID}, nil
}

// Areas is the resolver for the areas field.
func (r *discoverMotionResolver) Areas(ctx context.Context, obj *models.Motion) ([]*models.Area, error) {
	return lo.Map(obj.AreaIDs, func(v string, i int) *models.Area {
		return &models.Area{Code: v}
	}), nil
}

// PreferredPeriods is the resolver for the preferredPeriods field.
func (r *motionResolver) PreferredPeriods(ctx context.Context, obj *models.Motion) ([]models.DatePeriod, error) {
	return lo.Map(obj.PreferredPeriods, func(v string, i int) models.DatePeriod {
		return models.DatePeriod(v)
	}), nil
}

// Gender is the resolver for the gender field.
func (r *motionResolver) Gender(ctx context.Context, obj *models.Motion) (models.Gender, error) {
	return models.Gender(obj.Gender), nil
}

// Liked is the resolver for the liked field.
func (r *motionResolver) Liked(ctx context.Context, obj *models.Motion) (bool, error) {
	token := midacontext.GetClientToken(ctx)
	if token.IsUser() {
		thunk := midacontext.GetLoader[loader.Loader](ctx).UserLikeMotion.Load(ctx, obj.UserID)
		u, err := thunk()
		if err != nil {
			return false, err
		}
		return u.IsLike(obj.ID), nil
	}
	return false, nil
}

// ThumbsUp is the resolver for the thumbsUp field.
func (r *motionResolver) ThumbsUp(ctx context.Context, obj *models.Motion) (bool, error) {
	token := midacontext.GetClientToken(ctx)
	if token.IsUser() {
		thunk := midacontext.GetLoader[loader.Loader](ctx).UserThumbsUpMotion.Load(ctx, obj.UserID)
		u, err := thunk()
		if err != nil {
			return false, err
		}
		return u.ThumbsUp(obj.ID), nil
	}
	return false, nil
}

// Topic is the resolver for the topic field.
func (r *motionResolver) Topic(ctx context.Context, obj *models.Motion) (*models.Topic, error) {
	return &models.Topic{ID: obj.TopicID}, nil
}

// TopicOptionConfig is the resolver for the topicOptionConfig field.
func (r *motionResolver) TopicOptionConfig(ctx context.Context, obj *models.Motion) (*models.TopicOptionConfig, error) {
	return &models.TopicOptionConfig{TopicID: obj.TopicID}, nil
}

// User is the resolver for the user field.
func (r *motionResolver) User(ctx context.Context, obj *models.Motion) (*models.User, error) {
	return &models.User{ID: obj.UserID}, nil
}

// City is the resolver for the city field.
func (r *motionResolver) City(ctx context.Context, obj *models.Motion) (*models.Area, error) {
	return &models.Area{Code: obj.CityID}, nil
}

// Areas is the resolver for the areas field.
func (r *motionResolver) Areas(ctx context.Context, obj *models.Motion) ([]*models.Area, error) {
	return lo.Map(obj.AreaIDs, func(v string, i int) *models.Area {
		return &models.Area{Code: v}
	}), nil
}

// State is the resolver for the state field.
func (r *motionOfferRecordResolver) State(ctx context.Context, obj *models.MotionOfferRecord) (models.MotionOfferState, error) {
	return models.MotionOfferState(obj.State), nil
}

// Reviewed is the resolver for the reviewed field.
func (r *motionOfferRecordResolver) Reviewed(ctx context.Context, obj *models.MotionOfferRecord, userID *string) (bool, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() && !token.IsUser() {
		return false, nil
	}
	// 查询是否已经被评价时，需要提供评价者 id
	uid := graphqlutil.GetID(token, userID)
	if obj.UserID != uid && obj.ToUserID != uid {
		// id 不一致则无法查看
		return false, midacode.ErrNotPermitted
	}

	thunk := midacontext.GetLoader[loader.Loader](ctx).MotionReviewed.Load(ctx, obj.ID)
	motionReviewed, err := thunk()
	if err != nil {
		return false, err
	}
	return motionReviewed.IsReviewed(uid), nil
}

// ToMotion is the resolver for the toMotion field.
func (r *motionOfferRecordResolver) ToMotion(ctx context.Context, obj *models.MotionOfferRecord) (*models.Motion, error) {
	thunk := midacontext.GetLoader[loader.Loader](ctx).Motion.Load(ctx, obj.ToMotionID)
	return thunk()
}

// Motion is the resolver for the motion field.
func (r *motionOfferRecordResolver) Motion(ctx context.Context, obj *models.MotionOfferRecord) (*models.Motion, error) {
	thunk := midacontext.GetLoader[loader.Loader](ctx).Motion.Load(ctx, obj.MotionID)
	return thunk()
}

// GetAvailableMotionOffer is the resolver for the getAvailableMotionOffer field.
func (r *mutationResolver) GetAvailableMotionOffer(ctx context.Context, userID *string, targetMotionID string) (*models.AvailableMotionOffer, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsUser() && !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}
	thunk := midacontext.GetLoader[loader.Loader](ctx).Motion.Load(ctx, targetMotionID)
	motion, err := thunk()
	if err != nil {
		return nil, err
	}
	if !motion.Active {
		return nil, whalecode.ErrTheMotionIsNotActive
	}

	uid := graphqlutil.GetID(token, userID)
	db := dbutil.GetDB(ctx)
	Motion := dbquery.Use(db).Motion
	motion, err = Motion.WithContext(ctx).
		Where(Motion.UserID.Eq(uid), Motion.TopicID.Eq(motion.TopicID)).
		Where(Motion.Active.Is(true)).
		Take()
	if err != nil {
		if midacode.ItemIsNotFound(err) == midacode.ErrItemNotFound {
			return &models.AvailableMotionOffer{}, nil
		}
		return nil, err
	}

	next := time.Now()
	return &models.AvailableMotionOffer{Motion: motion, NextQuotaTime: &next}, nil
}

// CreateMotionOffer is the resolver for the createMotionOffer field.
func (r *mutationResolver) CreateMotionOffer(ctx context.Context, myMotionID string, targetMotionID string) (*models.CreateMotionOfferResult, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsUser() && !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}
	groupId, err := modelutil.CreateMotionOffer(ctx, token.UserID(), myMotionID, targetMotionID)
	if err != nil {
		return nil, err
	}
	return &models.CreateMotionOfferResult{ChatGroupID: groupId}, nil
}

// CancelMotionOffer is the resolver for the cancelMotionOffer field.
func (r *mutationResolver) CancelMotionOffer(ctx context.Context, myMotionID string, targetMotionID string) (*string, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsUser() && !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}
	panic("CancelMotionOffer not support now")
	//return nil, modelutil.CancelMotionOffer(ctx, token.UserID(), myMotionID, targetMotionID)
}

// AcceptMotionOffer is the resolver for the acceptMotionOffer field.
func (r *mutationResolver) AcceptMotionOffer(ctx context.Context, myMotionID string, targetMotionID string) (*string, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsUser() && !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}
	return nil, modelutil.AcceptMotionOffer(ctx, token.UserID(), myMotionID, targetMotionID)
}

// RejectMotionOffer is the resolver for the rejectMotionOffer field.
func (r *mutationResolver) RejectMotionOffer(ctx context.Context, myMotionID string, targetMotionID string) (*string, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsUser() && !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}
	return nil, modelutil.RejectMotionOffer(ctx, token.UserID(), myMotionID, targetMotionID)
}

// SendChatInOffer is the resolver for the sendChatInOffer field.
func (r *mutationResolver) SendChatInOffer(ctx context.Context, myMotionID string, targetMotionID string, sentence string) (*string, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsUser() && !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}
	return nil, modelutil.SendChatInOffer(ctx, token.UserID(), myMotionID, targetMotionID, sentence)
}

// FinishMotionOffer is the resolver for the finishMotionOffer field.
func (r *mutationResolver) FinishMotionOffer(ctx context.Context, fromMotionID string, toMotionID string) (*string, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsUser() && !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}
	return nil, modelutil.FinishMotionOffer(ctx, token.UserID(), fromMotionID, toMotionID)
}

// NotifyNewMotionOffer is the resolver for the notifyNewMotionOffer field.
func (r *mutationResolver) NotifyNewMotionOffer(ctx context.Context, param *models.NotifyNewMotionOfferMessageParam) (*string, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}

	if param == nil {
		startTime := carbon.Now().StartOfHour().SubHour()
		endTime := startTime.AddHour()
		return nil, modelutil.NotifyNewMotionOffer(ctx, startTime.ToStdTime(), endTime.ToStdTime())
	}
	return nil, modelutil.NotifyNewMotionOffer(ctx, param.Begin, param.End)
}

// SendMotionOfferAcceptMessage is the resolver for the sendMotionOfferAcceptMessage field.
func (r *mutationResolver) SendMotionOfferAcceptMessage(ctx context.Context, id int) (*string, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}
	db := dbutil.GetDB(ctx)
	MotionOfferRecord := dbquery.Use(db).MotionOfferRecord
	offer, err := MotionOfferRecord.WithContext(ctx).
		Where(MotionOfferRecord.ID.Eq(id)).
		Select(MotionOfferRecord.MotionID, MotionOfferRecord.UserID, MotionOfferRecord.ToUserID, MotionOfferRecord.ChatGroupID).
		Take()
	if err != nil {
		return nil, err
	}

	Motion := dbquery.Use(db).Motion
	motion, err := Motion.WithContext(ctx).Where(Motion.ID.Eq(offer.MotionID)).Select(Motion.TopicID).Take()
	if err != nil {
		return nil, err
	}

	return nil, modelutil.SendMotionOfferAcceptedMessage(ctx, motion.TopicID, offer)
}

// DiscoverCategoryMotions is the resolver for the discoverCategoryMotions field.
func (r *queryResolver) DiscoverCategoryMotions(ctx context.Context, userID *string, filter models.DiscoverTopicCategoryMotionFilter, topicCategoryID *string, nextToken *string) (*models.DiscoverMotionResult, error) {
	err := midacontext.GetLoader[loader.Loader](ctx).AllMotion.Load(ctx)
	if err != nil {
		return nil, err
	}
	next := ""
	if nextToken != nil {
		next = *nextToken
	}

	opt := loader.UserDiscoverMotionOpt{
		N:         6,
		NextToken: next,
		Gender:    models.GenderN,
		Type:      models.MotionTypeAll,
	}

	if filter.CategoryID != nil {
		opt.CategoryID = *filter.CategoryID
	} else {
		if topicCategoryID != nil {
			opt.CategoryID = *topicCategoryID
		}
	}

	if opt.CategoryID == "" {
		opt.CategoryID = loader.AllCategoryID
	}

	if filter.Type != nil {
		opt.Type = *filter.Type
	}

	if filter.CityID != nil {
		opt.CityID = *filter.CityID
	}
	if filter.Gender != nil {
		opt.Gender = *filter.Gender
	}
	if len(filter.TopicIds) > 0 {
		opt.TopicIDs = filter.TopicIds
	}
	opt.NextToken = next

	token := midacontext.GetClientToken(ctx)

	uid := graphqlutil.GetID(token, userID)
	var ids []string
	var retNext string
	if uid != "" {
		ids, retNext = midacontext.GetLoader[loader.Loader](ctx).AllMotion.LoadForUser(ctx, uid, opt)
	} else {
		ids = midacontext.GetLoader[loader.Loader](ctx).AllMotion.LoadForAnoumynous(ctx, opt)
	}
	thunk := midacontext.GetLoader[loader.Loader](ctx).Motion.LoadMany(ctx, ids)
	motions, err := utils.ReturnThunk(thunk)
	if err != nil {
		return nil, err
	}
	return &models.DiscoverMotionResult{Motions: motions, NextToken: retNext}, nil
}

// DiscoverLatestCategoryMotions is the resolver for the discoverLatestCategoryMotions field.
func (r *queryResolver) DiscoverLatestCategoryMotions(ctx context.Context, filter models.DiscoverTopicCategoryMotionFilter, topicCategoryID *string, lastID *string) ([]*models.Motion, error) {
	err := midacontext.GetLoader[loader.Loader](ctx).AllMotion.Load(ctx)
	if err != nil {
		return nil, err
	}
	opt := loader.UserDiscoverMotionOpt{
		N:      6,
		Gender: models.GenderN,
		Type:   models.MotionTypeAll,
	}

	if filter.CategoryID != nil {
		opt.CategoryID = *filter.CategoryID
	} else {
		if topicCategoryID != nil {
			opt.CategoryID = *topicCategoryID
		}
	}
	if opt.CategoryID == "" {
		opt.CategoryID = loader.AllCategoryID
	}

	if filter.Type != nil {
		opt.Type = *filter.Type
	}

	if opt.CategoryID == "" {
		opt.CategoryID = loader.AllCategoryID
	}

	if lastID != nil {
		opt.LastID = *lastID
	}
	if filter.CityID != nil {
		opt.CityID = *filter.CityID
	}
	if filter.Gender != nil {
		opt.Gender = *filter.Gender
	}
	if len(filter.TopicIds) > 0 {
		opt.TopicIDs = filter.TopicIds
	}

	var ids []string
	ids = midacontext.GetLoader[loader.Loader](ctx).AllMotion.GetOrderedMotions(ctx, opt)

	thunk := midacontext.GetLoader[loader.Loader](ctx).Motion.LoadMany(ctx, ids)
	motions, err := utils.ReturnThunk(thunk)
	if err != nil {
		return nil, err
	}

	token := midacontext.GetClientToken(ctx)

	if token.IsUser() {
		userDiscover := midacontext.GetLoader[loader.Loader](ctx).AllMotion.GenUserDiscoverMotion(ctx, token.String())
		viewedNum := userDiscover.Viewed(ids)

		// 如果查看的数量多于 300，就不再返回
		if viewedNum == 0 {
			return []*models.Motion{}, nil
		}
	}
	return motions, nil
}

// GetDiscoverMotion is the resolver for the getDiscoverMotion field.
func (r *queryResolver) GetDiscoverMotion(ctx context.Context, motionID string) (*models.Motion, error) {
	thunk := midacontext.GetLoader[loader.Loader](ctx).Motion.Load(ctx, motionID)
	return thunk()
}

// OutMotionOffers is the resolver for the outMotionOffers field.
func (r *queryResolver) OutMotionOffers(ctx context.Context, motionID string) ([]*models.MotionOfferRecord, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsUser() && !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}
	thunk := midacontext.GetLoader[loader.Loader](ctx).Motion.Load(ctx, motionID)
	motion, err := thunk()
	if err != nil {
		return nil, err
	}
	if token.IsUser() {
		if motion.UserID != token.String() {
			return nil, midacode.ErrNotPermitted
		}
	}
	offerThunk := midacontext.GetLoader[loader.Loader](ctx).OutMotionOfferRecord.Load(ctx, motionID)
	offers, err := offerThunk()
	if err != nil {
		return nil, err
	}
	return offers.Offers, nil
}

// InMotionOffers is the resolver for the inMotionOffers field.
func (r *queryResolver) InMotionOffers(ctx context.Context, motionID string) ([]*models.MotionOfferRecord, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsUser() && !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}
	thunk := midacontext.GetLoader[loader.Loader](ctx).Motion.Load(ctx, motionID)
	motion, err := thunk()
	if err != nil {
		return nil, err
	}
	if token.IsUser() {
		if motion.UserID != token.String() {
			return nil, midacode.ErrNotPermitted
		}
	}
	offerThunk := midacontext.GetLoader[loader.Loader](ctx).InMotionOfferRecord.Load(ctx, motionID)
	offers, err := offerThunk()
	if err != nil {
		return nil, err
	}
	return offers.Offers, nil
}

// GetMotionOffer is the resolver for the getMotionOffer field.
func (r *queryResolver) GetMotionOffer(ctx context.Context, motionID string, toMotionID string) (*models.MotionOfferRecord, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsUser() && !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}
	uid := token.UserID()
	db := dbutil.GetDB(ctx)
	MotionOfferRecord := dbquery.Use(db).MotionOfferRecord
	record, err := MotionOfferRecord.WithContext(ctx).Where(MotionOfferRecord.MotionID.Eq(motionID)).Where(MotionOfferRecord.ToMotionID.Eq(toMotionID)).Take()
	if err != nil {
		return nil, midacode.ItemMayNotFound(err)
	}
	if token.IsUser() && uid != record.UserID && uid != record.ToUserID {
		return nil, midacode.ErrNotPermitted
	}
	return record, nil
}

// DiscoverMotion returns DiscoverMotionResolver implementation.
func (r *Resolver) DiscoverMotion() DiscoverMotionResolver { return &discoverMotionResolver{r} }

// Motion returns MotionResolver implementation.
func (r *Resolver) Motion() MotionResolver { return &motionResolver{r} }

// MotionOfferRecord returns MotionOfferRecordResolver implementation.
func (r *Resolver) MotionOfferRecord() MotionOfferRecordResolver {
	return &motionOfferRecordResolver{r}
}

type discoverMotionResolver struct{ *Resolver }
type motionResolver struct{ *Resolver }
type motionOfferRecordResolver struct{ *Resolver }

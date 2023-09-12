package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.35

import (
	"context"
	"fmt"
	"whale/pkg/dbquery"
	"whale/pkg/loader"
	"whale/pkg/models"
	"whale/pkg/modelutil"
	"whale/pkg/utils"
	"whale/pkg/whalecode"

	"github.com/letjoy-club/mida-tool/dbutil"
	"github.com/letjoy-club/mida-tool/graphqlutil"
	"github.com/letjoy-club/mida-tool/keyer"
	"github.com/letjoy-club/mida-tool/logger"
	"github.com/letjoy-club/mida-tool/midacode"
	"github.com/letjoy-club/mida-tool/midacontext"
	"github.com/letjoy-club/mida-tool/redisutil"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"gorm.io/gen/field"
)

// CreateMotion is the resolver for the createMotion field.
func (r *mutationResolver) CreateMotion(ctx context.Context, userID *string, param models.CreateMotionParam) (*models.Motion, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() && !token.IsUser() {
		return nil, midacode.ErrNotPermitted
	}
	uid := graphqlutil.GetID(token, userID)
	if uid == "" {
		return nil, whalecode.ErrUserIDCannotBeEmpty
	}

	release, err := redisutil.LockAll(ctx, keyer.UserMatching(uid))
	if err != nil {
		return nil, err
	}
	defer release(ctx)
	return modelutil.CreateMotion(ctx, uid, &param)
}

// UpdateMotion is the resolver for the updateMotion field.
func (r *mutationResolver) UpdateMotion(ctx context.Context, id string, param models.UpdateMotionParam) (*models.Motion, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() && !token.IsUser() {
		return nil, midacode.ErrNotPermitted
	}
	if param.Remark == nil {
		emptyStr := ""
		param.Remark = &emptyStr
	}
	release, err := redisutil.LockAll(ctx, keyer.UserMotion(id))
	if err != nil {
		return nil, err
	}
	defer release(ctx)
	err = modelutil.UpdateMotion(ctx, id, &param)
	if err != nil {
		return nil, err
	}
	midacontext.GetLoader[loader.Loader](ctx).Motion.Clear(ctx, id)
	thunk := midacontext.GetLoader[loader.Loader](ctx).Motion.Load(ctx, id)
	return thunk()
}

// UserUpdateMotion is the resolver for the userUpdateMotion field.
func (r *mutationResolver) UserUpdateMotion(ctx context.Context, myMotionID string, param models.UserUpdateMotionParam) (*models.Motion, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() && !token.IsUser() {
		return nil, midacode.ErrNotPermitted
	}
	thunk := midacontext.GetLoader[loader.Loader](ctx).Motion.Load(ctx, myMotionID)
	_, err := thunk()
	if err != nil {
		return nil, err
	}

	db := dbutil.GetDB(ctx)
	Motion := dbquery.Use(db).Motion

	fields := []field.AssignExpr{}
	if param.Gender != nil {
		fields = append(fields, Motion.Gender.Value(*param.Gender))
	}
	if param.Remark != nil {
		fields = append(fields, Motion.Remark.Value(*param.Remark))
	}
	if param.AreaIds != nil {
		fields = append(fields, Motion.AreaIDs.Value(graphqlutil.ElementList[string](param.AreaIds)))
	}
	if param.DayRange != nil {
		fields = append(fields, Motion.DayRange.Value(graphqlutil.ElementList[string](param.DayRange)))
	}
	if param.Properties != nil {
		fields = append(fields, Motion.Properties.Value(graphqlutil.ElementList[*models.MotionPropertyParam](param.Properties)))
	}
	if param.PreferredPeriods != nil {
		fields = append(fields, Motion.Properties.Value(graphqlutil.ElementList[models.DatePeriod](param.PreferredPeriods)))
	}
	_, err = Motion.WithContext(ctx).Where(Motion.ID.Eq(myMotionID)).UpdateSimple(fields...)
	if err == nil {
		midacontext.GetLoader[loader.Loader](ctx).Motion.Clear(ctx, myMotionID)
	}
	return nil, err
}

// CloseMotion is the resolver for the closeMotion field.
func (r *mutationResolver) CloseMotion(ctx context.Context, id string) (*string, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() && !token.IsUser() {
		return nil, midacode.ErrNotPermitted
	}

	err := modelutil.CloseMotion(ctx, token.UserID(), id)
	return nil, err
}

// ReviewMotionOffer is the resolver for the reviewMotionOffer field.
func (r *mutationResolver) ReviewMotionOffer(ctx context.Context, userID *string, fromMotionID string, toMotionID string, param models.ReviewMotionParam) (*string, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() && !token.IsUser() {
		return nil, midacode.ErrNotPermitted
	}
	uid := graphqlutil.GetID(token, userID)
	// 校验Motion
	motionThunk := midacontext.GetLoader[loader.Loader](ctx).Motion.LoadMany(ctx, []string{fromMotionID, toMotionID})
	motions, errors := motionThunk()
	if errors != nil {
		return nil, multierr.Combine(errors...)
	}
	fromMotion, toMotion := motions[0], motions[1]

	// 必须要邀约参与方才能评价
	if fromMotion.UserID != uid && toMotion.UserID != uid {
		return nil, midacode.ErrNotPermitted
	}

	db := dbutil.GetDB(ctx)
	// 校验MotionOfferRecord
	MotionOfferRecord := dbquery.Use(db).MotionOfferRecord
	record, err := MotionOfferRecord.WithContext(ctx).Where(
		MotionOfferRecord.MotionID.Eq(fromMotionID),
		MotionOfferRecord.ToMotionID.Eq(toMotionID),
	).Take()
	if err != nil {
		return nil, err
	}
	if record.State != string(models.MotionOfferStateAccepted) && record.State != string(models.MotionOfferStateFinished) {
		return nil, whalecode.ErrMotionStateCannotReview
	}
	// 重复评价判断
	thunk := midacontext.GetLoader[loader.Loader](ctx).MotionReviewed.Load(ctx, record.ID)
	motionReviewed, err := thunk()
	if err != nil {
		return nil, err
	}
	if motionReviewed.IsReviewed(uid) {
		return nil, whalecode.ErrMotionReviewAlreadyExists
	}

	// 被评价者ID
	toUserId := fromMotion.UserID
	if fromMotion.UserID == uid {
		toUserId = toMotion.UserID
	}

	MotionReview := dbquery.Use(db).MotionReview
	if err = MotionReview.WithContext(ctx).Create(&models.MotionReview{
		MotionOfferID: record.ID,
		ReviewerID:    uid,
		ToUserID:      toUserId,
		TopicID:       fromMotion.TopicID,
		Score:         param.Score,
		Comment:       param.Comment,
	}); err != nil {
		return nil, err
	}

	if err := modelutil.PublishUserReviewEvent(ctx, uid, toUserId, string(models.ResultCreatedByOffer), fmt.Sprintf("%d", record.ID)); err != nil {
		logger.L.Error("ReviewMotionOffer - PublishUserReviewEvent error",
			zap.Error(err), zap.Any("motionOfferId", record.ID))
	}
	midacontext.GetLoader[loader.Loader](ctx).MotionReviewed.Clear(ctx, record.ID)

	return nil, err
}

// MarkMotionExpire is the resolver for the markMotionExpire field.
func (r *mutationResolver) MarkMotionExpire(ctx context.Context) (*string, error) {
	return nil, modelutil.MotionExpire(ctx)
}

// Motion is the resolver for the motion field.
func (r *queryResolver) Motion(ctx context.Context, id string) (*models.Motion, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() && !token.IsUser() {
		return nil, midacode.ErrNotPermitted
	}
	thunk := midacontext.GetLoader[loader.Loader](ctx).Motion.Load(ctx, id)
	motion, err := thunk()
	if err != nil {
		return nil, err
	}
	if token.IsUser() {
		if motion.UserID != token.UserID() {
			return nil, midacode.ErrNotPermitted
		}
	}
	return motion, nil
}

// UserMotions is the resolver for the userMotions field.
func (r *queryResolver) UserMotions(ctx context.Context, userID *string, paginator *graphqlutil.GraphQLPaginator) ([]*models.Motion, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() && !token.IsUser() {
		return nil, midacode.ErrNotPermitted
	}
	uid := graphqlutil.GetID(token, userID)
	if uid == "" {
		return nil, whalecode.ErrUserIDCannotBeEmpty
	}
	pager := graphqlutil.GetPager(paginator)

	db := dbutil.GetDB(ctx)
	Motion := dbquery.Use(db).Motion
	ids := []string{}

	err := Motion.WithContext(ctx).Where(Motion.UserID.Eq(uid)).Limit(pager.Limit()).Offset(pager.Offset()).Pluck(Motion.ID, &ids)
	if err != nil {
		return nil, err
	}
	thunk := midacontext.GetLoader[loader.Loader](ctx).Motion.LoadMany(ctx, ids)
	return utils.ReturnThunk(thunk)
}

// UserMotionsCount is the resolver for the userMotionsCount field.
func (r *queryResolver) UserMotionsCount(ctx context.Context, userID *string) (*models.Summary, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() && !token.IsUser() {
		return nil, midacode.ErrNotPermitted
	}
	uid := graphqlutil.GetID(token, userID)
	if uid == "" {
		return nil, whalecode.ErrUserIDCannotBeEmpty
	}

	db := dbutil.GetDB(ctx)
	Motion := dbquery.Use(db).Motion

	count, err := Motion.WithContext(ctx).Where(Motion.UserID.Eq(uid)).Count()
	return &models.Summary{Count: int(count)}, err
}

// ActiveMotions is the resolver for the activeMotions field.
func (r *queryResolver) ActiveMotions(ctx context.Context, userID *string) ([]*models.Motion, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() && !token.IsUser() {
		return nil, midacode.ErrNotPermitted
	}
	uid := graphqlutil.GetID(token, userID)
	if uid == "" {
		return nil, whalecode.ErrUserIDCannotBeEmpty
	}
	db := dbutil.GetDB(ctx)
	Motion := dbquery.Use(db).Motion
	ids := []string{}
	err := Motion.WithContext(ctx).
		Where(Motion.UserID.Eq(uid)).
		Where(Motion.Active.Is(true)).
		Order(Motion.ID.Desc()).
		Pluck(Motion.ID, &ids)
	if err != nil {
		return nil, err
	}
	thunk := midacontext.GetLoader[loader.Loader](ctx).Motion.LoadMany(ctx, ids)
	return utils.ReturnThunk(thunk)
}

// Motions is the resolver for the motions field.
func (r *queryResolver) Motions(ctx context.Context, filter *models.MotionFilter, paginator *graphqlutil.GraphQLPaginator) ([]*models.Motion, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}
	db := dbutil.GetDB(ctx)
	Motion := dbquery.Use(db).Motion
	ids := []string{}

	query := Motion.WithContext(ctx)
	if filter != nil {
		if filter.UserID != nil {
			query = query.Where(Motion.UserID.Eq(*filter.UserID))
		}
		if filter.Gender != nil {
			query = query.Where(Motion.Gender.Eq(filter.Gender.String()))
		}
		if filter.CityID != nil {
			query = query.Where(Motion.CityID.Eq(*filter.CityID))
		}
		if filter.TopicID != nil {
			query = query.Where(Motion.TopicID.Eq(*filter.TopicID))
		}
		if filter.Before != nil {
			query = query.Where(Motion.CreatedAt.Lt(*filter.Before))
		}
		if filter.After != nil {
			query = query.Where(Motion.CreatedAt.Gt(*filter.After))
		}
	}

	pager := graphqlutil.GetPager(paginator)

	err := query.Limit(pager.Limit()).Offset(pager.Offset()).Order(Motion.ID.Desc()).Pluck(Motion.ID, &ids)
	if err != nil {
		return nil, nil
	}
	thunk := midacontext.GetLoader[loader.Loader](ctx).Motion.LoadMany(ctx, ids)
	return utils.ReturnThunk(thunk)
}

// MotionsCount is the resolver for the motionsCount field.
func (r *queryResolver) MotionsCount(ctx context.Context, filter *models.MotionFilter) (*models.Summary, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}
	db := dbutil.GetDB(ctx)
	Motion := dbquery.Use(db).Motion

	query := Motion.WithContext(ctx)
	if filter != nil {
		if filter.UserID != nil {
			query = query.Where(Motion.UserID.Eq(*filter.UserID))
		}
		if filter.Gender != nil {
			query = query.Where(Motion.Gender.Eq(filter.Gender.String()))
		}
		if filter.CityID != nil {
			query = query.Where(Motion.CityID.Eq(*filter.CityID))
		}
		if filter.TopicID != nil {
			query = query.Where(Motion.TopicID.Eq(*filter.TopicID))
		}
		if filter.Before != nil {
			query = query.Where(Motion.CreatedAt.Lt(*filter.Before))
		}
		if filter.After != nil {
			query = query.Where(Motion.CreatedAt.Gt(*filter.After))
		}
	}

	count, err := query.Count()
	return &models.Summary{Count: int(count)}, err
}

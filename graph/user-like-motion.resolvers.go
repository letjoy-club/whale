package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.35

import (
	"context"
	"whale/pkg/dbquery"
	"whale/pkg/loader"
	"whale/pkg/models"
	"whale/pkg/whalecode"

	"github.com/letjoy-club/mida-tool/dbutil"
	"github.com/letjoy-club/mida-tool/graphqlutil"
	"github.com/letjoy-club/mida-tool/midacode"
	"github.com/letjoy-club/mida-tool/midacontext"
)

// LikeMotion is the resolver for the likeMotion field.
func (r *mutationResolver) LikeMotion(ctx context.Context, userID *string, motionID string) (int, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() && !token.IsUser() {
		return 0, midacode.ErrNotPermitted
	}
	uid := graphqlutil.GetID(token, userID)
	if uid == "" {
		return 0, whalecode.ErrUserIDCannotBeEmpty
	}
	thunk := midacontext.GetLoader[loader.Loader](ctx).Motion.Load(ctx, motionID)
	motion, err := thunk()
	if err != nil {
		return 0, err
	}

	userLikeMotionThunk := midacontext.GetLoader[loader.Loader](ctx).UserLikeMotion.Load(ctx, uid)
	likeMotion, err := userLikeMotionThunk()
	if err != nil {
		return 0, err
	}
	if likeMotion.IsLike(motionID) {
		return motion.LikeCount, nil
	}

	db := dbutil.GetDB(ctx)
	UserLikeMotion := dbquery.Use(db).UserLikeMotion
	Motion := dbquery.Use(db).Motion
	err = UserLikeMotion.WithContext(ctx).Create(&models.UserLikeMotion{
		ToMotionID: motionID,
		ToUserID:   motion.UserID,
		UserID:     uid,
	})
	if err != nil {
		return motion.LikeCount, nil
	}
	motion.LikeCount += 1
	if _, err := Motion.WithContext(ctx).Where(Motion.ID.Eq(motion.ID)).UpdateSimple(Motion.LikeCount.Add(1)); err != nil {
		return motion.LikeCount, err
	}
	midacontext.GetLoader[loader.Loader](ctx).UserLikeMotion.Clear(ctx, uid)
	midacontext.GetLoader[loader.Loader](ctx).Motion.Clear(ctx, motionID)
	return motion.LikeCount, err
}

// UnlikeMotion is the resolver for the unlikeMotion field.
func (r *mutationResolver) UnlikeMotion(ctx context.Context, userID *string, motionID string) (int, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() && !token.IsUser() {
		return 0, midacode.ErrNotPermitted
	}
	uid := graphqlutil.GetID(token, userID)
	if uid == "" {
		return 0, whalecode.ErrUserIDCannotBeEmpty
	}
	thunk := midacontext.GetLoader[loader.Loader](ctx).Motion.Load(ctx, motionID)
	motion, err := thunk()
	if err != nil {
		return 0, err
	}

	userLikeMotionThunk := midacontext.GetLoader[loader.Loader](ctx).UserLikeMotion.Load(ctx, uid)
	likeMotion, err := userLikeMotionThunk()
	if err != nil {
		return 0, err
	}
	if !likeMotion.IsLike(motionID) {
		return motion.LikeCount, nil
	}

	db := dbutil.GetDB(ctx)
	UserLikeMotion := dbquery.Use(db).UserLikeMotion
	Motion := dbquery.Use(db).Motion
	rx, err := UserLikeMotion.WithContext(ctx).Where(UserLikeMotion.UserID.Eq(uid), UserLikeMotion.ToMotionID.Eq(motionID)).Delete()
	if err != nil {
		return motion.LikeCount, err
	}
	if rx.RowsAffected == 0 {
		return motion.LikeCount, nil
	}
	if _, err := Motion.WithContext(ctx).Where(Motion.ID.Eq(motionID)).UpdateSimple(Motion.LikeCount.Add(-1)); err != nil {
		return motion.LikeCount, err
	}
	motion.LikeCount -= 1
	midacontext.GetLoader[loader.Loader](ctx).UserLikeMotion.Clear(ctx, uid)
	midacontext.GetLoader[loader.Loader](ctx).Motion.Clear(ctx, motionID)
	return motion.LikeCount, nil
}

// ThumbsUpMotion is the resolver for the thumbsUpMotion field.
func (r *mutationResolver) ThumbsUpMotion(ctx context.Context, userID *string, motionID string) (*string, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() && !token.IsUser() {
		return nil, midacode.ErrNotPermitted
	}
	uid := graphqlutil.GetID(token, userID)
	if uid == "" {
		return nil, whalecode.ErrUserIDCannotBeEmpty
	}

	thunk := midacontext.GetLoader[loader.Loader](ctx).Motion.Load(ctx, motionID)
	motion, err := thunk()
	if err != nil {
		return nil, err
	}

	userThumbsUpThunk := midacontext.GetLoader[loader.Loader](ctx).UserThumbsUpMotion.Load(ctx, uid)
	thumbsUp, err := userThumbsUpThunk()
	if err != nil {
		return nil, err
	}
	if thumbsUp.ThumbsUp(motionID) {
		return nil, nil
	}

	db := dbutil.GetDB(ctx)
	err = dbquery.Use(db).Transaction(func(tx *dbquery.Query) error {
		UserThumbsUpMotion := tx.UserThumbsUpMotion
		thumbsUpMotion := models.UserThumbsUpMotion{
			ToUserID:   motion.UserID,
			ToMotionID: motionID,
			UserID:     uid,
		}
		err := UserThumbsUpMotion.WithContext(ctx).Create(&thumbsUpMotion)
		if err != nil {
			return err
		}
		Motion := tx.Motion
		if _, err := Motion.WithContext(ctx).Where(Motion.ID.Eq(motionID)).UpdateSimple(Motion.ThumbsUpCount.Add(1)); err != nil {
			return err
		}
		thumbsUp.DoThumbsUp(motionID)
		midacontext.GetLoader[loader.Loader](ctx).Motion.Clear(ctx, motionID)
		return nil
	})
	return nil, err
}

// CancelThumbsUpMotion is the resolver for the cancelThumbsUpMotion field.
func (r *mutationResolver) CancelThumbsUpMotion(ctx context.Context, userID *string, motionID string) (*string, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() && !token.IsUser() {
		return nil, midacode.ErrNotPermitted
	}
	uid := graphqlutil.GetID(token, userID)
	if uid == "" {
		return nil, whalecode.ErrUserIDCannotBeEmpty
	}

	userThumbsUpThunk := midacontext.GetLoader[loader.Loader](ctx).UserThumbsUpMotion.Load(ctx, uid)
	thumbsUp, err := userThumbsUpThunk()
	if err != nil {
		return nil, err
	}
	if !thumbsUp.ThumbsUp(motionID) {
		return nil, nil
	}
	db := dbutil.GetDB(ctx)
	err = dbquery.Use(db).Transaction(func(tx *dbquery.Query) error {
		UserThumbsUpMotion := tx.UserThumbsUpMotion
		rx, err := UserThumbsUpMotion.WithContext(ctx).Where(UserThumbsUpMotion.UserID.Eq(uid), UserThumbsUpMotion.ToMotionID.Eq(motionID)).Delete()
		if err != nil {
			return err
		}
		if rx.RowsAffected == 1 {
			Motion := tx.Motion
			if _, err := Motion.WithContext(ctx).Where(Motion.ID.Eq(motionID)).UpdateSimple(Motion.ThumbsUpCount.Add(-1)); err != nil {
				return err
			}
			thumbsUp.UnThumbsUp(motionID)
			midacontext.GetLoader[loader.Loader](ctx).Motion.Clear(ctx, motionID)
		}
		return nil
	})
	return nil, err
}

// ThumbsUpMotions is the resolver for the thumbsUpMotions field.
func (r *mutationResolver) ThumbsUpMotions(ctx context.Context, userID *string, paginator *graphqlutil.GraphQLPaginator) ([]*models.UserThumbsUpMotion, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}
	uid := graphqlutil.GetID(token, userID)
	if uid == "" {
		return nil, whalecode.ErrUserIDCannotBeEmpty
	}
	pager := graphqlutil.GetPager(paginator)
	db := dbutil.GetDB(ctx)
	UserThumbsUpMotion := dbquery.Use(db).UserThumbsUpMotion
	thumbsUps, err := UserThumbsUpMotion.WithContext(ctx).Where(UserThumbsUpMotion.UserID.Eq(uid)).Offset(pager.Offset()).Limit(pager.Limit()).Find()
	if err != nil {
		return nil, err
	}
	return thumbsUps, nil
}

// ThumbsUpMotionsCount is the resolver for the thumbsUpMotionsCount field.
func (r *mutationResolver) ThumbsUpMotionsCount(ctx context.Context, userID *string) (*models.Summary, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}
	uid := graphqlutil.GetID(token, userID)
	if uid == "" {
		return nil, whalecode.ErrUserIDCannotBeEmpty
	}
	thunk := midacontext.GetLoader[loader.Loader](ctx).UserThumbsUpMotion.Load(ctx, uid)
	userThumbsUp, err := thunk()
	if err != nil {
		return nil, err
	}
	return &models.Summary{Count: userThumbsUp.Size()}, nil
}

// LikedMotions is the resolver for the likedMotions field.
func (r *queryResolver) LikedMotions(ctx context.Context, userID *string, paginator *graphqlutil.GraphQLPaginator) ([]*models.UserLikeMotion, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() && !token.IsUser() {
		return nil, midacode.ErrNotPermitted
	}
	uid := graphqlutil.GetID(token, userID)
	if uid == "" {
		return nil, whalecode.ErrUserIDCannotBeEmpty
	}
	db := dbutil.GetDB(ctx)
	UserLikeMotion := dbquery.Use(db).UserLikeMotion
	pager := graphqlutil.GetPager(paginator)
	if pager.IfExcceedLimit(token) {
		return []*models.UserLikeMotion{}, nil
	}
	likes, err := UserLikeMotion.WithContext(ctx).Where(UserLikeMotion.UserID.Eq(uid)).Offset(pager.Offset()).Limit(pager.Limit()).Order(UserLikeMotion.ID.Desc()).Find()
	return likes, err
}

// LikedMotionsCount is the resolver for the likedMotionsCount field.
func (r *queryResolver) LikedMotionsCount(ctx context.Context, userID *string) (*models.Summary, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() && !token.IsUser() {
		return nil, midacode.ErrNotPermitted
	}
	uid := graphqlutil.GetID(token, userID)
	if uid == "" {
		return nil, whalecode.ErrUserIDCannotBeEmpty
	}
	thunk := midacontext.GetLoader[loader.Loader](ctx).UserLikeMotion.Load(ctx, uid)
	likes, err := thunk()
	if err != nil {
		return nil, err
	}
	return &models.Summary{Count: len(likes.MotionIDs)}, nil
}

// Motion is the resolver for the motion field.
func (r *userLikeMotionResolver) Motion(ctx context.Context, obj *models.UserLikeMotion) (*models.Motion, error) {
	thunk := midacontext.GetLoader[loader.Loader](ctx).Motion.Load(ctx, obj.ToMotionID)
	motion, err := thunk()
	if err != nil {
		return nil, err
	}
	return motion, nil
}

// UserLikeMotion returns UserLikeMotionResolver implementation.
func (r *Resolver) UserLikeMotion() UserLikeMotionResolver { return &userLikeMotionResolver{r} }

type userLikeMotionResolver struct{ *Resolver }

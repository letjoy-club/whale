package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.28

import (
	"context"
	"time"
	"whale/pkg/dbquery"
	"whale/pkg/keyer"
	"whale/pkg/loader"
	"whale/pkg/matcher"
	"whale/pkg/models"
	"whale/pkg/modelutil"
	"whale/pkg/whalecode"

	"github.com/letjoy-club/mida-tool/dbutil"
	"github.com/letjoy-club/mida-tool/graphqlutil"
	"github.com/letjoy-club/mida-tool/midacode"
	"github.com/letjoy-club/mida-tool/midacontext"
	"github.com/samber/lo"
	"go.uber.org/multierr"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

// CreateMatching is the resolver for the createMatching field.
func (r *mutationResolver) CreateMatching(ctx context.Context, userID *string, param models.CreateMatchingParam) (*models.Matching, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() && !token.IsUser() {
		return nil, midacode.ErrNotPermitted
	}
	uid := graphqlutil.GetID(token, userID)
	if uid == "" {
		return nil, whalecode.ErrUserIDCannotBeEmpty
	}
	if param.Remark == nil {
		emptyStr := ""
		param.Remark = &emptyStr
	}
	release, err := keyer.LockAll(ctx, keyer.UserMatching(uid))
	if err != nil {
		return nil, err
	}
	defer release(ctx)
	return modelutil.CreateMatching(ctx, uid, param)
}

// UpdateMatching is the resolver for the updateMatching field.
func (r *mutationResolver) UpdateMatching(ctx context.Context, matchingID string, param models.UpdateMatchingParam) (*models.Matching, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}
	db := dbutil.GetDB(ctx)
	Matching := dbquery.Use(db).Matching
	fields := []field.AssignExpr{}
	if param.AreaIds != nil {
		fields = append(fields, Matching.AreaIDs.Value(graphqlutil.ElementList[string](param.AreaIds)))
	}
	if param.TopicID != nil {
		fields = append(fields, Matching.TopicID.Value(*param.TopicID))
	}
	if param.Gender != nil {
		fields = append(fields, Matching.Gender.Value(param.Gender.String()))
	}
	if param.Deadline != nil {
		fields = append(fields, Matching.Deadline.Value(*param.Deadline))
	}
	if param.Remark != nil {
		fields = append(fields, Matching.Remark.Value(*param.Remark))
	}
	if param.CityID != nil {
		fields = append(fields, Matching.CityID.Value(*param.CityID))
	}
	if param.CreatedAt != nil {
		fields = append(fields, Matching.CreatedAt.Value(*param.CreatedAt))
	}
	_, err := Matching.WithContext(ctx).Where(Matching.ID.Eq(matchingID)).UpdateSimple(fields...)
	if err != nil {
		return nil, err
	}
	midacontext.GetLoader[loader.Loader](ctx).Matching.Clear(ctx, matchingID)
	thunk := midacontext.GetLoader[loader.Loader](ctx).Matching.Load(ctx, matchingID)
	return thunk()
}

// UpdateMatchingQuota is the resolver for the updateMatchingQuota field.
func (r *mutationResolver) UpdateMatchingQuota(ctx context.Context, userID string, param models.UpdateMatchingQuotaParam) (string, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() {
		return "", midacode.ErrNotPermitted
	}
	db := dbutil.GetDB(ctx)
	MatchingQuota := dbquery.Use(db).MatchingQuota

	fields := []field.AssignExpr{}
	if param.Remain != nil {
		fields = append(fields, MatchingQuota.Remain.Value(*param.Remain))
	}
	if param.Total != nil {
		fields = append(fields, MatchingQuota.Total.Value(*param.Total))
	}

	_, err := MatchingQuota.WithContext(ctx).Where(MatchingQuota.UserID.Eq(userID)).UpdateSimple(fields...)
	if err != nil {
		return "", err
	}

	midacontext.GetLoader[loader.Loader](ctx).MatchingQuota.Clear(ctx, userID)
	return "", nil
}

// ConfirmMatchingResult is the resolver for the confirmMatchingResult field.
func (r *mutationResolver) ConfirmMatchingResult(ctx context.Context, userID *string, matchingID string, reject bool) (*string, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() && !token.IsUser() {
		return nil, midacode.ErrNotPermitted
	}

	matchingThunk := midacontext.GetLoader[loader.Loader](ctx).Matching.Load(ctx, matchingID)
	matching, err := matchingThunk()
	if err != nil {
		return nil, err
	}
	if matching == nil {
		return nil, midacode.ErrItemNotFound
	}

	if token.IsUser() {
		if matching.UserID != token.String() {
			return nil, midacode.ErrNotPermitted
		}
	}

	if matching.State != models.MatchingStateMatched.String() {
		return nil, whalecode.ErrMatchingStateShouldBeMatched
	}

	loader := midacontext.GetLoader[loader.Loader](ctx)

	thunk := loader.MatchingResult.Load(ctx, matching.ResultID)
	matchingResult, err := thunk()
	if err != nil {
		return nil, err
	}

	db := dbutil.GetDB(ctx)
	err = db.Transaction(func(tx *gorm.DB) error {
		if reject {
			Matching := dbquery.Use(tx).Matching
			matching.RejectedUserIDs = append(matching.RejectedUserIDs, matchingResult.OtherUserIDs(matching.UserID)...)
			_, err := Matching.WithContext(ctx).Select(Matching.RejectedUserIDs).Where(Matching.ID.Eq(matchingID)).Updates(matching)
			if err != nil {
				loader.Matching.Clear(ctx, matchingID)
				return err
			}
		}
		return modelutil.ConfirmMatching(ctx, tx, matchingResult, matchingID, token.String(), !reject)
	})
	if err == nil {
		err = modelutil.CheckMatchingResultAndCreateChatGroup(ctx, matchingResult)
	}
	for _, m := range matchingResult.MatchingIDs {
		loader.Matching.Clear(ctx, m)
	}
	loader.MatchingResult.Clear(ctx, matchingResult.ID)
	return nil, err
}

// ConfirmMatchingResultV2 is the resolver for the confirmMatchingResultV2 field.
func (r *mutationResolver) ConfirmMatchingResultV2(ctx context.Context, userID *string, matchingID string, confirm bool) (*string, error) {
	reject := !confirm
	return r.ConfirmMatchingResult(ctx, userID, matchingID, reject)
}

// CancelMatching is the resolver for the cancelMatching field.
func (r *mutationResolver) CancelMatching(ctx context.Context, matchingID string) (*string, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() && !token.IsUser() {
		return nil, midacode.ErrNotPermitted
	}

	matching, err := modelutil.GetMatchingAndCheckUser(ctx, matchingID, token.UserID())
	if err != nil {
		return nil, err
	}

	if matching.State != models.MatchingStateMatching.String() {
		return nil, whalecode.ErrMatchingStateShouldBeMatching
	}

	db := dbutil.GetDB(ctx)
	err = db.Transaction(func(tx *gorm.DB) error {
		Matching := dbquery.Use(db).Matching
		MatchingQuota := dbquery.Use(db).MatchingQuota
		_, err := Matching.WithContext(ctx).
			Where(Matching.ID.Eq(matchingID)).
			UpdateSimple(
				Matching.State.Value(models.MatchingStateCanceled.String()),
				Matching.MatchedAt.Null(),
			)
		if err != nil {
			return err
		}
		_, err = MatchingQuota.WithContext(ctx).
			Where(MatchingQuota.UserID.Eq(token.UserID())).
			UpdateSimple(MatchingQuota.Remain.Add(1))
		return err
	})
	loader := midacontext.GetLoader[loader.Loader](ctx)
	loader.Matching.Clear(ctx, matchingID)
	loader.MatchingQuota.Clear(ctx, matching.UserID)
	return nil, err
}

// StartMatching is the resolver for the startMatching field.
func (r *mutationResolver) StartMatching(ctx context.Context) (*string, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}

	matcher := matcher.Matcher{}
	err := matcher.Match(ctx)
	return nil, err
}

// FinishMatching is the resolver for the finishMatching field.
func (r *mutationResolver) FinishMatching(ctx context.Context, matchingID string) (*string, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() && !token.IsUser() {
		return nil, midacode.ErrNotPermitted
	}

	uid := ""
	if token.IsUser() {
		uid = token.String()
	}
	release, err := keyer.LockAll(ctx, keyer.UserMatching(uid))
	if err != nil {
		return nil, err
	}
	defer release(ctx)
	return nil, modelutil.FinishMatching(ctx, matchingID, uid)
}

// ReviewMatching is the resolver for the reviewMatching field.
func (r *mutationResolver) ReviewMatching(ctx context.Context, matchingID string, param models.ReviewMatchingParam) (*string, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() && !token.IsUser() {
		return nil, midacode.ErrNotPermitted
	}
	matchingThunk := midacontext.GetLoader[loader.Loader](ctx).Matching.Load(ctx, matchingID)
	matching, err := matchingThunk()
	if err != nil {
		return nil, err
	}
	if token.IsUser() {
		if matching.UserID != token.String() {
			return nil, midacode.ErrNotPermitted
		}
	}
	matchingResultThunk := midacontext.GetLoader[loader.Loader](ctx).MatchingResult.Load(ctx, matching.ResultID)
	result, err := matchingResultThunk()
	if err != nil {
		return nil, err
	}

	db := dbutil.GetDB(ctx)
	MatchingReview := dbquery.Use(db).MatchingReview
	_, err = MatchingReview.WithContext(ctx).Where(MatchingReview.MatchingID.Eq(matchingID)).Take()
	if err == nil {
		return nil, midacode.ErrRecordExists
	}
	peerMatchingID := ""
	for _, id := range result.MatchingIDs {
		if id != matchingID {
			peerMatchingID = id
			break
		}
	}
	err = MatchingReview.WithContext(ctx).Create(&models.MatchingReview{
		MatchingResultID: matching.ResultID,
		UserID:           matching.UserID,
		ToUserID:         param.ToUserID,
		TopicID:          matching.TopicID,
		Score:            param.Score,
		MatchingID:       matchingID,
		ToMatchingID:     peerMatchingID,
		Comment:          param.Comment,
	})
	return nil, err
}

// Matching is the resolver for the matching field.
func (r *queryResolver) Matching(ctx context.Context, id string) (*models.Matching, error) {
	token := midacontext.GetClientToken(ctx)
	if token.IsAnonymous() {
		return nil, midacode.ErrNotPermitted
	}
	thunk := midacontext.GetLoader[loader.Loader](ctx).Matching.Load(ctx, id)
	return thunk()
}

// UserMatchingQuota is the resolver for the userMatchingQuota field.
func (r *queryResolver) UserMatchingQuota(ctx context.Context, userID string) (*models.MatchingQuota, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}

	thunk := midacontext.GetLoader[loader.Loader](ctx).MatchingQuota.Load(ctx, userID)
	return thunk()
}

// UserMatchingCalendar is the resolver for the userMatchingCalendar field.
func (r *queryResolver) UserMatchingCalendar(ctx context.Context, userID *string, param models.UserMatchingCalenderParam) ([]*models.CalendarEvent, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsUser() && !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}
	uid := graphqlutil.GetID(token, userID)
	if uid == "" {
		return nil, whalecode.ErrUserIDCannotBeEmpty
	}
	if param.Before.Sub(param.After) > 64*24*time.Hour {
		return nil, whalecode.ErrQueryDurationTooLong
	}
	db := dbutil.GetDB(ctx).Debug()
	MatchingResult := dbquery.Use(db).MatchingResult
	matchingResultIDs := []int{}
	query := MatchingResult.WithContext(ctx).
		Where(
			dbutil.RawCond(dbutil.Contains(MatchingResult.UserIDs.ColumnName().String(), uid)),
			MatchingResult.Closed.Is(false),
			MatchingResult.ChatGroupCreatedAt.Between(param.After, param.Before.Add(-time.Second)),
		)
	if param.OtherUserID != nil {
		query = query.Where(dbutil.RawCond(dbutil.Contains(MatchingResult.UserIDs.ColumnName().String(), *param.OtherUserID)))
	}

	err := query.Limit(100).Pluck(MatchingResult.ID, &matchingResultIDs)
	if err != nil {
		return nil, err
	}
	thunk := midacontext.GetLoader[loader.Loader](ctx).MatchingResult.LoadMany(ctx, matchingResultIDs)
	matchingResults, errors := thunk()
	if len(errors) > 0 {
		return nil, multierr.Combine(errors...)
	}
	return lo.Map(matchingResults, func(m *models.MatchingResult, i int) *models.CalendarEvent {
		now := time.Now()
		if m.FinishedAt == nil {
			m.FinishedAt = &now
		}
		return &models.CalendarEvent{
			TopicID:            m.TopicID,
			MatchedAt:          m.CreatedAt,
			FinishedAt:         *m.FinishedAt,
			ChatGroupCreatedAt: m.ChatGroupCreatedAt,
		}
	}), nil
}

// UserMatchingsInTheDay is the resolver for the userMatchingsInTheDay field.
func (r *queryResolver) UserMatchingsInTheDay(ctx context.Context, userID *string, param models.UserMatchingInTheDayParam) ([]*models.MatchingResult, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsUser() && !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}
	uid := graphqlutil.GetID(token, userID)
	if uid == "" {
		return nil, whalecode.ErrUserIDCannotBeEmpty
	}
	t, err := time.Parse("20060102", param.DayStr)
	if err != nil {
		return nil, err
	}

	db := dbutil.GetDB(ctx)
	MatchingResult := dbquery.Use(db).MatchingResult
	matchingResultIDs := []int{}
	query := MatchingResult.WithContext(ctx).
		Where(
			dbutil.RawCond(dbutil.Contains(MatchingResult.UserIDs.ColumnName().String(), uid)),
			MatchingResult.Closed.Is(false),
			MatchingResult.ChatGroupCreatedAt.Between(t, t.Add(time.Hour*24-time.Second)),
		)
	if param.OtherUserID != nil {
		query = query.Where(dbutil.RawCond(dbutil.Contains(MatchingResult.UserIDs.ColumnName().String(), *param.OtherUserID)))
	}

	err = query.Limit(100).Pluck(MatchingResult.ID, &matchingResultIDs)
	if err != nil {
		return nil, err
	}
	thunk := midacontext.GetLoader[loader.Loader](ctx).MatchingResult.LoadMany(ctx, matchingResultIDs)

	matchingResults, errors := thunk()
	if len(errors) > 0 {
		return nil, multierr.Combine(errors...)
	}
	return matchingResults, nil
}

// MatchingResultByChatGroupID is the resolver for the matchingResultByChatGroupId field.
func (r *queryResolver) MatchingResultByChatGroupID(ctx context.Context, userID *string, chatGroupID string) (*models.MatchingResult, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsUser() && !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}
	uid := graphqlutil.GetID(token, userID)
	if uid == "" {
		return nil, whalecode.ErrUserIDCannotBeEmpty
	}

	db := dbutil.GetDB(ctx)
	MatchingResult := dbquery.Use(db).MatchingResult
	result, err := MatchingResult.WithContext(ctx).Where(MatchingResult.ChatGroupID.Eq(chatGroupID)).Take()
	if err != nil {
		return nil, midacode.ItemMayNotFound(err)
	}
	index := lo.IndexOf(result.UserIDs, uid)
	if index == -1 {
		return nil, midacode.ErrNotPermitted
	}
	return result, nil
}

// Matchings is the resolver for the matchings field.
func (r *queryResolver) Matchings(ctx context.Context, filter *models.MatchingFilter, paginator *graphqlutil.GraphQLPaginator) ([]*models.Matching, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}

	pager := graphqlutil.GetPager(paginator)

	db := dbutil.GetDB(ctx)
	Matching := dbquery.Use(db).Matching

	query := Matching.WithContext(ctx)
	if filter != nil {
		if filter.After != nil {
			query = query.Where(Matching.CreatedAt.Gt(*filter.After))
		}
		if filter.Before != nil {
			query = query.Where(Matching.CreatedAt.Lt(*filter.Before))
		}
		if filter.TopicID != nil {
			query = query.Where(Matching.TopicID.Eq(*filter.TopicID))
		}
		if filter.State != nil {
			query = query.Where(Matching.State.Eq(filter.State.String()))
		}
		if filter.UserID != nil {
			query = query.Where(Matching.UserID.Eq(*filter.UserID))
		}
		if filter.CityID != nil {
			query = query.Where(Matching.CityID.Eq(*filter.CityID))
		}
	}
	ids := []string{}
	err := query.Limit(pager.Limit()).Offset(pager.Offset()).Order(Matching.CreatedAt.Desc()).Pluck(Matching.ID, &ids)
	if err != nil {
		return nil, err
	}
	thunk := midacontext.GetLoader[loader.Loader](ctx).Matching.LoadMany(ctx, ids)
	matchings, errors := thunk()
	if len(errors) > 0 {
		return nil, multierr.Combine(errors...)
	}
	return matchings, nil
}

// MatchingsCount is the resolver for the matchingsCount field.
func (r *queryResolver) MatchingsCount(ctx context.Context, filter *models.MatchingFilter) (*models.Summary, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}

	db := dbutil.GetDB(ctx)
	Matching := dbquery.Use(db).Matching

	query := Matching.WithContext(ctx)
	if filter != nil {
		if filter.After != nil {
			query = query.Where(Matching.CreatedAt.Gt(*filter.After))
		}
		if filter.Before != nil {
			query = query.Where(Matching.CreatedAt.Lt(*filter.Before))
		}
		if filter.TopicID != nil {
			query = query.Where(Matching.TopicID.Eq(*filter.TopicID))
		}
		if filter.State != nil {
			query = query.Where(Matching.State.Eq(filter.State.String()))
		}
		if filter.UserID != nil {
			query = query.Where(Matching.UserID.Eq(*filter.UserID))
		}
		if filter.CityID != nil {
			query = query.Where(Matching.CityID.Eq(*filter.CityID))
		}
	}
	count, err := query.Count()
	return &models.Summary{Count: int(count)}, err
}

// MatchingResult is the resolver for the matchingResult field.
func (r *queryResolver) MatchingResult(ctx context.Context, id int) (*models.MatchingResult, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}

	thunk := midacontext.GetLoader[loader.Loader](ctx).MatchingResult.Load(ctx, id)
	return thunk()
}

// MatchingResults is the resolver for the matchingResults field.
func (r *queryResolver) MatchingResults(ctx context.Context, filter *models.MatchingResultFilter, paginator *graphqlutil.GraphQLPaginator) ([]*models.MatchingResult, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}

	pager := graphqlutil.GetPager(paginator)

	db := dbutil.GetDB(ctx)
	MatchingResult := dbquery.Use(db).MatchingResult

	query := MatchingResult.WithContext(ctx)
	if filter != nil {
		if filter.After != nil {
			query = query.Where(MatchingResult.CreatedAt.Gt(*filter.After))
		}
		if filter.Before != nil {
			query = query.Where(MatchingResult.CreatedAt.Lt(*filter.Before))
		}
		if filter.UserID != nil {
			query = query.Where(dbutil.RawCond(dbutil.Contains(MatchingResult.UserIDs.ColumnName().String(), *filter.UserID)))
		}
	}
	return query.Limit(pager.Limit()).Offset(pager.Offset()).Order(MatchingResult.CreatedAt.Desc()).Find()
}

// MatchingResultsCount is the resolver for the matchingResultsCount field.
func (r *queryResolver) MatchingResultsCount(ctx context.Context, filter *models.MatchingResultFilter) (*models.Summary, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}

	db := dbutil.GetDB(ctx)
	MatchingResult := dbquery.Use(db).MatchingResult

	query := MatchingResult.WithContext(ctx)
	if filter != nil {
		if filter.After != nil {
			query = query.Where(MatchingResult.CreatedAt.Gt(*filter.After))
		}
		if filter.Before != nil {
			query = query.Where(MatchingResult.CreatedAt.Lt(*filter.Before))
		}
		if filter.UserID != nil {
			query = query.Where(dbutil.RawCond(dbutil.Contains(MatchingResult.UserIDs.ColumnName().String(), *filter.UserID)))
		}
	}

	count, err := query.Count()
	return &models.Summary{Count: int(count)}, err
}

// UserMatchings is the resolver for the userMatchings field.
func (r *queryResolver) UserMatchings(ctx context.Context, userID *string, filter *models.UserMatchingFilter, paginator *graphqlutil.GraphQLPaginator) ([]*models.Matching, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsUser() && !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}

	pager := graphqlutil.GetPager(paginator)

	uid := graphqlutil.GetID(token, userID)
	if uid == "" {
		return nil, whalecode.ErrUserIDCannotBeEmpty
	}

	db := dbutil.GetDB(ctx)
	Matching := dbquery.Use(db).Matching

	query := Matching.WithContext(ctx).Where(Matching.UserID.Eq(uid))
	if filter != nil {
		if filter.State != nil {
			query = query.Where(Matching.State.Eq(filter.State.String()))
		}
		if len(filter.States) > 0 {
			states := lo.Map(filter.States, func(s models.MatchingState, i int) string {
				return s.String()
			})
			query = query.Where(Matching.State.In(states...))
		}
	}
	ids := []string{}
	err := query.Limit(pager.Limit()).Offset(pager.Offset()).Order(Matching.CreatedAt.Desc()).Pluck(Matching.ID, &ids)
	if err != nil {
		return nil, err
	}
	thunk := midacontext.GetLoader[loader.Loader](ctx).Matching.LoadMany(ctx, ids)
	matchings, errors := thunk()
	if len(errors) > 0 {
		return nil, multierr.Combine(errors...)
	}
	return matchings, nil
}

// UnconfirmedUserMatchings is the resolver for the unconfirmedUserMatchings field.
func (r *queryResolver) UnconfirmedUserMatchings(ctx context.Context, userID *string) ([]*models.Matching, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsUser() && !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}

	uid := graphqlutil.GetID(token, userID)
	if uid == "" {
		return nil, whalecode.ErrUserIDCannotBeEmpty
	}

	db := dbutil.GetDB(ctx)
	MatchingResult := dbquery.Use(db).MatchingResult
	matchingResults, err := MatchingResult.WithContext(ctx).Where(
		// 未关闭
		MatchingResult.Closed.Is(false),
		// 未结束
		MatchingResult.FinishedAt.IsNull(),
		// 未创建聊天
		MatchingResult.ChatGroupState.Eq(models.ChatGroupStateCreated.String()),
		dbutil.RawCond(dbutil.Contains(MatchingResult.UserIDs.ColumnName().String(), uid)),
	).Select(MatchingResult.MatchingIDs, MatchingResult.UserIDs).Limit(10).Find()
	if err != nil {
		return nil, err
	}
	matchingIDs := lo.Map(matchingResults, func(result *models.MatchingResult, i int) string {
		return result.GetMatchingID(uid)
	})
	thunk := midacontext.GetLoader[loader.Loader](ctx).Matching.LoadMany(ctx, matchingIDs)
	matchings, errors := thunk()
	if errors != nil {
		return nil, multierr.Combine(errors...)
	}
	return matchings, nil
}

// UserMatchingsCount is the resolver for the userMatchingsCount field.
func (r *queryResolver) UserMatchingsCount(ctx context.Context, userID *string, filter *models.UserMatchingFilter) (*models.Summary, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsUser() && !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}

	uid := graphqlutil.GetID(token, userID)
	if uid == "" {
		return nil, whalecode.ErrUserIDCannotBeEmpty
	}

	db := dbutil.GetDB(ctx)
	Matching := dbquery.Use(db).Matching

	query := Matching.WithContext(ctx).Where(Matching.UserID.Eq(uid))
	if filter != nil {
		if filter.State != nil {
			query = query.Where(Matching.State.Eq(filter.State.String()))
		}
	}
	count, err := query.Count()
	return &models.Summary{Count: int(count)}, err
}

// PreviewMatchingsOfTopic is the resolver for the previewMatchingsOfTopic field.
func (r *queryResolver) PreviewMatchingsOfTopic(ctx context.Context, cityID string, topicID string, limit *int) ([]*models.Matching, error) {
	db := dbutil.GetDB(ctx)
	Matching := dbquery.Use(db).Matching

	n := 4
	if limit != nil {
		if *limit <= 0 {
			*limit = 4
		}
		if *limit > 8 {
			*limit = 8
		}
		n = *limit
	}

	ids := []string{}
	err := Matching.WithContext(ctx).Where(
		Matching.CityID.Eq(cityID),
		Matching.TopicID.Eq(topicID),
	).Limit(n).Order(Matching.CreatedAt.Desc()).Order(Matching.CreatedAt.Desc()).Pluck(Matching.ID, &ids)
	if err != nil {
		return nil, err
	}
	thunk := midacontext.GetLoader[loader.Loader](ctx).Matching.LoadMany(ctx, ids)
	matchings, errors := thunk()
	if len(errors) > 0 {
		return nil, multierr.Combine(errors...)
	}
	return matchings, nil
}

// UnconfirmedInvitations is the resolver for the unconfirmedInvitations field.
func (r *queryResolver) UnconfirmedInvitations(ctx context.Context, userID *string) ([]*models.MatchingInvitation, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsUser() && !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}
	uid := graphqlutil.GetID(token, userID)
	if uid == "" {
		return nil, whalecode.ErrUserIDCannotBeEmpty
	}
	db := dbutil.GetDB(ctx)
	MatchingInvitation := dbquery.Use(db).MatchingInvitation
	ids := []string{}
	err := MatchingInvitation.WithContext(ctx).
		Where(MatchingInvitation.InviteeID.Eq(uid)).
		Where(MatchingInvitation.Closed.Is(false)).
		Where(MatchingInvitation.ConfirmState.Eq(models.InvitationConfirmStateUnconfirmed.String())).
		Order(MatchingInvitation.CreatedAt.Desc()).
		Pluck(MatchingInvitation.ID, &ids)
	if err != nil {
		return nil, err
	}
	thunk := midacontext.GetLoader[loader.Loader](ctx).MatchingInvitation.LoadMany(ctx, ids)
	invitaions, errors := thunk()
	if len(errors) > 0 {
		return nil, multierr.Combine(errors...)
	}
	return invitaions, nil
}

// UnconfirmedInvitationCount is the resolver for the unconfirmedInvitationCount field.
func (r *queryResolver) UnconfirmedInvitationCount(ctx context.Context, userID *string) (*models.Summary, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsUser() && !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}
	uid := graphqlutil.GetID(token, userID)
	if uid == "" {
		return nil, whalecode.ErrUserIDCannotBeEmpty
	}
	db := dbutil.GetDB(ctx)
	MatchingInvitation := dbquery.Use(db).MatchingInvitation
	count, err := MatchingInvitation.WithContext(ctx).
		Where(MatchingInvitation.InviteeID.Eq(uid)).
		Where(MatchingInvitation.Closed.Is(false)).
		Where(MatchingInvitation.ConfirmState.Eq(models.InvitationConfirmStateUnconfirmed.String())).
		Count()
	if err != nil {
		return nil, err
	}
	return &models.Summary{Count: int(count)}, nil
}

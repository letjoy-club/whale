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
	"whale/pkg/whalecode"

	"github.com/letjoy-club/mida-tool/dbutil"
	"github.com/letjoy-club/mida-tool/graphqlutil"
	"github.com/letjoy-club/mida-tool/midacode"
	"github.com/letjoy-club/mida-tool/midacontext"
	"github.com/samber/lo"
	"go.uber.org/multierr"
	"gorm.io/gorm/clause"
)

// CreateCityTopics is the resolver for the createCityTopics field.
func (r *mutationResolver) CreateCityTopics(ctx context.Context, param models.CreateCityTopicParam) (*models.CityTopics, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}

	db := dbutil.GetDB(ctx)
	CityTopics := dbquery.Use(db).CityTopics

	cityTopic := &models.CityTopics{
		CityID:   param.CityID,
		TopicIDs: param.TopicIds,
	}

	err := CityTopics.WithContext(ctx).Create(cityTopic)
	if err != nil {
		return nil, err
	}
	return cityTopic, nil
}

// UpdateCityTopics is the resolver for the updateCityTopics field.
func (r *mutationResolver) UpdateCityTopics(ctx context.Context, cityID string, param models.UpdateCityTopicParam) (*models.CityTopics, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}

	db := dbutil.GetDB(ctx)
	CityTopics := dbquery.Use(db).CityTopics

	_, err := CityTopics.WithContext(ctx).Where(CityTopics.CityID.Eq(cityID)).UpdateSimple(
		CityTopics.TopicIDs.Value(graphqlutil.ElementList[string](param.TopicIds)),
	)
	if err != nil {
		return nil, err
	}
	midacontext.GetLoader[loader.Loader](ctx).CityTopics.Clear(ctx, cityID)
	return CityTopics.WithContext(ctx).Where(CityTopics.CityID.Eq(cityID)).Take()
}

// UpdateHotTopicsInArea is the resolver for the updateHotTopicsInArea field.
func (r *mutationResolver) UpdateHotTopicsInArea(ctx context.Context, cityID string, param models.UpdateHotTopicParam) (*models.HotTopicsInArea, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}

	thunk := midacontext.GetLoader[loader.Loader](ctx).HotTopics.Load(ctx, cityID)
	hotTopics, err := thunk()
	if err != nil {
		return nil, err
	}

	db := dbutil.GetDB(ctx)
	HotTopicsInArea := dbquery.Use(db).HotTopicsInArea

	metrics := lo.Map(param.TopicMetrics, func(m *models.UpdateHotTopicMetricsParam, i int) models.TopicMetrics {
		return models.TopicMetrics{ID: m.TopicID, Matching: m.Matching, Matched: m.Matched, Heat: m.Heat}
	})

	hotTopics.TopicMetrics = metrics

	if _, err := HotTopicsInArea.WithContext(ctx).
		Where(HotTopicsInArea.CityID.Eq(cityID)).
		Select(HotTopicsInArea.TopicMetrics).
		Updates(hotTopics); err != nil {
		return nil, err
	}

	midacontext.GetLoader[loader.Loader](ctx).HotTopics.Clear(ctx, cityID)
	return hotTopics, nil
}

// UpdateUserJoinTopic is the resolver for the updateUserJoinTopic field.
func (r *mutationResolver) UpdateUserJoinTopic(ctx context.Context, id int, param models.UpdateUserJoinTopicParam) (*models.UserJoinTopic, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}
	thunk := midacontext.GetLoader[loader.Loader](ctx).UserJoinTopic.Load(ctx, id)
	userJoinTopic, err := thunk()
	if err != nil {
		return nil, err
	}
	matchingThunk := midacontext.GetLoader[loader.Loader](ctx).Matching.Load(ctx, userJoinTopic.LatestMatchingID)
	matching, err := matchingThunk()
	if err != nil {
		return nil, err
	}
	if userJoinTopic.CityID != matching.CityID || userJoinTopic.TopicID != matching.TopicID {
		return nil, whalecode.ErrMatchingNotMatchWithTopic
	}

	db := dbutil.GetDB(ctx)
	UserJoinTopic := dbquery.Use(db).UserJoinTopic
	_, err = UserJoinTopic.WithContext(ctx).Where(UserJoinTopic.ID.Eq(id)).UpdateSimple(UserJoinTopic.LatestMatchingID.Value(param.MatchingID))
	if err != nil {
		return nil, err
	}
	midacontext.GetLoader[loader.Loader](ctx).UserJoinTopic.Clear(ctx, id)
	return midacontext.GetLoader[loader.Loader](ctx).UserJoinTopic.Load(ctx, id)()
}

// CreateUserJoinTopic is the resolver for the createUserJoinTopic field.
func (r *mutationResolver) CreateUserJoinTopic(ctx context.Context, param models.CreateUserJoinTopicParam) (*models.UserJoinTopic, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}
	thunk := midacontext.GetLoader[loader.Loader](ctx).Matching.Load(ctx, param.MatchingID)
	matching, err := thunk()
	if err != nil {
		return nil, err
	}
	userJoinTopic := &models.UserJoinTopic{
		TopicID:          matching.TopicID,
		CityID:           matching.CityID,
		UserID:           matching.UserID,
		LatestMatchingID: matching.ID,
		Times:            1,
	}
	db := dbutil.GetDB(ctx)
	UserJoinTopic := dbquery.Use(db).UserJoinTopic
	if err := UserJoinTopic.WithContext(ctx).Clauses(
		clause.OnConflict{DoUpdates: clause.AssignmentColumns([]string{
			UserJoinTopic.LatestMatchingID.ColumnName().String(),
		})}).
		Create(userJoinTopic); err != nil {
		return nil, err
	}
	midacontext.GetLoader[loader.Loader](ctx).UserJoinTopic.Clear(ctx, userJoinTopic.ID)
	return userJoinTopic, nil
}

// CityTopics is the resolver for the cityTopics field.
func (r *queryResolver) CityTopics(ctx context.Context, cityID string) (*models.CityTopics, error) {
	thunk := midacontext.GetLoader[loader.Loader](ctx).CityTopics.Load(ctx, cityID)
	cityTopics, err := thunk()
	if err != nil {
		cityID = "310100"
		thunk := midacontext.GetLoader[loader.Loader](ctx).CityTopics.Load(ctx, cityID)
		return thunk()
	}
	return cityTopics, nil
}

// CitiesTopics is the resolver for the citiesTopics field.
func (r *queryResolver) CitiesTopics(ctx context.Context, filter *models.CitiesTopicsFilter, paginator *graphqlutil.GraphQLPaginator) ([]*models.CityTopics, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}
	db := dbutil.GetDB(ctx)
	CityTopics := dbquery.Use(db).CityTopics
	pager := graphqlutil.GetPager(paginator)

	cityIds := []string{}
	query := CityTopics.WithContext(ctx)
	if filter != nil {
		if filter.CityID != nil {
			query = query.Where(CityTopics.CityID.Eq(*filter.CityID))
		}
	}
	err := query.Limit(pager.Limit()).Offset(pager.Offset()).Pluck(CityTopics.CityID, &cityIds)
	if err != nil {
		return nil, err
	}
	thunk := midacontext.GetLoader[loader.Loader](ctx).CityTopics.LoadMany(ctx, cityIds)
	cityTopics, errors := thunk()
	if len(errors) > 0 {
		return nil, multierr.Combine(errors...)
	}
	return cityTopics, nil
}

// CitiesTopicsCount is the resolver for the citiesTopicsCount field.
func (r *queryResolver) CitiesTopicsCount(ctx context.Context, filter *models.CitiesTopicsFilter) (*models.Summary, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}
	db := dbutil.GetDB(ctx)
	CityTopics := dbquery.Use(db).CityTopics
	query := CityTopics.WithContext(ctx)
	if filter != nil {
		if filter.CityID != nil {
			query = query.Where(CityTopics.CityID.Eq(*filter.CityID))
		}
	}
	count, err := query.Count()
	return &models.Summary{Count: int(count)}, err
}

// HotTopics is the resolver for the hotTopics field.
func (r *queryResolver) HotTopicsInArea(ctx context.Context, cityID *string) (*models.HotTopicsInArea, error) {
	token := midacontext.GetClientToken(ctx)
	if token.IsUser() || token.IsAnonymous() {
		// 如果是用户，就获取v2版本的数据
		hotTopicsV2 := midacontext.GetLoader[loader.Loader](ctx).HotTopicsV2
		hotTopicsV2.Load(ctx)

		return &models.HotTopicsInArea{
			CityID:       "310100",
			TopicMetrics: hotTopicsV2.Metrics(),
			UpdatedAt:    time.Now(),
			CreatedAt:    time.Now(),
		}, nil
	}

	cid := "310100"
	if cityID != nil {
		cid = *cityID
	}
	thunk := midacontext.GetLoader[loader.Loader](ctx).HotTopics.Load(ctx, cid)
	hotTopics, err := thunk()
	if err != nil || len(hotTopics.TopicMetrics) == 0 {
		// 如果没有结果，就获取上海数据
		thunk := midacontext.GetLoader[loader.Loader](ctx).HotTopics.Load(ctx, "310100")
		return thunk()
	}
	return hotTopics, err
}

// HotTopics is the resolver for the hotTopics field.
func (r *queryResolver) HotTopics(ctx context.Context, filter *models.HotTopicsFilter, paginator *graphqlutil.GraphQLPaginator) ([]*models.HotTopicsInArea, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}
	pager := graphqlutil.GetPager(paginator)
	db := dbutil.GetDB(ctx)
	HotTopicsInArea := dbquery.Use(db).HotTopicsInArea
	query := HotTopicsInArea.WithContext(ctx)
	if filter != nil {
		if filter.CityID != nil {
			query = query.Where(HotTopicsInArea.CityID.Eq(*filter.CityID))
		}
	}

	cityIDs := []string{}
	err := query.Limit(pager.Limit()).Offset(pager.Offset()).Pluck(HotTopicsInArea.CityID, &cityIDs)
	if err != nil {
		return nil, err
	}

	thunk := midacontext.GetLoader[loader.Loader](ctx).HotTopics.LoadMany(ctx, cityIDs)
	hotTopics, errors := thunk()
	if len(errors) > 0 {
		return nil, multierr.Combine(errors...)
	}
	return hotTopics, nil
}

// HotTopicsCount is the resolver for the hotTopicsCount field.
func (r *queryResolver) HotTopicsCount(ctx context.Context, filter *models.HotTopicsFilter) (*models.Summary, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}
	db := dbutil.GetDB(ctx)
	HotTopicsInArea := dbquery.Use(db).HotTopicsInArea
	query := HotTopicsInArea.WithContext(ctx)
	if filter != nil {
		if filter.CityID != nil {
			query = query.Where(HotTopicsInArea.CityID.Eq(*filter.CityID))
		}
	}

	count, err := query.Count()
	return &models.Summary{Count: int(count)}, err
}

// UserJoinTopics is the resolver for the userJoinTopics field.
func (r *queryResolver) UserJoinTopics(ctx context.Context, filter *models.UserJoinTopicFilter, paginator *graphqlutil.GraphQLPaginator) ([]*models.UserJoinTopic, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}
	pager := graphqlutil.GetPager(paginator)
	db := dbutil.GetDB(ctx)
	UserJoinTopic := dbquery.Use(db).UserJoinTopic
	query := UserJoinTopic.WithContext(ctx)
	if filter != nil {
		if filter.UserID != nil {
			query = query.Where(UserJoinTopic.UserID.Eq(*filter.UserID))
		}
		if filter.CityID != nil {
			query = query.Where(UserJoinTopic.CityID.Eq(*filter.CityID))
		}
		if filter.TopicID != nil {
			query = query.Where(UserJoinTopic.TopicID.Eq(*filter.TopicID))
		}
	}
	userJoinTopicIDs := []int{}
	err := query.Limit(pager.Limit()).Offset(pager.Offset()).Order(UserJoinTopic.ID.Desc()).Pluck(UserJoinTopic.ID, &userJoinTopicIDs)
	if err != nil {
		return nil, err
	}
	thunk := midacontext.GetLoader[loader.Loader](ctx).UserJoinTopic.LoadMany(ctx, userJoinTopicIDs)
	userJoinTopics, errors := thunk()
	if len(errors) > 0 {
		return nil, multierr.Combine(errors...)
	}
	return userJoinTopics, nil
}

// UserJoinTopicsCount is the resolver for the userJoinTopicsCount field.
func (r *queryResolver) UserJoinTopicsCount(ctx context.Context, filter *models.UserJoinTopicFilter) (*models.Summary, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}
	db := dbutil.GetDB(ctx)
	UserJoinTopic := dbquery.Use(db).UserJoinTopic
	query := UserJoinTopic.WithContext(ctx)
	if filter != nil {
		if filter.UserID != nil {
			query = query.Where(UserJoinTopic.UserID.Eq(*filter.UserID))
		}
		if filter.CityID != nil {
			query = query.Where(UserJoinTopic.CityID.Eq(*filter.CityID))
		}
		if filter.TopicID != nil {
			query = query.Where(UserJoinTopic.TopicID.Eq(*filter.TopicID))
		}
	}
	count, err := query.Count()
	return &models.Summary{Count: int(count)}, err
}

// UserJoinTopic is the resolver for the userJoinTopic field.
func (r *queryResolver) UserJoinTopic(ctx context.Context, id int) (*models.UserJoinTopic, error) {
	thunk := midacontext.GetLoader[loader.Loader](ctx).UserJoinTopic.Load(ctx, id)
	return thunk()
}

// Topic is the resolver for the topic field.
func (r *recentMatchingResolver) Topic(ctx context.Context, obj *models.RecentMatching) (*models.Topic, error) {
	return &models.Topic{ID: obj.TopicID}, nil
}

// City is the resolver for the city field.
func (r *recentMatchingResolver) City(ctx context.Context, obj *models.RecentMatching) (*models.Area, error) {
	return &models.Area{Code: obj.CityID}, nil
}

// Matchings is the resolver for the matchings field.
func (r *recentMatchingResolver) Matchings(ctx context.Context, obj *models.RecentMatching) ([]*models.Matching, error) {
	thunk := midacontext.GetLoader[loader.Loader](ctx).Matching.LoadMany(ctx, obj.MatchingIDs)
	matchings, errors := thunk()
	if errors != nil {
		return nil, multierr.Combine(errors...)
	}
	return matchings, nil
}

// Topic is the resolver for the topic field.
func (r *userJoinTopicResolver) Topic(ctx context.Context, obj *models.UserJoinTopic) (*models.Topic, error) {
	return &models.Topic{ID: obj.TopicID}, nil
}

// City is the resolver for the city field.
func (r *userJoinTopicResolver) City(ctx context.Context, obj *models.UserJoinTopic) (*models.Area, error) {
	return &models.Area{Code: obj.CityID}, nil
}

// User is the resolver for the user field.
func (r *userJoinTopicResolver) User(ctx context.Context, obj *models.UserJoinTopic) (*models.User, error) {
	return &models.User{ID: obj.UserID}, nil
}

// Matching is the resolver for the matching field.
func (r *userJoinTopicResolver) Matching(ctx context.Context, obj *models.UserJoinTopic) (*models.Matching, error) {
	thunk := midacontext.GetLoader[loader.Loader](ctx).Matching.Load(ctx, obj.LatestMatchingID)
	return thunk()
}

// RecentMatching returns RecentMatchingResolver implementation.
func (r *Resolver) RecentMatching() RecentMatchingResolver { return &recentMatchingResolver{r} }

// UserJoinTopic returns UserJoinTopicResolver implementation.
func (r *Resolver) UserJoinTopic() UserJoinTopicResolver { return &userJoinTopicResolver{r} }

type recentMatchingResolver struct{ *Resolver }
type userJoinTopicResolver struct{ *Resolver }

package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.35

import (
	"context"
	"whale/pkg/dbquery"
	"whale/pkg/modelutil"

	"github.com/golang-module/carbon"
	"github.com/letjoy-club/mida-tool/dbutil"
	"github.com/letjoy-club/mida-tool/midacode"
	"github.com/letjoy-club/mida-tool/midacontext"
)

// YesterdayMatchingCount is the resolver for the yesterdayMatchingCount field.
func (r *queryResolver) YesterdayMatchingCount(ctx context.Context) (int, error) {
	db := dbutil.GetDB(ctx)
	Matching := dbquery.Use(db).Matching
	now := carbon.Now()
	start := now.StartOfDay().AddDays(-3)
	end := now.StartOfDay()
	count, err := Matching.WithContext(ctx).Where(Matching.CreatedAt.Between(start.ToStdTime(), end.ToStdTime())).Count()
	return int(count), err
}

// MotionSummary is the resolver for the motionSummary field.
func (r *queryResolver) MotionSummary(ctx context.Context) (map[string]interface{}, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() {
		return nil, midacode.ErrNotPermitted
	}

	summary, err := modelutil.GenMotionSummary(ctx)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"cities":    summary.Cities,
		"maleNum":   summary.MaleNum,
		"femaleNum": summary.FemaleNum,
	}, nil
}

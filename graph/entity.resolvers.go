package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.28

import (
	"context"
	"fmt"
	"whale/pkg/loader"
	"whale/pkg/models"

	"github.com/letjoy-club/mida-tool/midacode"
	"github.com/letjoy-club/mida-tool/midacontext"
)

// FindMatchingByID is the resolver for the findMatchingByID field.
func (r *entityResolver) FindMatchingByID(ctx context.Context, id string) (*models.Matching, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() && !token.IsUser() {
		return nil, midacode.ErrNotPermitted
	}
	thunk := midacontext.GetLoader[loader.Loader](ctx).Matching.Load(ctx, id)
	return thunk()
}

// FindMatchingQuotaByUserID is the resolver for the findMatchingQuotaByUserID field.
func (r *entityResolver) FindMatchingQuotaByUserID(ctx context.Context, userID string) (*models.MatchingQuota, error) {
	token := midacontext.GetClientToken(ctx)
	if !token.IsAdmin() && !token.IsUser() {
		return nil, midacode.ErrNotPermitted
	}
	if token.IsUser() {
		userID = token.UserID()
	}
	thunk := midacontext.GetLoader[loader.Loader](ctx).MatchingQuota.Load(ctx, userID)
	return thunk()
}

// FindUserByID is the resolver for the findUserByID field.
func (r *entityResolver) FindUserByID(ctx context.Context, id string) (*models.User, error) {
	panic(fmt.Errorf("not implemented: FindUserByID - findUserByID"))
}

// Entity returns EntityResolver implementation.
func (r *Resolver) Entity() EntityResolver { return &entityResolver{r} }

type entityResolver struct{ *Resolver }

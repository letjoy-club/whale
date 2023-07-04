package matcher

import (
	"context"
	"fmt"
	"whale/pkg/gqlient/hoopoe"

	"github.com/letjoy-club/mida-tool/midacontext"
)

func getBlacklistRelationship(ctx context.Context, ids []string) map[string]struct{} {
	relationships := map[string]struct{}{}
	resp, err := hoopoe.GetBlacklistRelationship(ctx, midacontext.GetServices(ctx).Hoopoe, ids)
	if err != nil {
		fmt.Println("failed to get blacklist relationship", err)
		return relationships
	}
	for _, tuple := range resp.BlacklistRelationship {
		if tuple.A > tuple.B {
			relationships[tuple.B+"-"+tuple.A] = struct{}{}
		} else {
			relationships[tuple.A+"-"+tuple.B] = struct{}{}
		}
	}
	return relationships
}

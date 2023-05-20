package keyer

import (
	"context"
	"fmt"
	"sort"
	"time"
	"whale/pkg/whalecode"

	"github.com/letjoy-club/mida-tool/redisutil"
)

func UserMatching(uid string) string {
	if uid == "" {
		return ""
	}
	return fmt.Sprintf("user:%s:matching", uid)
}

func Invitation(invitationId string) string {
	if invitationId == "" {
		return ""
	}
	return fmt.Sprintf("user:%s:matching", invitationId)
}

func MatchingResult(mid string) string {
	if mid == "" {
		return ""
	}
	return fmt.Sprintf("matching:%s:result", mid)
}

func LockAll(ctx context.Context, keys ...string) (release func(context.Context), err error) {
	sort.Strings(keys)
	locker := redisutil.GetLocker(ctx)

	releasers := []func(context.Context) error{}

	for _, key := range keys {
		if key == "" {
			continue
		}
		lock, err := locker.Obtain(ctx, key, time.Second*10, nil)
		if err != nil {
			if len(releasers) > 0 {
				for _, releaser := range releasers {
					releaser(ctx)
				}
			}
			return nil, whalecode.ErrResourceBusy
		}
		releasers = append(releasers, lock.Release)
	}

	return func(ctx context.Context) {
		for _, releaser := range releasers {
			releaser(ctx)
		}
	}, nil
}

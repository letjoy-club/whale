package modelutil

import (
	"context"
	"time"
	"whale/pkg/dbquery"

	"github.com/letjoy-club/mida-tool/dbutil"
)

func MotionExpire(ctx context.Context) error {
	db := dbutil.GetDB(ctx)
	Motion := dbquery.Use(db).Motion

	ids := []string{}

	err := Motion.WithContext(ctx).Where(Motion.Deadline.Lt(time.Now())).Pluck(Motion.ID, &ids)
	if err != nil {
		return err
	}

	for _, id := range ids {
		// 第二个参数是操作人，这里为空即可
		if err := CloseMotion(ctx, "", id); err != nil {
			return err
		}
	}
	return nil
}

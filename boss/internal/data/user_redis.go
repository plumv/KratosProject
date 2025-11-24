package data

import (
	"boss/internal/biz"
	"context"
	"encoding/json"
	"time"
)

func (repo userRepo) UpdateWithCache(ctx context.Context, id uint64, u *biz.User) error {
	us, err := repo.data.db.User(ctx).UpdateOneID(id).
		SetUsername(*u.Username).
		SetPassword(*u.Password).
		SetAge(*u.Age).Save(ctx)
	if j, err := json.Marshal(us); err == nil {
		_ = repo.data.rdb.Set(ctx, cacheKey(id), j, time.Minute*5).Err()
	}
	return err
}

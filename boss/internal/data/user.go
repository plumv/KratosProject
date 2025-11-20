package data

import (
	"boss/internal/biz"
	"boss/internal/data/ent"
	"boss/internal/data/ent/user"
	"context"
	"encoding/json"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserRepo .
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (u userRepo) Save(ctx context.Context, user *biz.User) (uuid.UUID, error) {
	us, err := u.data.db.User.Create().
		SetUsername(*user.Username).
		SetPassword(*user.Password).
		SetAge(*user.Age).Save(ctx)
	if j, err := json.Marshal(us); err == nil {
		_ = u.data.rdb.Set(ctx, cacheKey(us.ID), j, time.Minute*5).Err()
	}
	return us.ID, err
}

func (u userRepo) Update(ctx context.Context, id uuid.UUID, user *biz.User) error {
	_, err := u.data.db.User.Get(ctx, id)
	if err != nil {
		return err
	}
	us, err := u.data.db.User.UpdateOneID(id).
		SetUsername(*user.Username).
		SetPassword(*user.Password).
		SetAge(*user.Age).Save(ctx)
	if j, err := json.Marshal(us); err == nil {
		_ = u.data.rdb.Set(ctx, cacheKey(id), j, time.Minute*5).Err()
	}

	return err
}

func (u userRepo) FindByID(ctx context.Context, id uuid.UUID) (*biz.User, error) {
	// 1. 先读缓存
	key := cacheKey(id)
	if val, err := u.data.rdb.Get(ctx, key).Result(); err == nil {
		var us ent.User
		if err := json.Unmarshal([]byte(val), &us); err == nil {
			return &biz.User{
				ID:       &us.ID,
				Username: &us.Username,
				Age:      &us.Age,
			}, nil
		}
	}
	us, err := u.data.db.User.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	// 回填缓存
	if j, err := json.Marshal(us); err == nil {
		_ = u.data.rdb.Set(ctx, cacheKey(us.ID), j, time.Minute*5).Err()
	}
	return &biz.User{
		ID:       &us.ID,
		Username: &us.Username,
		Age:      &us.Age,
	}, nil
}

func (u userRepo) DeleteByID(ctx context.Context, id uuid.UUID) error {
	_ = u.data.rdb.Del(ctx, cacheKey(id)).Err()
	return u.data.db.User.DeleteOneID(id).Exec(ctx)
}

func (u userRepo) ListAll(ctx context.Context, f *biz.UserFilter, o *[]*biz.Order) ([]*biz.User, error) {
	query := u.data.db.User.Query()
	where(query, f)
	order(query, o)
	users, err := query.All(ctx)
	if err != nil {
		return nil, err
	}
	rv := make([]*biz.User, 0)
	for _, p := range users {
		rv = append(rv, &biz.User{
			ID:       &p.ID,
			Username: &p.Username,
			Age:      &p.Age,
		})
	}
	return rv, nil
}

func (u userRepo) PageAll(ctx context.Context, f *biz.UserFilter, p *biz.Page) ([]*biz.User, int, error) {
	q := u.data.db.User.Query()
	where(q, f)
	order(q, p.Orders)
	t, err := q.Count(ctx)
	if err != nil {
		return nil, t, err
	}
	page(q, p)
	users, err := q.All(ctx)
	if err != nil {
		return nil, t, err
	}
	rv := make([]*biz.User, 0)
	for _, p := range users {
		rv = append(rv, &biz.User{
			ID:       &p.ID,
			Username: &p.Username,
			Age:      &p.Age,
		})
	}
	return rv, t, nil
}

// where 动态查询条件拼接
func where(u *ent.UserQuery, f *biz.UserFilter) {
	// ---------- 条件翻译 ----------
	if v := f.Name; v != nil && *v != "" {
		u.Where(user.UsernameContains(*v))
	}
	if v := f.AgeEQ; v != nil {
		u.Where(user.AgeEQ(*v))
	}
	if v := f.AgeGTE; v != nil {
		u.Where(user.AgeGTE(*v))
	}
	if v := f.AgeLTE; v != nil {
		u.Where(user.AgeLTE(*v))
	}
	if f.IDIn != nil && len(*f.IDIn) > 0 {
		u.Where(user.IDIn(*f.IDIn...))
	}
	if v := f.CreatedAfter; v != nil {
		u.Where(user.CreatedAtGTE(*v))
	}
}

// order 排序方式
func order(u *ent.UserQuery, o *[]*biz.Order) {
	if len(*o) == 0 {
		u.Order(user.ByID(sql.OrderDesc()))
		return
	}
	for _, b := range *o {
		if user.ValidColumn(*b.Field) {
			asc := sql.OrderAsc()
			if *b.Desc {
				asc = sql.OrderDesc()
			}
			u.Order(sql.OrderByField(*b.Field, asc).ToFunc())
		}
	}
}

// page 分页条件
func page(q *ent.UserQuery, p *biz.Page) {
	q.Limit(int(*p.Limit)).
		Offset(int((*p.Page - 1) * *p.Limit))
}

func cacheKey(id uuid.UUID) string {
	return "cache:boss:user:" + id.String()
}

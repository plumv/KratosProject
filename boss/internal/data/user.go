package data

import (
	"boss/internal/biz"
	"boss/internal/data/ent"
	"boss/internal/data/ent/user"
	"context"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
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

func (repo userRepo) Save(ctx context.Context, u *biz.User) (uint64, error) {
	us, err := repo.data.db.User(ctx).Create().
		SetUsername(*u.Username).
		SetPassword(*u.Password).
		SetAge(*u.Age).Save(ctx)
	repo.data.cdb.Add(ctx, cacheKey(us.ID), us, time.Minute*5)
	return us.ID, err
}

func (repo userRepo) Update(ctx context.Context, id uint64, u *biz.User) error {
	us, err := repo.data.db.User(ctx).UpdateOneID(id).
		SetUsername(*u.Username).
		SetPassword(*u.Password).
		SetAge(*u.Age).Save(ctx)

	repo.data.cdb.Update(ctx, cacheKey(id), us, time.Minute*5)

	return err
}

func (repo userRepo) FindByID(ctx context.Context, id uint64) (*biz.User, error) {
	us := new(ent.User)
	err := repo.data.cdb.Get(ctx, cacheKey(id), us)
	if err == nil {
		return &biz.User{
			ID:       &id,
			Username: &us.Username,
			Age:      &us.Age,
		}, nil
	}
	us, err = repo.data.db.User(ctx).Get(ctx, id)
	if err != nil {
		return nil, err
	}
	repo.data.cdb.Add(ctx, cacheKey(id), us, time.Minute*5)
	return &biz.User{
		ID:       &us.ID,
		Username: &us.Username,
		Age:      &us.Age,
	}, nil
}

func (repo userRepo) DeleteByID(ctx context.Context, id uint64) error {
	repo.data.cdb.Delete(ctx, cacheKey(id))
	return repo.data.db.User(ctx).DeleteOneID(id).Exec(ctx)
}

func (repo userRepo) ListAll(ctx context.Context, f *biz.UserFilter, o *[]*biz.Order) ([]*biz.User, error) {
	query := repo.data.db.User(ctx).Query()
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

func (repo userRepo) PageAll(ctx context.Context, f *biz.UserFilter, p *biz.Page) ([]*biz.User, int, error) {
	q := repo.data.db.User(ctx).Query()
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

func cacheKey(id uint64) string {
	return fmt.Sprintf("cache:boss:user:%v", id)
}

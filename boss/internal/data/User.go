package data

import (
	"boss/internal/biz"
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// userRepo 是biz.UserRepo接口的实现，实现了数据库的具体操作
type userRepo struct {
	data *Data
	log  *log.Helper
}

func (r *userRepo) FindAll(ctx context.Context) ([]*biz.User, error) {
	panic("implement me")
}

func (r *userRepo) FindPage(ctx context.Context) ([]*biz.User, error) {
	panic("implement me")
}

func (r *userRepo) Save(ctx context.Context, g *biz.User) (*biz.User, error) {
	return g, nil
}

func (r *userRepo) Update(ctx context.Context, g *biz.User) (*biz.User, error) {
	return g, nil
}

func (r *userRepo) FindByID(context.Context, int64) (*biz.User, error) {
	return nil, nil
}

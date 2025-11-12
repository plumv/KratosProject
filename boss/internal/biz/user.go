package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// User 定义数据库中对于表映射的结构体
type User struct {
	Id       int64
	Name     string
	Age      int32
	Email    string
	Password string
}

// UserRepo 定义数据库可以提供什么操作的接口
type UserRepo interface {
	Save(context.Context, *User) (*User, error)
	Update(context.Context, *User) (*User, error)
	FindByID(context.Context, int64) (*User, error)
	FindAll(context.Context) ([]*User, error)
	FindPage(context.Context) ([]*User, error)
}

// UserUsecase 表示可以如何操作数据库，底层是通过UserRepo来操作数据库的
type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
}

func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *UserUsecase) CreateUser(ctx context.Context, g *User) (*User, error) {
	uc.log.WithContext(ctx).Infof("CreateUser: %v", g.Name)
	return uc.repo.Save(ctx, g)
}

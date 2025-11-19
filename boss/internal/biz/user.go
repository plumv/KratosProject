package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

// User is a User model.
type User struct {
	// 主键
	ID uuid.UUID
	// 用户名
	Username string
	// 密码
	Password string
	// 年龄
	Age int32
}

// UserFilter 用户查询
type UserFilter struct {
	Name         *string // 模糊
	AgeEQ        *int32  // 精确
	AgeGTE       *int32  // 范围
	AgeLTE       *int32
	IDIn         []uuid.UUID // IN
	CreatedAfter *time.Time
}

// UserRepo is a Greater repo.
type UserRepo interface {
	Save(context.Context, *User) (uuid.UUID, error)
	Update(context.Context, uuid.UUID, *User) error
	FindByID(context.Context, uuid.UUID) (*User, error)
	DeleteByID(context.Context, uuid.UUID) error
	ListAll(context.Context, *UserFilter, *[]*Order) ([]*User, error)
	PageAll(context.Context, *UserFilter, *Page) ([]*User, int, error)
}

// UserUsecase is a User usecase.
type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
}

// NewUserUsecase new a User usecase.
func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateUser creates a User, and returns the new User.
func (uc *UserUsecase) CreateUser(ctx context.Context, g *User) (uuid.UUID, error) {
	return uc.repo.Save(ctx, g)
}
func (uc *UserUsecase) UpdateUser(ctx context.Context, id uuid.UUID, g *User) error {
	return uc.repo.Update(ctx, id, g)
}
func (uc *UserUsecase) FindUser(ctx context.Context, id uuid.UUID) (*User, error) {
	return uc.repo.FindByID(ctx, id)
}
func (uc *UserUsecase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return uc.repo.DeleteByID(ctx, id)
}

func (uc *UserUsecase) ListUser(ctx context.Context, filter *UserFilter, order *[]*Order) ([]*User, error) {
	return uc.repo.ListAll(ctx, filter, order)
}

func (uc *UserUsecase) PageUser(ctx context.Context, filter *UserFilter, page *Page) ([]*User, int, error) {
	return uc.repo.PageAll(ctx, filter, page)
}

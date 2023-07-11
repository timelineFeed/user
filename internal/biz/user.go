package biz

import (
	"context"
	v1 "user/api/user/v1"
	"user/internal/data/model"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

type User struct {
}

// UserRepo is a user repo.
type UserRepo interface {
	Create(context.Context, *model.User) error
	Update(context.Context, *model.User) (*model.User, error)
	FindByID(context.Context, int64) (*model.User, error)
	List(context.Context, []int64) ([]*model.User, error)
}

// UserRepo is a Greeter usecase.
type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
}

// NewUserRepoUsecase new a Greeter usecase.
func NewUserRepoUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateUser creates a Greeter, and returns the new Greeter.
func (uc *UserUsecase) CreateUser(ctx context.Context, g *UserRepo) (*UserRepo, error) {
	//TODO implement me
	panic("implement me")
}

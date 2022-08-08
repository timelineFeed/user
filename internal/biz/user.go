package biz

import (
	"context"
	"github.com/timelineFeed/user/internal/model"

	v1 "github.com/timelineFeed/user/api/user/v1"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

// User is a User model.
type User struct {
	model.User
}

// UserRepo is a Greater repo.
type UserRepo interface {
	Create(context.Context, *User) error
	FindPasswordByID(context.Context, uint64) (*User, error)
	FindByID(context.Context, uint64) (*User, error)
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

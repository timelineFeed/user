package service

import (
	"context"
	"github.com/timelineFeed/user/internal/model"
	"time"

	v1 "github.com/timelineFeed/user/api/user/v1"
	"github.com/timelineFeed/user/internal/biz"
)

// GreeterService is a greeter service.
type UserService struct {
	v1.UnimplementedUserServer
	uc *biz.UserUsecase
}

// NewUserService new a greeter service.
func NewUserService(uc *biz.UserUsecase) *UserService {
	return &UserService{uc: uc}
}

func (u *UserService) Register(ctx context.Context, in *v1.RegisterRequest) (*v1.RegisterReply, error) {
	err := u.uc.Register(ctx, &biz.User{
		User: &model.User{
			Name:      in.Name,
			Password:  in.Password,
			Telephone: in.Telephone,
			Email:     in.Email,
			Image:     in.Image,
			Status:    0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	})
	if err != nil {
		return nil, err
	}
	return &v1.RegisterReply{}, err
}
func (u *UserService) Login(ctx context.Context, in *v1.LoginRequest) (*v1.LoginReply, error) {
	login, err := u.uc.Login(ctx, in.Email, in.Password)
	if err != nil {
		return nil, err
	}
	return &v1.LoginReply{
		Name:  login.Name,
		Image: login.Image,
		Token: login.Token,
	}, nil
}
func (u *UserService) UserDetail(ctx context.Context, in *v1.UserDetailRequest) (*v1.UserDetailReply, error) {
	user, err := u.uc.UserDetail(ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	return &v1.UserDetailReply{
		Name:      user.Name,
		Telephone: user.Telephone,
		Email:     user.Email,
		Image:     user.Image,
	}, nil
}

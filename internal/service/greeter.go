package service

import (
	"context"

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

// SayHello implements user.GreeterServer.
func (s *GreeterService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	g, err := s.uc.CreateGreeter(ctx, &biz.Greeter{Hello: in.Name})
	if err != nil {
		return nil, err
	}
	return &v1.HelloReply{Message: "Hello " + g.Hello}, nil
}

func (u *UserService) Register(ctx context.Context, in *v1.RegisterRequest) (*v1.RegisterReply, error) {
	panic("impl")
}
func (u *UserService) Login(ctx context.Context, in *v1.LoginRequest) (*v1.LoginReply, error) {
	panic("impl")
}
func (u *UserService) UserDetail(ctx context.Context, in *v1.UserDetailRequest) (*v1.UserDetailReply, error) {
	panic("impl")
}

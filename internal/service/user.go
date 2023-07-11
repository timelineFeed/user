package service

import (
	"context"

	v1 "user/api/user/v1"
	"user/internal/biz"
)

// UserService is a greeter service.
type UserService struct {
	v1.UnimplementedUserServer
	uc *biz.UserUsecase
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.UserUsecase) *UserService {
	return &UserService{uc: uc}
}

func (s *UserService) Register(ctx context.Context, in *v1.RegisterRequest) (*v1.RegisterReply, error) {

}

// Login user login api
func (s *UserService) Login(ctx context.Context, in *v1.LoginRequest) (*v1.LoginReply, error) {

}

// GetUser get one user info
func (s *UserService) GetUser(ctx context.Context, in *v1.GetUserRequest) (*v1.GetUserReply, error) {

}

// GetUserList get user info list
func (s *UserService) GetUserList(ctx context.Context, in *v1.GetUserListRequest) (*v1.GetUserListReply, error) {

}

// Update update user info
func (s *UserService) Update(ctx context.Context, in *v1.UpdateRequest) (*v1.UpdateReply, error) {

}

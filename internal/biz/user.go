package biz

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"gorm.io/gorm"
	"time"
	v1 "user/api/user/v1"
	"user/internal/data/model"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound  = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
	ErrInternal      = errors.InternalServer(v1.ErrorReason_INTERNAL_ERROR.String(), "service internal error")
	ErrBdOperate     = errors.InternalServer(v1.ErrorReason_DB_OPERATE_ERROR.String(), "db operate error")
	ErrPasswordCheck = errors.InternalServer(v1.ErrorReason_USER_PASSWD_CHECK_FAILED.String(), "passwd check failed")
)

type User struct {
}

// UserRepo is a user repo.
type UserRepo interface {
	Create(context.Context, *model.User) error
	Update(context.Context, *model.User) (*model.User, error)
	FindByID(context.Context, uint64) (*model.User, error)
	List(context.Context, []uint64) ([]*model.User, error)
}

// UserUseCase is a user usecase.
type UserUseCase struct {
	repo UserRepo
	log  *log.Helper
}

// NewUserUseCase new a Greeter usecase.
func NewUserUseCase(repo UserRepo, logger log.Logger) *UserUseCase {
	return &UserUseCase{repo: repo, log: log.NewHelper(logger)}
}

// Register user register
func (uc *UserUseCase) Register(ctx context.Context, in *v1.RegisterRequest) error {
	// 对密码进行加密
	password, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("generate password err=%+v", err)
		return ErrInternal
	}
	user := &model.User{
		Name:      in.GetName(),
		Password:  string(password),
		Telephone: in.GetPhone(),
		Email:     in.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = uc.repo.Create(ctx, user)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("create user err=%+v", err)
		return ErrBdOperate
	}
	return nil
}

func (uc *UserUseCase) Login(ctx context.Context, in *v1.LoginRequest) (*v1.LoginReply, error) {
	user, err := uc.repo.FindByID(ctx, in.Uid)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("find user by id err=%+v", err)
		return nil, ErrUserNotFound
	}
	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.Password))
	if err != nil {
		uc.log.WithContext(ctx).Errorf("compare password err=%+v", err)
		return nil, ErrPasswordCheck
	}
	// todo jwt

	return &v1.LoginReply{
		User: &v1.UserInfo{
			Uid:   user.ID,
			Phone: user.Telephone,
			Email: user.Email,
			Name:  user.Name,
			Extra: convertMap(user.Extra),
		},
		Token: "",
	}, nil
}

// convertMap todo test
func convertMap(originalMap map[string]interface{}) map[string]*anypb.Any {
	convertedMap := make(map[string]*anypb.Any)
	for key, value := range originalMap {
		pv, ok := value.(proto.Message)
		if !ok {
			panic("value can not assert proto.Message")
		}
		anyValue, err := anypb.New(pv)
		if err != nil {
			panic(fmt.Sprintf("anypb new err=%+v", err))
		}
		convertedMap[key] = anyValue
	}
	return convertedMap
}

// GetUser get one user info
func (uc *UserUseCase) GetUser(ctx context.Context, in *v1.GetUserRequest) (*v1.GetUserReply, error) {
	user, err := uc.repo.FindByID(ctx, in.GetUid())
	if err != nil {
		uc.log.WithContext(ctx).Errorf("find user user_id=%d,err=%+v", in.GetUid(), err)
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}
		return nil, ErrBdOperate
	}
	return &v1.GetUserReply{User: &v1.UserInfo{
		Uid:   user.ID,
		Phone: user.Telephone,
		Email: user.Email,
		Name:  user.Name,
		Extra: convertMap(user.Extra),
	}}, nil
}

// GetUserList get user info list
func (uc *UserUseCase) GetUserList(ctx context.Context, in *v1.GetUserListRequest) (*v1.GetUserListReply, error) {
	users, err := uc.repo.List(ctx, in.GetUidList())
	if err != nil {
		uc.log.WithContext(ctx).Errorf("get user list user_ids=%+v,err=%+v", in.GetUidList(), err)
		return nil, ErrBdOperate
	}
	if len(users) != len(in.GetUidList()) {
		uc.log.WithContext(ctx).Errorf("get user list failed,len(users)=%d,len(in.GetUidList())=%d",
			len(users), len(in.GetUidList()))
		return nil, ErrInternal
	}
	resUsers := make([]*v1.UserInfo, 0, len(users))
	for _, user := range users {
		resUser := &v1.UserInfo{
			Uid:   user.ID,
			Phone: user.Telephone,
			Email: user.Email,
			Name:  user.Name,
			Extra: convertMap(user.Extra),
		}
		resUsers = append(resUsers, resUser)
	}
	return &v1.GetUserListReply{Users: resUsers}, nil
}

// Update update user info
func (uc *UserUseCase) Update(ctx context.Context, in *v1.UpdateRequest) (*v1.UpdateReply, error) {

	// todo extra pb any to go any
	user := &model.User{
		ID:        in.User.GetUid(),
		Name:      in.User.GetName(),
		Telephone: in.User.GetPhone(),
		Email:     in.User.GetEmail(),
		UpdatedAt: time.Now(),
	}
	user, err := uc.repo.Update(ctx, user)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("update user err=%+v", err)
		return nil, ErrBdOperate
	}
	return &v1.UpdateReply{}, nil
}

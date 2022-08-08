package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/timelineFeed/user/internal/model"
	"regexp"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/crypto/bcrypt"

	jwtv4 "github.com/golang-jwt/jwt/v4"
	v1 "github.com/timelineFeed/user/api/user/v1"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

// User is a User model.
type User struct {
	model.User
	Token string
}

// UserRepo is a Greater repo.
type UserRepo interface {
	Create(context.Context, *User) error
	FindByPhoneOrEmail(context.Context, string, string) (*User, error)
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

func (u *UserUsecase) Register(ctx context.Context, user *User) error {
	// 密码转hash存储
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(passwordHash)
	return u.repo.Create(ctx, user)
}

func (u *UserUsecase) Login(ctx context.Context, emailOrPhone, password string) (*User, error) {
	phone := ""
	email := ""
	// 做正则表达式校验输入
	ok, err := regexp.Match("^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$", []byte(emailOrPhone))
	if err != nil {
		return nil, err
	}
	if ok {
		email = emailOrPhone
	}
	ok2, err := regexp.Match("^(13[0-9]|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18[0-9]|19[0-35-9])\\d{8}$",
		[]byte(emailOrPhone))
	if err != nil {
		return nil, err
	}
	if !ok2 {
		phone = emailOrPhone
	}
	if !(ok && ok2) {
		return nil, errors.New(int(v1.ErrorReason_INCORRECT_INPUT_PARAMETERS),
			"输入的参数有误", "输入的参数有误")
	}
	// db 查询数据
	user, err := u.repo.FindByPhoneOrEmail(ctx, phone, email)
	if err != nil {
		return nil, errors.New(int(v1.ErrorReason_SERVICE_INTERNAL_ERROR), "查询数据出错", "查询数据出错")
	}
	//校验密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New(int(v1.ErrorReason_SERVICE_INTERNAL_ERROR), "密码错误", "密码错误")
	}
	// 用jwt 生成token
	jwt.WithSigningMethod(jwtv4.SigningMethodHS256)

	token := jwtv4.NewWithClaims(jwtv4.SigningMethodHS256, jwtv4.StandardClaims{})

}

func (u *UserUsecase) UserDetail(ctx context.Context, id uint64) (*User, error) {
	return u.repo.FindByID(ctx, id)
}

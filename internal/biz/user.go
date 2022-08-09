package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/timelineFeed/user/internal/data"
	"github.com/timelineFeed/user/internal/model"
	"golang.org/x/sync/singleflight"
	"regexp"
	"strconv"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	jwtv4 "github.com/golang-jwt/jwt/v4"
	"github.com/jellydator/ttlcache/v2"
	v1 "github.com/timelineFeed/user/api/user/v1"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
	jwtKey          = []byte("key")
)

const (
	TokenExpireDuration     = 24 * time.Hour
	UserRedisExpireDuration = 60 * time.Second
	UserTTLExpireDuration   = 10 * time.Second
	Issuer                  = "user service"
)

// User is a User model.
type User struct {
	*model.User
	Token string
}

type AuthClaim struct {
	UserId uint64 `json:"user_id"`
	jwtv4.StandardClaims
}

// UserRepo is a Greater repo.
type UserRepo interface {
	Create(context.Context, *User) error
	FindByPhoneOrEmail(context.Context, string, string) (*User, error)
	FindByID(context.Context, uint64) (*User, error)
	SetUser(context.Context, string, string, time.Duration) error
	GetUser(context.Context, string) (string, error)
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

// Register 注册
func (u *UserUsecase) Register(ctx context.Context, user *User) error {
	// 密码转hash存储
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(passwordHash)
	return u.repo.Create(ctx, user)
}

// Login 用户登录
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
	c := AuthClaim{
		UserId: user.Id,
		StandardClaims: jwtv4.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			Issuer:    Issuer,
		},
	}
	token := jwtv4.NewWithClaims(jwtv4.SigningMethodHS256, c)
	user.Token, err = token.SignedString(jwtKey)
	if err != nil {
		return nil, errors.New(int(v1.ErrorReason_SERVICE_INTERNAL_ERROR), "生成token失败", "生成token失败")
	}
	return user, nil

}

// UserDetail 查询用户详情
func (u *UserUsecase) UserDetail(ctx context.Context, id uint64) (*User, error) {
	// 场景分为 主态查询、客态查询
	// 主态查询直接返回全部数据
	// 客态查询需要判断信息是否公开

	// TODO 客态实现
	// 主态实现
	var sf singleflight.Group
	v, err, _ := sf.Do(strconv.Itoa(int(id)), u.genUserDetail(ctx, id))
	if err != nil {
		return nil, err
	}
	ret, ok := v.(*User)
	if !ok {
		return nil, errors.New(int(v1.ErrorReason_SERVICE_INTERNAL_ERROR), "反射失败", "服务内部错误")
	}
	return ret, nil
}

// genUserDetail 合并查询用户信息的请求
func (u *UserUsecase) genUserDetail(ctx context.Context, userID uint64) func() (interface{}, error) {
	return func() (interface{}, error) {
		// 先查询内存缓存、再查询redis、再查询db
		key := fmt.Sprintf(data.UserRedisKey, userID)
		var ur model.User
		cache := ttlcache.NewCache()
		value, err := cache.Get(key)
		if err != nil && err != ttlcache.ErrNotFound {
			return nil, err
		}
		// ttl 查询到的处理
		if err == nil {
			v, ok := value.(string)
			if !ok {
				return nil, errors.New(int(v1.ErrorReason_SERVICE_INTERNAL_ERROR), "断言失败", "服务内部错误")
			}
			err = json.Unmarshal([]byte(v), &ur)
			if err != nil {
				return nil, errors.New(int(v1.ErrorReason_SERVICE_INTERNAL_ERROR), "序列化失败", "服务内部错误")
			}
			return User{User: &ur}, nil
		}

		userStr, err := u.repo.GetUser(ctx, key)
		if err != nil && err != redis.Nil {
			return nil, errors.New(int(v1.ErrorReason_SERVICE_INTERNAL_ERROR), "查询redis出错", "查询出错")
		}
		// redis查询到数据
		if err == nil {

			err = json.Unmarshal([]byte(userStr), &u)
			if err != nil {
				return nil, errors.New(int(v1.ErrorReason_SERVICE_INTERNAL_ERROR), "redis获取的值序列化失败", "服务内部错误")
			}
			return &User{User: &ur}, nil
		}
		// 查询db,写入redis
		user, err := u.repo.FindByID(ctx, userID)
		if err != nil {
			return nil, errors.New(int(v1.ErrorReason_SERVICE_INTERNAL_ERROR), "查询db出错", "服务内部错误")
		}
		// 回写数据
		go func() {
			userByte, err := json.Marshal(&user)
			if err != nil {
				log.Errorf("序列化user失败,err=%+v", err)
				return
			}
			//写内存缓存 ,写入失败不返回
			err = cache.SetWithTTL(key, string(userByte), UserTTLExpireDuration)
			if err != nil {
				log.Errorf("ttl写入失败,err=%+v", err)
			}
			//写redis
			err = u.repo.SetUser(ctx, key, string(userByte), UserRedisExpireDuration)
			if err != nil {
				log.Errorf("redis写入失败,err=%+v", err)
			}
		}()
		return user, nil
	}
}

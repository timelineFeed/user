package data

import (
	"context"
	"user/internal/data/model"

	"user/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserRepo .
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (u userRepo) Create(ctx context.Context, user *model.User) error {
	return u.data.DB.WithContext(ctx).Create(user).Error
}

func (u userRepo) Update(ctx context.Context, user *model.User) (*model.User, error) {
	return user, u.data.DB.WithContext(ctx).Where("id = ?", user.ID).Updates(user).Error
}

func (u userRepo) FindByID(ctx context.Context, id int64) (*model.User, error) {
	user := &model.User{}
	return user, u.data.DB.WithContext(ctx).Where("id = ?", id).First(user).Error
}

func (u userRepo) List(ctx context.Context, ids []int64) ([]*model.User, error) {
	userList := make([]*model.User, 0, len(ids))
	return userList, u.data.DB.WithContext(ctx).Where("id in ?", ids).Find(&userList).Error
}

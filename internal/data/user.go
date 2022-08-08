package data

import (
	"context"
	"github.com/timelineFeed/user/internal/model"
	"gorm.io/gorm"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/timelineFeed/user/internal/biz"
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

func (r *userRepo) GetTableName() *gorm.DB {
	return r.data.db.Table((&model.User{}).TableName())
}

func (r *userRepo) Create(ctx context.Context, g *biz.User) error {
	return r.GetTableName().Create(g).Error
}

func (r *userRepo) FindPasswordByID(cxt context.Context, id uint64) (*biz.User, error) {
	user := new(biz.User)
	return user, r.GetTableName().Select("password").
		Where("id = ? and status = ?", id, model.StatusNormal).Find(&user).Error
}

func (r *userRepo) FindByID(ctx context.Context, id uint64) (*biz.User, error) {
	user := new(biz.User)
	return user, r.GetTableName().Select("name", "image", "telephone", "email").
		Where("id = ? and status = ?", id, model.StatusNormal).Find(&user).Error
}

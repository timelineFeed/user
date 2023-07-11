package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type User struct {
	ID        uint64    `gorm:"column:id"`
	Name      string    `gorm:"column:name"`       // 用户昵称
	Password  string    `gorm:"column:password"`   // 用户密码hash
	Telephone string    `gorm:"column:telephone"`  // 用户电话号码
	Email     string    `gorm:"column:email"`      //  用户邮箱号
	Status    int       `gorm:"column:status"`     // 状态，10-删除
	Extra     ExtraMap  `gorm:"column:extra"`      // 额外配置
	CreatedAt time.Time `gorm:"column:created_at"` // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at"` // 更新时间
}

func (m *User) TableName() string {
	return "timeline.user"
}

type ExtraMap map[string]interface{}

func (e ExtraMap) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("extra scan failed")
	}
	return json.Unmarshal(bytes, &e)
}
func (e ExtraMap) Value() (driver.Value, error) {
	return json.Marshal(e)
}

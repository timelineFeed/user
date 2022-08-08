package model

import "time"

type User struct {
	Id        uint64    `gorm:"column:id" json:"id"`
	Name      string    `gorm:"column:name" json:"name"`
	Password  string    `gorm:"column:password" json:"password"`
	Telephone string    `gorm:"column:telephone" json:"telephone"`
	Email     string    `gorm:"column:email" json:"email"`
	Image     string    `gorm:"column:image" json:"image"`
	Status    uint      `gorm:"column:status" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (u *User) TableName() string {
	return "timeline.user"
}

type StatusType uint

const (
	StatusNormal StatusType = 0  //用户正常状态
	StatusDelete StatusType = 10 //用户被删除
)

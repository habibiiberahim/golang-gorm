package belajar_golang_gorm

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID           string    `gorm:"primary_key; column:id;<-:create" json:"id,omitempty"`
	Password     string    `gorm:"column:password" json:"password,omitempty"`
	Name         Name      `gorm:"embedded"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime;<-:create" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoCreateTime" json:"updated_at"`
	Information  string    `gorm:"-"`
	Wallet       Wallet    `gorm:"foreignKey:user_id;references:id"`
	Addresses    []Address `gorm:"foreignKey:user_id;references:id"`
	LikeProducts []Product `gorm:"many2many:user_like_product;foreignKey:id;joinForeignKey:user_id;references:id;joinReferences:product_ide"`
}

type Name struct {
	FirstName  string `gorm:"column:first_name"`
	MiddleName string `gorm:"column:middle_name"`
	LastName   string `gorm:"column:last_name"`
}

func (u *User) TableName() string {
	return "users"
}

type UserLog struct {
	ID        int    `gorm:"primary_key;column:id;autoIncrement" json:"id,omitempty"`
	UserId    string `gorm:"column:user_id" json:"user_id,omitempty"`
	Action    string `gorm:"column:action" json:"action,omitempty"`
	CreatedAt int64  `gorm:"column:created_at;autoCreateTime:milli;" json:"created_at"`
	UpdatedAt int64  `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli" json:"updated_at"`
}

func (ul *UserLog) TableName() string {
	return "user_logs"
}

func (u *User) BeforeCreate(db *gorm.DB) error {
	if u.ID == "" {
		u.ID = "user-" + time.Now().Format("20230102110105")
	}
	return nil
}

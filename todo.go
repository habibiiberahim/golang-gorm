package belajar_golang_gorm

import "gorm.io/gorm"

type Todo struct {
	gorm.Model
	UserId      string `gorm:"column:user_id;"json:"user_id,omitempty"`
	Title       string `gorm:"column:title;"json:"title,omitempty"`
	Description string `gorm:"column:description;"json:"description,omitempty"`
}

func (t *Todo) TableName() string {
	return "todos"
}

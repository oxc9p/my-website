package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID `json:"id" gorm:"id; unique"`
	Username    string    `json:"username" gorm:"username; unique; primaryKey"`
	Password    string    `json:"password" gorm:"password"`
	VisibleName string    `json:"visible_name" gorm:"visible_name"`
	Permission  int       `json:"permission" gorm:"permission; default(0)"`
	Image       string    `json:"image" gorm:"image; default:'default.jpg'"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}

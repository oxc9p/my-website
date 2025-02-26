package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID          uuid.UUID `json:"id" gorm:"id; unique"`
	Username    string    `json:"username" gorm:"username; unique; primaryKey"`
	Password    string    `json:"password" gorm:"password"`
	Permission  int       `json:"permission" gorm:"permission; default(0)"`
	Image       string    `json:"image" gorm:"image'"`
	DateCreated time.Time `json:"date_created" gorm:"date_created"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	u.DateCreated = time.Now()
	u.Image = "static/images/default.jpg"
	return
}

package models

import "github.com/google/uuid"

type User struct {
	ID          uuid.UUID `json:"id" gorm:"id"`
	Username    string    `json:"username" gorm:"username; unique; primaryKey"`
	Password    string    `json:"password" gorm:"password"`
	VisibleName string    `json:"visible_name" gorm:"visible_name"`
	Permission  int       `json:"permission" gorm:"permission"`
	Image       string    `json:"image" gorm:"image; default:'default.jpg'"`
}

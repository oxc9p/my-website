package models

type Article struct {
	Image       string `json:"image" gorm:"primaryKey"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Id          string `json:"id"`
}

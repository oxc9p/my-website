package models

type Article struct {
	Image       string `json:"image"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Author      string `json:"author"`
	Id          string `json:"id" gorm:"primaryKey"`
}

package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"myPage/models"
)

func Init() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	e := db.AutoMigrate(&models.Article{}, &models.User{}, &models.Session{})
	if e != nil {
		panic("failed to migrate database")
	}
	return db
}

func CreateArticle(db *gorm.DB, article models.Article) {
	db.Create(&article)
}

func GetArticles(db *gorm.DB) []models.Article {
	var articles []models.Article
	db.Find(&articles)
	return articles
}

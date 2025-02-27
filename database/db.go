package database

import (
	"errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"myPage/models"
	"strings"
)

func Init() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	e := db.AutoMigrate(&models.Article{}, &models.User{}, &models.Session{}, &models.Project{})
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

func CreateProject(db *gorm.DB, project models.Project) {
	db.Create(project)
}

func GetProjects(db *gorm.DB) []models.Project {
	var projects []models.Project
	db.Find(&projects)
	return projects
}

// FindUserByUsername retrieves a user from the database by username.
func FindUserByUsername(db *gorm.DB, username string) (models.User, error) {
	var user models.User
	if err := db.First(&user, "username = ?", username).Error; err != nil {
		return user, err
	}
	return user, nil
}

// CreateUser creates a new user in the database.
func CreateUser(db *gorm.DB, user *models.User) error {
	if err := db.Create(user).Error; err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return errors.New("username already exists")
		}
		return errors.New("unable to create user")
	}
	return nil
}

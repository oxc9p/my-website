package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"myPage/database"
	"os"
)

func ArticleHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("article", fiber.Map{
			"WebLink":  os.Getenv("WEB_URL"),
			"Articles": database.GetArticles(db),
		})
	}
}

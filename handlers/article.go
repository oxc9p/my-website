package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"myPage/database"
	"myPage/models"
	"myPage/tools"
)

func ArticleHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var articles []models.Article
		return c.Render("article", fiber.Map{
			"WebLink":  tools.WebLink,
			"Articles": database.Get(db, &articles),
			"IsLogin":  tools.IsSessionExist(db, c),
		})
	}
}

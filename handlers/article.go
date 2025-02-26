package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"myPage/database"
	"myPage/tools"
)

func ArticleHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("article", fiber.Map{
			"WebLink":  tools.WebLink,
			"Articles": database.GetArticles(db),
			"IsLogin":  tools.IsSessionExist(db, c),
		})
	}
}

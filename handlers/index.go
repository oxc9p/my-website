package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"myPage/tools"
	"os"
)

func IndexHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"WebLink": os.Getenv("WEB_URL"),
			"IsLogin": tools.IsSessionExist(db, c),
		})
	}
}

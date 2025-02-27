package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"myPage/tools"
)

func AboutHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("about", fiber.Map{
			"WebLink": tools.WebLink,
			"IsLogin": tools.IsSessionExist(db, c),
		})
	}
}

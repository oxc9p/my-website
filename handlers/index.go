package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"myPage/database"
	"myPage/tools"
)

func IndexHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"WebLink":  tools.WebLink,
			"IsLogin":  tools.IsSessionExist(db, c),
			"Projects": database.GetProjects(db),
		})
	}
}

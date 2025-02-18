package handlers

import (
	"github.com/gofiber/fiber/v2"
	"os"
)

func IndexHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"WebLink": os.Getenv("WEB_URL"),
		})
	}
}

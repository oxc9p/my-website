package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"myPage/models"
	"myPage/tools"
)

func DashboardHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		session, err := tools.CheckSession(db, c)
		if err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		username := session.UserName
		var user models.User
		if err := db.First(&user, "username = ?", username).Error; err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		return c.SendString("Welcome to the dashboard, " + user.Username)
	}
}

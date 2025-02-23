package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"myPage/models"
	"myPage/tools"
	"time"
)

func LoginHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var user models.User
		if err := c.BodyParser(&user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}
		var enteredPassword = user.Password
		// Search user by username
		if err := db.First(&user, "username = ?", user.Username).Error; err != nil {
			return c.Status(401).SendString("Invalid username or password")
		}

		//Checking if the password is correct
		if tools.CheckPasswordHash(enteredPassword, user.Password) == false {
			return c.Status(401).SendString("Invalid username or password")
		}

		// Successfully authenticated
		sessionID := tools.GenerateSessionID()
		session := models.Session{
			SessionID:     sessionID,
			UserID:        user.ID,
			UserName:      user.Username,
			Authenticated: true,
		}
		if err := db.Save(&session).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Unable to save session",
			})
		}

		c.Cookie(&fiber.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Expires:  time.Now().Add(720 * time.Hour),
			HTTPOnly: true,
		})

		return c.SendString("Login successful")
	}
}

package tools

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"myPage/models"
)

func GenerateSessionID() string {
	return uuid.New().String()
}

func CheckSession(db *gorm.DB, c *fiber.Ctx) (*models.Session, error) {
	sessionID := c.Cookies("session_id")
	if sessionID == "" {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Session not found")
	}

	var session models.Session
	if err := db.First(&session, "session_id = ?", sessionID).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Session is invalid")
	}

	return &session, nil
}

func IsSessionExist(db *gorm.DB, c *fiber.Ctx) (bool, error) {
	sessionID := c.Cookies("session_id")
	if sessionID == "" {
		return false, fiber.NewError(fiber.StatusUnauthorized, "Session not found")
	}

	var session models.Session
	if err := db.First(&session, "session_id = ?", sessionID).Error; err != nil {
		return false, fiber.NewError(fiber.StatusUnauthorized, "Session is invalid")
	}

	return true, nil
}

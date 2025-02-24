package tools

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"myPage/models"
	"time"
)

func GenerateSessionID() string {
	return uuid.New().String()
}

func CreateSession(db *gorm.DB, user models.User, c *fiber.Ctx) error {
	sessionID := GenerateSessionID()
	session := models.Session{
		SessionID:     sessionID,
		UserID:        user.ID,
		UserName:      user.Username,
		Authenticated: true,
	}
	if err := db.Save(&session).Error; err != nil {
		log.Printf("Error saving session: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to save session",
		})
	}

	// Set the session cookie.
	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  time.Now().Add(720 * time.Hour),
		HTTPOnly: true,
		Secure:   true,                        //  IMPORTANT:  Set Secure to true in production (HTTPS only).
		SameSite: fiber.CookieSameSiteLaxMode, //  Good practice for CSRF protection.
	})
	return nil
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

func IsSessionExist(db *gorm.DB, c *fiber.Ctx) bool {
	session, err := CheckSession(db, c)
	if err != nil || session == nil || db == nil {
		return false
	}

	username := session.UserName
	var user models.User

	if err := db.First(&user, "username = ?", username).Error; err != nil {
		return false
	}

	return true
}

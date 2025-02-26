package handlers

import (
	"errors"
	"fmt"
	"log"
	"myPage/database"
	"myPage/models"
	"myPage/tools"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// LoginHandler handles user login.
func LoginHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := c.FormValue("username")
		password := c.FormValue("password")

		if err := tools.ValidateCredentials(username, password); err != nil {
			return tools.HandleUserError(c, fiber.StatusBadRequest, err.Error())
		}

		user, err := database.FindUserByUsername(db, username)
		if err != nil {
			if strings.Contains(err.Error(), "record not found") {
				return tools.HandleUserError(c, fiber.StatusConflict, "Username does not exist")
			}
			return tools.HandleUserError(c, fiber.StatusUnauthorized, "Invalid username or password")
		}

		if !tools.CheckPasswordHash(password, user.Password) {
			return tools.HandleUserError(c, fiber.StatusUnauthorized, "Invalid username or password")
		}

		if err := tools.CreateSession(db, user, c); err != nil {
			return err
		}

		return c.Status(fiber.StatusCreated).Redirect(tools.WebLink + "/dashboard")
	}
}

// RegisterHandler handles user registration.
func RegisterHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var user models.User
		user.Username = c.FormValue("username")
		user.Password = c.FormValue("password")

		if err := tools.ValidateCredentials(user.Username, user.Password); err != nil {
			return tools.HandleUserError(c, fiber.StatusBadRequest, err.Error())
		}

		hashedPassword, err := tools.HashPassword(user.Password)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			return tools.HandleUserError(c, fiber.StatusInternalServerError, "Unable to hash password")
		}
		user.Password = hashedPassword

		if err := database.CreateUser(db, &user); err != nil {
			return tools.HandleUserError(c, fiber.StatusConflict, err.Error())

		}

		// Create a session
		if err := tools.CreateSession(db, user, c); err != nil {
			return err
		}
		tools.CreateDirectories(user.Username)

		return c.Status(fiber.StatusCreated).Redirect(tools.WebLink + "/dashboard")
	}
}

// handleSessionRetrieval retrieves and validates a session, handling errors appropriately.
func handleSessionRetrieval(db *gorm.DB, c *fiber.Ctx) (models.Session, error) {
	sessionID := c.Cookies("session_id")
	if sessionID == "" {
		return models.Session{}, errors.New("unauthorized: Session not found")
	}

	var session models.Session
	if err := db.Where("session_id = ?", sessionID).First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Session{}, err // Return the error directly for 404
		}
		log.Printf("Error finding session: %v", err)
		return models.Session{}, fmt.Errorf("unable to find session")
	}
	return session, nil
}

// LogoutHandler handles user logout.
func LogoutHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		session, err := handleSessionRetrieval(db, c)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// If session is not found, redirect to login
				return c.Status(fiber.StatusNotFound).Redirect(tools.WebLink + "/login")
			}
			// For other errors, return 500
			return tools.HandleUserError(c, fiber.StatusInternalServerError, err.Error())
		}
		// Deleting a session
		if err := db.Where("session_id = ?", session.SessionID).Delete(&models.Session{}).Error; err != nil {
			log.Printf("Error deleting session: %v", err)
			return tools.HandleUserError(c, fiber.StatusInternalServerError, "Unable to delete session")
		}

		// Clearing session cookies
		clearSessionCookie(c)

		// Redirect to login page
		return c.Redirect(tools.WebLink + "/login")
	}
}

// clearSessionCookie clears the session cookie.
func clearSessionCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // Set the expiration date in the past
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteLaxMode,
	})
}

func RenderLoginHandler(db *gorm.DB) fiber.Handler {

	return func(c *fiber.Ctx) error {
		return tools.RenderWithSessionCheck(db, c, "login", false)
	}
}

func RenderRegisterHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return tools.RenderWithSessionCheck(db, c, "register", false)
	}
}

func RenderLogoutHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return tools.RenderWithSessionCheck(db, c, "logout", true)
	}
}

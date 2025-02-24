package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"log"
	"myPage/models"
	"myPage/tools"
	"os"
	"strings"
	"time"
)

func LoginHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := c.FormValue("username")
		password := c.FormValue("password")

		if username == "" || password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Username and password are required",
			})
		}

		var user models.User
		if err := db.First(&user, "username = ?", username).Error; err != nil {
			if strings.Contains(err.Error(), "record not found") {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{
					"error": "Username does not exist",
				})
			}
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid username or password",
			})
		}

		if tools.CheckPasswordHash(password, user.Password) == false {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid username or password",
			})
		}

		if err := tools.CreateSession(db, user, c); err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message":  "Login successful",
			"username": user.Username,
		})
	}
}

func RegisterHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var user models.User

		user.Username = c.FormValue("username")
		user.Password = c.FormValue("password")
		user.VisibleName = c.FormValue("visible_name")

		if user.Username == "" || user.Password == "" || user.VisibleName == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Username and password and visible name are required",
			})
		}

		// Validate password length.
		if len(user.Password) > 72 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Password is too long, the length must not exceed 72 characters",
			})
		}
		if len(user.Password) < 8 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Password is too short, must be at least 8 characters",
			})
		}

		// Hash the password
		hashedPassword, err := tools.HashPassword(user.Password)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Unable to hash password",
			})
		}
		user.Password = hashedPassword

		// Create the user
		if err := db.Create(&user).Error; err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{
					"error": "Username already exists",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Unable to create user",
			})
		}

		// Create a session
		if err := tools.CreateSession(db, user, c); err != nil {
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message":  "Registration successful",
			"username": user.Username,
		})
	}
}

func LogoutHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionID := c.Cookies("session_id")
		if sessionID == "" {
			// If there is no session cookie, return error 401
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: Session not found",
			})
		}

		// Trying to find a session in the database
		var session models.Session
		if err := db.Where("session_id = ?", sessionID).First(&session).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).Redirect(os.Getenv("WEB_URL") + "/login")
			} else {
				// Session search error
				log.Printf("Error finding session: %v", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Unable to find session",
				})
			}
		}

		// Deleting a session
		if err := db.Where("session_id = ?", sessionID).Delete(&models.Session{}).Error; err != nil {
			log.Printf("Error deleting session: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Unable to delete session",
			})
		}

		// Clearing session cookies
		c.Cookie(&fiber.Cookie{
			Name:     "session_id",
			Value:    "",
			Expires:  time.Now().Add(-time.Hour), // Set the expiration date in the past
			HTTPOnly: true,
			Secure:   true,
			SameSite: fiber.CookieSameSiteLaxMode,
		})

		// Redirect to login page
		return c.Redirect(os.Getenv("WEB_URL") + "/login")
	}
}

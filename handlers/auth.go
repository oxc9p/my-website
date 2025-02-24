package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"log"
	"myPage/models"
	"myPage/tools"
	"strings"
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
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid username or password",
			})
		}

		//Checking if the password is correct
		if tools.CheckPasswordHash(enteredPassword, user.Password) == false {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid username or password",
			})
		}

		// Successfully authenticated
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

		// Parse the request body
		if err := c.BodyParser(&user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}
		// Set the permission after parse
		user.Permission = 0

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
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		// Delete the session from the database
		var session models.Session
		if err := db.Where("session_id = ?", sessionID).First(&session).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Session not found, but proceed with clearing cookie
			} else {
				log.Printf("Error finding session: %v", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Unable to logout",
				})
			}
		}

		// Delete session
		if err := db.Delete(&session).Error; err != nil {
			log.Printf("Error deleting session: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal Server Error",
			})
		}

		// Clear the cookie Set expiry to the past
		c.Cookie(&fiber.Cookie{
			Name:     "session_id",
			Value:    "",
			Expires:  time.Now().Add(-time.Hour),
			HTTPOnly: true,
			Secure:   true,
			SameSite: fiber.CookieSameSiteLaxMode,
		})

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Logout successful",
		}) // Redirect to login after logout
	}
}

package tools

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
)

var WebLink = os.Getenv("WEB_URL")

// HandleUserError returns a Fiber handler that sends a JSON error response.
func HandleUserError(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"error": message,
	})
}

// ValidateCredentials checks if username and password are provided and if the password meets length requirements.
func ValidateCredentials(username, password string) error {
	if username == "" || password == "" {
		return errors.New("username and password are required")
	}
	if len(password) > 72 {
		return errors.New("password is too long, the length must not exceed 72 characters")
	}
	if len(password) < 8 {
		return errors.New("password is too short, must be at least 8 characters")
	}
	return nil
}

func CreateDirectories(username string) {
	baseDir := "userfiles/" + username
	dirs := []string{"", "/md", "/img"} //Using empty string to create base directory

	for _, dir := range dirs {
		if err := os.Mkdir(baseDir+dir, os.ModePerm); err != nil && !os.IsExist(err) {
			log.Println(err)
		}
	}
}

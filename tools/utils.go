package tools

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"log"
	"myPage/database"
	"myPage/models"
	"os"
	"path/filepath"
	"time"
)

type FileInfo struct {
	FileDir       string
	FileExtension []string
}

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

func AuthenticateAndGetUser(db *gorm.DB, c *fiber.Ctx) (*models.User, error) {
	session, err := CheckSession(db, c)
	if err != nil {
		return nil, c.Status(fiber.StatusUnauthorized).Redirect(WebLink + "/login")
	}

	user, err := database.FindUserByUsername(db, session.UserName)
	if err != nil {
		return nil, c.Status(fiber.StatusUnauthorized).Redirect(WebLink + "/login")
	}

	return &user, nil
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

func UploadFile(c *fiber.Ctx, user *models.User, formFieldName string, fileInfo FileInfo) (string, error) {
	file, err := c.FormFile(formFieldName)
	if err != nil || file == nil {
		return "", c.Status(fiber.StatusBadRequest).SendString("Error retrieving file")
	}

	found := false
	for _, e := range fileInfo.FileExtension {
		if ext := filepath.Ext(file.Filename); e == ext {
			found = true
			break
		}
	}
	if !found {
		return "", c.Status(fiber.StatusBadRequest).SendString("Invalid file format. Allowed formats: " + fmt.Sprint(fileInfo.FileExtension))
	}

	// Generate filename
	filename := fmt.Sprintf("%s%s", time.Now().Format("15:04:05:04-2006.01.02"), filepath.Ext(file.Filename))

	// Create user directory if it doesn't exist
	userDir := fmt.Sprintf("./userfiles/%s/%s/", user.Username, fileInfo.FileDir)
	if err := os.MkdirAll(userDir, os.ModePerm); err != nil {
		return "", c.Status(fiber.StatusInternalServerError).SendString("Error creating directory")
	}

	// Save file
	filePath := filepath.Join(userDir, filename)
	if err := c.SaveFile(file, filePath); err != nil {
		return "", c.Status(fiber.StatusInternalServerError).SendString("Error saving file")
	}

	return filepath.Join("userfiles", user.Username, fileInfo.FileDir, filename), nil
}

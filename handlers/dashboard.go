package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"log"
	"myPage/database"
	"myPage/models"
	"myPage/tools"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func DashboardHandler(db *gorm.DB) fiber.Handler {

	return func(c *fiber.Ctx) error {

		session, err := tools.CheckSession(db, c)
		if err != nil || session == nil || db == nil {
			return c.Status(fiber.StatusUnauthorized).Redirect(tools.WebLink + "/login")
		}

		var user models.User

		if err := db.First(&user, "username = ?", session.UserName).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).Redirect(tools.WebLink + "/login")
		}

		return tools.RenderWithSessionCheck(db, c, "dashboard", true, fiber.Map{
			"Username":    user.Username,
			"Image":       user.Image,
			"DateCreated": user.DateCreated.Format("2006-01-02"),
			"WebLink":     tools.WebLink,
		})
	}
}

// UploadImageHandler handles image uploads.
func UploadImageHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		session, err := tools.CheckSession(db, c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).Redirect(tools.WebLink + "/login")
		}

		user, err := database.FindUserByUsername(db, session.UserName)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).Redirect(tools.WebLink + "/login")
		}

		return handleImageUpload(c, db, &user)
	}
}

// handleImageUpload processes the image upload, validates it, and updates the user record.
func handleImageUpload(c *fiber.Ctx, db *gorm.DB, user *models.User) error {
	file, err := c.FormFile("image")
	if err != nil || file == nil {
		return c.Status(fiber.StatusBadRequest).SendString("Error retrieving file")
	}

	if ext := filepath.Ext(file.Filename); ext != ".jpg" && ext != ".jpeg" {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid file format. Only .jpg or .jpeg allowed")
	}

	// Generate filename
	filename := fmt.Sprintf("%d.jpg", time.Now().UnixNano())

	// Create user directory if it doesn't exist
	userDir := fmt.Sprintf("./userfiles/%s/img/", user.Username)
	if err := os.MkdirAll(userDir, os.ModePerm); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error creating directory")
	}

	// Save file
	filePath := filepath.Join(userDir, filename)
	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error saving file")
	}

	// Remove old image if it's not the default one
	if !strings.Contains(user.Image, "static/images/default.jpg") {
		if err := os.Remove(user.Image); err != nil && !os.IsNotExist(err) { // Check if error is not "file does not exist"
			log.Println("Error removing old image:", err)
		}
	}

	// Update image path in the database
	return db.Model(user).Update("image", filepath.Join("userfiles", user.Username, "img", filename)).Error
}

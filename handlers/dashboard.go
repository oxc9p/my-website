package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"myPage/database"
	"myPage/models"
	"myPage/tools"
	"os"
	"strings"
)

func DashboardHandler(db *gorm.DB) fiber.Handler {

	return func(c *fiber.Ctx) error {

		user, err := tools.AuthenticateAndGetUser(db, c)
		if err != nil {
			log.Println("Error getting user:", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		md, e := os.ReadDir("userfiles" + "/" + user.Username + "/md")
		if e != nil {
			log.Println("Error reading directory:", e)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		var files []string
		for _, file := range md {
			files = append(files, file.Name())
		}

		isWriter := false
		if user.Permission > 0 {
			isWriter = true
		}

		isAdmin := false
		if user.Permission > 1 {
			isAdmin = true
		}

		return tools.RenderWithSessionCheck(db, c, "dashboard", true, fiber.Map{
			"Username":      user.Username,
			"Image":         user.Image,
			"DateCreated":   user.DateCreated.Format("2006-01-02"),
			"WebLink":       tools.WebLink,
			"IsWriter":      isWriter,
			"IsAdmin":       isAdmin,
			"MarkdownFiles": files,
		})
	}
}

func UploadAvatarHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := tools.AuthenticateAndGetUser(db, c)
		if err != nil {
			return err
		}

		return handleAvatarUpload(c, db, user)
	}
}

func UploadImageHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := tools.AuthenticateAndGetUser(db, c)
		if err != nil {
			return err
		}

		return handleImageUpload(c, user)
	}
}

func UploadMdHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := tools.AuthenticateAndGetUser(db, c)
		if err != nil {
			return err
		}

		return handleMdUpload(c, user)
	}
}

func DeleteMdHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := tools.AuthenticateAndGetUser(db, c)
		if err != nil {
			return err
		}
		return handleMdDelete(c, user)
	}
}

func DeleteImageHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := tools.AuthenticateAndGetUser(db, c)
		if err != nil {
			return err
		}

		return handleImageDelete(c, user)
	}
}

func handleAvatarUpload(c *fiber.Ctx, db *gorm.DB, user *models.User) error {
	filePath, err := tools.UploadFile(c, user, "avatar", tools.FileInfo{FileDir: "img", FileExtension: []string{".jpg", ".jpeg", ".png"}})
	if err != nil {
		return err
	}

	// Remove old image if it's not the default one
	if !strings.Contains(user.Image, "static/images/default.jpg") {
		if err := os.Remove(user.Image); err != nil && !os.IsNotExist(err) { // Check if error is not "file does not exist"
			log.Println("Error removing old image:", err)
		}
	}

	// Update image path in the database
	return db.Model(user).Update("image", filePath).Error
}

func handleImageUpload(c *fiber.Ctx, user *models.User) error {
	_, err := tools.UploadFile(c, user, "image", tools.FileInfo{FileDir: "img", FileExtension: []string{".jpg", ".jpeg", ".png"}})
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).SendString("Image uploaded successfully")
}

func handleMdUpload(c *fiber.Ctx, user *models.User) error {
	_, err := tools.UploadFile(c, user, "md", tools.FileInfo{FileDir: "md", FileExtension: []string{".md"}})
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).SendString("Markdown uploaded successfully")
}

func handleMdDelete(c *fiber.Ctx, user *models.User) error {
	if err := tools.DeleteFile(c, user, "md"); err != nil {
		log.Println(err)
		return err
	}
	return c.Status(fiber.StatusOK).SendString("Markdown removed successfully")
}

func handleImageDelete(c *fiber.Ctx, user *models.User) error {
	if err := tools.DeleteFile(c, user, "img"); err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).SendString("Image removed successfully")
}

func UploadArticleHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := tools.AuthenticateAndGetUser(db, c)
		if err != nil {
			return err
		}
		if user.Permission == 0 {
			return c.Status(fiber.StatusUnauthorized).Redirect(tools.WebLink + "/dashboard")
		}
		article := models.Article{
			Image:       c.FormValue("image"),
			Title:       c.FormValue("title"),
			Description: c.FormValue("description"),
			Link:        c.FormValue("link"),
			Author:      user.Username,
			Id:          uuid.NewString(),
		}
		database.Create(db, &article)
		return c.Status(fiber.StatusOK).SendString("Article added successfully")
	}
}

func UploadProjectHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := tools.AuthenticateAndGetUser(db, c)
		if err != nil {
			return err
		}
		if user.Permission != 2 {
			return c.Status(fiber.StatusUnauthorized).Redirect(tools.WebLink + "/dashboard")
		}
		project := models.Project{
			Image:       c.FormValue("image"),
			Title:       c.FormValue("title"),
			Description: c.FormValue("description"),
			Label:       c.FormValue("label"),
			Link:        c.FormValue("link"),
		}
		database.Create(db, &project)
		return c.Status(fiber.StatusOK).SendString("Project added successfully")
	}
}

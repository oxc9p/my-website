package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"myPage/models"
	"myPage/tools"
	"os"
)

func RenderLoginHandler(db *gorm.DB) fiber.Handler {

	return func(c *fiber.Ctx) error {
		session, err := tools.CheckSession(db, c)
		if err != nil {
			return c.Render("login", fiber.Map{
				"WebLink": os.Getenv("WEB_URL"),
			})
		}

		username := session.UserName
		var user models.User
		if err := db.First(&user, "username = ?", username).Error; err != nil {
			return c.Render("login", fiber.Map{
				"WebLink": os.Getenv("WEB_URL"),
			})
		}

		return c.Redirect(os.Getenv("WEB_URL"))
	}
}

func RenderRegisterHandler(db *gorm.DB) fiber.Handler {

	return func(c *fiber.Ctx) error {
		session, err := tools.CheckSession(db, c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		username := session.UserName
		var user models.User
		if err := db.First(&user, "username = ?", username).Error; err != nil {
			return c.Render("register", fiber.Map{
				"WebLink": os.Getenv("WEB_URL"),
			})
		}

		return c.Redirect(os.Getenv("WEB_URL"))
	}
}

func RenderLogoutHandler(db *gorm.DB) fiber.Handler {

	return func(c *fiber.Ctx) error {
		session, err := tools.CheckSession(db, c)
		if err != nil {
			return c.Redirect(os.Getenv("WEB_URL"))
		}

		username := session.UserName
		var user models.User
		if err := db.First(&user, "username = ?", username).Error; err == nil {
			return c.Render("logout", fiber.Map{
				"WebLink": os.Getenv("WEB_URL"),
			})
		}

		return c.Redirect(os.Getenv("WEB_URL"))
	}
}

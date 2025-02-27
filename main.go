package main

import (
	"github.com/gofiber/fiber/v2"
	tmpl "github.com/gofiber/template/html/v2"
	"log"
	"myPage/database"
	"myPage/handlers"
	"os"
)

func main() {
	db := database.Init()
	engine := tmpl.New("./templates", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/static", "./templates/static")
	app.Static("/userfiles", "./userfiles")

	// Creating api
	api := app.Group("/api")
	api.Post("/login", handlers.LoginHandler(db))
	api.Post("/register", handlers.RegisterHandler(db))
	api.Post("/logout", handlers.LogoutHandler(db))
	upload := api.Group("/upload")
	upload.Post("/article", handlers.UploadArticleHandler(db))
	upload.Post("/avatar", handlers.UploadAvatarHandler(db))
	upload.Post("/image", handlers.UploadImageHandler(db))
	upload.Post("/md", handlers.UploadMdHandler(db))
	upload.Post("/project", handlers.UploadProjectHandler(db))

	// Connecting handlers using the get method
	app.Get("/", handlers.IndexHandler(db))
	app.Get("/markdown/:filename", handlers.MarkdownHandler())
	app.Get("/blog", handlers.ArticleHandler(db))
	app.Get("/dashboard", handlers.DashboardHandler(db))
	app.Get("/login", handlers.RenderLoginHandler(db))
	app.Get("/register", handlers.RenderRegisterHandler(db))
	app.Get("/logout", handlers.RenderLogoutHandler(db))
	app.Get("/about", handlers.AboutHandler(db))

	users := app.Group("/users")
	users.Get("/:username/:filename", handlers.UserMarkdownHandler())

	// Run webpage
	log.Fatal(app.Listen(os.Getenv("IP") + ":8080"))
}

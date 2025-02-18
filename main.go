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

	// Connect handlers using method get
	app.Get("/", handlers.IndexHandler())
	app.Get("/markdown/:filename", handlers.MarkdownHandler())
	app.Get("/blog", handlers.ArticleHandler(db))

	log.Fatal(app.Listen(os.Getenv("IP") + ":80"))
}

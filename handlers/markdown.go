package handlers

import (
	"github.com/gofiber/fiber/v2"
	"myPage/tools"
)

func MarkdownHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return tools.RenderMd(c, "markdown/"+c.Params("filename")+".md")
	}
}

package handlers

import (
	"github.com/gofiber/fiber/v2"
	"myPage/tools"
)

func UserMarkdownHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return tools.RenderMd(c, "userfiles/"+c.Params("username")+"/md/"+c.Params("filename"))
	}
}

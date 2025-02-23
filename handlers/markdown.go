package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"myPage/tools"
)

func MarkdownHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		byteArrayChunked, err := tools.ParseFileToByteArray("markdown/" + c.Params("filename") + ".md")
		if err != nil {
			fmt.Println(err.Error())
			return c.Status(500).SendString("Error loading file")
		}
		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.Send([]byte(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Markdown</title>
			<link rel="stylesheet" href="/static/style.css">
		</head>
		<body>
			<div class="container">` + string(tools.MdToHTML(byteArrayChunked)) + `</div></body>
		</html>`))
	}
}

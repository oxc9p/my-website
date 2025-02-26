package tools

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"io/ioutil"
	"os"
)

// ParseFileToByteArray reads a file and returns its contents as a byte array.
func ParseFileToByteArray(filePath string) ([]byte, error) {
	fileContents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return fileContents, nil
}

// ParseFileToByteArrayChunked is an alternative implementation that reads the
// file in chunks, suitable for very large files.
func ParseFileToByteArrayChunked(filePath string, chunkSize int) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	var byteArray []byte
	buffer := make([]byte, chunkSize)

	for {
		n, err := file.Read(buffer)
		if err != nil {
			if errors.Is(err, os.ErrClosed) {
				return nil, fmt.Errorf("error reading file (file closed): %w", err)
			}
			if err.Error() == "EOF" {
				break
			}
			return nil, fmt.Errorf("error reading file chunk: %w", err)
		}
		byteArray = append(byteArray, buffer[:n]...)
	}

	return byteArray, nil
}

func MdToHTML(md []byte) []byte {

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func RenderMd(c *fiber.Ctx, path string) error {
	byteArrayChunked, err := ParseFileToByteArray(path)
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
			<div class="container">` + string(MdToHTML(byteArrayChunked)) + `</div></body>
		</html>`))
}

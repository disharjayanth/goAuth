package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

func indexHandler(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title":   "Go fiber app",
		"Body":    "Hello world!",
		"Message": "Today is a very good day!",
		"Info":    "This page is built with Tailwindcss and backend in Golang",
	})
}

func main() {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", indexHandler)

	log.Fatal(app.Listen(":8000"))
}

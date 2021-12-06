package main

import (
	"github.com/disharjayanth/goAuth/11_exercises/02_jwt/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

func main() {
	engine := html.New("views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", handlers.SignUpPageHandler)
	app.Get("/login", handlers.SignInPageHandler)

	app.Post("/signup", handlers.SignUpPostHandler)
	app.Post("/signin", handlers.SignInPostHandler)

	app.Get("/success", handlers.SuccessPageHandler)
	app.Get("/logout", handlers.LogOutHandler)

	app.Listen(":8000")
}

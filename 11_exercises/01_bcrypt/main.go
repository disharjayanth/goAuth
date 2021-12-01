package main

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"golang.org/x/crypto/bcrypt"
)

type UserDB map[string]string

var users = UserDB{}

func formHandler(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "User Details",
	})
}

func registerHandler(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "" || password == "" {
		return c.Redirect("/", http.StatusSeeOther)
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error while generating password: %w", err)
	}

	users[username] = string(hashPassword)

	for k, v := range users {
		fmt.Println(k, ":", v)
	}

	return c.Redirect("/", http.StatusSeeOther)
}

func main() {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", formHandler)
	app.Post("/register", registerHandler)

	app.Listen(":8000")
}

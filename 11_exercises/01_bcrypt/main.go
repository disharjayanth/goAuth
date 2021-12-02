package main

import (
	"fmt"
	"net/http"

	"github.com/disharjayanth/goAuth/11_exercises/01_bcrypt/usersdb"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"golang.org/x/crypto/bcrypt"
)

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

	user := &usersdb.User{
		Name:     username,
		Password: string(hashPassword),
	}

	trueOrFalse, err := user.Store()
	if err != nil {
		return c.SendString(err.Error())
	}

	if trueOrFalse {
		return c.Redirect("/", http.StatusSeeOther)
	}

	return c.SendString("error while storing user datas")
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

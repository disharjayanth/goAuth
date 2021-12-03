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
		"Title":           "User Sign Up",
		"FormName":        "Sign Up",
		"Link":            "/login",
		"LinkTitle":       "Already have an account? Sign In",
		"FormLinkHandler": "/register",
	})
}

func loginPageHandler(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title":           "User Login",
		"FormName":        "Sign In",
		"Link":            "/",
		"LinkTitle":       "Don't have an account? Sign Up",
		"FormLinkHandler": "/login",
	})
}

func loginHandler(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "" || password == "" {
		return c.Redirect("/login", http.StatusSeeOther)
	}

	user, err := usersdb.RetrieveUser(username)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Printf("error while comparing password: %w", err)
		return c.Render("index", fiber.Map{
			"Title":           "Login Failed",
			"FormName":        "Login Failed",
			"FormLinkHandler": "/login",
		})
	}

	if user.Name == username && err == nil {
		return c.Render("index", fiber.Map{
			"Title":           "Login successfull",
			"FormName":        "Login successfull",
			"FormLinkHandler": "/login",
		})
	}

	fmt.Println(username, password)
	fmt.Println(user.Name, string(user.Password))

	return nil
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
	app.Get("/login", loginPageHandler)
	app.Post("/login", loginHandler)

	app.Listen(":8000")
}

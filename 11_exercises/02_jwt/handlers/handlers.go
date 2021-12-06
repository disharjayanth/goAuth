package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/disharjayanth/goAuth/11_exercises/02_jwt/utilities"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func SignUpPageHandler(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title":          "Sign Up",
		"FormName":       "Sign Up",
		"FormSubmitLink": "/signup",
		"Question":       "Already have an account? Sign In!",
		"QuestionLink":   "/login",
	})
}

func SignInPageHandler(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title":          "Sign In",
		"FormName":       "Sign In",
		"FormSubmitLink": "/signin",
		"Question":       "Don't have an account? Sign Up!",
		"QuestionLink":   "/",
	})
}

func SignUpPostHandler(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "" || password == "" {
		return c.Redirect("/", http.StatusSeeOther)
	}

	sid := uuid.NewString()
	token, err := utilities.CreateToken(sid)
	if err != nil {
		return fmt.Errorf("error while creating token: %w", err)
	}

	c1 := fiber.Cookie{
		Name:    "session",
		Value:   token,
		Expires: time.Now().Add(1 * time.Minute),
	}

	c.Cookie(&c1)

	return c.Redirect("/success", http.StatusSeeOther)
}

func SignInPostHandler(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "" || password == "" {
		return c.Redirect("/login", http.StatusSeeOther)
	}

	return nil
}

func SuccessPageHandler(c *fiber.Ctx) error {
	cookie := c.Cookies("session")

	if cookie == "" {
		return errors.New("cookie is empty")
	}

	sid, err := utilities.ParseToken(cookie)
	if err != nil {
		return err
	}

	return c.Render("success", fiber.Map{
		"sid": sid,
	})
}

func LogOutHandler(c *fiber.Ctx) error {
	c.ClearCookie("session")
	return c.SendString("Successfully logged out!")
}

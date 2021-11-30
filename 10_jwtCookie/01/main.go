package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/golang-jwt/jwt"
)

var privateKey string = "secret string!"

type CustomClaims struct {
	jwt.StandardClaims
	Email string
}

func getJWT(email string) (string, error) {
	claims := &CustomClaims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	sToken, err := token.SignedString([]byte(privateKey))
	if err != nil {
		return "", fmt.Errorf("could'nt sign jwt token: %w", err)
	}

	return sToken, nil
}

func indexHandler(c *fiber.Ctx) error {
	cookies := c.Cookies("session")
	if cookies == "" {
		cookie := &fiber.Cookie{
			Name:    "session",
			Expires: time.Now().Add(1 * time.Minute),
		}
		cookies = cookie.Value
	}

	return c.Render("index", fiber.Map{
		"Title":   "Go fiber app",
		"Body":    "Hello world!",
		"Message": "Today is a very good day!",
		"Info":    "This page is built with Tailwindcss and backend in Golang",
		"Cookie":  cookies,
	})
}

func submitHandler(c *fiber.Ctx) error {
	email := c.FormValue("email")
	if email == "" {
		c.Redirect("/", http.StatusSeeOther)
	}

	jwt, err := getJWT(email)
	if err != nil {
		return c.SendString(err.Error())
	}

	cookie := &fiber.Cookie{
		Name:  "session",
		Value: jwt,
	}

	c.Cookie(cookie)

	return c.Redirect("/", http.StatusSeeOther)
}

func main() {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", indexHandler)
	app.Post("/submit", submitHandler)

	log.Fatal(app.Listen(":8000"))
}

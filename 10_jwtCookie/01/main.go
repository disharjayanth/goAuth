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
			ExpiresAt: time.Now().Add(1 * time.Minute).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	sToken, err := token.SignedString([]byte(privateKey))
	if err != nil {
		return "", fmt.Errorf("could'nt sign jwt token: %w", err)
	}

	return sToken, nil
}

func parseJWT(token string) (string, string, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		// optional to check signing algorithm!
		if t.Method.Alg() != jwt.SigningMethodHS256.Name {
			return nil, fmt.Errorf("algorithms are not matching")
		}
		return []byte(privateKey), nil
	})

	message := "Not signed in"
	if err != nil {
		return "", message, fmt.Errorf("error while parsing jwt token: %w", err)
	}

	if claims, ok := parsedToken.Claims.(*CustomClaims); ok && parsedToken.Valid {
		message = "Signed in"
		return claims.Email, message, nil
	} else {
		return "", message, fmt.Errorf("invalid token")
	}
}

func indexHandler(c *fiber.Ctx) error {
	cookies := c.Cookies("session")

	var message string
	var email string
	var err error

	message = "Not signed in"
	if cookies != "" {
		email, message, err = parseJWT(cookies)
		if err != nil {
			fmt.Println(err)
		}
	}

	return c.Render("index", fiber.Map{
		"Title":         "Go fiber app",
		"Body":          "Hello world!",
		"Message":       "Today is a very good day!",
		"Info":          "This page is built with Tailwindcss and backend in Golang",
		"Cookie":        cookies,
		"SignedInOrNot": message,
	})
}

func submitHandler(c *fiber.Ctx) error {
	email := c.FormValue("email")
	if email == "" {
		return c.Redirect("/", http.StatusSeeOther)
	}

	jwt, err := getJWT(email)
	if err != nil {
		return c.SendString(err.Error())
	}

	cookie := &fiber.Cookie{
		Name:    "session",
		Value:   jwt,
		Expires: time.Now().Add(1 * time.Minute),
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

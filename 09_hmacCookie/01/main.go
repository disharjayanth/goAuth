package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

func createHmacHash(msg string) string {
	h := hmac.New(sha256.New, []byte("secret string!"))
	io.WriteString(h, msg)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func htmlHandler(c *fiber.Ctx) error {
	cookie := c.Cookies("session")

	sb := strings.SplitN(cookie, "|", 2)

	isEqual := true
	if len(sb) == 2 {
		hash := sb[0]
		email := sb[1]

		hmacHash := createHmacHash(email)

		isEqual = hmac.Equal([]byte(hash), []byte(hmacHash))
	}

	message := "Not logged in"

	if isEqual {
		message = "Logged in"
	}

	return c.Render("index", fiber.Map{
		"cookie":             cookie,
		"loggedInOutMessage": message,
	})
}

func cookieHandler(c *fiber.Ctx) error {
	if c.Method() != http.MethodPost {
		c.Redirect("/", http.StatusSeeOther)
		return nil
	}

	email := c.FormValue("email")
	if email == "" {
		c.Redirect("/", http.StatusSeeOther)
		return nil
	}

	// hash | what we store
	hash := createHmacHash(email)

	cookie := &fiber.Cookie{
		Name:  "session",
		Value: hash + "|" + email,
	}

	c.Cookie(cookie)

	return c.Redirect("/", http.StatusSeeOther)
}

func main() {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", htmlHandler)
	app.Post("/submit", cookieHandler)

	log.Fatal(app.Listen(":8000"))
}

package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var githubOAuthConfig = &oauth2.Config{}

func intialPageHandler(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "OAuth2 with Github.",
	})
}

func startOAuthProcessHandler(fiberC *fiber.Ctx) error {
	redirectURL := githubOAuthConfig.AuthCodeURL("0000")
	fmt.Println("clientID:", githubOAuthConfig.ClientID, "secret:", githubOAuthConfig.ClientSecret)
	fmt.Println("clientID", os.Getenv("ClientID"), "secret:", os.Getenv("ClientSecret"))
	fmt.Println(redirectURL)
	return fiberC.Redirect(redirectURL)
}

func init() {
	godotenv.Load()
	githubOAuthConfig = &oauth2.Config{
		ClientID:     os.Getenv("ClientID"),
		ClientSecret: os.Getenv("ClientSecret"),
		Endpoint:     github.Endpoint,
	}
}

func callBackCodeHandler(c *fiber.Ctx) error {
	code := c.Query("code")
	fmt.Println("code:", code)
	return nil
}

func main() {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", intialPageHandler)
	app.Post("/oauth/github", startOAuthProcessHandler)
	app.Get("/oauth2/receive", callBackCodeHandler)

	app.Listen(":8000")
}

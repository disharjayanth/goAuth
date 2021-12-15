package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOAuthConfig = &oauth2.Config{}

func intialPageHandler(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "Google OAuth",
	})
}

func startOAuthProcessHandler(c *fiber.Ctx) error {
	redirectURL := googleOAuthConfig.AuthCodeURL("0000")
	fmt.Println("redirectURL:", redirectURL)
	return c.Redirect(redirectURL)
}

func callBackCodeHandler(c *fiber.Ctx) error {
	code := c.Query("code")
	// state := c.Query("state")

	token, err := googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return err
	}

	tokenSource := googleOAuthConfig.TokenSource(c.Context(), token)

	client := oauth2.NewClient(c.Context(), tokenSource)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v1/userinfo?alt=json")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	sb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return c.SendString(string(sb))
}

func init() {
	godotenv.Load()
	googleOAuthConfig = &oauth2.Config{
		ClientID:     os.Getenv("ClientID"),
		ClientSecret: os.Getenv("ClientSecret"),
		Endpoint:     google.Endpoint,
		Scopes:       []string{"email", "profile"},
		RedirectURL:  os.Getenv("RedirectURL"),
	}
}

func main() {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", intialPageHandler)

	app.Post("/oauth/google", startOAuthProcessHandler)

	app.Get("/authorized", callBackCodeHandler)

	app.Listen(":8000")
}

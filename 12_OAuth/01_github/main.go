package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// Key is githubID, value is userID
var gitHubConnections = make(map[string]string)

var githubOAuthConfig = &oauth2.Config{}

// JSON Layout: {"data":{"viewer":{"id":".."}}}
type GitHubResponse struct {
	Data struct {
		Viewer struct {
			ID string `json:"id"`
		} `json:"viewer"`
	} `json:"data"`
}

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
	state := c.Query("state")
	fmt.Println("code:", code, "state:", state)

	token, err := githubOAuthConfig.Exchange(c.Context(), code)
	if err != nil {
		return fmt.Errorf("Couldn't login: %w", err)
	}

	tokenSource := githubOAuthConfig.TokenSource(c.Context(), token)

	client := oauth2.NewClient(c.Context(), tokenSource)

	const query = `{ 
		"query": "query { viewer { id }}" 
	}`

	queryReader := strings.NewReader(query)
	resp, err := client.Post("https://api.github.com/graphql", "application/json", queryReader)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var githubResponse GitHubResponse
	err = json.NewDecoder(resp.Body).Decode(&githubResponse)
	if err != nil {
		return err
	}

	gitHubID := githubResponse.Data.Viewer.ID
	userID, ok := gitHubConnections[gitHubID]
	if !ok {
		// new user create account
		// depends on project
		gitHubConnections[gitHubID] = uuid.NewString()
		return c.SendString(gitHubConnections[gitHubID])
	}

	fmt.Println("userID:", githubResponse)
	// Login user and send JWT

	return c.SendString(userID)
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

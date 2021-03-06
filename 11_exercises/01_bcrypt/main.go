package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/disharjayanth/goAuth/11_exercises/01_bcrypt/sessiondb"
	"github.com/disharjayanth/goAuth/11_exercises/01_bcrypt/usersdb"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var privateKey []byte = []byte{239, 25, 78, 56, 68, 194, 108, 94, 228, 87, 231, 160, 160, 112, 184, 189, 189, 97, 77, 74, 43, 241, 248, 184, 205, 97, 127, 233, 197, 17, 241, 232, 99, 195, 116, 162, 3, 30, 6, 91, 103, 238, 131, 206, 240, 41, 30, 216, 115, 96, 239, 123, 254, 167, 60, 102, 206, 96, 144, 120, 137, 133, 13, 127}

func formHandler(c *fiber.Ctx) error {
	cookie := c.Cookies("session")
	username := ""
	if cookie != "" {
		sid, err := parseHMACToken(cookie)
		if err != nil {
			return c.SendString("error while parsing token from cookie")
		}

		username, err = sessiondb.Get(sid)
		if err != nil {
			return c.SendString("error while getting user's session id")
		}

		return c.Render("index", fiber.Map{
			"Title":           "User Sign Up",
			"FormName":        "Sign Up",
			"Link":            "/login",
			"LinkTitle":       "Already have an account? Sign In",
			"FormLinkHandler": "/register",
			"Username":        username,
		})
	}
	return c.Render("index", fiber.Map{
		"Title":           "User Sign Up",
		"FormName":        "Sign Up",
		"Link":            "/login",
		"LinkTitle":       "Already have an account? Sign In",
		"FormLinkHandler": "/register",
		"Username":        username,
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
		return c.Render("index", fiber.Map{
			"Title":           "Login Failed",
			"FormName":        "Login Failed",
			"FormLinkHandler": "/login",
			"LinkTitle":       "Don't have an account? Sign Up",
		})
	}

	// Store user's session id to sessiondb
	session := sessiondb.Session{
		Name: username,
		SID:  uuid.NewString(),
	}

	success, err := session.Store()
	if err != nil {
		return c.SendString("error while storing user's session id" + err.Error())
	}

	token, err := createHMACToken(session.SID)
	if err != nil {
		return c.SendString("error while create HMAC token")
	}

	cookie := fiber.Cookie{
		Name:    "session",
		Value:   token,
		Expires: time.Now().Add(1 * time.Minute),
	}

	c.Cookie(&cookie)

	if user.Name == username && err == nil && success {
		return c.Redirect("/", http.StatusSeeOther)
	}

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

	// Store user to userdb
	user := &usersdb.User{
		Name:     username,
		Password: string(hashPassword),
	}

	trueOrFalse, err := user.Store()
	if err != nil {
		return c.SendString(err.Error())
	}

	// Store user's session id to sessiondb
	session := sessiondb.Session{
		Name: username,
		SID:  uuid.NewString(),
	}

	success, err := session.Store()
	if err != nil {
		return c.SendString("error while storing user's session id")
	}

	token, err := createHMACToken(session.SID)
	if err != nil {
		return c.SendString("error while create HMAC token")
	}

	cookie := fiber.Cookie{
		Name:    "session",
		Value:   token,
		Expires: time.Now().Add(1 * time.Minute),
	}

	c.Cookie(&cookie)

	if trueOrFalse && success {
		return c.Redirect("/", http.StatusSeeOther)
	}

	return c.SendString("error while storing user datas")
}

func createHMACToken(sid string) (string, error) {
	h := hmac.New(sha512.New, privateKey)
	if _, err := h.Write([]byte(sid)); err != nil {
		return "", fmt.Errorf("error in createHMACToken function while writing sid to hash: %w", err)
	}

	// hex representation
	// signedHMACHex := fmt.Sprintf("%x", h.Sum(nil))

	// base64 representation
	signedHMACBase64 := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return signedHMACBase64 + "|" + sid, nil
}

func parseHMACToken(ss string) (string, error) {
	xs := strings.SplitN(ss, "|", 2)
	token := xs[0]
	sid := xs[1]

	sb, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return "", fmt.Errorf("error in parseHMACToken function while decoding base64 to hmac: %w", err)
	}

	h := hmac.New(sha512.New, privateKey)
	if _, err := h.Write([]byte(sid)); err != nil {
		return "", fmt.Errorf("error in parseHMACToken function while writing sid to hash: %w", err)
	}

	if hmac.Equal(sb, h.Sum(nil)) {
		return sid, nil
	} else {
		return "", fmt.Errorf("error in parseHMACToken function and hmacs are'nt equal: %w", err)
	}
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

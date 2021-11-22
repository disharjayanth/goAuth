package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

// private key
var signKey []byte

// public key
var verifyKey []byte

type UserClaims struct {
	jwt.StandardClaims
	SessionID int64
}

func (u *UserClaims) Valid() error {
	if !u.VerifyExpiresAt(time.Now().Unix(), true) {
		return fmt.Errorf("Token has expired")
	}

	if u.SessionID == 0 {
		return fmt.Errorf("Invalid session ID")
	}

	return nil
}

func createToken(user string) (string, error) {
	t := jwt.New(jwt.SigningMethodHS256)

	t.Claims = &UserClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Minute).Unix(),
		},
		1234,
	}

	return t.SignedString(signKey)
}

func verifyToken(token string) {
	t, err := jwt.ParseWithClaims(token, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	if err != nil {
		fmt.Println("Error while verifying token:", err)
		return
	}

	claims := t.Claims.(*UserClaims)
	fmt.Println(claims.SessionID)
}

func main() {
	token, err := createToken("hello world")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println("Token:", token)

	verifyToken(token)
}

package utilities

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

var privateKey string = "secret string!"

type CustomClaims struct {
	jwt.StandardClaims
	Sid string
}

func CreateToken(sid string) (string, error) {
	claims := CustomClaims{
		Sid: sid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &claims)
	sToken, err := token.SignedString([]byte(privateKey))
	if err != nil {
		return "", fmt.Errorf("error while signing jwt token: %w", err)
	}

	return sToken, nil
}

func ParseToken(token string) (string, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		// optional to check if algorithms are matching
		if t.Method.Alg() != jwt.SigningMethodHS512.Name {
			return nil, fmt.Errorf("algorithms are'nt matching")
		}

		return []byte(privateKey), nil
	})

	if err != nil {
		return "", fmt.Errorf("error while parsing token: %w", err)
	}

	if claims, ok := parsedToken.Claims.(*CustomClaims); ok && parsedToken.Valid {
		fmt.Println("parsedToken:", claims.Sid, ok)
		return claims.Sid, nil
	} else {
		return "", fmt.Errorf("invalid token")
	}
}

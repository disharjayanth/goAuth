package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "123456789"
	hashedPass, err := hashpassword(password)
	if err != nil {
		panic(err)
	}

	if err = comparePassword(password, hashedPass); err != nil {
		fmt.Println("Not logged in!")
		return
	}

	fmt.Println("Login successful")
}

func hashpassword(password string) ([]byte, error) {
	sb, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error while creating hash for given password: %w", err)
	}
	return sb, nil
}

func comparePassword(password string, hashedPass []byte) error {
	if err := bcrypt.CompareHashAndPassword(hashedPass, []byte(password)); err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}
	return nil
}

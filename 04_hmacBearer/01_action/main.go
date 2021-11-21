package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"errors"
	"fmt"
	"log"
)

// private key used to sign cryptographic key also used to validate other key
var key []byte

func signMessage(msg []byte) ([]byte, error) {
	h := hmac.New(sha512.New, key)
	if _, err := h.Write(msg); err != nil {
		return nil, fmt.Errorf("error in signedMessage function while writing msg: %w", err)
	}

	signature := h.Sum(nil)
	return signature, nil
}

func checkSignature(msg, sign []byte) (bool, error) {
	newSign, err := signMessage(msg)
	if err != nil {
		return false, fmt.Errorf("error in checkSignature func while getting signature of message: %w", err)
	}

	if hmac.Equal(newSign, sign) {
		return true, nil
	} else {
		return false, errors.New("signature doesnt match")
	}
}

func main() {
	for i := 1; i <= 64; i++ {
		key = append(key, byte(i))
	}

	signature, err := signMessage([]byte("hello world"))
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(string(signature))

	pass, err := checkSignature([]byte("hello worl"), signature)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(pass)
}

package main

import (
	"encoding/base64"
	"fmt"
)

func main() {
	msg := "Hello world!"
	fmt.Println("message:", msg)
	encodedString := encode(msg)
	fmt.Println("Encoded message", encodedString)

	decodedString, err := decode(encodedString)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("Decoded message:", decodedString)
}

func encode(msg string) string {
	encodedString := base64.URLEncoding.EncodeToString([]byte(msg))
	return encodedString
}

func decode(encodedString string) (string, error) {
	sb, err := base64.URLEncoding.DecodeString(encodedString)
	if err != nil {
		return "", fmt.Errorf("error while decoding encoded string:", err)
	}
	return string(sb), nil
}

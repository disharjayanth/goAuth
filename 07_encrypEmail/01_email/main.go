package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func encodeAndDecode(key []byte, msg string) ([]byte, error) {
	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("couldn't new cipher: %w", err)
	}

	sb := make([]byte, aes.BlockSize)

	cipherStream := cipher.NewCTR(cipherBlock, sb)

	buff := &bytes.Buffer{}

	sw := cipher.StreamWriter{
		S: cipherStream,
		W: buff,
	}

	if _, err := sw.Write([]byte(msg)); err != nil {
		return nil, fmt.Errorf("error while writing message to stream writer buffer: %w", err)
	}

	output := buff.Bytes()

	return output, nil
}

func main() {
	msg := "hello world!"
	password := "1234"

	sb, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		fmt.Println("error creating password:", err)
		panic(err)
	}

	sb = sb[:16]

	result, err := encodeAndDecode(sb, msg)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println(string(result))

	result2, err := encodeAndDecode(sb, string(result))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println(string(result2))
}

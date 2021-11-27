package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// writer can be buffer or response writer
func encryptEncodeDecode(wtr io.Writer, key []byte) (io.Writer, error) {
	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("couldn't new cipher: %w", err)
	}

	// intialize vector
	iv := make([]byte, aes.BlockSize)

	cipherStream := cipher.NewCTR(cipherBlock, iv)

	return cipher.StreamWriter{
		S: cipherStream,
		W: wtr,
	}, nil
}

func main() {
	msg := "hello world!"
	password := "1234"

	sb, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalln("error while generating password:", err)
	}
	sb = sb[:16]

	encWtr := &bytes.Buffer{}

	encWriter, err := encryptEncodeDecode(encWtr, sb)
	if err != nil {
		log.Fatalln(err)
	}

	io.WriteString(encWriter, msg)

	encrypted := encWtr.String()
	fmt.Println("encrypted:", encrypted)

	decWriter := &bytes.Buffer{}

	decryptedWriter, err := encryptEncodeDecode(decWriter, sb)
	if err != nil {
		log.Fatalln(err)
	}

	io.WriteString(decryptedWriter, encrypted)

	decrypted := decWriter.String()

	fmt.Println("decrypted:", decrypted)
}

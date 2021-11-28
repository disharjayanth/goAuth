package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	f, err := os.Open("08_sha256File/info.txt")
	if err != nil {
		log.Fatalln("Failed to open info.txt file:", err)
	}
	defer f.Close()

	h := sha256.New()

	io.Copy(h, f)

	fmt.Println("sha256:", string(h.Sum(nil)))
	fmt.Println("base64:", base64.StdEncoding.EncodeToString(h.Sum(nil)))
	fmt.Println("hex:", hex.EncodeToString(h.Sum(nil)))
	fmt.Println("hex:", hex.EncodeToString(h.Sum([]byte{1, 2})))
}

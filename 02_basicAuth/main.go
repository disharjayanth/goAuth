package main

import (
	"encoding/base64"
	"fmt"
)

func main() {
	// Request header container Authoriztion user:pass
	fmt.Println(base64.StdEncoding.EncodeToString([]byte("user:pass")))
}

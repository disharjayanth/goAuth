package main

import (
	"encoding/base64"
	"fmt"
)

func main() {
	// Request header container Authoriztion user:pass
	// Format is Authorization username:password
	fmt.Println(base64.StdEncoding.EncodeToString([]byte("user:pass")))
}

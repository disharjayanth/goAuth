package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type person struct {
	First string `json:"first"`
}

func encode(w http.ResponseWriter, r *http.Request) {
	p1 := person{
		First: "James",
	}

	if err := json.NewEncoder(w).Encode(&p1); err != nil {
		fmt.Println("Error encoding to json:", err)
		return
	}
}

func decode(w http.ResponseWriter, r *http.Request) {
	p1 := person{}

	if err := json.NewDecoder(r.Body).Decode(&p1); err != nil {
		fmt.Println("Error decoding from json:", err)
		return
	}

	fmt.Printf("%+v\n", p1)
}

func main() {
	http.HandleFunc("/encode", encode)
	http.HandleFunc("/decode", decode)

	log.Println("Server listening at port@3000")
	http.ListenAndServe(":3000", nil)
}

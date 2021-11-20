package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type person struct {
	First string `json:"first"`
	Last  string `json:"last"`
	Age   int    `json:"age"`
	Admin bool   `json:"admin"`
}

func encode(w http.ResponseWriter, r *http.Request) {
	sliceOfPerson := []person{
		{
			First: "Jamie",
			Last:  "Young",
			Age:   22,
			Admin: true,
		},
		{
			First: "Joe",
			Last:  "Rogan",
			Age:   52,
			Admin: true,
		},
	}

	json.NewEncoder(w).Encode(&sliceOfPerson)
}

func decode(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		sliceOfPerson := []person{}
		json.NewDecoder(r.Body).Decode(&sliceOfPerson)
		fmt.Printf("%+v\n", sliceOfPerson)
	case http.MethodGet:
		w.Write([]byte("Get request"))
	}
}

func main() {
	http.HandleFunc("/encode", encode)
	http.HandleFunc("/decode", decode)

	log.Println("Server listening at port:3000")
	http.ListenAndServe(":3000", nil)
}

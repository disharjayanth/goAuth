package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type person struct {
	First string
}

func main() {
	p1 := person{
		First: "Jenny",
	}

	p2 := person{
		First: "James",
	}

	xp := []person{p1, p2}

	sb, err := json.Marshal(xp)
	if err != nil {
		log.Panic("Error marshalling to JSON:", err)
		return
	}

	fmt.Println(string(sb))

	xp2 := []person{}
	if err = json.Unmarshal(sb, &xp2); err != nil {
		log.Panic(err)
	}
	fmt.Println(xp2)
	fmt.Printf("%+v\n", xp2)
}

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

	p3 := person{}
	p4 := person{}

	xp2 := []person{p3, p4}
	json.Unmarshal(sb, &xp2)
	fmt.Println(xp2)
}

package main

import (
	"encoding/json"
	"log"
	"os"
)

func main() {
	r, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	syx, err := New(r)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.NewEncoder(os.Stdout).Encode(syx); err != nil {
		log.Fatal(err)
	}
}

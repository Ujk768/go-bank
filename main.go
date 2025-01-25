package main

import (
	"fmt"
	"log"
)

func main() {
	store, err := NewPostGresStorage()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Store: %v\n", store)
	server := NewAPIServer(":3000", store)
	server.Run()
}

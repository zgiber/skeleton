package main

import (
	"log"
	"os"
)

func main() {
	transport := NewTransport(&service{}, ":8000")
	initRoutes(transport)

	if err := transport.Run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

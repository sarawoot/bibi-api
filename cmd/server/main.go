package main

import (
	"api/internal/server"
	"log"
)

func main() {
	if err := server.RunServer(); err != nil {
		log.Fatal(err)
	}
}

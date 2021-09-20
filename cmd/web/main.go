package main

import (
	"log"
	"net/http"

	"github.com/VishalHasija/ChatApp/internal/handlers"
)

func main() {
	routes := routes()

	log.Println("starting Channel listener")
	go handlers.ListenToWsChannel()

	log.Println("Starting web server on port 8080")
	_ = http.ListenAndServe(":8080", routes)
}

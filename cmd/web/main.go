package main

import (
	"log"
	"net/http"
)

func main() {
	routes := routes()

	log.Println("Starting web server on port 8080")

	_ = http.ListenAndServe(":8080", routes)
}

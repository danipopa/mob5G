package main

import (
	"log"
	"net/http"
	"github.com/danipopa/nrf/internal/api"
)

func main() {
	// Initialize the router
	router := api.SetupRouter()

	// Start the server
	log.Println("Starting NRF server on port 8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}


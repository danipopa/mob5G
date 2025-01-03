package main

import (
	"log"
	"net/http"


	"github.com/gorilla/mux"

        "github.com/danipopa/mob5g/ausf/src"
)

func main() {
	// Initialize Redis storage
	storage := ausf.NewRedisStorage("redis-service:6379")

	// Create AUSF handlers
	handlers := ausf.NewHandlers(storage)

	// Setup HTTP routes
	router := mux.NewRouter()
	router.HandleFunc("/nausf-auth/v1/authenticate/{ueId}", handlers.Authenticate).Methods("POST")
	router.HandleFunc("/nausf-auth/v1/verify/{ueId}", handlers.Verify).Methods("POST")

	// Start HTTP server
	log.Println("Starting AUSF service on port 8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start AUSF service: %v", err)
	}
}


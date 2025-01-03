package main

import (
	"log"
	"net/http"

	"github.com/danipopa/mob5g/udm/udm"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize Redis storage
	storage := udm.NewRedisStorage("redis-service:6379")

	// Create UDM handlers
	handlers := udm.NewHandlers(storage)

	// Setup HTTP routes
	router := mux.NewRouter()
	router.HandleFunc("/nudm-sdm/v1/subscription-data/{ueId}", handlers.GetSubscriptionData).Methods("GET")
	router.HandleFunc("/nudm-auth/v1/auth-vectors/{ueId}", handlers.GetAuthVector).Methods("GET")

	// Start HTTP server
	log.Println("Starting UDM service on port 8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start UDM service: %v", err)
	}
}


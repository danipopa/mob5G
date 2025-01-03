package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/danipopa/mob5g/pcf/src"
)

func main() {
	// Initialize storage (optional if policies are pre-configured)
	storage := pcf.NewRedisStorage("redis-service:6379")

	// Create PCF handlers
	handlers := pcf.NewHandlers(storage)

	// Setup HTTP routes
	router := mux.NewRouter()
	router.HandleFunc("/npcf-am/v1/mobility-policy", handlers.N15Handler).Methods("POST")
	router.HandleFunc("/npcf-am/v1/mobility-policy/update", handlers.N15UpdateHandler).Methods("PUT")

	// Start HTTP server
	log.Println("Starting PCF service on port 8081...")
	if err := http.ListenAndServe(":8081", router); err != nil {
		log.Fatalf("Failed to start PCF service: %v", err)
	}
}


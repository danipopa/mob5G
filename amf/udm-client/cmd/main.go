package main

import (
	"log"
	"net/http"
	"os"

	"github.com/danipopa/mob5g/amf/udm-client/udmclient"
	"github.com/danipopa/mob5g/amf/udm-client/udmservice"
)

func main() {
	// Fetch UDM base URL from an environment variable or default
	udmBaseURL := os.Getenv("UDM_BASE_URL")
	if udmBaseURL == "" {
		udmBaseURL = "http://udm-service:8080"
	}

	// Initialize the UDM service
	service := udmservice.NewUDMService(udmBaseURL)
	router := service.SetupRouter()

	// Start the HTTP server
	log.Printf("Starting UDM Client Service on port 8081 (UDM base URL: %s)...", udmBaseURL)
	if err := http.ListenAndServe(":8081", router); err != nil {
		log.Fatalf("Failed to start UDM Client Service: %v", err)
	}
}


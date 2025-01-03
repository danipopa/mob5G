package main

import (
	"log"
	"net/http"

	"github.com/danipopa/mob5g/smf/smf-n10/src"
)

func main() {
	// Initialize API handlers
	handlers := smfn10.NewHandlers("http://udm-service:8084")

	// Define routes
	http.HandleFunc("/n10/subscription-data", handlers.GetSubscriptionDataHandler)

	// Start the server
	log.Println("Starting SMF-N10 service on port 8082...")
	log.Fatal(http.ListenAndServe(":8082", nil))
}


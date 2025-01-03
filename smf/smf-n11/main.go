package main

import (
	"log"
	"net/http"

	"github.com/danipopa/mob5g/smf/smf-n11/src"
)

func main() {
	handlers := smfn11.NewHandlers()

	// Define unique endpoints
	http.HandleFunc("/sm-contexts", handlers.CreateSessionHandler)                // POST for session creation
	http.HandleFunc("/sm-contexts/modify", handlers.ModifySessionHandler)         // PUT for session modification
	http.HandleFunc("/sm-contexts/release", handlers.ReleaseSessionHandler)       // DELETE for session release
	http.HandleFunc("/sm-contexts/notify", handlers.HandleSMFEventNotification)   // POST for notifications

	log.Println("Starting SMF-N11 service on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}


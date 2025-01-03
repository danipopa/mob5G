package main

import (
	"log"
	"net/http"

	"github.com/danipopa/mob5g/smf/smf-n7/src"
)

func main() {
	handlers := smfn7.NewHandlers("http://pcf-service:8080") // PCF base URL

	http.HandleFunc("/n7/policy-association", handlers.HandlePolicyAssociation) // POST
	http.HandleFunc("/n7/policy-termination", handlers.HandlePolicyTermination) // DELETE
	http.HandleFunc("/n7/event-reporting", handlers.HandleEventReporting) // POST for event reporting
	http.HandleFunc("/n7/policy-update", handlers.HandlePolicyUpdate) // PUT for policy updates

	log.Println("Starting SMF-N7 service on port 8082...")
	log.Fatal(http.ListenAndServe(":8082", nil))
}


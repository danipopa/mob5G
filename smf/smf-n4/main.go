package main

import (
	"log"
	"time"
	"net/http"
	"github.com/danipopa/mob5g/smf/smf-n4/src"
)

func main() {
	// Initialize components
	pfcpClient := &src.PFCPClient{UPFAddress: "127.0.0.1:8805"}
	sessionHandler := &src.SessionHandler{PFCPClient: pfcpClient}
	webHandlers := &src.WebHandlers{SessionHandler: sessionHandler}

	// Expose APIs for SMF Session Manager
	http.HandleFunc("/n4/establish-session", webHandlers.HandleEstablishSession) // POST
	http.HandleFunc("/n4/modify-session", webHandlers.HandleModifySession)       // PUT
	http.HandleFunc("/n4/release-session", webHandlers.HandleReleaseSession)     // DELETE
	http.HandleFunc("/n4/usage-report", webHandlers.HandleUsageReport)           // POST

	// Start heartbeat goroutine
	go src.StartHeartbeat(pfcpClient, 10*time.Second)

	// Start HTTP server
	log.Println("Starting SMF-N4 service on port 8084...")
	log.Fatal(http.ListenAndServe(":8084", nil))
}


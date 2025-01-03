package main

import (
	"log"
	"net/http"
	"github.com/danipopa/mob5g/smf/smf-session-controller/src"
)

func main() {
	// Initialize Redis client
	redisClient := src.NewRedisClient("redis-service:6379")

	// Initialize Session Manager
	sessionManager := src.NewSessionManager(redisClient)

	// Initialize Handlers
	handlers := src.NewHandlers(sessionManager)

	// Define HTTP routes
	http.HandleFunc("/sessions", handlers.CreateSessionHandler)  // POST
	http.HandleFunc("/sessions/", handlers.ModifySessionHandler) // PUT
	http.HandleFunc("/sessions/", handlers.DeleteSessionHandler) // DELETE

	// Start HTTP server
	log.Println("Starting SMF Session Controller on port 8085...")
	log.Fatal(http.ListenAndServe(":8085", nil))
}


package main

import (
	"net/http"

	"github.com/danipopa/mob5g/udr/src"
)

func main() {
	// Initialize Redis storage
	storage := udr.NewRedisStorage("redis-service:6379")

	// Create UDR handlers
	handlers := udr.NewHandlers(storage)

	// Setup routes
	router := http.NewServeMux()
	router.HandleFunc("/nudr-dr/v1/subscriptions", handlers.GetSubscription)
	router.HandleFunc("/nudr-dr/v1/subscriptions/save", handlers.SaveSubscription)

	// Start HTTPS server
	udr.StartServer(router)
}


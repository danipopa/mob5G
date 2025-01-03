package udr

import (
	"log"
	"net/http"
)

func StartServer(router *http.ServeMux) {
	log.Println("Starting UDR service on port 8080...")
	err := http.ListenAndServe(":8080", router) // Use plain HTTP
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

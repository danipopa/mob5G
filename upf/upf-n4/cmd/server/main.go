package main

import (
	"upf-n4/internal/pfcp"
	"upf-n4/internal/transport"
	"log"
)

func main() {
	log.Println("Starting UPF-N4 PFCP server...")

	// Initialize Redis
	pfcp.InitializeRedis("localhost:6379")

	// Start listening for PFCP messages
	err := transport.ListenForPFCPMessages(8805, pfcp.HandleMessage)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

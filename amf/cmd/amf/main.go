package main

import (
	"log"

	"github.com/danipopa/mob5g/amf/pkg/nrf-agent"
)

func main() {
	nrfURL := "http://nrf"
	amfInstance := "amf-instance-1"
	servicePorts := map[string]int{
		"n1": 38412,
		"n2": 38413,
	}

	nrfClient := nrfagent.NewNRFClient(nrfURL, amfInstance, servicePorts)

	// Register with NRF
	if err := nrfClient.Register(); err != nil {
		log.Fatalf("Failed to register with NRF: %v", err)
	}

	// Start heartbeat
	go nrfClient.Heartbeat()

	// Simulate application lifecycle
	select {}
}


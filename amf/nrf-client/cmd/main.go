package main

import (
	"log"
	"time"

	"github.com/danipopa/mob5g/amf/nrf-client/internal/api"
)

func main() {
	nrfURL := "http://10.107.50.172:80"
	amfID := "amf-instance-1"

	nrfClient := nrfagent.NewNRFClient(nrfURL, amfID)

	profile := nrfagent.NFProfile{
		NFID:           amfID,
		NFInstanceID:   amfID,
		NFType:         "AMF",
		Status:         "ACTIVE",
		FQDN:           "amf1.example.com",
		IPAddresses:    []string{"192.168.1.10"},
		ServiceURLs:    []string{"http://amf1:8080"},
		HeartbeatTimer: 30,
		PLMNID:         nrfagent.PLMN{MCC: "310", MNC: "001"},
		SNSSAIs:        []nrfagent.SNSSAI{{SST: "01", SD: "abc123"}},
		AreaID:         "area1",
		DNNs:           []string{"internet"},
	}

	// Retry registration in a loop
	for {
		err := nrfClient.Register(profile)
		if err != nil {
			log.Printf("Failed to register AMF with NRF: %v. Retrying in 1 second...", err)
			time.Sleep(1 * time.Second)
			continue
		}
		log.Println("Successfully registered AMF with NRF")
		break
	}

	// Start heartbeat in a separate goroutine
	go nrfClient.Heartbeat(profile)

	// Keep the application running
	select {}
}

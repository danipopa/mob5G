package src

import (
	"fmt"
	"time"
)

func StartHeartbeat(client *PFCPClient, interval time.Duration) {
	for {
		header := CreatePFCPHeader(1, 100) // MessageType: Heartbeat Request
		_, err := client.SendPFCPMessage(header)
		if err != nil {
			fmt.Printf("Heartbeat failed: %v\n", err)
		} else {
			fmt.Println("Heartbeat successful")
		}
		time.Sleep(interval)
	}
}


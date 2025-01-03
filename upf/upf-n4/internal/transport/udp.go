package transport

import (
	"log"
	"net"
)

// ListenForPFCPMessages listens for incoming PFCP messages on the specified UDP port.
func ListenForPFCPMessages(port int, handler func([]byte, string)) error {
	addr := net.UDPAddr{Port: port}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	log.Printf("Listening for PFCP messages on port %d...", port)
	buf := make([]byte, 1500)

	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("Error reading UDP packet: %v", err)
			continue
		}
		go handler(buf[:n], remoteAddr.String())
	}
}


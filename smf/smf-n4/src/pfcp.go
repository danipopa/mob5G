package src

import (
	"encoding/binary"
	"fmt"
	"net"
)

// PFCPClient handles PFCP communication with the UPF
type PFCPClient struct {
	UPFAddress string
}

// SendPFCPMessage sends a PFCP message to the UPF and returns the response
func (c *PFCPClient) SendPFCPMessage(message []byte) ([]byte, error) {
	conn, err := net.Dial("udp", c.UPFAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to UPF: %w", err)
	}
	defer conn.Close()

	_, err = conn.Write(message)
	if err != nil {
		return nil, fmt.Errorf("failed to send PFCP message: %w", err)
	}

	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to read PFCP response: %w", err)
	}

	return buffer[:n], nil
}

// CreatePFCPHeader creates a PFCP message header
func CreatePFCPHeader(messageType uint8, sequenceNumber uint32) []byte {
	header := make([]byte, 8)
	header[0] = 0x20 // Version 1
	header[1] = messageType
	binary.BigEndian.PutUint16(header[2:4], 8) // Message Length (placeholder)
	binary.BigEndian.PutUint32(header[4:8], sequenceNumber)
	return header
}

// HandlePFCPUsageReport processes usage reports received from the UPF
func (c *PFCPClient) HandlePFCPUsageReport(data []byte) (*UsageReport, error) {
	// Decode the binary PFCP message into a UsageReport structure
	if len(data) < 8 {
		return nil, fmt.Errorf("invalid PFCP usage report length")
	}

	report := &UsageReport{
		SessionID:  "session-id-placeholder", // Decode from the message
		VolumeMB:   12345,                    // Example decoded value
		DurationMS: 60000,                    // Example decoded value
	}
	return report, nil
}


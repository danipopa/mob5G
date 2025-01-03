package src

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// SessionHandler handles session-related requests
type SessionHandler struct {
	PFCPClient *PFCPClient
}

// EstablishSession sends a PFCP Session Establishment Request
func (h *SessionHandler) EstablishSession(pdr PDR) error {
	header := CreatePFCPHeader(50, 1) // MessageType: PFCP Session Establishment Request
	// Append PDR or other payload data as required
	message := append(header, encodePDR(pdr)...)

	response, err := h.PFCPClient.SendPFCPMessage(message)
	if err != nil {
		return fmt.Errorf("failed to establish session: %w", err)
	}

	fmt.Printf("PFCP Session Establishment Response: %x\n", response)
	return nil
}

// ModifySession sends a PFCP Session Modification Request
func (h *SessionHandler) ModifySession(pdr PDR) error {
	header := CreatePFCPHeader(51, 2) // MessageType: PFCP Session Modification Request
	// Append PDR or other payload data as required
	message := append(header, encodePDR(pdr)...)

	response, err := h.PFCPClient.SendPFCPMessage(message)
	if err != nil {
		return fmt.Errorf("failed to modify session: %w", err)
	}

	fmt.Printf("PFCP Session Modification Response: %x\n", response)
	return nil
}

// ReleaseSession sends a PFCP Session Deletion Request
func (h *SessionHandler) ReleaseSession(sessionID uint32) error {
	header := CreatePFCPHeader(52, 3) // MessageType: PFCP Session Deletion Request
	// Append session ID or other payload data as required
	message := header

	response, err := h.PFCPClient.SendPFCPMessage(message)
	if err != nil {
		return fmt.Errorf("failed to release session: %w", err)
	}

	fmt.Printf("PFCP Session Deletion Response: %x\n", response)
	return nil
}

// Helper function to encode a PDR (example, not complete)
func encodePDR(pdr PDR) []byte {
	// Serialize PDR fields into binary format
	return []byte{} // Placeholder
}

// ForwardUsageReportToSessionManager forwards the usage report to the SMF Session Manager
func (h *SessionHandler) ForwardUsageReportToSessionManager(report UsageReport) error {
	url := "http://smf-session-manager:8080/usage-report" // SMF Session Manager API
	payload, err := json.Marshal(report)
	if err != nil {
		return fmt.Errorf("failed to encode usage report: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to forward usage report to session manager: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("session manager returned non-OK status: %d", resp.StatusCode)
	}

	return nil
}


package src

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// N4Client communicates with SMF-N4-Service
type N4Client struct {
	BaseURL string
}

// NewN4Client initializes an N4 client
func NewN4Client(baseURL string) *N4Client {
	return &N4Client{BaseURL: baseURL}
}

// CreateSessionInUPF sends a session creation request to SMF-N4
func (c *N4Client) CreateSessionInUPF(session *Session) error {
	url := fmt.Sprintf("%s/sm-contexts", c.BaseURL)
	payload, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to encode session for SMF-N4: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to communicate with SMF-N4: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("SMF-N4 returned non-OK status: %d", resp.StatusCode)
	}

	return nil
}

// ReleaseSessionInUPF sends a session deletion request to the UPF via SMF-N4
func (c *N4Client) ReleaseSessionInUPF(sessionID string) error {
	url := fmt.Sprintf("%s/sm-contexts/%s", c.BaseURL, sessionID)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request for session deletion: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to communicate with SMF-N4: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("SMF-N4 returned non-OK status: %d", resp.StatusCode)
	}

	return nil
}

// ModifySessionInUPF sends a session modification request to the UPF via SMF-N4
func (c *N4Client) ModifySessionInUPF(session *Session) error {
	url := fmt.Sprintf("%s/sm-contexts/%s", c.BaseURL, session.SessionID)
	payload, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to encode session for modification: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request for session modification: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to communicate with SMF-N4: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("SMF-N4 returned non-OK status: %d", resp.StatusCode)
	}

	return nil
}

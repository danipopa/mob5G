package src

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// N11Client handles communication with SMF-N11-Service
type N11Client struct {
	BaseURL string
}

// NewN11Client initializes a new N11 client
func NewN11Client(baseURL string) *N11Client {
	return &N11Client{BaseURL: baseURL}
}

// NotifyAMF sends a notification to the AMF via smf-n11-service
func (c *N11Client) NotifyAMF(notification *AMFNotification) error {
	url := fmt.Sprintf("%s/amf-notify", c.BaseURL)
	payload, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to encode notification: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to send notification to AMF: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("SMF-N11 returned non-OK status: %d", resp.StatusCode)
	}

	return nil
}

// QueryAMF queries the AMF for specific session-related information
func (c *N11Client) QueryAMF(sessionID string) (*AMFResponse, error) {
	url := fmt.Sprintf("%s/sm-contexts/%s", c.BaseURL, sessionID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to query AMF: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("SMF-N11 returned non-OK status: %d", resp.StatusCode)
	}

	var response AMFResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode AMF response: %w", err)
	}

	return &response, nil
}


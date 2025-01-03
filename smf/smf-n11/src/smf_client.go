package smfn11

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type SMFClient struct {
	BaseURL string
}

// ForwardRequest forwards a request to the SMF with timeout handling
func (c *SMFClient) ForwardRequest(request PDURequest) (*PDUResponse, error) {
	url := fmt.Sprintf("%s/sm-contexts", c.BaseURL)
	payload, _ := json.Marshal(request)

	client := &http.Client{Timeout: 5 * time.Second} // Add timeout
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to contact SMF: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("SMF returned non-OK status: %d", resp.StatusCode)
	}

	var response PDUResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to parse SMF response: %w", err)
	}

	return &response, nil
}


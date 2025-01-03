package src

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// N7Client communicates with SMF-N7-Service
type N7Client struct {
	BaseURL string
}

// NewN7Client initializes an N7 client
func NewN7Client(baseURL string) *N7Client {
	return &N7Client{BaseURL: baseURL}
}

// FetchPolicyData retrieves policy data from SMF-N7
func (c *N7Client) FetchPolicyData(session *Session) error {
	url := fmt.Sprintf("%s/policies", c.BaseURL)
	payload, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to encode session for SMF-N7: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to fetch policy data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("SMF-N7 returned non-OK status: %d", resp.StatusCode)
	}

	return nil
}


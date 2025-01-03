//SMF-N7 Logic and PCF Communication

package smfn7

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type PCFClient struct {
	BaseURL string
}

// SendPolicyRequest sends a policy request to the PCF
func (c *PCFClient) SendPolicyRequest(contextData SMPolicyContextData) (*SMPolicyDecision, error) {
	url := fmt.Sprintf("%s/sm-policies", c.BaseURL)
	payload, _ := json.Marshal(contextData)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to contact PCF: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("PCF returned non-OK status: %d", resp.StatusCode)
	}

	var decision SMPolicyDecision
	if err := json.NewDecoder(resp.Body).Decode(&decision); err != nil {
		return nil, fmt.Errorf("failed to parse PCF response: %w", err)
	}

	return &decision, nil
}


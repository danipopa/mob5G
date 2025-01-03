package smfn10

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// UDMClient represents the client for UDM interactions
type UDMClient struct {
	BaseURL string // UDM Base URL
}

// GetSubscriptionData retrieves subscription data for a given UEID from the UDM
func (c *UDMClient) GetSubscriptionData(ueID string) (*SubscriptionData, error) {
	url := fmt.Sprintf("%s/subscription-data/%s", c.BaseURL, ueID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to contact UDM: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("UDM returned non-OK status: %d", resp.StatusCode)
	}

	var data SubscriptionData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to parse UDM response: %w", err)
	}

	return &data, nil
}


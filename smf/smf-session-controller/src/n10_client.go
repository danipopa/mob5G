package src

import (
	"fmt"
	"net/http"
)

// N10Client communicates with SMF-N10-Service
type N10Client struct {
	BaseURL string
}

// NewN10Client initializes an N10 client
func NewN10Client(baseURL string) *N10Client {
	return &N10Client{BaseURL: baseURL}
}

// FetchSubscriptionData retrieves subscription data from SMF-N10
func (c *N10Client) FetchSubscriptionData(ueID string) error {
	url := fmt.Sprintf("%s/subscriptions/%s", c.BaseURL, ueID)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch subscription data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("SMF-N10 returned non-OK status: %d", resp.StatusCode)
	}

	return nil
}


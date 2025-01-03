package pcf

import (
	"fmt"
	"net/http"
	"encoding/json"
)

type UserSubscription struct {
	UEID        string `json:"ue_id"`
	Mobility    string `json:"mobility"`    // HANDOVER_ALLOWED, ATTACH_ALLOWED
	DataQuota   int    `json:"data_quota"` // Remaining data quota
}

// RetrieveSubscriptionData queries UDR/UDM for subscription details.
func RetrieveSubscriptionData(ueID string) (*UserSubscription, error) {
	// Simulate an HTTP request to UDR/UDM
	resp, err := http.Get(fmt.Sprintf("http://udr-service:8080/subscription/%s", ueID))
	if err != nil {
		return nil, fmt.Errorf("failed to query UDR/UDM: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("UDR/UDM returned status: %d", resp.StatusCode)
	}

	var subscription UserSubscription
	if err := json.NewDecoder(resp.Body).Decode(&subscription); err != nil {
		return nil, fmt.Errorf("failed to decode UDR/UDM response: %w", err)
	}
	return &subscription, nil
}


package udmclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// UDMClient represents the client for interacting with the UDM
type UDMClient struct {
	udmBaseURL string
	httpClient *http.Client
}

// NewUDMClient initializes a new UDM client
func NewUDMClient(baseURL string) *UDMClient {
	return &UDMClient{
		udmBaseURL: baseURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// SubscriptionData represents the UE subscription data
type SubscriptionData struct {
	UEID                string   `json:"ue_id"`
	AllowedPLMNs        []PLMN   `json:"allowed_plmns"`
	AllowedSNSSAIs      []SNSSAI `json:"allowed_snssais"`
	MobilityRestrictions string   `json:"mobility_restrictions"`
}

// AuthVector represents the authentication vector for a UE
type AuthVector struct {
	RAND  string `json:"rand"`
	AUTN  string `json:"autn"`
	XRES  string `json:"xres"`
	KASME string `json:"kasme"`
}

// PLMN represents a Public Land Mobile Network identifier
type PLMN struct {
	MCC string `json:"mcc"`
	MNC string `json:"mnc"`
}

// SNSSAI represents the Single Network Slice Selection Assistance Information
type SNSSAI struct {
	SST string `json:"sst"`
	SD  string `json:"sd"`
}

// GetSubscriptionData fetches subscription data for a given UE ID
func (c *UDMClient) GetSubscriptionData(ueID string) (*SubscriptionData, error) {
	url := fmt.Sprintf("%s/nudm-sdm/v1/subscription-data/%s", c.udmBaseURL, ueID)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch subscription data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response from UDM: %s", resp.Status)
	}

	var data SubscriptionData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode subscription data: %w", err)
	}

	return &data, nil
}

// GetAuthVector fetches the authentication vector for a given UE ID
func (c *UDMClient) GetAuthVector(ueID string) (*AuthVector, error) {
	url := fmt.Sprintf("%s/nudm-auth/v1/auth-vectors/%s", c.udmBaseURL, ueID)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch authentication vector: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response from UDM: %s", resp.Status)
	}

	var vector AuthVector
	if err := json.NewDecoder(resp.Body).Decode(&vector); err != nil {
		return nil, fmt.Errorf("failed to decode authentication vector: %w", err)
	}

	return &vector, nil
}


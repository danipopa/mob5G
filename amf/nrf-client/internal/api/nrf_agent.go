package nrfagent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"strings" // Import strings for utility functions like Contains

)

// NRFClient represents the NRF client for interacting with the NRF
type NRFClient struct {
	nrfURL     string
	amfID      string
	httpClient *http.Client
}

// NFProfile represents the NF profile for registration and updates
type NFProfile struct {
	NFID           string   `json:"nf_id"`
	NFInstanceID   string   `json:"nf_instance_id"`
	NFType         string   `json:"nf_type"`
	Status         string   `json:"status"`
	FQDN           string   `json:"fqdn"`
	IPAddresses    []string `json:"ip_addresses"`
	ServiceURLs    []string `json:"service_urls"`
	HeartbeatTimer int      `json:"heartbeat_timer"`
	PLMNID         PLMN     `json:"plmn_id"`
	SNSSAIs        []SNSSAI `json:"snssais"`
	AreaID         string   `json:"area_id"`
	DNNs           []string `json:"dnns"`
}

// PLMN represents the Public Land Mobile Network identifier
type PLMN struct {
	MCC string `json:"mcc"`
	MNC string `json:"mnc"`
}

// SNSSAI represents the Single Network Slice Selection Assistance Information
type SNSSAI struct {
	SST string `json:"sst"`
	SD  string `json:"sd"`
}

// NewNRFClient initializes a new NRF client
func NewNRFClient(nrfURL, amfID string) *NRFClient {
	return &NRFClient{
		nrfURL:     nrfURL,
		amfID:      amfID,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// Register registers the AMF with the NRF
func (n *NRFClient) Register(profile NFProfile) error {
	url := fmt.Sprintf("%s/nnrf-nfm/v1/nf-instances/%s", n.nrfURL, n.amfID)
	return n.sendRequest(http.MethodPut, url, profile)
}

// Deregister deregisters the AMF from the NRF
func (n *NRFClient) Deregister() error {
	url := fmt.Sprintf("%s/nnrf-nfm/v1/nf-instances/%s", n.nrfURL, n.amfID)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create deregistration request: %w", err)
	}
	resp, err := n.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send deregistration request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected response from NRF during deregistration: %s", resp.Status)
	}
	fmt.Println("Successfully deregistered AMF from NRF")
	return nil
}

// Update updates the AMF profile in the NRF
func (n *NRFClient) Update(updatedFields map[string]interface{}) error {
	url := fmt.Sprintf("%s/nnrf-nfm/v1/nf-instances/%s", n.nrfURL, n.amfID)
	return n.sendRequest(http.MethodPatch, url, updatedFields)
}

// Heartbeat sends a heartbeat to the NRF to indicate liveness
func (n *NRFClient) Heartbeat(profile NFProfile) {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        url := fmt.Sprintf("%s/nnrf-nfm/v1/nf-instances/%s", n.nrfURL, n.amfID)
        payload := map[string]interface{}{
            "status": "ACTIVE",
        }
        err := n.sendRequest(http.MethodPatch, url, payload)
        if err != nil {
            if isNotFoundError(err) {
                // Handle NF profile not found
                handleNFNotFound(n, profile)
            } else {
                fmt.Printf("Failed to send heartbeat: %v\n", err)
            }
        } else {
            fmt.Println("Heartbeat sent successfully")
        }
    }
}

func isNotFoundError(err error) bool {
    return err != nil && strings.Contains(err.Error(), "404")
}

func handleNFNotFound(nrfClient *NRFClient, profile NFProfile) {
    fmt.Printf("NF profile not found in NRF. Attempting re-registration...\n")
    retryCount := 0
    for {
        err := nrfClient.Register(profile)
        if err != nil {
            retryCount++
            waitTime := time.Duration(retryCount*2) * time.Second
            fmt.Printf("Re-registration failed. Retrying in %v... (attempt %d)\n", waitTime, retryCount)
            time.Sleep(waitTime)
        } else {
            fmt.Println("Successfully re-registered with NRF.")
            break
        }
    }
}

// NotifyOperators sends an alert to the operator in case of persistent issues
func NotifyOperators(message string) {
	// Simulate sending an alert (e.g., email, SMS, or external API)
	fmt.Printf("ALERT: %s\n", message)
}

// Discover retrieves NF instances of a specific type from the NRF
func (n *NRFClient) Discover(nfType string) ([]NFProfile, error) {
	url := fmt.Sprintf("%s/nnrf-disc/v1/nfs?nf-type=%s", n.nrfURL, nfType)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create discovery request: %w", err)
	}
	resp, err := n.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to query NRF: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response from NRF: %s", resp.Status)
	}
	var nfInstances []NFProfile
	if err := json.NewDecoder(resp.Body).Decode(&nfInstances); err != nil {
		return nil, fmt.Errorf("failed to decode NRF response: %w", err)
	}
	return nfInstances, nil
}

// Subscribe subscribes to notifications from the NRF
func (n *NRFClient) Subscribe(subscriptionPayload map[string]interface{}) error {
	url := fmt.Sprintf("%s/nnrf-nfm/v1/subscriptions", n.nrfURL)
	return n.sendRequest(http.MethodPost, url, subscriptionPayload)
}

// HandleNotification processes incoming notifications from the NRF
func (n *NRFClient) HandleNotification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}
	var notification map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
		http.Error(w, "Failed to parse notification payload", http.StatusBadRequest)
		return
	}
	fmt.Printf("Received notification from NRF: %+v\n", notification)
	w.WriteHeader(http.StatusOK)
}

// sendRequest is a helper function to send HTTP requests
func (n *NRFClient) sendRequest(method, url string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}
	req, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := n.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected response: %s", resp.Status)
	}
	fmt.Printf("Successfully sent %s request to %s\n", method, url)
	return nil
}


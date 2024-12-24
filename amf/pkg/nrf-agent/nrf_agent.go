package nrfagent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type NRFClient struct {
	nrfURL       string
	amfInstance  string
	nfType       string
	nfStatus     string
	servicePorts map[string]int
	httpClient   *http.Client
	heartbeatInterval time.Duration
}

type NFInstance struct {
	NFInstanceID string            `json:"nfInstanceId"`
	NFType       string            `json:"nfType"`
	NFStatus     string            `json:"nfStatus"`
	Services     map[string]int    `json:"services"`
}

// NewNRFClient initializes a new NRF Client
func NewNRFClient(nrfURL, amfInstance string, servicePorts map[string]int) *NRFClient {
	return &NRFClient{
		nrfURL:       nrfURL,
		amfInstance:  amfInstance,
		nfType:       "AMF",
		nfStatus:     "REGISTERED",
		servicePorts: servicePorts,
		httpClient:   &http.Client{},
		heartbeatInterval: 30 * time.Second, // Default 30s heartbeat interval
	}
}

// Register registers the AMF with the NRF
func (n *NRFClient) Register() error {
	payload := NFInstance{
		NFInstanceID: n.amfInstance,
		NFType:       n.nfType,
		NFStatus:     n.nfStatus,
		Services:     n.servicePorts,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal NFInstance: %w", err)
	}

	url := fmt.Sprintf("%s/nnrf-nfm/v1/nf-instances/%s", n.nrfURL, n.amfInstance)
	req, err := http.NewRequest("PUT", url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to register with NRF: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected response from NRF: %s", resp.Status)
	}

	fmt.Println("Successfully registered with NRF")
	return nil
}

// Deregister deregisters the AMF from the NRF
func (n *NRFClient) Deregister() error {
	url := fmt.Sprintf("%s/nnrf-nfm/v1/nf-instances/%s", n.nrfURL, n.amfInstance)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := n.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to deregister from NRF: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected response from NRF: %s", resp.Status)
	}

	fmt.Println("Successfully deregistered from NRF")
	return nil
}

// Heartbeat sends periodic heartbeats to the NRF
func (n *NRFClient) Heartbeat() {
	for {
		time.Sleep(n.heartbeatInterval)
		url := fmt.Sprintf("%s/nnrf-nfm/v1/nf-instances/%s/heartbeat", n.nrfURL, n.amfInstance)
		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			fmt.Printf("Failed to create heartbeat request: %v\n", err)
			continue
		}

		resp, err := n.httpClient.Do(req)
		if err != nil {
			fmt.Printf("Failed to send heartbeat: %v\n", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Unexpected response from NRF during heartbeat: %s\n", resp.Status)
		} else {
			fmt.Println("Heartbeat sent successfully")
		}
	}
}

// Update updates the AMF registration in the NRF
func (n *NRFClient) Update(updatedFields map[string]interface{}) error {
	url := fmt.Sprintf("%s/nnrf-nfm/v1/nf-instances/%s", n.nrfURL, n.amfInstance)

	data, err := json.Marshal(updatedFields)
	if err != nil {
		return fmt.Errorf("failed to marshal update payload: %w", err)
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create update request: %w", err)
	}
	req.Header.Set("Content-Type", "application/merge-patch+json")

	resp, err := n.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update NRF registration: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response from NRF during update: %s", resp.Status)
	}

	fmt.Println("Successfully updated AMF registration with NRF")
	return nil
}

// HandleNotification processes notifications from the NRF
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

	// Process notification logic here (e.g., update local state)
	w.WriteHeader(http.StatusOK)
}

// Subscribe subscribes to NRF notifications
func (n *NRFClient) Subscribe(subscriptionPayload map[string]interface{}) error {
	url := fmt.Sprintf("%s/nnrf-sub/v1/subscriptions", n.nrfURL)

	data, err := json.Marshal(subscriptionPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal subscription payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create subscription request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to subscribe to NRF: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected response from NRF during subscription: %s", resp.Status)
	}

	fmt.Println("Successfully subscribed to NRF notifications")
	return nil
}


// Discover queries the NRF for other NFs
func (n *NRFClient) Discover(nfType string) ([]NFInstance, error) {
	url := fmt.Sprintf("%s/nnrf-disc/v1/nfs?nf-type=%s", n.nrfURL, nfType)
	req, err := http.NewRequest("GET", url, nil)
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

	var nfInstances []NFInstance
	if err := json.NewDecoder(resp.B


package smfn11

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type AMFClient struct {
	BaseURL string
}

// SendResponse sends a response back to the AMF
func (c *AMFClient) SendResponse(response *PDUResponse) error {
	url := fmt.Sprintf("%s/amf/response", c.BaseURL)
	payload, _ := json.Marshal(response)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to send response to AMF: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("AMF returned non-OK status: %d", resp.StatusCode)
	}

	return nil
}

// NotifyUE sends session-related updates to the UE
func (c *AMFClient) NotifyUE(notification map[string]interface{}) error {
	url := fmt.Sprintf("%s/ue-notification", c.BaseURL)
	payload, _ := json.Marshal(notification)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to send notification to UE: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("UE returned non-OK status: %d", resp.StatusCode)
	}

	return nil
}

// SendNotification sends notifications to the AMF
func (c *AMFClient) SendNotification(notification *GenericRequest) error {
	url := fmt.Sprintf("%s/notify", c.BaseURL)
	payload, _ := json.Marshal(notification)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to send notification to AMF: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("AMF returned non-OK status: %d", resp.StatusCode)
	}

	return nil
}


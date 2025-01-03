//N7 Service Logic

package smfn7

import (
	"fmt"
	"net/http"
)

type N7Service struct {
	PCFClient *PCFClient
}

// CreatePolicyAssociation handles policy association creation
func (s *N7Service) CreatePolicyAssociation(contextData SMPolicyContextData) (*SMPolicyDecision, error) {
	return s.PCFClient.SendPolicyRequest(contextData)
}

// DeletePolicyAssociation handles policy termination
func (s *N7Service) DeletePolicyAssociation(policyID string) error {
	url := fmt.Sprintf("%s/sm-policies/%s", s.PCFClient.BaseURL, policyID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to contact PCF: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("PCF returned non-OK status: %d", resp.StatusCode)
	}

	return nil
}


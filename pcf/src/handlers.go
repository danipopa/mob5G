package pcf

import (
	"encoding/json"
	"net/http"
	"fmt" // Add this line
)

type Handlers struct {
	storage Storage
}

// NewHandlers initializes the Handlers struct
func NewHandlers(storage Storage) *Handlers {
	return &Handlers{storage: storage}
}

func (h *Handlers) N15Handler(w http.ResponseWriter, r *http.Request) {
	var request AMPolicyRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid Mobility Policy Request JSON", http.StatusBadRequest)
		return
	}

	// Evaluate the mobility policy
	engine := PolicyDecisionEngine{}
	response := engine.EvaluateMobilityPolicy(request)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) N15UpdateHandler(w http.ResponseWriter, r *http.Request) {
	var updateRequest AMPolicyRequest
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		http.Error(w, "Invalid Policy Update JSON", http.StatusBadRequest)
		return
	}

	// Update the policy in storage
	newPolicy := AMPolicyResponse{
		PolicyID:       "policy-" + updateRequest.UEID,
		AccessAllowed:  true,
		HandoverAllowed: updateRequest.MobilityType == "HANDOVER",
		QoSPriority:    1, // Example new priority
		Message:        "Policy updated successfully",
	}

	if err := h.storage.SaveMobilityPolicy(newPolicy); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update policy: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newPolicy)
}



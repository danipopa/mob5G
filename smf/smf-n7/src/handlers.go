//RESTful APIs for SMF-PCF Communication (TS 29.512)
package smfn7

import (
	"encoding/json"
	"net/http"
)

type Handlers struct {
	N7Service *N7Service
}

func NewHandlers(pcfBaseURL string) *Handlers {
	return &Handlers{
		N7Service: &N7Service{
			PCFClient: &PCFClient{BaseURL: pcfBaseURL},
		},
	}
}

// HandlePolicyAssociation handles policy association creation
func (h *Handlers) HandlePolicyAssociation(w http.ResponseWriter, r *http.Request) {
	var contextData SMPolicyContextData
	if err := json.NewDecoder(r.Body).Decode(&contextData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	decision, err := h.N7Service.CreatePolicyAssociation(contextData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(decision)
}

// HandlePolicyTermination handles policy termination
func (h *Handlers) HandlePolicyTermination(w http.ResponseWriter, r *http.Request) {
	policyID := r.URL.Query().Get("policy_id")
	if policyID == "" {
		http.Error(w, "Policy ID is required", http.StatusBadRequest)
		return
	}

	if err := h.N7Service.DeletePolicyAssociation(policyID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandleEventReporting handles event reporting to PCF
func (h *Handlers) HandleEventReporting(w http.ResponseWriter, r *http.Request) {
    policyID := r.URL.Query().Get("policy_id")
    if policyID == "" {
        http.Error(w, "Policy ID is required", http.StatusBadRequest)
        return
    }

    var eventReport map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&eventReport); err != nil {
        http.Error(w, "Invalid event report payload", http.StatusBadRequest)
        return
    }

    // Log the event (or forward it to another system if required)
    eventReport["policyId"] = policyID
    w.WriteHeader(http.StatusNoContent)
}

// HandlePolicyUpdate handles PCF-initiated policy updates
func (h *Handlers) HandlePolicyUpdate(w http.ResponseWriter, r *http.Request) {
    policyID := r.URL.Query().Get("policy_id")
    if policyID == "" {
        http.Error(w, "Policy ID is required", http.StatusBadRequest)
        return
    }

    var updateData SMPolicyDecision
    if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
        http.Error(w, "Invalid policy update payload", http.StatusBadRequest)
        return
    }

    // Apply the policy update (stub logic for now)
    updateData.PolicyID = policyID
    w.WriteHeader(http.StatusNoContent)
}


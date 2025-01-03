package smfn10

import (
	"encoding/json"
	"net/http"
)

// Handlers contains the logic for handling N10 API requests
type Handlers struct {
	UDMClient *UDMClient // Client for UDM interactions
}

// NewHandlers initializes a new Handlers instance
func NewHandlers(udmBaseURL string) *Handlers {
	return &Handlers{
		UDMClient: &UDMClient{BaseURL: udmBaseURL},
	}
}

// GetSubscriptionDataHandler handles requests to retrieve subscription data
func (h *Handlers) GetSubscriptionDataHandler(w http.ResponseWriter, r *http.Request) {
	ueID := r.URL.Query().Get("ue_id")
	if ueID == "" {
		http.Error(w, "Missing UEID parameter", http.StatusBadRequest)
		return
	}

	// Fetch data from UDM
	data, err := h.UDMClient.GetSubscriptionData(ueID)
	if err != nil {
		http.Error(w, "Failed to retrieve subscription data", http.StatusInternalServerError)
		return
	}

	// Return the data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}


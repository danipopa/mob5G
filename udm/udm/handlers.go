package udm

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Handlers struct {
	storage Storage
}

// NewHandlers initializes UDM HTTP handlers
func NewHandlers(storage Storage) *Handlers {
	return &Handlers{storage: storage}
}

// GetSubscriptionData handles requests for subscription data
func (h *Handlers) GetSubscriptionData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ueID := vars["ueId"]

	data, err := h.storage.GetSubscriptionData(ueID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// GetAuthVector handles requests for authentication vectors
func (h *Handlers) GetAuthVector(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ueID := vars["ueId"]

	vector, err := h.storage.GetAuthVector(ueID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vector)
}


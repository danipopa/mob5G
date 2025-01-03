package udr

import (
	"encoding/json"
	"fmt"
	"net/http"
)


// Handlers structure
type Handlers struct {
	storage *RedisStorage
}

// NewHandlers initializes API handlers
func NewHandlers(storage *RedisStorage) *Handlers {
	return &Handlers{storage: storage}
}

// GetSubscription retrieves a subscription by UEID
func (h *Handlers) GetSubscription(w http.ResponseWriter, r *http.Request) {
	ueID := r.URL.Query().Get("ue_id")
	subscription, err := h.storage.GetSubscription(ueID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve subscription: %v", err), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(subscription)
}

// SaveSubscription stores a subscription
func (h *Handlers) SaveSubscription(w http.ResponseWriter, r *http.Request) {
	var subscription Subscription
	if err := json.NewDecoder(r.Body).Decode(&subscription); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if err := h.storage.SaveSubscription(subscription); err != nil {
		http.Error(w, fmt.Sprintf("Failed to save subscription: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// Similar methods for Policies and Sessions can be added here


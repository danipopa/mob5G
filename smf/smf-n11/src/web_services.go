package smfn11

import (
	"encoding/json"
	"net/http"
)

// NotifyAMFHandler handles requests from SMF Session Manager to notify the AMF
func (h *Handlers) NotifyAMFHandler(w http.ResponseWriter, r *http.Request) {
	var notification GenericRequest
	if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
		http.Error(w, "Invalid notification payload", http.StatusBadRequest)
		return
	}

	// Forward the notification to AMF
	err := h.AMFClient.SendNotification(&notification)
	if err != nil {
		http.Error(w, "Failed to notify AMF", http.StatusInternalServerError)
		return
	}

	// Acknowledge the notification
	w.WriteHeader(http.StatusOK)
}


package smfn11

import (
	"encoding/json"
	"net/http"
)

type Handlers struct {
	SessionManager *SessionManager
	AMFClient      *AMFClient // Add the AMFClient field
}

func NewHandlers() *Handlers {
	return &Handlers{
		SessionManager: &SessionManager{},
		AMFClient:      &AMFClient{BaseURL: "http://amf-service:8080"}, // Initialize AMFClient with base URLvi
	}
}

// CreateSessionHandler handles session creation requests
func (h *Handlers) CreateSessionHandler(w http.ResponseWriter, r *http.Request) {
	var request PDURequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	response := h.SessionManager.CreateSession(request)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) ModifySessionHandler(w http.ResponseWriter, r *http.Request) {
	smContextID := r.URL.Query().Get("session_id") // Extract session ID from query parameter
	if smContextID == "" {
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	var request PDURequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	request.SessionID = smContextID
	response := h.SessionManager.ModifySession(request)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) ReleaseSessionHandler(w http.ResponseWriter, r *http.Request) {
	smContextID := r.URL.Query().Get("session_id") // Extract session ID from query parameter
	if smContextID == "" {
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	response := h.SessionManager.ReleaseSession(smContextID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) HandleSMFEventNotification(w http.ResponseWriter, r *http.Request) {
	var notification map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
		http.Error(w, "Invalid notification payload", http.StatusBadRequest)
		return
	}

	// Forward the notification to the UE via AMF
	if err := h.AMFClient.NotifyUE(notification); err != nil {
		http.Error(w, "Failed to notify UE", http.StatusInternalServerError)
		return
	}

	// Acknowledge the notification
	w.WriteHeader(http.StatusNoContent)
}

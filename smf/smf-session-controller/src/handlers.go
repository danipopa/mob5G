package src

import (
	"encoding/json"
	"net/http"
)

// Handlers provides HTTP handlers for session management
type Handlers struct {
	sessionManager *SessionManager
}

// NewHandlers initializes Handlers
func NewHandlers(manager *SessionManager) *Handlers {
	return &Handlers{sessionManager: manager}
}

// CreateSessionHandler handles session creation requests
func (h *Handlers) CreateSessionHandler(w http.ResponseWriter, r *http.Request) {
	var req SessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	session, err := h.sessionManager.CreateSession(&req)
	if err != nil {
		http.Error(w, "Failed to create session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

// ModifySessionHandler handles session modification requests
func (h *Handlers) ModifySessionHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Path[len("/sessions/"):]
	var req SessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	session, err := h.sessionManager.ModifySession(sessionID, &req)
	if err != nil {
		http.Error(w, "Failed to modify session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

// DeleteSessionHandler handles session deletion requests
func (h *Handlers) DeleteSessionHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Path[len("/sessions/"):]
	if err := h.sessionManager.DeleteSession(sessionID); err != nil {
		http.Error(w, "Failed to delete session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}


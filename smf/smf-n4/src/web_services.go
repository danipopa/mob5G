package src

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type WebHandlers struct {
	SessionHandler *SessionHandler
}

// HandleEstablishSession handles session establishment requests from the SMF Session Manager
func (h *WebHandlers) HandleEstablishSession(w http.ResponseWriter, r *http.Request) {
	var request PDR
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := h.SessionHandler.EstablishSession(request)
	if err != nil {
		http.Error(w, "Failed to establish session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Session established successfully"))
}

// HandleModifySession handles session modification requests from the SMF Session Manager
func (h *WebHandlers) HandleModifySession(w http.ResponseWriter, r *http.Request) {
	var request PDR
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := h.SessionHandler.ModifySession(request)
	if err != nil {
		http.Error(w, "Failed to modify session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Session modified successfully"))
}

func (h *WebHandlers) HandleReleaseSession(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	// Convert sessionID to uint32 if necessary
	id, err := strconv.ParseUint(sessionID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	err = h.SessionHandler.ReleaseSession(uint32(id))
	if err != nil {
		http.Error(w, "Failed to release session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Session released successfully"))
}

func (h *WebHandlers) HandleUsageReport(w http.ResponseWriter, r *http.Request) {
	var usageReport UsageReport
	if err := json.NewDecoder(r.Body).Decode(&usageReport); err != nil {
		http.Error(w, "Invalid usage report payload", http.StatusBadRequest)
		return
	}

	// Forward the usage report to the SMF Session Manager
	err := h.SessionHandler.ForwardUsageReportToSessionManager(usageReport)
	if err != nil {
		http.Error(w, "Failed to forward usage report: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Usage report processed successfully"))
}

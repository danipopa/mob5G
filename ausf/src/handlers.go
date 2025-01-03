package ausf

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Handlers struct {
	storage Storage
}

// NewHandlers initializes AUSF HTTP handlers
func NewHandlers(storage Storage) *Handlers {
	return &Handlers{storage: storage}
}

// Authenticate handles authentication requests
func (h *Handlers) Authenticate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ueID := vars["ueId"]

	// Generate authentication vector
	authVector := GenerateAuthVector(ueID)

	// Save authentication session
	if err := h.storage.SaveAuthSession(ueID, *authVector); err != nil {
		http.Error(w, fmt.Sprintf("Failed to save auth session: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authVector)
}

// Verify handles authentication verification requests
func (h *Handlers) Verify(w http.ResponseWriter, r *http.Request) {
	var verifyRequest VerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&verifyRequest); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	authData, err := h.storage.GetAuthSession(verifyRequest.UEID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve auth session: %v", err), http.StatusNotFound)
		return
	}

	result := VerifyAuthResponse(authData.XRES, verifyRequest.Response)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(VerifyResponse{Result: result})
}


package udmservice
import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/danipopa/mob5g/amf/udm-client/udmclient"	

	"github.com/gorilla/mux"

)

// UDMService represents the microservice for the UDM client
type UDMService struct {
	Client *udmclient.UDMClient
}

// NewUDMService initializes the UDMService with a UDM client
func NewUDMService(udmBaseURL string) *UDMService {
	return &UDMService{
		Client: udmclient.NewUDMClient(udmBaseURL),
	}
}

// GetSubscriptionDataHandler handles requests for subscription data
func (s *UDMService) GetSubscriptionDataHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ueID := vars["ueId"]

	data, err := s.Client.GetSubscriptionData(ueID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve subscription data: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// GetAuthVectorHandler handles requests for authentication vectors
func (s *UDMService) GetAuthVectorHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ueID := vars["ueId"]

	vector, err := s.Client.GetAuthVector(ueID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve authentication vector: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vector)
}

// SetupRouter initializes the routes for the microservice
func (s *UDMService) SetupRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/subscription-data/{ueId}", s.GetSubscriptionDataHandler).Methods("GET")
	router.HandleFunc("/auth-vectors/{ueId}", s.GetAuthVectorHandler).Methods("GET")
	return router
}


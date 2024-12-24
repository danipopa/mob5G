package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var redisClient *redis.Client

// NFProfile represents a Network Function profile
type NFProfile struct {
    // Basic NF Information
    NFID         string   `json:"nf_id"`          // Unique NF identifier
    NFInstanceID string   `json:"nf_instance_id"` // Instance ID of the NF
    NFType       string   `json:"nf_type"`        // Type of NF (e.g., AMF, SMF, PCF)
    Status       string   `json:"status"`         // Registration status (e.g., REGISTERED, DEREGISTERED)
    FQDN         string   `json:"fqdn"`           // Fully Qualified Domain Name
    IPAddresses  []string `json:"ip_addresses"`   // List of IP addresses (IPv4/IPv6)
    ServiceURLs  []string `json:"service_urls"`   // URLs for NF services
    HeartbeatTimer int     `json:"heartbeat_timer"` // Heartbeat interval in seconds

    // PLMN Information
    PLMNID struct {                             // Public Land Mobile Network Identifier
        MNC string `json:"mnc"`                 // Mobile Network Code
        MCC string `json:"mcc"`                 // Mobile Country Code
    } `json:"plmn_id"`

    // Network Slicing
    SNssais []struct {                          // Slice/Service Type and Slice Differentiator
        SST string `json:"sst"`                 // Slice/Service Type (e.g., eMBB, URLLC)
        SD  string `json:"sd"`                  // Slice Differentiator
    } `json:"snssais"`

    // Area and DNN Information
    AreaID string   `json:"area_id"`            // Area ID for NF coverage
    DNNs   []string `json:"dnns"`               // List of supported Data Network Names

    // Security and Integrity
    SecurityFeatures []string `json:"security_features"` // List of supported security features
    IntegrityAlgorithms []string `json:"integrity_algorithms"` // Message integrity algorithms
    EncryptionAlgorithms []string `json:"encryption_algorithms"` // Encryption algorithms

    // NF Services
    NFServices []struct {                       // Services provided by the NF
        ServiceName     string   `json:"service_name"`     // Name of the service
        ServiceURLs     []string `json:"service_urls"`     // URLs for accessing the service
        ServiceStatus   string   `json:"service_status"`   // Status of the service (e.g., ACTIVE, INACTIVE)
        AllowedPlmns    []string `json:"allowed_plmns"`    // Allowed PLMNs for the service
        SupportedFeatures []string `json:"supported_features"` // Supported features for the service
    } `json:"nf_services"`

    // Load and Capacity
    Capacity       int    `json:"capacity"`     // NF capacity metric
    LoadLevelInfo  int    `json:"load_level"`   // Current load level (percentage)
    MaxCapacity    int    `json:"max_capacity"` // Maximum capacity of the NF

    // Discovery and Management
    RecoveryTime      string            `json:"recovery_time"`      // Last recovery timestamp
    AllowedPlmnList   []string          `json:"allowed_plmn_list"`  // List of allowed PLMNs
    PreferredNFs      []string          `json:"preferred_nfs"`      // List of preferred NFs for interactions
    ConfigurationInfo map[string]string `json:"configuration_info"` // Additional configuration details

    // Additional Features
    IWKPCOs           []string          `json:"iwk_pcos"`           // Interworking Packet Core Optimizations
    SupportedProtocols []string          `json:"supported_protocols"` // Protocols supported by the NF
    AdditionalInfo     map[string]string `json:"additional_info"`    // Vendor-specific extensions
}
// Initialize Redis connection
func init() {
	log.Println("Initializing Redis client...")
	redisClient = redis.NewClient(&redis.Options{
		Addr: "redis-service:6379",
	})
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Redis client initialized successfully.")
}

// Register a new NF
func RegisterNF(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling NF registration request...")
	var profile NFProfile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Log received profile
	log.Printf("Received NF Profile: %+v\n", profile)

	// Store NF profile in Redis
	data, err := json.Marshal(profile)
	if err != nil {
		log.Printf("Error marshaling NF profile: %v", err)
		http.Error(w, "Failed to process NF profile", http.StatusInternalServerError)
		return
	}

	if err := redisClient.Set(ctx, profile.NFID, data, 0).Err(); err != nil {
		log.Printf("Error saving NF profile to Redis: %v", err)
		http.Error(w, "Failed to save NF profile", http.StatusInternalServerError)
		return
	}

	log.Printf("NF %s registered successfully.", profile.NFID)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("NF registered successfully"))
}

// Discover NFs by type
func DiscoverNFs(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling NF discovery request...")
	nfType := r.URL.Query().Get("nf_type")
	log.Printf("NF type filter: %s\n", nfType)

	var results []NFProfile
	iter := redisClient.Scan(ctx, 0, "*", 0).Iterator()
	for iter.Next(ctx) {
		val, err := redisClient.Get(ctx, iter.Val()).Result()
		if err != nil {
			log.Printf("Error retrieving NF profile from Redis: %v", err)
			continue
		}

		var profile NFProfile
		if err := json.Unmarshal([]byte(val), &profile); err == nil {
			if nfType == "" || profile.NFType == nfType {
				results = append(results, profile)
			}
		} else {
			log.Printf("Error unmarshaling NF profile: %v", err)
		}
	}

	if err := iter.Err(); err != nil {
		log.Printf("Error scanning Redis: %v", err)
		http.Error(w, "Failed to discover NFs", http.StatusInternalServerError)
		return
	}

	log.Printf("Discovered %d NFs matching type '%s'.\n", len(results), nfType)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// Deregister an NF
func DeregisterNF(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling NF deregistration request...")
	vars := mux.Vars(r)
	nfID := vars["nf_id"]
	log.Printf("Deregistering NF ID: %s\n", nfID)

	if err := redisClient.Del(ctx, nfID).Err(); err != nil {
		log.Printf("Error deleting NF profile from Redis: %v", err)
		http.Error(w, "Failed to deregister NF", http.StatusInternalServerError)
		return
	}

	log.Printf("NF %s deregistered successfully.", nfID)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("NF deregistered successfully"))
}

// SetupRouter initializes the router and routes
func SetupRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/nrf/register", RegisterNF).Methods("POST")
	router.HandleFunc("/nrf/discover", DiscoverNFs).Methods("GET")
	router.HandleFunc("/nrf/deregister/{nf_id}", DeregisterNF).Methods("DELETE")
	return router
}


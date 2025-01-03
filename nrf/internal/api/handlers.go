package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time" // Import the time package for date and time operations

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

var ctx = context.Background()
var redisClient *redis.Client

// NFProfile represents a Network Function profile
type NFProfile struct {
	NFID              string            `json:"nf_id"`
	NFInstanceID      string            `json:"nf_instance_id"`
	NFType            string            `json:"nf_type"`
	Status            string            `json:"status"`
	FQDN              string            `json:"fqdn"`
	IPAddresses       []string          `json:"ip_addresses"`
	ServiceURLs       []string          `json:"service_urls"`
	HeartbeatTimer    int               `json:"heartbeat_timer"`
	PLMNID            map[string]string `json:"plmn_id"`
	SNssais           []map[string]string `json:"snssais"`
	AdditionalInfo    map[string]string `json:"additional_info"`
	Subscriptions     []string          `json:"subscriptions"`
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

// RegisterNF handles NF registration or replacement
func RegisterNF(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nfInstanceID := vars["nfInstanceID"]

	var profile NFProfile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if profile.NFID != nfInstanceID {
		http.Error(w, "Mismatch between URL and profile NFID", http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(profile)
	if err != nil {
		http.Error(w, "Failed to process profile", http.StatusInternalServerError)
		return
	}

	if err := redisClient.Set(ctx, nfInstanceID, data, 0).Err(); err != nil {
		http.Error(w, "Failed to save profile", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("NF registered successfully"))
}

// DiscoverNFs retrieves NFs based on filters
func DiscoverNFs(w http.ResponseWriter, r *http.Request) {
	nfType := r.URL.Query().Get("nf_type")
	iter := redisClient.Scan(ctx, 0, "*", 0).Iterator()

	var results []NFProfile
	for iter.Next(ctx) {
		val, err := redisClient.Get(ctx, iter.Val()).Result()
		if err != nil {
			continue
		}

		var profile NFProfile
		if err := json.Unmarshal([]byte(val), &profile); err == nil {
			if nfType == "" || profile.NFType == nfType {
				results = append(results, profile)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// DeregisterNF removes an NF profile
func DeregisterNF(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nfID := vars["nf_id"]

	if err := redisClient.Del(ctx, nfID).Err(); err != nil {
		http.Error(w, "Failed to deregister NF", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("NF deregistered successfully"))
}

// NFHeartBeat updates the heartbeat of an NF
func NFHeartBeat(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    nfInstanceID := vars["nfInstanceID"]

    existingProfileJSON, err := redisClient.Get(ctx, nfInstanceID).Result()
    if err != nil {
        http.Error(w, "NF not found", http.StatusNotFound)
        return
    }

    var profile NFProfile
    if err := json.Unmarshal([]byte(existingProfileJSON), &profile); err != nil {
        http.Error(w, "Failed to process profile", http.StatusInternalServerError)
        return
    }

    // Initialize AdditionalInfo if it is nil
    if profile.AdditionalInfo == nil {
        profile.AdditionalInfo = make(map[string]string)
    }

    // Use server-side timestamp for the heartbeat
    serverTime := time.Now().UTC().Format(time.RFC3339)
    profile.AdditionalInfo["last_heartbeat"] = serverTime

    data, err := json.Marshal(profile)
    if err != nil {
        http.Error(w, "Failed to update profile", http.StatusInternalServerError)
        return
    }

    if err := redisClient.Set(ctx, nfInstanceID, data, 0).Err(); err != nil {
        http.Error(w, "Failed to save profile", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}


// CreateSubscription adds a new subscription
func CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var subscription map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&subscription); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	subscriptionID := subscription["subscriptionId"].(string)
	data, err := json.Marshal(subscription)
	if err != nil {
		http.Error(w, "Failed to process subscription", http.StatusInternalServerError)
		return
	}

	if err := redisClient.Set(ctx, subscriptionID, data, 0).Err(); err != nil {
		http.Error(w, "Failed to save subscription", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Subscription created successfully"))
}

// NotifySubscribers sends notifications to subscribers
func NotifySubscribers(w http.ResponseWriter, r *http.Request) {
	var notification map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	iter := redisClient.Scan(ctx, 0, "subscription:*", 0).Iterator()
	for iter.Next(ctx) {
		subscriptionJSON, err := redisClient.Get(ctx, iter.Val()).Result()
		if err != nil {
			continue
		}

		var subscription map[string]interface{}
		if err := json.Unmarshal([]byte(subscriptionJSON), &subscription); err == nil {
			notificationURL := subscription["notificationUri"].(string)
			go sendNotification(notificationURL, notification)
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Notifications sent"))
}

func sendNotification(notificationURL string, notification map[string]interface{}) {
	data, _ := json.Marshal(notification)
	http.Post(notificationURL, "application/json", bytes.NewBuffer(data))
}

// SetupRouter initializes the router and routes
func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	// NF Management
	router.HandleFunc("/nnrf-nfm/v1/nf-instances/{nfInstanceID}", RegisterNF).Methods("PUT")
	router.HandleFunc("/nnrf-nfm/v1/nf-instances/{nfInstanceID}", NFHeartBeat).Methods("PATCH")
	router.HandleFunc("/nnrf-nfm/v1/nf-instances/{nf_id}", DeregisterNF).Methods("DELETE")

	// NF Discovery
	router.HandleFunc("/nnrf-disc/v1/nfs", DiscoverNFs).Methods("GET")

	// Subscription Management
	router.HandleFunc("/nnrf-sub/v1/subscriptions", CreateSubscription).Methods("POST")
	router.HandleFunc("/nnrf-notify/v1/notifications", NotifySubscribers).Methods("POST")

	return router
}


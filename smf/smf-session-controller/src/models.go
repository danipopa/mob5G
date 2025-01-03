package src

// Session represents a PDU session
type Session struct {
	SessionID   string `json:"session_id"`
	UEID        string `json:"ue_id"`
	DNN         string `json:"dnn"`
	Slice       string `json:"slice"`
	QoSProfile  string `json:"qos_profile"`
	UPFAddress  string `json:"upf_address"`
	PCFPolicyID string `json:"pcf_policy_id"`
}

// SessionRequest represents a request for session creation or modification
type SessionRequest struct {
	UEID       string `json:"ue_id"`
	DNN        string `json:"dnn"`
	Slice      string `json:"slice"`
	QoSProfile string `json:"qos_profile"`
}


// AMFNotification represents a notification sent to the AMF via SMF-N11
type AMFNotification struct {
	SessionID  string `json:"session_id"`  // Unique ID of the session
	Event      string `json:"event"`       // Event type (e.g., "SessionEstablished", "SessionReleased")
	Details    string `json:"details"`     // Additional details about the event
	Subscriber string `json:"subscriber"` // Subscriber ID (e.g., UEID)
}

// AMFResponse represents a response received from the AMF
type AMFResponse struct {
	SessionID string `json:"session_id"` // Unique ID of the session
	Status    string `json:"status"`     // Status of the request (e.g., "SUCCESS", "FAILED")
	Message   string `json:"message"`    // Additional message details
}

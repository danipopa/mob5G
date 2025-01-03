//Shared Data Models (TS 29.503)

package smfn11

// PDURequest represents a request for PDU session management
type PDURequest struct {
	UEID       string `json:"ue_id"`        // User Equipment ID
	DNN        string `json:"dnn"`          // Data Network Name
	SST        string `json:"sst"`          // Slice/Service Type
	SD         string `json:"sd"`           // Slice Differentiator
	SessionID  string `json:"session_id"`   // PDU Session ID
	QoSProfile string `json:"qos_profile"`  // Quality of Service Profile
}

// PDUResponse represents a response for PDU session management
type PDUResponse struct {
	SessionID string `json:"session_id"` // PDU Session ID
	Status    string `json:"status"`     // Status of the operation
	Message   string `json:"message"`    // Additional information
}

// PolicyData represents shared policy information
type PolicyData struct {
	QoSProfile string `json:"qos_profile"` // Quality of Service Profile
}

// GenericRequest represents a basic notification or request format
type GenericRequest struct {
	Type          string `json:"type"`          // Notification or request type
	Payload       string `json:"payload"`       // Serialized payload
	TransactionID string `json:"transaction_id"` // Transaction ID for tracking
}

package udr

// Subscription model
type Subscription struct {
	UEID       string `json:"ue_id"`       // User Equipment Identifier
	AccessType string `json:"access_type"` // e.g., 4G, 5G
	QoS        string `json:"qos"`        // QoS Profile
}

// Policy model
type Policy struct {
	PolicyID   string `json:"policy_id"`
	UEID       string `json:"ue_id"`
	QoSPriority int    `json:"qos_priority"`
	Action      string `json:"action"` // e.g., ALLOW, DENY
}

// Session model
type Session struct {
	SessionID string `json:"session_id"`
	UEID      string `json:"ue_id"`
	IPAddress string `json:"ip_address"`
}


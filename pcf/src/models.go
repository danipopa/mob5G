package pcf

type Policy struct {
	ID          string `json:"id"`
	Description string `json:"description"` // Example field for policy description
	QoSPriority int    `json:"qos_priority"`
}

type AMPolicyRequest struct {
	UEID       string `json:"ue_id"`       // User Equipment Identifier
	AccessPoint string `json:"access_point"` // Target access point
	MobilityType string `json:"mobility_type"` // e.g., HANDOVER, INITIAL_ATTACH
}

type AMPolicyResponse struct {
	PolicyID       string `json:"policy_id"`       // Policy Identifier
	AccessAllowed  bool   `json:"access_allowed"`  // Whether access is allowed
	HandoverAllowed bool  `json:"handover_allowed"` // Whether handover is allowed
	QoSPriority     int    `json:"qos_priority"`    // QoS Priority Level
	Message         string `json:"message"`        // Additional information
}


//Shared Data Models for Policies, Charging, and QoS (TS 29.503)


package smfn7

// SMPolicyContextData represents the data sent from SMF to PCF
type SMPolicyContextData struct {
	Supi         string  `json:"supi"`          // Subscription Permanent Identifier
	DNN          string  `json:"dnn"`           // Data Network Name
	PDUSessionID int     `json:"pdu_session_id"`// PDU Session Identifier
	SliceInfo    Slice   `json:"slice_info"`    // Slice/Service information
	QoSProfile   QoS     `json:"qos_profile"`   // QoS requirements
	UELocation   Location `json:"ue_location"`  // UE location information
}

// SMPolicyDecision represents the policies sent from PCF to SMF
type SMPolicyDecision struct {
	PolicyID     string         `json:"policy_id"`      // Policy Identifier
	QoSRules     []QoSRule      `json:"qos_rules"`      // QoS rules
	ChargingRule ChargingPolicy `json:"charging_policy"`// Charging rules
}

// Slice represents slice information for the session
type Slice struct {
	SST int    `json:"sst"` // Slice/Service Type
	SD  string `json:"sd"`  // Slice Differentiator
}

// QoS represents Quality of Service requirements
type QoS struct {
	FiveQI int `json:"5qi"` // 5G QoS Identifier
	ARP    int `json:"arp"` // Allocation and Retention Priority
}

// Location represents UE location information
type Location struct {
	Latitude  float64 `json:"latitude"`  // Latitude
	Longitude float64 `json:"longitude"` // Longitude
}

// QoSRule represents a single QoS rule
type QoSRule struct {
	RuleID         string `json:"rule_id"`         // Rule Identifier
	FlowDescription string `json:"flow_description"`// Description of flow
	QoSParameters  QoS    `json:"qos_parameters"`  // QoS parameters
}

// ChargingPolicy represents the charging policy
type ChargingPolicy struct {
	Online bool `json:"online"` // Online charging enabled
	Quota  int  `json:"quota"`  // Data quota in MB
}


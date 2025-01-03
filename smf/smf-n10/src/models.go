package smfn10

// SubscriptionData represents the subscription data retrieved from the UDM
type SubscriptionData struct {
	UEID           string   `json:"ue_id"`            // User Equipment ID
	AllowedNetworks []string `json:"allowed_networks"` // List of allowed DNNs (Data Network Names)
	QoSProfile     string   `json:"qos_profile"`      // Default QoS Profile
	NSSAI          []NSSAI  `json:"nssai"`           // Allowed Network Slice Selection Assistance Information
}

// NSSAI represents a single Network Slice Selection Assistance Information entry
type NSSAI struct {
	SST string `json:"sst"` // Slice/Service Type
	SD  string `json:"sd"`  // Slice Differentiator
}


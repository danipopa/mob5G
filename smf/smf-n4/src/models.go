package src

// PFCPHeader represents the PFCP message header
type PFCPHeader struct {
	Version        uint8
	MessageType    uint8
	MessageLength  uint16
	SequenceNumber uint32
}

// PDR (Packet Detection Rule)
type PDR struct {
	RuleID         uint32
	MatchCriteria  string
	QoSParameters  QoS
	ForwardingRule FAR
}

// FAR (Forwarding Action Rule)
type FAR struct {
	Action      string
	Destination string
}

// QoS represents Quality of Service parameters
type QoS struct {
	FiveQI  uint8
	GBR     uint32
	MBR     uint32
	Priority uint8
}

// UsageReport represents usage data sent from UPF to SMF
type UsageReport struct {
	SessionID  string `json:"session_id"`
	VolumeMB   uint64 `json:"volume_mb"`    // Data volume in MB
	DurationMS uint64 `json:"duration_ms"` // Duration in milliseconds
}

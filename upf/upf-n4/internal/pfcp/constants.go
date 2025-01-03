package pfcp

// PFCP Message Types
const (
	PFCPAssociationSetupRequest        uint8 = 1
	PFCPAssociationSetupResponse       uint8 = 2
	PFCPAssociationUpdateRequest       uint8 = 3
	PFCPAssociationUpdateResponse      uint8 = 4
	PFCPAssociationReleaseRequest      uint8 = 5
	PFCPAssociationReleaseResponse     uint8 = 6
	PFCPSessionEstablishmentRequest    uint8 = 50
	PFCPSessionEstablishmentResponse   uint8 = 51
	PFCPSessionModificationRequest     uint8 = 52
	PFCPSessionModificationResponse    uint8 = 53
	PFCPSessionDeletionRequest         uint8 = 54
	PFCPSessionDeletionResponse        uint8 = 55
	PFCPHeartbeatRequest               uint8 = 56
	PFCPHeartbeatResponse              uint8 = 57
)

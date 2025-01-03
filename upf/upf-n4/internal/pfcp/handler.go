package pfcp

import (
	"log"
	"encoding/json"
)

// Association data structure
type Association struct {
	NodeID   string `json:"node_id"`
	NodeType string `json:"node_type"`
}

// HandleMessage processes incoming PFCP messages.
func HandleMessage(data []byte, addr string) {
	msg, err := DeserializePFCPMessage(data)
	if err != nil {
		log.Printf("Error parsing PFCP message from %s: %v", addr, err)
		return
	}

	switch msg.MessageType {
	case PFCPAssociationSetupRequest:
		handleAssociationSetupRequest(msg, addr)
	case PFCPAssociationUpdateRequest:
		handleAssociationUpdateRequest(msg, addr)
	case PFCPAssociationReleaseRequest:
		handleAssociationReleaseRequest(msg, addr)
	case PFCPSessionEstablishmentRequest:
		handleSessionEstablishmentRequest(msg, addr)
	case PFCPSessionModificationRequest:
		handleSessionModificationRequest(msg, addr)
	case PFCPSessionDeletionRequest:
		handleSessionDeletionRequest(msg, addr)
	case PFCPHeartbeatRequest:
		handleHeartbeatRequest(msg, addr)
	default:
		log.Printf("Unknown PFCP message type %d from %s", msg.MessageType, addr)
	}
}

// Handlers for each message type

func handleAssociationSetupRequest(msg *PFCPMessage, addr string) {
	log.Printf("Handling PFCP Association Setup Request from %s", addr)

	// Parse Node ID from the payload (simplified example)
	nodeID := string(msg.Payload[:8]) // Assume Node ID is 8 bytes
	association := Association{
		NodeID:   nodeID,
		NodeType: "UPF",
	}

	// Save to Redis
	data, _ := json.Marshal(association)
	err := SaveAssociation(nodeID, string(data))
	if err != nil {
		log.Printf("Failed to save association: %v", err)
	}

	// Respond with PFCP Association Setup Response
	response := CreateAssociationSetupResponse(msg.SequenceNumber)
	sendResponse(response, addr)
}

func handleAssociationUpdateRequest(msg *PFCPMessage, addr string) {
	log.Printf("Handling PFCP Association Update Request from %s", addr)

	// Example: Update association in Redis
	nodeID := string(msg.Payload[:8]) // Assume Node ID is 8 bytes
	_, err := GetAssociation(nodeID)
	if err != nil {
		log.Printf("No existing association found for Node ID %s", nodeID)
		return
	}

	// Update logic (this is simplified, real implementation may parse new data)
	newData := Association{
		NodeID:   nodeID,
		NodeType: "Updated-UPF",
	}
	data, _ := json.Marshal(newData)
	SaveAssociation(nodeID, string(data))

	log.Printf("Updated association for Node ID %s", nodeID)
}

func handleAssociationReleaseRequest(msg *PFCPMessage, addr string) {
	log.Printf("Handling PFCP Association Release Request from %s", addr)

	// Example: Delete association from Redis
	nodeID := string(msg.Payload[:8]) // Assume Node ID is 8 bytes
	err := DeleteAssociation(nodeID)
	if err != nil {
		log.Printf("Failed to delete association for Node ID %s", nodeID)
		return
	}

	log.Printf("Deleted association for Node ID %s", nodeID)
}

func handleSessionEstablishmentRequest(msg *PFCPMessage, addr string) {
	log.Printf("Handling PFCP Session Establishment Request from %s", addr)
	// TODO: Implement session establishment logic
}

func handleSessionModificationRequest(msg *PFCPMessage, addr string) {
	log.Printf("Handling PFCP Session Modification Request from %s", addr)
	// TODO: Implement session modification logic
}

func handleSessionDeletionRequest(msg *PFCPMessage, addr string) {
	log.Printf("Handling PFCP Session Deletion Request from %s", addr)
	// TODO: Implement session deletion logic
}

func handleHeartbeatRequest(msg *PFCPMessage, addr string) {
	log.Printf("Handling PFCP Heartbeat Request from %s", addr)
	response := CreateHeartbeatResponse(msg.SequenceNumber)
	sendResponse(response, addr)
}


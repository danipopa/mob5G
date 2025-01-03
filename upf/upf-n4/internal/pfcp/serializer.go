package pfcp

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// PFCPMessage represents a PFCP message.
type PFCPMessage struct {
	Version       uint8
	MessageType   uint8
	SequenceNumber uint32
	MessageLength uint16
	Payload       []byte
}

// SerializePFCPMessage serializes a PFCP message into a byte slice.
func SerializePFCPMessage(msg *PFCPMessage) ([]byte, error) {
	buf := new(bytes.Buffer)

	// Serialize header
	buf.WriteByte(msg.Version)
	buf.WriteByte(msg.MessageType)
	binary.Write(buf, binary.BigEndian, msg.SequenceNumber)
	binary.Write(buf, binary.BigEndian, msg.MessageLength)

	// Serialize payload
	if len(msg.Payload) > 0 {
		buf.Write(msg.Payload)
	}

	return buf.Bytes(), nil
}

// DeserializePFCPMessage deserializes a byte slice into a PFCPMessage.
func DeserializePFCPMessage(data []byte) (*PFCPMessage, error) {
	if len(data) < 8 {
		return nil, errors.New("data too short for PFCP message")
	}

	buf := bytes.NewReader(data)
	msg := &PFCPMessage{}

	// Read header
	binary.Read(buf, binary.BigEndian, &msg.Version)
	binary.Read(buf, binary.BigEndian, &msg.MessageType)
	binary.Read(buf, binary.BigEndian, &msg.SequenceNumber)
	binary.Read(buf, binary.BigEndian, &msg.MessageLength)

	// Read payload
	msg.Payload = make([]byte, buf.Len())
	buf.Read(msg.Payload)

	return msg, nil
}


func CreateAssociationSetupResponse(sequenceNumber uint32) *PFCPMessage {
	return &PFCPMessage{
		Version:       1,
		MessageType:   PFCPAssociationSetupResponse,
		SequenceNumber: sequenceNumber,
		MessageLength: uint16(8), // Example length
		Payload:       []byte{0x01, 0x02}, // Dummy payload
	}
}

func CreateHeartbeatResponse(sequenceNumber uint32) *PFCPMessage {
	return &PFCPMessage{
		Version:       1,
		MessageType:   PFCPHeartbeatResponse,
		SequenceNumber: sequenceNumber,
		MessageLength: uint16(8),
		Payload:       []byte{},
	}
}

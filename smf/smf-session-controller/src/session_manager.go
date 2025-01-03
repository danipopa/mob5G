package src

import (
	"fmt"
)

// SessionManager handles session lifecycle management
type SessionManager struct {
	redisClient *RedisClient
	n4Client    *N4Client
	n10Client   *N10Client
	n7Client    *N7Client
	n11Client   *N11Client
}

// NewSessionManager initializes a new SessionManager
func NewSessionManager(redisClient *RedisClient) *SessionManager {
	return &SessionManager{
		redisClient: redisClient,
		n4Client:    NewN4Client("http://smf-n4-service:8084"),
		n10Client:   NewN10Client("http://smf-n10-service:8086"),
		n7Client:    NewN7Client("http://smf-n7-service:8087"),
		n11Client:   NewN11Client("http://smf-n11-service:8086"),
	}
}

// CreateSession handles session creation
func (m *SessionManager) CreateSession(req *SessionRequest) (*Session, error) {
	// Step 1: Create session context
	session := &Session{
		SessionID:   fmt.Sprintf("session-%s-%s", req.UEID, req.DNN),
		UEID:        req.UEID,
		DNN:         req.DNN,
		Slice:       req.Slice,
		QoSProfile:  req.QoSProfile,
		UPFAddress:  "upf-placeholder",
		PCFPolicyID: "pcf-placeholder",
	}

	// Step 2: Interact with SMF-N10 for subscription data
	if err := m.n10Client.FetchSubscriptionData(req.UEID); err != nil {
		return nil, fmt.Errorf("failed to fetch subscription data: %w", err)
	}

	// Step 3: Interact with SMF-N7 for policies
	if err := m.n7Client.FetchPolicyData(session); err != nil {
		return nil, fmt.Errorf("failed to fetch policy data: %w", err)
	}

	// Step 4: Interact with SMF-N4 for user plane setup
	if err := m.n4Client.CreateSessionInUPF(session); err != nil {
		return nil, fmt.Errorf("failed to create session in UPF: %w", err)
	}

	// Step 5: Save session to Redis
	if err := m.redisClient.SaveSession(session.SessionID, session); err != nil {
		return nil, fmt.Errorf("failed to save session: %w", err)
	}

	return session, nil
}

// ModifySession and DeleteSession can be implemented similarly.

// ModifySession handles session modification
func (m *SessionManager) ModifySession(sessionID string, req *SessionRequest) (*Session, error) {
	// Step 1: Retrieve existing session from Redis
	session, err := m.redisClient.GetSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve session: %w", err)
	}

	// Step 2: Update session context with new data
	session.QoSProfile = req.QoSProfile
	session.Slice = req.Slice

	// Step 3: Interact with SMF-N7 for updated policies
	if err := m.n7Client.FetchPolicyData(session); err != nil {
		return nil, fmt.Errorf("failed to fetch policy data: %w", err)
	}

	// Step 4: Interact with SMF-N4 for updating user plane resources
	if err := m.n4Client.ModifySessionInUPF(session); err != nil {
		return nil, fmt.Errorf("failed to modify session in UPF: %w", err)
	}

	// Step 5: Save updated session to Redis
	if err := m.redisClient.SaveSession(sessionID, session); err != nil {
		return nil, fmt.Errorf("failed to save updated session: %w", err)
	}

	return session, nil
}

// DeleteSession handles session deletion
func (m *SessionManager) DeleteSession(sessionID string) error {
	// Step 1: Retrieve session from Redis to ensure it exists
	session, err := m.redisClient.GetSession(sessionID)
	if err != nil {
		return fmt.Errorf("failed to retrieve session: %w", err)
	}

	// Step 2: Interact with SMF-N4 to release user plane resources
	if err := m.n4Client.ReleaseSessionInUPF(sessionID); err != nil {
		return fmt.Errorf("failed to release session in UPF: %w", err)
	}

	// Step 3: Delete session from Redis
	if err := m.redisClient.DeleteSession(sessionID); err != nil {
		return fmt.Errorf("failed to delete session from Redis: %w", err)
	}

	// Step 4: Notify AMF about session deletion
	notification := &AMFNotification{
		SessionID:  sessionID,
		Event:      "SessionReleased",
		Details:    "Session successfully deleted",
		Subscriber: session.UEID,
	}

	if err := m.n11Client.NotifyAMF(notification); err != nil {
		return fmt.Errorf("failed to notify AMF: %w", err)
	}

	return nil
}

// Example integration of N11Client in SessionManager
func (m *SessionManager) NotifyAMFOnSessionCreation(session *Session) error {
	notification := &AMFNotification{
		SessionID:  session.SessionID,
		Event:      "SessionEstablished",
		Details:    "Session successfully created",
		Subscriber: session.UEID,
	}

	err := m.n11Client.NotifyAMF(notification)
	if err != nil {
		return fmt.Errorf("failed to notify AMF: %w", err)
	}

	return nil
}

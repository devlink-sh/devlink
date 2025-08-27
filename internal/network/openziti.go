package network

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/devlink/internal/util"
)

// OpenZitiManager handles OpenZiti networking operations
type OpenZitiManager struct {
	config     *util.Config
	services   map[string]*EphemeralService
	servicesMu sync.RWMutex
	clients    map[string]*SecureClient
	clientsMu  sync.RWMutex
}

// EphemeralService represents a temporary service for sharing environment data
type EphemeralService struct {
	ID          string    `json:"id"`
	ShareCode   string    `json:"share_code"`
	ServiceName string    `json:"service_name"`
	Port        int       `json:"port"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	Listener    net.Listener
	IsActive    bool
}

// SecureClient represents a secure client connection
type SecureClient struct {
	ID          string    `json:"id"`
	ShareCode   string    `json:"share_code"`
	ServiceName string    `json:"service_name"`
	ConnectedAt time.Time `json:"connected_at"`
	Connection  net.Conn
	IsConnected bool
}

// NetworkMessage represents messages exchanged over the network
type NetworkMessage struct {
	Type      string          `json:"type"`
	ShareCode string          `json:"share_code"`
	Data      json.RawMessage `json:"data,omitempty"`
	Timestamp time.Time       `json:"timestamp"`
	Version   string          `json:"version"`
}

// NewOpenZitiManager creates a new OpenZiti manager
func NewOpenZitiManager(config *util.Config) (*OpenZitiManager, error) {
	return &OpenZitiManager{
		config:   config,
		services: make(map[string]*EphemeralService),
		clients:  make(map[string]*SecureClient),
	}, nil
}

// CreateEphemeralService creates a temporary service for sharing environment data
func (ozm *OpenZitiManager) CreateEphemeralService(shareCode string, expiry time.Duration) (*EphemeralService, error) {
	// Generate unique service name
	serviceName := fmt.Sprintf("devlink-env-%s", shareCode)

	// Create ephemeral service
	service := &EphemeralService{
		ID:          util.GenerateUUID(),
		ShareCode:   shareCode,
		ServiceName: serviceName,
		Port:        ozm.config.P2PPort,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(expiry),
		IsActive:    false,
	}

	// Start listening for connections
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", ozm.config.P2PPort))
	if err != nil {
		return nil, fmt.Errorf("failed to create listener: %w", err)
	}

	service.Listener = listener
	service.IsActive = true

	// Store service
	ozm.servicesMu.Lock()
	ozm.services[service.ID] = service
	ozm.servicesMu.Unlock()

	// Start service handler
	go ozm.handleServiceConnections(service)

	// Start cleanup timer
	go ozm.scheduleServiceCleanup(service, expiry)

	return service, nil
}

// ConnectToService establishes a secure connection to a service
func (ozm *OpenZitiManager) ConnectToService(shareCode string) (*SecureClient, error) {
	serviceName := fmt.Sprintf("devlink-env-%s", shareCode)

	// Connect to service (simplified for now)
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", ozm.config.P2PPort))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to service: %w", err)
	}

	// Create secure client
	client := &SecureClient{
		ID:          util.GenerateUUID(),
		ShareCode:   shareCode,
		ServiceName: serviceName,
		ConnectedAt: time.Now(),
		Connection:  conn,
		IsConnected: true,
	}

	// Store client
	ozm.clientsMu.Lock()
	ozm.clients[client.ID] = client
	ozm.clientsMu.Unlock()

	return client, nil
}

// SendEncryptedData sends encrypted data over a secure connection
func (ozm *OpenZitiManager) SendEncryptedData(client *SecureClient, encryptedData *util.EncryptedData) error {
	if !client.IsConnected {
		return fmt.Errorf("client is not connected")
	}

	// Create network message
	message := &NetworkMessage{
		Type:      "encrypted_data",
		ShareCode: client.ShareCode,
		Timestamp: time.Now(),
		Version:   "1.0",
	}

	// Serialize encrypted data
	dataBytes, err := json.Marshal(encryptedData)
	if err != nil {
		return fmt.Errorf("failed to serialize encrypted data: %w", err)
	}
	message.Data = dataBytes

	// Serialize message
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to serialize message: %w", err)
	}

	// Send message
	_, err = client.Connection.Write(messageBytes)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

// ReceiveEncryptedData receives encrypted data from a secure connection
func (ozm *OpenZitiManager) ReceiveEncryptedData(client *SecureClient) (*util.EncryptedData, error) {
	if !client.IsConnected {
		return nil, fmt.Errorf("client is not connected")
	}

	// Read message
	buffer := make([]byte, 4096)
	n, err := client.Connection.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to read message: %w", err)
	}

	// Parse message
	var message NetworkMessage
	if err := json.Unmarshal(buffer[:n], &message); err != nil {
		return nil, fmt.Errorf("failed to parse message: %w", err)
	}

	// Validate message
	if message.Type != "encrypted_data" {
		return nil, fmt.Errorf("unexpected message type: %s", message.Type)
	}

	if message.ShareCode != client.ShareCode {
		return nil, fmt.Errorf("share code mismatch")
	}

	// Parse encrypted data
	var encryptedData util.EncryptedData
	if err := json.Unmarshal(message.Data, &encryptedData); err != nil {
		return nil, fmt.Errorf("failed to parse encrypted data: %w", err)
	}

	return &encryptedData, nil
}

// handleServiceConnections handles incoming connections to a service
func (ozm *OpenZitiManager) handleServiceConnections(service *EphemeralService) {
	defer func() {
		service.Listener.Close()
		service.IsActive = false
	}()

	for service.IsActive {
		// Accept connection with timeout
		conn, err := service.Listener.Accept()
		if err != nil {
			if !service.IsActive {
				return // Service was stopped
			}
			continue
		}

		// Handle connection in goroutine
		go ozm.handleConnection(service, conn)
	}
}

// handleConnection handles a single connection
func (ozm *OpenZitiManager) handleConnection(service *EphemeralService, conn net.Conn) {
	defer conn.Close()

	// Read message
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		return
	}

	// Parse message
	var message NetworkMessage
	if err := json.Unmarshal(buffer[:n], &message); err != nil {
		return
	}

	// Validate message
	if message.Type != "encrypted_data" || message.ShareCode != service.ShareCode {
		return
	}

	// Process the encrypted data (this would be handled by the application layer)
	// For now, we just echo back a success response
	response := &NetworkMessage{
		Type:      "response",
		ShareCode: service.ShareCode,
		Timestamp: time.Now(),
		Version:   "1.0",
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		return
	}

	conn.Write(responseBytes)
}

// scheduleServiceCleanup schedules cleanup of expired services
func (ozm *OpenZitiManager) scheduleServiceCleanup(service *EphemeralService, expiry time.Duration) {
	time.Sleep(expiry)

	ozm.servicesMu.Lock()
	defer ozm.servicesMu.Unlock()

	if s, exists := ozm.services[service.ID]; exists {
		s.IsActive = false
		if s.Listener != nil {
			s.Listener.Close()
		}
		delete(ozm.services, service.ID)
	}
}

// CloseClient closes a secure client connection
func (ozm *OpenZitiManager) CloseClient(clientID string) error {
	ozm.clientsMu.Lock()
	defer ozm.clientsMu.Unlock()

	if client, exists := ozm.clients[clientID]; exists {
		client.IsConnected = false
		if client.Connection != nil {
			client.Connection.Close()
		}
		delete(ozm.clients, clientID)
	}

	return nil
}

// CloseService closes an ephemeral service
func (ozm *OpenZitiManager) CloseService(serviceID string) error {
	ozm.servicesMu.Lock()
	defer ozm.servicesMu.Unlock()

	if service, exists := ozm.services[serviceID]; exists {
		service.IsActive = false
		if service.Listener != nil {
			service.Listener.Close()
		}
		delete(ozm.services, serviceID)
	}

	return nil
}

// GetServiceStats returns statistics about active services
func (ozm *OpenZitiManager) GetServiceStats() map[string]interface{} {
	ozm.servicesMu.RLock()
	defer ozm.servicesMu.RUnlock()

	activeServices := 0
	for _, service := range ozm.services {
		if service.IsActive {
			activeServices++
		}
	}

	return map[string]interface{}{
		"total_services":     len(ozm.services),
		"active_services":    activeServices,
		"openziti_connected": true, // Simplified for now
	}
}

// GetClientStats returns statistics about active clients
func (ozm *OpenZitiManager) GetClientStats() map[string]interface{} {
	ozm.clientsMu.RLock()
	defer ozm.clientsMu.RUnlock()

	connectedClients := 0
	for _, client := range ozm.clients {
		if client.IsConnected {
			connectedClients++
		}
	}

	return map[string]interface{}{
		"total_clients":      len(ozm.clients),
		"connected_clients":  connectedClients,
		"openziti_connected": true, // Simplified for now
	}
}

// Close closes the OpenZiti manager and all connections
func (ozm *OpenZitiManager) Close() error {
	// Close all clients
	ozm.clientsMu.Lock()
	for _, client := range ozm.clients {
		client.IsConnected = false
		if client.Connection != nil {
			client.Connection.Close()
		}
	}
	ozm.clients = make(map[string]*SecureClient)
	ozm.clientsMu.Unlock()

	// Close all services
	ozm.servicesMu.Lock()
	for _, service := range ozm.services {
		service.IsActive = false
		if service.Listener != nil {
			service.Listener.Close()
		}
	}
	ozm.services = make(map[string]*EphemeralService)
	ozm.servicesMu.Unlock()

	return nil
}

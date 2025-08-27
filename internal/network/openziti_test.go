package network

import (
	"testing"
	"time"

	"github.com/devlink/internal/util"
)

func TestNewOpenZitiManager(t *testing.T) {
	config := util.DefaultConfig()
	ozm, err := NewOpenZitiManager(config)

	if err != nil {
		t.Fatalf("NewOpenZitiManager failed: %v", err)
	}

	if ozm == nil {
		t.Fatal("NewOpenZitiManager() returned nil")
	}

	if ozm.config == nil {
		t.Error("OpenZitiManager should have config")
	}

	if ozm.services == nil {
		t.Error("OpenZitiManager should have services map")
	}

	if ozm.clients == nil {
		t.Error("OpenZitiManager should have clients map")
	}
}

func TestCreateEphemeralService(t *testing.T) {
	config := util.DefaultConfig()
	ozm, err := NewOpenZitiManager(config)
	if err != nil {
		t.Fatalf("Failed to create OpenZiti manager: %v", err)
	}
	defer ozm.Close()

	shareCode := "blue-whale-42"
	expiry := 1 * time.Hour

	service, err := ozm.CreateEphemeralService(shareCode, expiry)
	if err != nil {
		t.Fatalf("CreateEphemeralService failed: %v", err)
	}

	if service == nil {
		t.Fatal("CreateEphemeralService returned nil")
	}

	if service.ID == "" {
		t.Error("Service should have an ID")
	}

	if service.ShareCode != shareCode {
		t.Errorf("Service share code mismatch. Expected: %s, Got: %s", shareCode, service.ShareCode)
	}

	if !service.IsActive {
		t.Error("Service should be active")
	}

	if service.Listener == nil {
		t.Error("Service should have a listener")
	}

	// Check that service is stored
	stats := ozm.GetServiceStats()
	if stats["active_services"].(int) != 1 {
		t.Error("Service should be counted as active")
	}
}

func TestConnectToService(t *testing.T) {
	config := util.DefaultConfig()
	ozm, err := NewOpenZitiManager(config)
	if err != nil {
		t.Fatalf("Failed to create OpenZiti manager: %v", err)
	}
	defer ozm.Close()

	shareCode := "red-dragon-123"
	expiry := 1 * time.Hour

	// Create service first
	service, err := ozm.CreateEphemeralService(shareCode, expiry)
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}

	// Connect to service
	client, err := ozm.ConnectToService(shareCode)
	if err != nil {
		t.Fatalf("ConnectToService failed: %v", err)
	}

	if client == nil {
		t.Fatal("ConnectToService returned nil")
	}

	if client.ID == "" {
		t.Error("Client should have an ID")
	}

	if client.ShareCode != shareCode {
		t.Errorf("Client share code mismatch. Expected: %s, Got: %s", shareCode, client.ShareCode)
	}

	if !client.IsConnected {
		t.Error("Client should be connected")
	}

	if client.Connection == nil {
		t.Error("Client should have a connection")
	}

	// Check that client is stored
	stats := ozm.GetClientStats()
	if stats["connected_clients"].(int) != 1 {
		t.Error("Client should be counted as connected")
	}

	// Clean up
	ozm.CloseClient(client.ID)
	ozm.CloseService(service.ID)
}

func TestSendReceiveEncryptedData(t *testing.T) {
	config := util.DefaultConfig()
	ozm, err := NewOpenZitiManager(config)
	if err != nil {
		t.Fatalf("Failed to create OpenZiti manager: %v", err)
	}
	defer ozm.Close()

	shareCode := "green-forest-999"
	expiry := 1 * time.Hour

	// Create service
	_, err = ozm.CreateEphemeralService(shareCode, expiry)
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}

	// Connect to service
	client, err := ozm.ConnectToService(shareCode)
	if err != nil {
		t.Fatalf("Failed to connect to service: %v", err)
	}

	// Create test encrypted data
	encryptedData := &util.EncryptedData{
		Data:      "test-encrypted-data",
		Nonce:     "test-nonce",
		Salt:      "test-salt",
		Version:   "1.0",
		Algorithm: "AES-256-GCM",
	}

	// Send encrypted data
	err = ozm.SendEncryptedData(client, encryptedData)
	if err != nil {
		t.Fatalf("SendEncryptedData failed: %v", err)
	}

	// Clean up
	ozm.CloseClient(client.ID)
}

func TestServiceCleanup(t *testing.T) {
	config := util.DefaultConfig()
	ozm, err := NewOpenZitiManager(config)
	if err != nil {
		t.Fatalf("Failed to create OpenZiti manager: %v", err)
	}
	defer ozm.Close()

	shareCode := "yellow-sun-777"
	expiry := 100 * time.Millisecond // Short expiry for testing

	// Create service
	_, err = ozm.CreateEphemeralService(shareCode, expiry)
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}

	// Check that service is active
	stats := ozm.GetServiceStats()
	if stats["active_services"].(int) != 1 {
		t.Error("Service should be active")
	}

	// Wait for cleanup
	time.Sleep(200 * time.Millisecond)

	// Check that service is cleaned up
	stats = ozm.GetServiceStats()
	if stats["active_services"].(int) != 0 {
		t.Error("Service should be cleaned up")
	}

	if stats["total_services"].(int) != 0 {
		t.Error("Service should be removed from map")
	}
}

func TestCloseClient(t *testing.T) {
	config := util.DefaultConfig()
	ozm, err := NewOpenZitiManager(config)
	if err != nil {
		t.Fatalf("Failed to create OpenZiti manager: %v", err)
	}
	defer ozm.Close()

	shareCode := "purple-moon-555"
	expiry := 1 * time.Hour

	// Create service and client
	service, err := ozm.CreateEphemeralService(shareCode, expiry)
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}

	client, err := ozm.ConnectToService(shareCode)
	if err != nil {
		t.Fatalf("Failed to connect to service: %v", err)
	}

	// Check that client is connected
	stats := ozm.GetClientStats()
	if stats["connected_clients"].(int) != 1 {
		t.Error("Client should be connected")
	}

	// Close client
	err = ozm.CloseClient(client.ID)
	if err != nil {
		t.Fatalf("CloseClient failed: %v", err)
	}

	// Check that client is closed
	stats = ozm.GetClientStats()
	if stats["connected_clients"].(int) != 0 {
		t.Error("Client should be closed")
	}

	// Clean up
	ozm.CloseService(service.ID)
}

func TestCloseService(t *testing.T) {
	config := util.DefaultConfig()
	ozm, err := NewOpenZitiManager(config)
	if err != nil {
		t.Fatalf("Failed to create OpenZiti manager: %v", err)
	}
	defer ozm.Close()

	shareCode := "orange-star-333"
	expiry := 1 * time.Hour

	// Create service
	service, err := ozm.CreateEphemeralService(shareCode, expiry)
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}

	// Check that service is active
	stats := ozm.GetServiceStats()
	if stats["active_services"].(int) != 1 {
		t.Error("Service should be active")
	}

	// Close service
	err = ozm.CloseService(service.ID)
	if err != nil {
		t.Fatalf("CloseService failed: %v", err)
	}

	// Check that service is closed
	stats = ozm.GetServiceStats()
	if stats["active_services"].(int) != 0 {
		t.Error("Service should be closed")
	}
}

func TestGetServiceStats(t *testing.T) {
	config := util.DefaultConfig()
	ozm, err := NewOpenZitiManager(config)
	if err != nil {
		t.Fatalf("Failed to create OpenZiti manager: %v", err)
	}
	defer ozm.Close()

	// Check initial stats
	stats := ozm.GetServiceStats()

	expectedFields := []string{"total_services", "active_services", "openziti_connected"}
	for _, field := range expectedFields {
		if _, exists := stats[field]; !exists {
			t.Errorf("Stats should contain field: %s", field)
		}
	}

	if stats["total_services"].(int) != 0 {
		t.Error("Initial total services should be 0")
	}

	if stats["active_services"].(int) != 0 {
		t.Error("Initial active services should be 0")
	}

	if !stats["openziti_connected"].(bool) {
		t.Error("OpenZiti should be connected")
	}
}

func TestGetClientStats(t *testing.T) {
	config := util.DefaultConfig()
	ozm, err := NewOpenZitiManager(config)
	if err != nil {
		t.Fatalf("Failed to create OpenZiti manager: %v", err)
	}
	defer ozm.Close()

	// Check initial stats
	stats := ozm.GetClientStats()

	expectedFields := []string{"total_clients", "connected_clients", "openziti_connected"}
	for _, field := range expectedFields {
		if _, exists := stats[field]; !exists {
			t.Errorf("Stats should contain field: %s", field)
		}
	}

	if stats["total_clients"].(int) != 0 {
		t.Error("Initial total clients should be 0")
	}

	if stats["connected_clients"].(int) != 0 {
		t.Error("Initial connected clients should be 0")
	}

	if !stats["openziti_connected"].(bool) {
		t.Error("OpenZiti should be connected")
	}
}

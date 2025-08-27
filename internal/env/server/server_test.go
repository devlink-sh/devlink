package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/devlink/internal/util"
)

func TestNewServer(t *testing.T) {
	config := util.DefaultConfig()
	server, err := NewServer(config)

	if err != nil {
		t.Fatalf("NewServer failed: %v", err)
	}

	if server == nil {
		t.Fatal("NewServer() returned nil")
	}

	if server.config == nil {
		t.Error("Server should have config")
	}

	if server.encryption == nil {
		t.Error("Server should have encryption manager")
	}

	if server.shares == nil {
		t.Error("Server should have shares map")
	}
}

func TestCreateShare(t *testing.T) {
	config := util.DefaultConfig()
	server, err := NewServer(config)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Create test environment file
	envFile := &util.EnvFile{
		Variables: []util.EnvVariable{
			{
				Key:         "DATABASE_URL",
				Value:       "postgresql://localhost:5432/mydb",
				IsSensitive: true,
				LineNumber:  1,
			},
			{
				Key:         "API_KEY",
				Value:       "sk-1234567890abcdef",
				IsSensitive: true,
				LineNumber:  2,
			},
		},
		RawContent: "DATABASE_URL=postgresql://localhost:5432/mydb\nAPI_KEY=sk-1234567890abcdef",
		FilePath:   "test.env",
		TotalLines: 2,
		ValidLines: 2,
	}

	shareCode := "blue-whale-42"
	expiry := 1 * time.Hour

	share, err := server.CreateShare(envFile, shareCode, expiry, false)
	if err != nil {
		t.Fatalf("CreateShare failed: %v", err)
	}

	if share == nil {
		t.Fatal("CreateShare returned nil")
	}

	if share.ID == "" {
		t.Error("Share should have an ID")
	}

	if share.ShareCode != shareCode {
		t.Errorf("Share code mismatch. Expected: %s, Got: %s", shareCode, share.ShareCode)
	}

	if share.EncryptedData == nil {
		t.Error("Share should have encrypted data")
	}

	if share.AccessCount != 0 {
		t.Error("New share should have 0 access count")
	}

	if share.MaxAccesses != 1 {
		t.Error("Share should be single-use by default")
	}

	// Check that share is stored
	stats := server.GetStats()
	if stats["active_shares"].(int) != 1 {
		t.Error("Share should be counted as active")
	}
}

func TestGetShare(t *testing.T) {
	config := util.DefaultConfig()
	server, err := NewServer(config)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Create test environment file
	envFile := &util.EnvFile{
		Variables: []util.EnvVariable{
			{
				Key:         "DATABASE_URL",
				Value:       "postgresql://localhost:5432/mydb",
				IsSensitive: true,
				LineNumber:  1,
			},
		},
		RawContent: "DATABASE_URL=postgresql://localhost:5432/mydb",
		FilePath:   "test.env",
		TotalLines: 1,
		ValidLines: 1,
	}

	shareCode := "red-dragon-123"
	expiry := 1 * time.Hour

	// Create share
	_, err = server.CreateShare(envFile, shareCode, expiry, false)
	if err != nil {
		t.Fatalf("Failed to create share: %v", err)
	}

	// Get share
	share, err := server.GetShare(shareCode)
	if err != nil {
		t.Fatalf("GetShare failed: %v", err)
	}

	if share == nil {
		t.Fatal("GetShare returned nil")
	}

	if share.ShareCode != shareCode {
		t.Errorf("Share code mismatch. Expected: %s, Got: %s", shareCode, share.ShareCode)
	}

	if share.AccessCount != 1 {
		t.Error("Share should have access count of 1")
	}

	// Try to get the same share again (should fail for single-use)
	_, err = server.GetShare(shareCode)
	if err == nil {
		t.Error("Second access to single-use share should fail")
	}
}

func TestGetShareExpired(t *testing.T) {
	config := util.DefaultConfig()
	server, err := NewServer(config)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Create test environment file
	envFile := &util.EnvFile{
		Variables: []util.EnvVariable{
			{
				Key:         "DATABASE_URL",
				Value:       "postgresql://localhost:5432/mydb",
				IsSensitive: true,
				LineNumber:  1,
			},
		},
		RawContent: "DATABASE_URL=postgresql://localhost:5432/mydb",
		FilePath:   "test.env",
		TotalLines: 1,
		ValidLines: 1,
	}

	shareCode := "green-forest-999"
	expiry := 1 * time.Millisecond // Very short expiry

	// Create share
	_, err = server.CreateShare(envFile, shareCode, expiry, false)
	if err != nil {
		t.Fatalf("Failed to create share: %v", err)
	}

	// Wait for expiry
	time.Sleep(10 * time.Millisecond)

	// Try to get expired share
	_, err = server.GetShare(shareCode)
	if err == nil {
		t.Error("Access to expired share should fail")
	}

	if err.Error() != "share expired: "+shareCode {
		t.Errorf("Expected expired error, got: %v", err)
	}
}

func TestHandleHealth(t *testing.T) {
	config := util.DefaultConfig()
	server, err := NewServer(config)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Create request
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	server.handleHealth(rr, req)

	// Check status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	// Check response body
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got %s", response["status"])
	}

	if response["version"] != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got %s", response["version"])
	}
}

func TestHandleShare(t *testing.T) {
	config := util.DefaultConfig()
	server, err := NewServer(config)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Create request data
	requestData := map[string]interface{}{
		"share_code": "blue-whale-42",
		"data":       "DATABASE_URL=postgresql://localhost:5432/mydb\nAPI_KEY=sk-1234567890abcdef",
		"expiry":     "1h",
		"read_only":  false,
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	// Create request
	req, err := http.NewRequest("POST", "/share", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	server.handleShare(rr, req)

	// Check status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	// Check response body
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if !response["success"].(bool) {
		t.Error("Expected success to be true")
	}

	if response["share_code"] != "blue-whale-42" {
		t.Errorf("Expected share code 'blue-whale-42', got %s", response["share_code"])
	}
}

func TestHandleGet(t *testing.T) {
	config := util.DefaultConfig()
	server, err := NewServer(config)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Create test environment file
	envFile := &util.EnvFile{
		Variables: []util.EnvVariable{
			{
				Key:         "DATABASE_URL",
				Value:       "postgresql://localhost:5432/mydb",
				IsSensitive: true,
				LineNumber:  1,
			},
		},
		RawContent: "DATABASE_URL=postgresql://localhost:5432/mydb",
		FilePath:   "test.env",
		TotalLines: 1,
		ValidLines: 1,
	}

	shareCode := "red-dragon-123"
	expiry := 1 * time.Hour

	// Create share
	_, err = server.CreateShare(envFile, shareCode, expiry, false)
	if err != nil {
		t.Fatalf("Failed to create share: %v", err)
	}

	// Create request
	req, err := http.NewRequest("GET", "/get/"+shareCode, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	server.handleGet(rr, req)

	// Check status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	// Check response body
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if !response["success"].(bool) {
		t.Error("Expected success to be true")
	}

	if response["share_code"] != shareCode {
		t.Errorf("Expected share code '%s', got %s", shareCode, response["share_code"])
	}
}

func TestHandleStats(t *testing.T) {
	config := util.DefaultConfig()
	server, err := NewServer(config)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Create request
	req, err := http.NewRequest("GET", "/stats", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	server.handleStats(rr, req)

	// Check status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	// Check response body
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response["version"] != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got %s", response["version"])
	}
}

func TestGetStats(t *testing.T) {
	config := util.DefaultConfig()
	server, err := NewServer(config)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Check initial stats
	stats := server.GetStats()

	expectedFields := []string{"total_shares", "active_shares", "expired_shares", "total_accesses", "server_port", "encryption_ready"}
	for _, field := range expectedFields {
		if _, exists := stats[field]; !exists {
			t.Errorf("Stats should contain field: %s", field)
		}
	}

	if stats["total_shares"].(int) != 0 {
		t.Error("Initial total shares should be 0")
	}

	if stats["active_shares"].(int) != 0 {
		t.Error("Initial active shares should be 0")
	}

	if !stats["encryption_ready"].(bool) {
		t.Error("Encryption should be ready")
	}
}

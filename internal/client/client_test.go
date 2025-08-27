package client

import (
	"testing"
	"time"

	"github.com/devlink/internal/util"
)

func TestNewClient(t *testing.T) {
	config := util.DefaultConfig()
	serverURL := "http://localhost:8080"

	client, err := NewClient(config, serverURL)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	if client == nil {
		t.Fatal("NewClient() returned nil")
	}

	if client.config == nil {
		t.Error("Client should have config")
	}

	if client.encryption == nil {
		t.Error("Client should have encryption manager")
	}

	if client.serverURL != serverURL {
		t.Errorf("Server URL mismatch. Expected: %s, Got: %s", serverURL, client.serverURL)
	}

	if client.httpClient == nil {
		t.Error("Client should have HTTP client")
	}
}

func TestValidateShareCode(t *testing.T) {
	config := util.DefaultConfig()
	client, err := NewClient(config, "http://localhost:8080")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Valid share codes
	validCodes := []string{
		"blue-whale-42",
		"red-dragon-137",
		"green-forest-999",
	}

	for _, code := range validCodes {
		if err := client.ValidateShareCode(code); err != nil {
			t.Errorf("Valid share code '%s' failed validation: %v", code, err)
		}
	}

	// Invalid share codes
	invalidCodes := []string{
		"",                    // Empty
		"blue-whale",          // Missing number
		"blue-whale-42-extra", // Too many parts
		"BLUE-whale-42",       // Uppercase
		"blue-whale-0",        // Zero number
		"blue-whale-1000",     // Number too large
	}

	for _, code := range invalidCodes {
		if err := client.ValidateShareCode(code); err == nil {
			t.Errorf("Invalid share code '%s' should fail validation", code)
		}
	}
}

func TestCreateShare(t *testing.T) {
	config := util.DefaultConfig()
	client, err := NewClient(config, "http://localhost:8080")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
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

	// This will fail because there's no server running, but we can test the request preparation
	_, err = client.CreateShare(envFile, shareCode, expiry, false)
	if err == nil {
		t.Error("Expected error when server is not running")
	}

	// The error should be related to connection failure, not request preparation
	if err != nil && !contains(err.Error(), "connection") && !contains(err.Error(), "refused") {
		t.Errorf("Expected connection error, got: %v", err)
	}
}

func TestHealthCheck(t *testing.T) {
	config := util.DefaultConfig()
	client, err := NewClient(config, "http://localhost:8080")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// This will fail because there's no server running
	err = client.HealthCheck()
	if err == nil {
		t.Error("Expected error when server is not running")
	}

	// The error should be related to connection failure
	if err != nil && !contains(err.Error(), "connection") && !contains(err.Error(), "refused") {
		t.Errorf("Expected connection error, got: %v", err)
	}
}

func TestGetServerStats(t *testing.T) {
	config := util.DefaultConfig()
	client, err := NewClient(config, "http://localhost:8080")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// This will fail because there's no server running
	_, err = client.GetServerStats()
	if err == nil {
		t.Error("Expected error when server is not running")
	}

	// The error should be related to connection failure
	if err != nil && !contains(err.Error(), "connection") && !contains(err.Error(), "refused") {
		t.Errorf("Expected connection error, got: %v", err)
	}
}

func TestShareFile(t *testing.T) {
	config := util.DefaultConfig()
	client, err := NewClient(config, "http://localhost:8080")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	shareCode := "blue-whale-42"
	expiry := 1 * time.Hour

	// This will fail because the file doesn't exist
	_, err = client.ShareFile("nonexistent.env", shareCode, expiry, false)
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}

	// The error should be related to file not found
	if err != nil && !contains(err.Error(), "no such file") && !contains(err.Error(), "not found") {
		t.Errorf("Expected file not found error, got: %v", err)
	}
}

func TestGetShareToFile(t *testing.T) {
	config := util.DefaultConfig()
	client, err := NewClient(config, "http://localhost:8080")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	shareCode := "blue-whale-42"
	outputPath := "/tmp/test_output.env"

	// This will fail because there's no server running
	err = client.GetShareToFile(shareCode, outputPath)
	if err == nil {
		t.Error("Expected error when server is not running")
	}

	// The error should be related to connection failure
	if err != nil && !contains(err.Error(), "connection") && !contains(err.Error(), "refused") {
		t.Errorf("Expected connection error, got: %v", err)
	}
}

func TestGetShareInfo(t *testing.T) {
	config := util.DefaultConfig()
	client, err := NewClient(config, "http://localhost:8080")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	shareCode := "blue-whale-42"

	// This will fail because there's no server running
	_, err = client.GetShareInfo(shareCode)
	if err == nil {
		t.Error("Expected error when server is not running")
	}

	// The error should be related to connection failure
	if err != nil && !contains(err.Error(), "connection") && !contains(err.Error(), "refused") {
		t.Errorf("Expected connection error, got: %v", err)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

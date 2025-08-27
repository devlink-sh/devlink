package util

import (
	"strings"
	"testing"
	"time"
)

func TestNewValidationManager(t *testing.T) {
	vm := NewValidationManager()
	if vm == nil {
		t.Fatal("NewValidationManager() returned nil")
	}

	// Check that configuration was loaded
	if vm.maxAttemptsPerMinute <= 0 {
		t.Error("Max attempts should be positive")
	}

	if vm.blockDuration <= 0 {
		t.Error("Block duration should be positive")
	}

	if vm.codePattern == nil {
		t.Error("Code pattern should be initialized")
	}

	if vm.config == nil {
		t.Error("ValidationManager should have config")
	}

	if vm.securityConfig == nil {
		t.Error("ValidationManager should have security config")
	}
}

func TestValidateShareCode_ValidCodes(t *testing.T) {
	vm := NewValidationManager()

	validCodes := []string{
		"blue-whale-42",
		"red-dragon-1",
		"green-forest-999",
		"bright-river-137",
		"quiet-ocean-7",
	}

	for _, code := range validCodes {
		if err := vm.ValidateShareCode(code, "test-client"); err != nil {
			t.Errorf("Valid code '%s' failed validation: %v", code, err)
		}
	}
}

func TestValidateShareCode_InvalidFormat(t *testing.T) {
	vm := NewValidationManager()

	invalidCodes := []string{
		"",                    // Empty
		"blue-whale",          // Missing number
		"blue-whale-42-extra", // Too many parts
		"BLUE-whale-42",       // Uppercase
		"blue-WHALE-42",       // Mixed case
		"blue-whale-0",        // Zero number
		"blue-whale-1000",     // Number too large
		"a-whale-42",          // Adjective too short
		"blue-b-42",           // Noun too short
		"blue-whale-abc",      // Non-numeric number
		"blue whale 42",       // Wrong separator
		"blue-whale-",         // Missing number
		"-whale-42",           // Missing adjective
		"blue--42",            // Missing noun
	}

	for _, code := range invalidCodes {
		if err := vm.ValidateShareCode(code, "test-client"); err == nil {
			t.Errorf("Invalid code '%s' should fail validation", code)
		}
	}
}

func TestValidateShareCode_OffensiveWords(t *testing.T) {
	vm := NewValidationManager()

	// Use words that are in the default security config
	offensiveCodes := []string{
		"test-whale-42",
		"blue-demo-42",
		"temp-dragon-42",
		"fake-forest-42",
	}

	for _, code := range offensiveCodes {
		if err := vm.ValidateShareCode(code, "test-client"); err == nil {
			t.Errorf("Offensive code '%s' should fail validation", code)
		}
	}
}

func TestValidateShareCode_SequentialPatterns(t *testing.T) {
	vm := NewValidationManager()

	sequentialCodes := []string{
		"blue-whale-123", // Sequential numbers
		"blue-whale-321", // Reverse sequential
	}

	for _, code := range sequentialCodes {
		if err := vm.ValidateShareCode(code, "test-client"); err == nil {
			t.Errorf("Sequential code '%s' should fail validation", code)
		}
	}
}

func TestValidateShareCode_CommonPatterns(t *testing.T) {
	vm := NewValidationManager()

	commonPatterns := []string{
		"test-test-1",
		"demo-demo-1",
		"temp-temp-1",
		"fake-fake-1",
		"blue-blue-1",
		"red-red-1",
		"green-green-1",
	}

	for _, code := range commonPatterns {
		if err := vm.ValidateShareCode(code, "test-client"); err == nil {
			t.Errorf("Common pattern code '%s' should fail validation", code)
		}
	}
}

func TestValidateShareCode_RateLimit(t *testing.T) {
	vm := NewValidationManager()
	clientID := "test-client"

	// Make multiple validation attempts
	for i := 0; i < 10; i++ {
		err := vm.ValidateShareCode("blue-whale-42", clientID)
		if err != nil {
			t.Errorf("Validation %d should succeed: %v", i, err)
		}
	}

	// The 11th attempt should fail due to rate limit
	err := vm.ValidateShareCode("blue-whale-42", clientID)
	if err == nil {
		t.Error("11th validation should fail due to rate limit")
	}

	if !strings.Contains(err.Error(), "rate limit exceeded") {
		t.Errorf("Expected rate limit error, got: %v", err)
	}
}

func TestValidateShareCode_RateLimitReset(t *testing.T) {
	vm := NewValidationManager()
	clientID := "test-client"

	// Make some validation attempts
	for i := 0; i < 5; i++ {
		vm.ValidateShareCode("blue-whale-42", clientID)
	}

	// Manually reset the rate limit by setting old timestamp
	vm.rateMutex.Lock()
	if entry, exists := vm.rateLimitMap[clientID]; exists {
		entry.lastAttempt = time.Now().Add(-2 * time.Minute)
		entry.attempts = 0
	}
	vm.rateMutex.Unlock()

	// Should be able to make more attempts
	for i := 0; i < 5; i++ {
		err := vm.ValidateShareCode("blue-whale-42", clientID)
		if err != nil {
			t.Errorf("Validation after reset should succeed: %v", err)
		}
	}
}

func TestValidateShareCode_Blocking(t *testing.T) {
	vm := NewValidationManager()
	clientID := "test-client"

	// Exceed rate limit to get blocked
	for i := 0; i < 11; i++ {
		vm.ValidateShareCode("blue-whale-42", clientID)
	}

	// Should be blocked
	err := vm.ValidateShareCode("blue-whale-42", clientID)
	if err == nil {
		t.Error("Client should be blocked")
	}

	if !strings.Contains(err.Error(), "client is blocked") {
		t.Errorf("Expected blocking error, got: %v", err)
	}
}

func TestValidateShareCode_BlockExpiration(t *testing.T) {
	vm := NewValidationManager()
	clientID := "test-client"

	// Exceed rate limit to get blocked
	for i := 0; i < 11; i++ {
		vm.ValidateShareCode("blue-whale-42", clientID)
	}

	// Manually set block to expire
	vm.rateMutex.Lock()
	if entry, exists := vm.rateLimitMap[clientID]; exists {
		entry.blockUntil = time.Now().Add(-1 * time.Minute)
	}
	vm.rateMutex.Unlock()

	// Should be able to validate again
	err := vm.ValidateShareCode("blue-whale-42", clientID)
	if err != nil {
		t.Errorf("Validation after block expiration should succeed: %v", err)
	}
}

func TestValidateWord(t *testing.T) {
	vm := NewValidationManager()

	// Valid words
	validWords := []string{"blue", "red", "green", "whale", "dragon", "forest"}
	for _, word := range validWords {
		if err := vm.validateWord(word, "adjective"); err != nil {
			t.Errorf("Valid word '%s' should pass validation: %v", word, err)
		}
	}

	// Invalid words
	invalidWords := []string{
		"",                          // Empty
		"a",                         // Too short
		"verylongwordthatistoolong", // Too long
		"BLUE",                      // Uppercase
		"Blue",                      // Mixed case
		"blue123",                   // Contains numbers
		"blue-whale",                // Contains hyphen
	}

	for _, word := range invalidWords {
		if err := vm.validateWord(word, "adjective"); err == nil {
			t.Errorf("Invalid word '%s' should fail validation", word)
		}
	}
}

func TestValidateNumber(t *testing.T) {
	vm := NewValidationManager()

	// Valid numbers
	validNumbers := []string{"1", "42", "123", "999"}
	for _, num := range validNumbers {
		if err := vm.validateNumber(num); err != nil {
			t.Errorf("Valid number '%s' should pass validation: %v", num, err)
		}
	}

	// Invalid numbers
	invalidNumbers := []string{
		"",     // Empty
		"0",    // Zero
		"1000", // Too large
		"abc",  // Non-numeric
		"12a",  // Mixed
		"1.5",  // Decimal
		"-1",   // Negative
	}

	for _, num := range invalidNumbers {
		if err := vm.validateNumber(num); err == nil {
			t.Errorf("Invalid number '%s' should fail validation", num)
		}
	}
}

func TestIsSequentialPattern(t *testing.T) {
	vm := NewValidationManager()

	// Sequential patterns
	sequentialPatterns := []string{
		"blue-whale-123",
		"blue-whale-321",
		"blue-whale-1234",
	}

	for _, code := range sequentialPatterns {
		if !vm.isSequentialPattern(code) {
			t.Errorf("Sequential pattern '%s' should be detected", code)
		}
	}

	// Non-sequential patterns
	nonSequentialPatterns := []string{
		"blue-whale-42",
		"blue-whale-137",
		"blue-whale-999",
	}

	for _, code := range nonSequentialPatterns {
		if vm.isSequentialPattern(code) {
			t.Errorf("Non-sequential pattern '%s' should not be detected", code)
		}
	}
}

func TestHasExcessiveRepetition(t *testing.T) {
	vm := NewValidationManager()

	// Excessive repetition
	repetitiveCodes := []string{
		"blue-whale-1111",
		"blue-whale-aaaa",
		"blue-whale-xxxx",
	}

	for _, code := range repetitiveCodes {
		if !vm.hasExcessiveRepetition(code) {
			t.Errorf("Excessive repetition in '%s' should be detected", code)
		}
	}

	// Normal codes
	normalCodes := []string{
		"blue-whale-42",
		"blue-whale-123",
		"blue-whale-999",
	}

	for _, code := range normalCodes {
		if vm.hasExcessiveRepetition(code) {
			t.Errorf("Normal code '%s' should not have excessive repetition", code)
		}
	}
}

func TestIsCommonPattern(t *testing.T) {
	vm := NewValidationManager()

	// Common patterns
	commonPatterns := []string{
		"test-test-1",
		"demo-demo-1",
		"temp-temp-1",
		"fake-fake-1",
		"blue-blue-1",
		"red-red-1",
		"green-green-1",
	}

	for _, code := range commonPatterns {
		if !vm.isCommonPattern(code) {
			t.Errorf("Common pattern '%s' should be detected", code)
		}
	}

	// Non-common patterns
	nonCommonPatterns := []string{
		"blue-whale-42",
		"red-dragon-1",
		"green-forest-999",
	}

	for _, code := range nonCommonPatterns {
		if vm.isCommonPattern(code) {
			t.Errorf("Non-common pattern '%s' should not be detected", code)
		}
	}
}

func TestGetRateLimitStats(t *testing.T) {
	vm := NewValidationManager()

	// Make some validation attempts
	for i := 0; i < 5; i++ {
		vm.ValidateShareCode("blue-whale-42", "client1")
	}

	stats := vm.GetRateLimitStats()

	// Check that stats contain expected fields
	expectedFields := []string{"active_clients", "blocked_clients", "max_attempts", "block_duration_s"}
	for _, field := range expectedFields {
		if _, exists := stats[field]; !exists {
			t.Errorf("Stats should contain field: %s", field)
		}
	}

	// Check specific values
	if stats["active_clients"].(int) < 1 {
		t.Error("Should have at least one active client")
	}

	// Max attempts should be configurable and positive
	if stats["max_attempts"].(int) <= 0 {
		t.Error("Max attempts should be positive")
	}

	// Block duration should be configurable and positive
	if stats["block_duration_s"].(float64) <= 0 {
		t.Error("Block duration should be positive")
	}
}

func TestValidationCleanup(t *testing.T) {
	vm := NewValidationManager()

	// Add some clients
	vm.ValidateShareCode("blue-whale-42", "client1")
	vm.ValidateShareCode("red-dragon-1", "client2")

	// Manually set one client to be old
	vm.rateMutex.Lock()
	if entry, exists := vm.rateLimitMap["client1"]; exists {
		entry.lastAttempt = time.Now().Add(-2 * time.Hour)
	}
	vm.rateMutex.Unlock()

	// Run cleanup
	vm.Cleanup()

	// Check that old client was removed
	vm.rateMutex.RLock()
	_, exists := vm.rateLimitMap["client1"]
	vm.rateMutex.RUnlock()

	if exists {
		t.Error("Old client should be removed by cleanup")
	}

	// Check that recent client is still there
	vm.rateMutex.RLock()
	_, exists = vm.rateLimitMap["client2"]
	vm.rateMutex.RUnlock()

	if !exists {
		t.Error("Recent client should not be removed by cleanup")
	}
}

func TestConcurrentValidation(t *testing.T) {
	vm := NewValidationManager()
	results := make(chan error, 10)

	// Run validations concurrently
	for i := 0; i < 10; i++ {
		go func() {
			err := vm.ValidateShareCode("blue-whale-42", "concurrent-client")
			results <- err
		}()
	}

	// Collect results
	for i := 0; i < 10; i++ {
		select {
		case err := <-results:
			// Some may fail due to rate limiting, which is expected
			if err != nil && !strings.Contains(err.Error(), "rate limit exceeded") {
				t.Errorf("Unexpected error in concurrent validation: %v", err)
			}
		case <-time.After(5 * time.Second):
			t.Error("Concurrent validation test timed out")
		}
	}
}

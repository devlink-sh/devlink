package util

import (
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestNewTokenGenerator(t *testing.T) {
	tg := NewTokenGenerator()
	if tg == nil {
		t.Fatal("NewTokenGenerator() returned nil")
	}

	if len(tg.adjectives) == 0 {
		t.Error("TokenGenerator should have adjectives")
	}

	if len(tg.nouns) == 0 {
		t.Error("TokenGenerator should have nouns")
	}

	// Rate limit should be configurable and positive
	if tg.rateLimit <= 0 {
		t.Error("Rate limit should be positive")
	}

	if tg.config == nil {
		t.Error("TokenGenerator should have config")
	}
}

func TestGenerateShareCode(t *testing.T) {
	tg := NewTokenGenerator()

	// Test basic generation
	code, err := tg.GenerateShareCode()
	if err != nil {
		t.Fatalf("GenerateShareCode failed: %v", err)
	}

	if code == "" {
		t.Error("Generated code should not be empty")
	}

	// Test format: word-word-number
	parts := strings.Split(code, "-")
	if len(parts) != 3 {
		t.Errorf("Generated code should have 3 parts, got %d: %s", len(parts), code)
	}

	// Test adjective
	if !tg.isValidWord(parts[0], tg.adjectives) {
		t.Errorf("Generated adjective '%s' is not valid", parts[0])
	}

	// Test noun
	if !tg.isValidWord(parts[1], tg.nouns) {
		t.Errorf("Generated noun '%s' is not valid", parts[1])
	}

	// Test number
	if !tg.isValidNumber(parts[2]) {
		t.Errorf("Generated number '%s' is not valid", parts[2])
	}
}

func TestGenerateShareCode_Uniqueness(t *testing.T) {
	// Use a test configuration with disabled rate limiting
	config := DefaultConfig()
	config.TokenRateLimit = "0ms"
	tg := NewTokenGeneratorWithConfig(config)

	generatedCodes := make(map[string]bool)

	// Generate multiple codes and check uniqueness
	for i := 0; i < 10; i++ {
		code, err := tg.GenerateShareCode()
		if err != nil {
			t.Fatalf("GenerateShareCode failed on attempt %d: %v", i, err)
		}

		if generatedCodes[code] {
			t.Errorf("Duplicate code generated: %s", code)
		}
		generatedCodes[code] = true
	}
}

func TestGenerateShareCode_RateLimit(t *testing.T) {
	tg := NewTokenGenerator()

	// First generation should succeed
	_, err := tg.GenerateShareCode()
	if err != nil {
		t.Fatalf("First generation failed: %v", err)
	}

	// Immediate second generation should fail due to rate limit
	_, err = tg.GenerateShareCode()
	if err == nil {
		t.Error("Second generation should fail due to rate limit")
	}

	if !strings.Contains(err.Error(), "rate limit exceeded") {
		t.Errorf("Expected rate limit error, got: %v", err)
	}

	// Wait for rate limit to expire
	time.Sleep(150 * time.Millisecond)

	// Third generation should succeed
	_, err = tg.GenerateShareCode()
	if err != nil {
		t.Errorf("Third generation should succeed after rate limit: %v", err)
	}
}

func TestValidateShareCode(t *testing.T) {
	tg := NewTokenGenerator()

	// Valid codes
	validCodes := []string{
		"blue-whale-42",
		"red-dragon-1",
		"green-forest-999",
		"bright-river-137",
		"quiet-ocean-7",
	}

	for _, code := range validCodes {
		if err := tg.ValidateShareCode(code); err != nil {
			t.Errorf("Valid code '%s' failed validation: %v", code, err)
		}
	}

	// Invalid codes
	invalidCodes := []string{
		"",                    // Empty
		"blue-whale",          // Missing number
		"blue-whale-42-extra", // Too many parts
		"blue-whale-0",        // Zero number
		"blue-whale-1000",     // Number too large
		"a-whale-42",          // Adjective too short
		"blue-b-42",           // Noun too short
		"blue-whale-abc",      // Non-numeric number
		"blue whale 42",       // Wrong separator
	}

	for _, code := range invalidCodes {
		if err := tg.ValidateShareCode(code); err == nil {
			t.Errorf("Invalid code '%s' should fail validation", code)
		}
	}
}

func TestIsValidWord(t *testing.T) {
	tg := NewTokenGenerator()

	// Valid words
	validWords := []string{"blue", "red", "green", "whale", "dragon", "forest"}
	for _, word := range validWords {
		if !tg.isValidWord(word, tg.adjectives) && !tg.isValidWord(word, tg.nouns) {
			t.Errorf("Valid word '%s' should be valid", word)
		}
	}

	// Invalid words
	invalidWords := []string{"", "a", "invalidword", "blue123"}
	for _, word := range invalidWords {
		if tg.isValidWord(word, tg.adjectives) || tg.isValidWord(word, tg.nouns) {
			t.Errorf("Invalid word '%s' should not be valid", word)
		}
	}
}

func TestIsValidNumber(t *testing.T) {
	tg := NewTokenGenerator()

	// Valid numbers
	validNumbers := []string{"1", "42", "123", "999"}
	for _, num := range validNumbers {
		if !tg.isValidNumber(num) {
			t.Errorf("Valid number '%s' should be valid", num)
		}
	}

	// Invalid numbers
	invalidNumbers := []string{"", "0", "1000", "abc", "12a", "1.5", "-1"}
	for _, num := range invalidNumbers {
		if tg.isValidNumber(num) {
			t.Errorf("Invalid number '%s' should not be valid", num)
		}
	}
}

func TestCollisionHandling(t *testing.T) {
	tg := NewTokenGenerator()

	// Mark a code as used
	testCode := "blue-whale-42"
	tg.markCodeAsUsed(testCode)

	// Check that it's marked as used
	if !tg.isCodeUsed(testCode) {
		t.Error("Code should be marked as used")
	}

	// Check that a different code is not used
	if tg.isCodeUsed("red-dragon-1") {
		t.Error("Different code should not be marked as used")
	}
}

func TestCodeExpiration(t *testing.T) {
	tg := NewTokenGenerator()

	// Mark a code as used
	testCode := "blue-whale-42"
	tg.markCodeAsUsed(testCode)

	// Manually set the used time to be old (older than 24 hours)
	tg.mutex.Lock()
	tg.usedCodes[testCode] = time.Now().Add(-25 * time.Hour)
	tg.mutex.Unlock()

	// Check that it's no longer marked as used (should be cleaned up)
	if tg.isCodeUsed(testCode) {
		t.Error("Expired code should not be marked as used")
	}
}

func TestGetStats(t *testing.T) {
	// Use a test configuration with disabled rate limiting
	config := DefaultConfig()
	config.TokenRateLimit = "0ms"
	tg := NewTokenGeneratorWithConfig(config)

	// Generate some codes to populate stats
	for i := 0; i < 5; i++ {
		tg.GenerateShareCode()
	}

	stats := tg.GetStats()

	// Check that stats contain expected fields
	expectedFields := []string{"total_adjectives", "total_nouns", "active_codes", "rate_limit_ms"}
	for _, field := range expectedFields {
		if _, exists := stats[field]; !exists {
			t.Errorf("Stats should contain field: %s", field)
		}
	}

	// Check specific values
	if stats["active_codes"].(int) < 1 {
		t.Error("Should have at least one active code")
	}

	if stats["total_adjectives"].(int) < 10 {
		t.Error("Should have reasonable number of adjectives")
	}

	if stats["total_nouns"].(int) < 10 {
		t.Error("Should have reasonable number of nouns")
	}
}

func TestCleanup(t *testing.T) {
	tg := NewTokenGenerator()

	// Mark some codes as used
	codes := []string{"blue-whale-42", "red-dragon-1", "green-forest-999"}
	for _, code := range codes {
		tg.markCodeAsUsed(code)
	}

	// Manually set one code to be old
	tg.mutex.Lock()
	tg.usedCodes[codes[0]] = time.Now().Add(-25 * time.Hour)
	tg.mutex.Unlock()

	// Run cleanup
	tg.Cleanup()

	// Check that old code was removed
	if tg.isCodeUsed(codes[0]) {
		t.Error("Old code should be removed by cleanup")
	}

	// Check that recent codes are still there
	if !tg.isCodeUsed(codes[1]) {
		t.Error("Recent code should not be removed by cleanup")
	}
}

func TestCodeFormatRegex(t *testing.T) {
	// Test the regex pattern used in validation
	pattern := regexp.MustCompile(`^[a-z]+-[a-z]+-\d{1,3}$`)

	validCodes := []string{
		"blue-whale-42",
		"red-dragon-1",
		"green-forest-999",
	}

	for _, code := range validCodes {
		if !pattern.MatchString(code) {
			t.Errorf("Valid code '%s' should match pattern", code)
		}
	}

	invalidCodes := []string{
		"BLUE-whale-42",
		"blue-WHALE-42",
		"blue-whale-1000",
		"blue-whale-abc",
		"blue-whale",
		"blue-whale-42-extra",
	}

	for _, code := range invalidCodes {
		if pattern.MatchString(code) {
			t.Errorf("Invalid code '%s' should not match pattern", code)
		}
	}
}

func TestConcurrentGeneration(t *testing.T) {
	// Use a test configuration with disabled rate limiting
	config := DefaultConfig()
	config.TokenRateLimit = "0ms"
	tg := NewTokenGeneratorWithConfig(config)

	results := make(chan string, 10)
	errors := make(chan error, 10)

	// Generate codes concurrently
	for i := 0; i < 10; i++ {
		go func() {
			code, err := tg.GenerateShareCode()
			if err != nil {
				errors <- err
				return
			}
			results <- code
		}()
	}

	// Collect results
	generatedCodes := make(map[string]bool)
	for i := 0; i < 10; i++ {
		select {
		case code := <-results:
			if generatedCodes[code] {
				t.Errorf("Duplicate code generated in concurrent test: %s", code)
			}
			generatedCodes[code] = true
		case err := <-errors:
			// Rate limiting errors are expected in concurrent tests
			if !strings.Contains(err.Error(), "rate limit exceeded") {
				t.Errorf("Unexpected error in concurrent test: %v", err)
			}
		case <-time.After(5 * time.Second):
			t.Error("Concurrent generation test timed out")
		}
	}
}

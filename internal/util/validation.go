package util

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"
)

type ValidationManager struct {
	rateLimitMap         map[string]*rateLimitEntry
	rateMutex            sync.RWMutex
	codePattern          *regexp.Regexp
	maxAttemptsPerMinute int
	blockDuration        time.Duration
	cleanupInterval      time.Duration
	securityConfig       *SecurityConfig
	config               *Config
}

type rateLimitEntry struct {
	attempts    int
	lastAttempt time.Time
	blocked     bool
	blockUntil  time.Time
}

func NewValidationManager() *ValidationManager {
	config, err := LoadConfig()
	if err != nil {
		config = DefaultConfig()
	}

	securityConfig := NewSecurityConfig()
	securityConfig.LoadFromEnvironment()

	return &ValidationManager{
		rateLimitMap:         make(map[string]*rateLimitEntry),
		codePattern:          regexp.MustCompile(`^[a-z]+-[a-z]+-\d{1,3}$`),
		maxAttemptsPerMinute: config.ValidationMaxAttempts,
		blockDuration:        config.GetValidationBlockDuration(),
		cleanupInterval:      config.GetValidationCleanupInterval(),
		securityConfig:       securityConfig,
		config:               config,
	}
}

func NewValidationManagerWithConfig(config *Config) *ValidationManager {
	if config == nil {
		config = DefaultConfig()
	}

	securityConfig := NewSecurityConfig()
	securityConfig.LoadOffensiveWords(config.OffensiveWordsFile)

	return &ValidationManager{
		rateLimitMap:         make(map[string]*rateLimitEntry),
		codePattern:          regexp.MustCompile(`^[a-z]+-[a-z]+-\d{1,3}$`),
		maxAttemptsPerMinute: config.ValidationMaxAttempts,
		blockDuration:        config.GetValidationBlockDuration(),
		cleanupInterval:      config.GetValidationCleanupInterval(),
		securityConfig:       securityConfig,
		config:               config,
	}
}

func (vm *ValidationManager) ValidateShareCode(code, clientID string) error {
	if err := vm.checkRateLimit(clientID); err != nil {
		return fmt.Errorf("rate limit exceeded: %w", err)
	}

	if err := vm.validateFormat(code); err != nil {
		vm.recordAttempt(clientID)
		return fmt.Errorf("invalid format: %w", err)
	}

	if err := vm.validateContent(code); err != nil {
		vm.recordAttempt(clientID)
		return fmt.Errorf("invalid content: %w", err)
	}

	if err := vm.validateSecurity(code); err != nil {
		vm.recordAttempt(clientID)
		return fmt.Errorf("security check failed: %w", err)
	}

	vm.recordAttempt(clientID)
	return nil
}

func (vm *ValidationManager) validateFormat(code string) error {
	if code == "" {
		return fmt.Errorf("share code cannot be empty")
	}

	if len(code) < 8 || len(code) > 25 {
		return fmt.Errorf("share code must be between 8 and 25 characters")
	}

	if !vm.codePattern.MatchString(code) {
		return fmt.Errorf("share code must match pattern: word-word-number")
	}

	return nil
}

func (vm *ValidationManager) validateContent(code string) error {
	parts := strings.Split(code, "-")
	if len(parts) != 3 {
		return fmt.Errorf("invalid number of parts")
	}

	if err := vm.validateWord(parts[0], "adjective"); err != nil {
		return err
	}

	if err := vm.validateWord(parts[1], "noun"); err != nil {
		return err
	}

	if err := vm.validateNumber(parts[2]); err != nil {
		return err
	}

	return nil
}

func (vm *ValidationManager) validateWord(word, wordType string) error {
	if len(word) < 2 || len(word) > 12 {
		return fmt.Errorf("%s must be between 2 and 12 characters", wordType)
	}

	if !regexp.MustCompile(`^[a-z]+$`).MatchString(word) {
		return fmt.Errorf("%s must contain only lowercase letters", wordType)
	}

	if vm.isOffensiveWord(word) {
		return fmt.Errorf("%s contains inappropriate content", wordType)
	}

	return nil
}

func (vm *ValidationManager) validateNumber(numStr string) error {
	if len(numStr) < 1 || len(numStr) > 3 {
		return fmt.Errorf("number must be between 1 and 3 digits")
	}

	if !regexp.MustCompile(`^\d+$`).MatchString(numStr) {
		return fmt.Errorf("number must contain only digits")
	}

	if numStr == "0" {
		return fmt.Errorf("number cannot be zero")
	}

	if len(numStr) == 3 {
		if numStr[0] > '9' || (numStr[0] == '9' && (numStr[1] > '9' || (numStr[1] == '9' && numStr[2] > '9'))) {
			return fmt.Errorf("number must be between 1 and 999")
		}
	}

	return nil
}

func (vm *ValidationManager) validateSecurity(code string) error {
	if vm.isSequentialPattern(code) {
		return fmt.Errorf("share code contains predictable pattern")
	}

	if vm.hasExcessiveRepetition(code) {
		return fmt.Errorf("share code contains excessive repetition")
	}

	if vm.isCommonPattern(code) {
		return fmt.Errorf("share code matches common pattern")
	}

	return nil
}

func (vm *ValidationManager) isOffensiveWord(word string) bool {
	return vm.securityConfig.IsOffensiveWord(word)
}

func (vm *ValidationManager) isSequentialPattern(code string) bool {
	if strings.Contains(code, "aaa") || strings.Contains(code, "bbb") {
		return true
	}

	parts := strings.Split(code, "-")
	if len(parts) == 3 {
		numStr := parts[2]
		if len(numStr) >= 3 {
			for i := 0; i < len(numStr)-2; i++ {
				if numStr[i+1] == numStr[i]+1 && numStr[i+2] == numStr[i+1]+1 {
					return true
				}
				if numStr[i+1] == numStr[i]-1 && numStr[i+2] == numStr[i+1]-1 {
					return true
				}
			}
		}
	}

	return false
}

func (vm *ValidationManager) hasExcessiveRepetition(code string) bool {
	for i := 0; i < len(code)-3; i++ {
		if code[i] == code[i+1] && code[i] == code[i+2] && code[i] == code[i+3] {
			return true
		}
	}

	return false
}

func (vm *ValidationManager) isCommonPattern(code string) bool {
	commonPatterns := []string{
		"test-test-1", "demo-demo-1", "temp-temp-1", "fake-fake-1",
		"blue-blue-1", "red-red-1", "green-green-1",
	}

	code = strings.ToLower(code)
	for _, pattern := range commonPatterns {
		if code == pattern {
			return true
		}
	}

	return false
}

// checkRateLimit checks if a client has exceeded rate limits
func (vm *ValidationManager) checkRateLimit(clientID string) error {
	vm.rateMutex.Lock()
	defer vm.rateMutex.Unlock()

	entry, exists := vm.rateLimitMap[clientID]
	if !exists {
		vm.rateLimitMap[clientID] = &rateLimitEntry{
			attempts:    0,
			lastAttempt: time.Now(),
			blocked:     false,
		}
		return nil
	}

	// Check if currently blocked
	if entry.blocked {
		if time.Now().Before(entry.blockUntil) {
			return fmt.Errorf("client is blocked until %v", entry.blockUntil)
		}
		// Unblock if time has passed
		entry.blocked = false
		entry.attempts = 0
	}

	// Check if we need to reset attempts (new minute)
	if time.Since(entry.lastAttempt) > time.Minute {
		entry.attempts = 0
	}

	// Check if rate limit exceeded
	if entry.attempts >= vm.maxAttemptsPerMinute {
		entry.blocked = true
		entry.blockUntil = time.Now().Add(vm.blockDuration)
		return fmt.Errorf("rate limit exceeded, blocked for %v", vm.blockDuration)
	}

	return nil
}

// recordAttempt records a validation attempt for rate limiting
func (vm *ValidationManager) recordAttempt(clientID string) {
	vm.rateMutex.Lock()
	defer vm.rateMutex.Unlock()

	entry, exists := vm.rateLimitMap[clientID]
	if !exists {
		vm.rateLimitMap[clientID] = &rateLimitEntry{
			attempts:    1,
			lastAttempt: time.Now(),
			blocked:     false,
		}
		return
	}

	entry.attempts++
	entry.lastAttempt = time.Now()
}

// GetRateLimitStats returns rate limiting statistics
func (vm *ValidationManager) GetRateLimitStats() map[string]interface{} {
	vm.rateMutex.RLock()
	defer vm.rateMutex.RUnlock()

	activeClients := 0
	blockedClients := 0

	for _, entry := range vm.rateLimitMap {
		if time.Since(entry.lastAttempt) <= vm.cleanupInterval {
			activeClients++
		}
		if entry.blocked && time.Now().Before(entry.blockUntil) {
			blockedClients++
		}
	}

	return map[string]interface{}{
		"active_clients":   activeClients,
		"blocked_clients":  blockedClients,
		"max_attempts":     vm.maxAttemptsPerMinute,
		"block_duration_s": vm.blockDuration.Seconds(),
	}
}

// Cleanup removes old rate limit entries
func (vm *ValidationManager) Cleanup() {
	vm.rateMutex.Lock()
	defer vm.rateMutex.Unlock()

	now := time.Now()
	for clientID, entry := range vm.rateLimitMap {
		// Remove entries older than configured cleanup interval
		if now.Sub(entry.lastAttempt) > vm.cleanupInterval {
			delete(vm.rateLimitMap, clientID)
		}
	}
}

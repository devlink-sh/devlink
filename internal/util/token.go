package util

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"
)

type TokenGenerator struct {
	adjectives      []string
	nouns           []string
	usedCodes       map[string]time.Time
	mutex           sync.RWMutex
	lastGeneration  time.Time
	rateLimit       time.Duration
	cleanupInterval time.Duration
	config          *Config
}

func NewTokenGenerator() *TokenGenerator {
	config, err := LoadConfig()
	if err != nil {
		config = DefaultConfig()
	}

	return &TokenGenerator{
		adjectives: []string{
			"blue", "red", "green", "yellow", "purple", "orange", "pink", "brown",
			"black", "white", "gray", "silver", "gold", "navy", "teal", "coral",
			"lime", "indigo", "violet", "maroon", "olive", "cyan", "magenta",
			"fast", "slow", "big", "small", "tall", "short", "wide", "narrow",
			"bright", "dark", "light", "heavy", "soft", "hard", "smooth", "rough",
			"warm", "cool", "fresh", "old", "new", "young", "ancient", "modern",
			"quiet", "loud", "calm", "wild", "gentle", "fierce", "brave", "shy",
		},
		nouns: []string{
			"whale", "dolphin", "shark", "turtle", "seahorse", "octopus", "jellyfish",
			"crab", "lobster", "starfish", "clam", "oyster", "mussel", "coral",
			"eagle", "hawk", "owl", "falcon", "raven", "crow", "sparrow", "robin",
			"lion", "tiger", "bear", "wolf", "fox", "deer", "rabbit", "squirrel",
			"elephant", "giraffe", "zebra", "rhino", "hippo", "camel", "llama",
			"dragon", "phoenix", "unicorn", "griffin", "pegasus", "centaur",
			"mountain", "river", "ocean", "forest", "desert", "island", "valley",
			"castle", "tower", "bridge", "temple", "palace", "cottage", "cabin",
			"diamond", "ruby", "emerald", "sapphire", "pearl", "crystal", "gem",
		},
		usedCodes:       make(map[string]time.Time),
		rateLimit:       config.GetTokenRateLimit(),
		cleanupInterval: config.GetTokenCleanupInterval(),
		config:          config,
	}
}

func NewTokenGeneratorWithConfig(config *Config) *TokenGenerator {
	if config == nil {
		config = DefaultConfig()
	}

	return &TokenGenerator{
		adjectives: []string{
			"blue", "red", "green", "yellow", "purple", "orange", "pink", "brown",
			"black", "white", "gray", "silver", "gold", "navy", "teal", "coral",
			"lime", "indigo", "violet", "maroon", "olive", "cyan", "magenta",
			"fast", "slow", "big", "small", "tall", "short", "wide", "narrow",
			"bright", "dark", "light", "heavy", "soft", "hard", "smooth", "rough",
			"warm", "cool", "fresh", "old", "new", "young", "ancient", "modern",
			"quiet", "loud", "calm", "wild", "gentle", "fierce", "brave", "shy",
		},
		nouns: []string{
			"whale", "dolphin", "shark", "turtle", "seahorse", "octopus", "jellyfish",
			"crab", "lobster", "starfish", "clam", "oyster", "mussel", "coral",
			"eagle", "hawk", "owl", "falcon", "raven", "crow", "sparrow", "robin",
			"lion", "tiger", "bear", "wolf", "fox", "deer", "rabbit", "squirrel",
			"elephant", "giraffe", "zebra", "rhino", "hippo", "camel", "llama",
			"dragon", "phoenix", "unicorn", "griffin", "pegasus", "centaur",
			"mountain", "river", "ocean", "forest", "desert", "island", "valley",
			"castle", "tower", "bridge", "temple", "palace", "cottage", "cabin",
			"diamond", "ruby", "emerald", "sapphire", "pearl", "crystal", "gem",
		},
		usedCodes:       make(map[string]time.Time),
		rateLimit:       config.GetTokenRateLimit(),
		cleanupInterval: config.GetTokenCleanupInterval(),
		config:          config,
	}
}

func (tg *TokenGenerator) GenerateShareCode() (string, error) {
	tg.mutex.Lock()
	if time.Since(tg.lastGeneration) < tg.rateLimit {
		tg.mutex.Unlock()
		return "", fmt.Errorf("rate limit exceeded, please wait before generating another code")
	}
	tg.lastGeneration = time.Now()
	tg.mutex.Unlock()

	maxAttempts := 100
	for attempt := 0; attempt < maxAttempts; attempt++ {
		code, err := tg.generateCode()
		if err != nil {
			return "", fmt.Errorf("failed to generate code: %w", err)
		}

		if !tg.isCodeUsed(code) {
			tg.markCodeAsUsed(code)
			return code, nil
		}
	}

	return "", fmt.Errorf("unable to generate unique code after %d attempts", maxAttempts)
}

func (tg *TokenGenerator) generateCode() (string, error) {
	adjIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(tg.adjectives))))
	if err != nil {
		return "", fmt.Errorf("failed to generate random adjective index: %w", err)
	}

	nounIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(tg.nouns))))
	if err != nil {
		return "", fmt.Errorf("failed to generate random noun index: %w", err)
	}

	number, err := rand.Int(rand.Reader, big.NewInt(999))
	if err != nil {
		return "", fmt.Errorf("failed to generate random number: %w", err)
	}

	adjective := tg.adjectives[adjIndex.Int64()]
	noun := tg.nouns[nounIndex.Int64()]
	numberStr := fmt.Sprintf("%d", number.Int64()+1)

	return fmt.Sprintf("%s-%s-%s", adjective, noun, numberStr), nil
}

func (tg *TokenGenerator) isCodeUsed(code string) bool {
	tg.mutex.RLock()
	defer tg.mutex.RUnlock()

	usedTime, exists := tg.usedCodes[code]
	if !exists {
		return false
	}

	if time.Since(usedTime) > tg.cleanupInterval {
		delete(tg.usedCodes, code)
		return false
	}

	return true
}

func (tg *TokenGenerator) markCodeAsUsed(code string) {
	tg.mutex.Lock()
	defer tg.mutex.Unlock()
	tg.usedCodes[code] = time.Now()
}

func (tg *TokenGenerator) ValidateShareCode(code string) error {
	if code == "" {
		return fmt.Errorf("share code cannot be empty")
	}

	parts := strings.Split(code, "-")
	if len(parts) != 3 {
		return fmt.Errorf("share code must be in format: word-word-number")
	}

	if !tg.isValidWord(strings.ToLower(parts[0]), tg.adjectives) {
		return fmt.Errorf("invalid adjective in share code")
	}

	if !tg.isValidWord(strings.ToLower(parts[1]), tg.nouns) {
		return fmt.Errorf("invalid noun in share code")
	}

	if !tg.isValidNumber(parts[2]) {
		return fmt.Errorf("invalid number in share code (must be 1-999)")
	}

	return nil
}

func (tg *TokenGenerator) isValidWord(word string, wordList []string) bool {
	word = strings.ToLower(strings.TrimSpace(word))
	for _, validWord := range wordList {
		if validWord == word {
			return true
		}
	}
	return false
}

func (tg *TokenGenerator) isValidNumber(numStr string) bool {
	if len(numStr) == 0 || len(numStr) > 3 {
		return false
	}

	for _, char := range numStr {
		if char < '0' || char > '9' {
			return false
		}
	}

	if numStr == "0" || len(numStr) > 3 {
		return false
	}

	return true
}

func (tg *TokenGenerator) GetStats() map[string]interface{} {
	tg.mutex.RLock()
	defer tg.mutex.RUnlock()

	activeCodes := 0
	for _, usedTime := range tg.usedCodes {
		if time.Since(usedTime) <= tg.cleanupInterval {
			activeCodes++
		}
	}

	return map[string]interface{}{
		"total_adjectives": len(tg.adjectives),
		"total_nouns":      len(tg.nouns),
		"active_codes":     activeCodes,
		"rate_limit_ms":    tg.rateLimit.Milliseconds(),
	}
}

func (tg *TokenGenerator) Cleanup() {
	tg.mutex.Lock()
	defer tg.mutex.Unlock()

	now := time.Now()
	for code, usedTime := range tg.usedCodes {
		if now.Sub(usedTime) > tg.cleanupInterval {
			delete(tg.usedCodes, code)
		}
	}
}

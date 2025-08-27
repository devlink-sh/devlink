package util

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type SecurityConfig struct {
	offensiveWords        map[string]bool
	mutex                 sync.RWMutex
	defaultOffensiveWords []string
}

func NewSecurityConfig() *SecurityConfig {
	return &SecurityConfig{
		offensiveWords: make(map[string]bool),
		defaultOffensiveWords: []string{
			"test", "demo", "temp", "fake", "invalid", "null", "undefined",
		},
	}
}

func (sc *SecurityConfig) LoadOffensiveWords(configPath string) error {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	sc.offensiveWords = make(map[string]bool)

	if configPath != "" {
		if err := sc.loadFromFile(configPath); err != nil {
			fmt.Printf("Warning: Could not load offensive words from %s: %v. Using defaults.\n", configPath, err)
		}
	}

	for _, word := range sc.defaultOffensiveWords {
		sc.offensiveWords[strings.ToLower(word)] = true
	}

	return nil
}

func (sc *SecurityConfig) loadFromFile(filePath string) error {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("invalid file path: %w", err)
	}

	if strings.Contains(absPath, "..") {
		return fmt.Errorf("path traversal not allowed")
	}

	file, err := os.Open(absPath)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
		if lineCount > 10000 {
			return fmt.Errorf("config file too large (max 10000 lines)")
		}

		word := strings.TrimSpace(strings.ToLower(scanner.Text()))
		if word != "" && !strings.HasPrefix(word, "#") {
			if len(word) > 50 {
				continue
			}
			sc.offensiveWords[word] = true
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	return nil
}

func (sc *SecurityConfig) LoadFromEnvironment() error {
	configPath := os.Getenv("DEVLINK_OFFENSIVE_WORDS_FILE")
	return sc.LoadOffensiveWords(configPath)
}

func (sc *SecurityConfig) IsOffensiveWord(word string) bool {
	sc.mutex.RLock()
	defer sc.mutex.RUnlock()

	word = strings.ToLower(strings.TrimSpace(word))
	return sc.offensiveWords[word]
}

func (sc *SecurityConfig) GetOffensiveWordsCount() int {
	sc.mutex.RLock()
	defer sc.mutex.RUnlock()
	return len(sc.offensiveWords)
}

func (sc *SecurityConfig) AddOffensiveWord(word string) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	word = strings.ToLower(strings.TrimSpace(word))
	if word != "" && len(word) <= 50 {
		sc.offensiveWords[word] = true
	}
}

func (sc *SecurityConfig) RemoveOffensiveWord(word string) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	word = strings.ToLower(strings.TrimSpace(word))
	delete(sc.offensiveWords, word)
}

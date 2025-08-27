package util

import (
	"fmt"
	"os"
	"testing"
)

func TestNewSecurityConfig(t *testing.T) {
	sc := NewSecurityConfig()
	if sc == nil {
		t.Fatal("NewSecurityConfig() returned nil")
	}

	if sc.offensiveWords == nil {
		t.Error("SecurityConfig should have offensiveWords map")
	}

	if len(sc.defaultOffensiveWords) == 0 {
		t.Error("SecurityConfig should have default offensive words")
	}
}

func TestLoadOffensiveWords_NoFile(t *testing.T) {
	sc := NewSecurityConfig()

	// Test with empty path - should use defaults
	err := sc.LoadOffensiveWords("")
	if err != nil {
		t.Errorf("LoadOffensiveWords with empty path should not error: %v", err)
	}

	// Should have at least the default words
	if sc.GetOffensiveWordsCount() == 0 {
		t.Error("Should have default offensive words loaded")
	}

	// Test default words
	for _, word := range sc.defaultOffensiveWords {
		if !sc.IsOffensiveWord(word) {
			t.Errorf("Default word '%s' should be offensive", word)
		}
	}
}

func TestLoadOffensiveWords_FromFile(t *testing.T) {
	sc := NewSecurityConfig()

	// Create a temporary file with test words
	tmpFile, err := os.CreateTemp("", "offensive_words_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write test words
	testWords := "badword1\nbadword2\n# This is a comment\nbadword3\n\n"
	_, err = tmpFile.WriteString(testWords)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Load from file
	err = sc.LoadOffensiveWords(tmpFile.Name())
	if err != nil {
		t.Errorf("LoadOffensiveWords should not error: %v", err)
	}

	// Check that words were loaded
	if !sc.IsOffensiveWord("badword1") {
		t.Error("badword1 should be loaded as offensive")
	}

	if !sc.IsOffensiveWord("badword2") {
		t.Error("badword2 should be loaded as offensive")
	}

	if !sc.IsOffensiveWord("badword3") {
		t.Error("badword3 should be loaded as offensive")
	}

	// Comments should be ignored
	if sc.IsOffensiveWord("# This is a comment") {
		t.Error("Comments should not be loaded as offensive words")
	}
}

func TestLoadOffensiveWords_PathTraversal(t *testing.T) {
	sc := NewSecurityConfig()

	// Test path traversal attack - these should gracefully handle failures
	// and fall back to defaults rather than returning errors
	maliciousPaths := []string{
		"../../../etc/passwd",
		"..\\..\\..\\windows\\system32\\config\\sam",
		"./config/../../../etc/shadow",
	}

	for _, path := range maliciousPaths {
		err := sc.LoadOffensiveWords(path)
		// The LoadOffensiveWords function logs warnings but doesn't return errors
		// It falls back to defaults for security
		if err != nil {
			t.Errorf("LoadOffensiveWords should handle errors gracefully: %s", path)
		}

		// Should still have default words loaded
		if sc.GetOffensiveWordsCount() == 0 {
			t.Errorf("Should have default words after path traversal attempt: %s", path)
		}
	}
}

func TestLoadOffensiveWords_FileSize(t *testing.T) {
	sc := NewSecurityConfig()

	// Create a file with too many lines
	tmpFile, err := os.CreateTemp("", "large_offensive_words_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write many lines to trigger size limit
	for i := 0; i < 10001; i++ {
		tmpFile.WriteString("word\n")
	}
	tmpFile.Close()

	// Should handle oversized file gracefully and fall back to defaults
	err = sc.LoadOffensiveWords(tmpFile.Name())
	if err != nil {
		t.Error("LoadOffensiveWords should handle oversized files gracefully")
	}

	// Should still have default words loaded
	if sc.GetOffensiveWordsCount() == 0 {
		t.Error("Should have default words after oversized file attempt")
	}
}

func TestLoadFromEnvironment(t *testing.T) {
	sc := NewSecurityConfig()

	// Test without environment variable
	err := sc.LoadFromEnvironment()
	if err != nil {
		t.Errorf("LoadFromEnvironment should not error when env var is not set: %v", err)
	}

	// Should have default words
	if sc.GetOffensiveWordsCount() == 0 {
		t.Error("Should have default words when no env file specified")
	}
}

func TestIsOffensiveWord(t *testing.T) {
	sc := NewSecurityConfig()
	sc.LoadOffensiveWords("")

	// Test case insensitivity
	sc.AddOffensiveWord("TESTWORD")
	if !sc.IsOffensiveWord("testword") {
		t.Error("Offensive word check should be case insensitive")
	}

	if !sc.IsOffensiveWord("TESTWORD") {
		t.Error("Offensive word check should work with uppercase")
	}

	if !sc.IsOffensiveWord("TestWord") {
		t.Error("Offensive word check should work with mixed case")
	}

	// Test with whitespace
	if !sc.IsOffensiveWord("  testword  ") {
		t.Error("Offensive word check should trim whitespace")
	}

	// Test non-offensive word
	if sc.IsOffensiveWord("nonoffensiveword") {
		t.Error("Non-offensive word should not be flagged")
	}
}

func TestAddOffensiveWord(t *testing.T) {
	sc := NewSecurityConfig()
	sc.LoadOffensiveWords("")

	// Add a new word
	sc.AddOffensiveWord("newbadword")
	if !sc.IsOffensiveWord("newbadword") {
		t.Error("Added word should be offensive")
	}

	// Test length limit
	longWord := "verylongwordthatexceedsthelimitofcharactersandshouldbeignored"
	sc.AddOffensiveWord(longWord)
	if sc.IsOffensiveWord(longWord) {
		t.Error("Overly long words should not be added")
	}

	// Test empty word
	sc.AddOffensiveWord("")
	if sc.IsOffensiveWord("") {
		t.Error("Empty words should not be added")
	}
}

func TestRemoveOffensiveWord(t *testing.T) {
	sc := NewSecurityConfig()
	sc.LoadOffensiveWords("")

	// Add a word
	sc.AddOffensiveWord("removeme")
	if !sc.IsOffensiveWord("removeme") {
		t.Error("Word should be offensive after adding")
	}

	// Remove the word
	sc.RemoveOffensiveWord("removeme")
	if sc.IsOffensiveWord("removeme") {
		t.Error("Word should not be offensive after removing")
	}
}

func TestGetOffensiveWordsCount(t *testing.T) {
	sc := NewSecurityConfig()
	sc.LoadOffensiveWords("")

	initialCount := sc.GetOffensiveWordsCount()
	if initialCount == 0 {
		t.Error("Should have at least default words")
	}

	// Add a word
	sc.AddOffensiveWord("testword")
	newCount := sc.GetOffensiveWordsCount()
	if newCount != initialCount+1 {
		t.Error("Count should increase after adding word")
	}

	// Remove a word
	sc.RemoveOffensiveWord("testword")
	finalCount := sc.GetOffensiveWordsCount()
	if finalCount != initialCount {
		t.Error("Count should return to original after removing word")
	}
}

func TestConcurrentAccess(t *testing.T) {
	sc := NewSecurityConfig()
	sc.LoadOffensiveWords("")

	// Test concurrent read/write operations
	done := make(chan bool, 10)

	// Concurrent readers
	for i := 0; i < 5; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				sc.IsOffensiveWord("test")
			}
			done <- true
		}()
	}

	// Concurrent writers
	for i := 0; i < 5; i++ {
		go func(id int) {
			for j := 0; j < 10; j++ {
				word := fmt.Sprintf("testword%d_%d", id, j)
				sc.AddOffensiveWord(word)
				sc.RemoveOffensiveWord(word)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

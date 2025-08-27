package env

import (
	"os"
	"testing"
)

func TestNewParser(t *testing.T) {
	parser := NewParser()
	if parser == nil {
		t.Fatal("NewParser() returned nil")
	}
}

func TestParseContent_ValidFile(t *testing.T) {
	content := `# Test file
DATABASE_URL=postgresql://localhost:5432/mydb
API_KEY=sk-1234567890abcdef
NODE_ENV=development

# Empty line above
REDIS_URL=redis://localhost:6379`

	parser := NewParser()
	envFile, err := parser.ParseContent(content, "test.env")
	if err != nil {
		t.Fatalf("ParseContent failed: %v", err)
	}

	// Check statistics
	if envFile.TotalLines != 7 {
		t.Errorf("Expected 7 total lines, got %d", envFile.TotalLines)
	}
	if envFile.ValidLines != 4 {
		t.Errorf("Expected 4 valid lines, got %d", envFile.ValidLines)
	}
	if envFile.CommentLines != 2 {
		t.Errorf("Expected 2 comment lines, got %d", envFile.CommentLines)
	}
	if envFile.EmptyLines != 1 {
		t.Errorf("Expected 1 empty line, got %d", envFile.EmptyLines)
	}

	// Check variables
	if len(envFile.Variables) != 4 {
		t.Errorf("Expected 4 variables, got %d", len(envFile.Variables))
	}

	// Check specific variables
	expectedVars := map[string]string{
		"DATABASE_URL": "postgresql://localhost:5432/mydb",
		"API_KEY":      "sk-1234567890abcdef",
		"NODE_ENV":     "development",
		"REDIS_URL":    "redis://localhost:6379",
	}

	for _, variable := range envFile.Variables {
		if expectedValue, exists := expectedVars[variable.Key]; exists {
			if variable.Value != expectedValue {
				t.Errorf("Expected %s=%s, got %s", variable.Key, expectedValue, variable.Value)
			}
		} else {
			t.Errorf("Unexpected variable: %s", variable.Key)
		}
	}
}

func TestParseContent_WithQuotes(t *testing.T) {
	content := `QUOTED_STRING="Hello World"
SINGLE_QUOTED='Single quoted string'
ESCAPED_QUOTES="He said \"Hello\""
MIXED_QUOTES='He said "Hello"'`

	parser := NewParser()
	envFile, err := parser.ParseContent(content, "test.env")
	if err != nil {
		t.Fatalf("ParseContent failed: %v", err)
	}

	expectedVars := map[string]string{
		"QUOTED_STRING":  "Hello World",
		"SINGLE_QUOTED":  "Single quoted string",
		"ESCAPED_QUOTES": "He said \"Hello\"",
		"MIXED_QUOTES":   "He said \"Hello\"",
	}

	for _, variable := range envFile.Variables {
		if expectedValue, exists := expectedVars[variable.Key]; exists {
			if variable.Value != expectedValue {
				t.Errorf("Expected %s=%s, got %s", variable.Key, expectedValue, variable.Value)
			}
		}
	}
}

func TestParseContent_WithComments(t *testing.T) {
	content := `DATABASE_URL=postgresql://localhost:5432/mydb # Database connection
API_KEY=sk-1234567890abcdef # API key for external service
NODE_ENV=development # Environment`

	parser := NewParser()
	envFile, err := parser.ParseContent(content, "test.env")
	if err != nil {
		t.Fatalf("ParseContent failed: %v", err)
	}

	// Check that comments are extracted
	for _, variable := range envFile.Variables {
		if variable.Key == "DATABASE_URL" && variable.Comment != "# Database connection" {
			t.Errorf("Expected comment '# Database connection', got '%s'", variable.Comment)
		}
		if variable.Key == "API_KEY" && variable.Comment != "# API key for external service" {
			t.Errorf("Expected comment '# API key for external service', got '%s'", variable.Comment)
		}
	}
}

func TestParseContent_InvalidLines(t *testing.T) {
	content := `VALID_VAR=value
INVALID_LINE_NO_EQUALS
ANOTHER_VALID=value
=no_key
KEY_ONLY=
`

	parser := NewParser()
	envFile, err := parser.ParseContent(content, "test.env")
	if err != nil {
		t.Fatalf("ParseContent failed: %v", err)
	}

	// Should have parse errors
	if len(envFile.ParseErrors) == 0 {
		t.Error("Expected parse errors, got none")
	}

	// Should still have valid variables
	if len(envFile.Variables) != 3 {
		t.Errorf("Expected 3 valid variables, got %d", len(envFile.Variables))
	}
}

func TestParseContent_SensitiveDetection(t *testing.T) {
	content := `DATABASE_URL=postgresql://localhost:5432/mydb
API_KEY=sk-1234567890abcdef
JWT_SECRET=my-secret-key
PASSWORD=password123
NODE_ENV=development
REDIS_PASSWORD=redis_pass`

	parser := NewParser()
	envFile, err := parser.ParseContent(content, "test.env")
	if err != nil {
		t.Fatalf("ParseContent failed: %v", err)
	}

	// Check sensitive variable detection
	sensitiveCount := 0
	for _, variable := range envFile.Variables {
		if variable.IsSensitive {
			sensitiveCount++
		}
	}

	if sensitiveCount < 4 {
		t.Errorf("Expected at least 4 sensitive variables, got %d", sensitiveCount)
	}
}

func TestParseFile_FileNotFound(t *testing.T) {
	parser := NewParser()
	_, err := parser.ParseFile("nonexistent.env")
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

func TestParseFile_ValidFile(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "test*.env")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := `DATABASE_URL=postgresql://localhost:5432/mydb
API_KEY=sk-1234567890abcdef`

	_, err = tmpFile.WriteString(content)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	parser := NewParser()
	envFile, err := parser.ParseFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	if len(envFile.Variables) != 2 {
		t.Errorf("Expected 2 variables, got %d", len(envFile.Variables))
	}
}

func TestValidateKey(t *testing.T) {
	parser := NewParser()

	// Valid keys
	validKeys := []string{"KEY", "key", "Key123", "_KEY", "KEY_123"}
	for _, key := range validKeys {
		if err := parser.validateKey(key); err != nil {
			t.Errorf("Key '%s' should be valid: %v", key, err)
		}
	}

	// Invalid keys
	invalidKeys := []string{"", "123KEY", "KEY-123", "KEY.123", "KEY 123"}
	for _, key := range invalidKeys {
		if err := parser.validateKey(key); err == nil {
			t.Errorf("Key '%s' should be invalid", key)
		}
	}
}

func TestParseValue(t *testing.T) {
	parser := NewParser()

	testCases := []struct {
		input    string
		expected string
		comment  string
		hasError bool
	}{
		{"", "", "", false},
		{"value", "value", "", false},
		{"value # comment", "value", "# comment", false},
		{`"quoted value"`, "quoted value", "", false},
		{`'single quoted'`, "single quoted", "", false},
		{`"quoted # not comment"`, "quoted # not comment", "", false},
		{`"unclosed quote`, "", "", true},
		{`value\#notcomment`, "value\\#notcomment", "", false},
	}

	for _, tc := range testCases {
		value, comment, err := parser.parseValue(tc.input)
		if tc.hasError && err == nil {
			t.Errorf("Expected error for input '%s', got nil", tc.input)
		}
		if !tc.hasError && err != nil {
			t.Errorf("Unexpected error for input '%s': %v", tc.input, err)
		}
		if !tc.hasError {
			if value != tc.expected {
				t.Errorf("For input '%s', expected value '%s', got '%s'", tc.input, tc.expected, value)
			}
			if comment != tc.comment {
				t.Errorf("For input '%s', expected comment '%s', got '%s'", tc.input, tc.comment, comment)
			}
		}
	}
}

func TestIsSensitiveVariable(t *testing.T) {
	parser := NewParser()

	sensitiveKeys := []string{
		"PASSWORD", "API_KEY", "JWT_SECRET", "DATABASE_URL",
		"REDIS_PASSWORD", "AWS_SECRET_KEY", "PRIVATE_KEY",
	}

	for _, key := range sensitiveKeys {
		if !parser.isSensitiveVariable(key) {
			t.Errorf("Key '%s' should be detected as sensitive", key)
		}
	}

	nonSensitiveKeys := []string{
		"NODE_ENV", "DEBUG", "PORT", "HOST", "LOG_LEVEL",
	}

	for _, key := range nonSensitiveKeys {
		if parser.isSensitiveVariable(key) {
			t.Errorf("Key '%s' should not be detected as sensitive", key)
		}
	}
}

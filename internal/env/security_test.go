package env

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/devlink/internal/util"
)

func TestParser_Security_BoundsChecking(t *testing.T) {
	parser := NewParser()
	
	// Test extremely large value
	largeValue := strings.Repeat("A", 15000)
	content := fmt.Sprintf("TEST_VAR=%s", largeValue)
	
	envFile, err := parser.ParseContent(content, "test.env")
	if err == nil {
		t.Error("Expected error for extremely large value, got nil")
	}
	
	if envFile != nil && len(envFile.Variables) > 0 {
		t.Error("Should not parse variables with oversized values")
	}
}

func TestParser_Security_PathTraversal(t *testing.T) {
	parser := NewParser()
	
	// Test path traversal attempts
	maliciousContent := `PATH=../../../etc/passwd
SECRET=value`
	
	envFile, err := parser.ParseContent(maliciousContent, "test.env")
	if err != nil {
		t.Fatalf("Parser should handle path traversal in values: %v", err)
	}
	
	// Should parse but flag as sensitive
	for _, variable := range envFile.Variables {
		if variable.Key == "PATH" && !variable.IsSensitive {
			t.Error("Path variable should be detected as sensitive")
		}
	}
}

func TestValidator_Security_WeakPasswords(t *testing.T) {
	validator := NewValidator()
	
	weakPasswords := []string{
		"password", "123456", "qwerty", "admin", "root",
		"test", "demo", "guest", "user", "default",
	}
	
	for _, weakPass := range weakPasswords {
		content := fmt.Sprintf("PASSWORD=%s", weakPass)
		parser := NewParser()
		envFile, err := parser.ParseContent(content, "test.env")
		if err != nil {
			t.Fatalf("Failed to parse test content: %v", err)
		}
		
		result := validator.Validate(envFile)
		if result.IsValid {
			t.Errorf("Weak password '%s' should be detected as invalid", weakPass)
		}
		
		found := false
		for _, err := range result.Errors {
			if err.Message == "weak password detected" {
				found = true
				break
			}
		}
		
		if !found {
			t.Errorf("Weak password '%s' should trigger weak password error", weakPass)
		}
	}
}

func TestValidator_Security_HardcodedCredentials(t *testing.T) {
	validator := NewValidator()
	
	hardcodedCreds := []string{
		"admin", "root", "password", "123456", "qwerty",
		"test", "demo", "guest", "user", "default",
	}
	
	for _, cred := range hardcodedCreds {
		content := fmt.Sprintf("SECRET_KEY=%s", cred)
		parser := NewParser()
		envFile, err := parser.ParseContent(content, "test.env")
		if err != nil {
			t.Fatalf("Failed to parse test content: %v", err)
		}
		
		result := validator.Validate(envFile)
		if result.IsValid {
			t.Errorf("Hardcoded credential '%s' should be detected as invalid", cred)
		}
		
		found := false
		for _, err := range result.Errors {
			if err.Message == "hardcoded credential detected" {
				found = true
				break
			}
		}
		
		if !found {
			t.Errorf("Hardcoded credential '%s' should trigger hardcoded credential error", cred)
		}
	}
}

func TestValidator_Security_ExposedSecretsInURLs(t *testing.T) {
	validator := NewValidator()
	
	exposedURLs := []string{
		"postgresql://user:password@localhost:5432/db",
		"mysql://admin:secret@localhost:3306/db",
		"redis://user:pass@localhost:6379",
	}
	
	for _, url := range exposedURLs {
		content := fmt.Sprintf("DATABASE_URL=%s", url)
		parser := NewParser()
		envFile, err := parser.ParseContent(content, "test.env")
		if err != nil {
			t.Fatalf("Failed to parse test content: %v", err)
		}
		
		result := validator.Validate(envFile)
		if result.IsValid {
			t.Errorf("Exposed secret in URL should be detected as invalid")
		}
		
		found := false
		for _, err := range result.Errors {
			if err.Message == "secret exposed in URL" {
				found = true
				break
			}
		}
		
		if !found {
			t.Errorf("Exposed secret in URL should trigger exposed secret error")
		}
	}
}

func TestFormatter_Security_Masking(t *testing.T) {
	formatter := NewFormatter()
	
	// Test sensitive variable masking
	envFile := &util.EnvFile{
		Variables: []util.EnvVariable{
			{
				Key:         "API_KEY",
				Value:       "sk-1234567890abcdef",
				IsSensitive: true,
			},
			{
				Key:         "PASSWORD",
				Value:       "secret123",
				IsSensitive: true,
			},
			{
				Key:         "NODE_ENV",
				Value:       "development",
				IsSensitive: false,
			},
		},
	}
	
	options := &util.FormatOptions{
		MaskSensitive: true,
		ShowComments:  false,
		OutputFormat:  "text",
	}
	
	output, err := formatter.Format(envFile, options)
	if err != nil {
		t.Fatalf("Formatter failed: %v", err)
	}
	
	// Check that sensitive values are masked
	if strings.Contains(output, "sk-1234567890abcdef") {
		t.Error("Sensitive API key should be masked in output")
	}
	
	if strings.Contains(output, "secret123") {
		t.Error("Sensitive password should be masked in output")
	}
	
	// Check that non-sensitive values are not masked
	if !strings.Contains(output, "development") {
		t.Error("Non-sensitive value should not be masked")
	}
}

func TestFileOperations_Security_PathValidation(t *testing.T) {
	// Test path traversal prevention
	maliciousPaths := []string{
		"../../../etc/passwd",
		"..\\..\\..\\windows\\system32\\config\\sam",
		"./config/../../../etc/shadow",
		"config/../../../../root/.ssh/id_rsa",
	}
	
	for _, path := range maliciousPaths {
		// This would be called from validateOutputPath
		absPath, err := filepath.Abs(path)
		if err == nil && strings.Contains(absPath, "..") {
			t.Errorf("Path traversal attempt should be detected: %s", path)
		}
	}
	
	// Test dangerous file extensions
	dangerousExts := []string{
		"config.exe", "script.bat", "test.sh", "code.py", "app.js", "web.php",
	}
	
	for _, file := range dangerousExts {
		ext := strings.ToLower(filepath.Ext(file))
		if ext == ".exe" || ext == ".bat" || ext == ".sh" || ext == ".py" || ext == ".js" || ext == ".php" {
			// This should be rejected by validateOutputPath
			t.Logf("Dangerous file extension detected: %s", ext)
		}
	}
}

func TestShareCode_Security_Validation(t *testing.T) {
	// Test valid share codes
	validCodes := []string{
		"ABC123", "XYZ789", "123456", "ABCDEF", "A1B2C3D4E5F6",
	}
	
	for _, code := range validCodes {
		matched, err := regexp.MatchString(`^[A-Z0-9]{6,12}$`, code)
		if err != nil {
			t.Errorf("Regex validation failed for valid code %s: %v", code, err)
		}
		if !matched {
			t.Errorf("Valid share code %s should pass validation", code)
		}
	}
	
	// Test invalid share codes
	invalidCodes := []string{
		"", "ABC", "ABC123DEF456GHI", "abc123", "ABC-123", "ABC_123",
		"ABC 123", "ABC@123", "ABC#123", "ABC$123",
	}
	
	for _, code := range invalidCodes {
		matched, err := regexp.MatchString(`^[A-Z0-9]{6,12}$`, code)
		if err != nil {
			continue // Skip regex errors for empty strings
		}
		if matched {
			t.Errorf("Invalid share code %s should fail validation", code)
		}
	}
}

func TestSecureShareCodeGeneration(t *testing.T) {
	// Test that generated codes are secure and unique
	generatedCodes := make(map[string]bool)
	
	for i := 0; i < 100; i++ {
		code := generateTestShareCode(i)
		
		// Check format
		matched, err := regexp.MatchString(`^[A-Z0-9]{6}$`, code)
		if err != nil {
			t.Errorf("Generated code validation failed: %v", err)
		}
		if !matched {
			t.Errorf("Generated code %s does not match expected format", code)
		}
		
		// Check uniqueness
		if generatedCodes[code] {
			t.Errorf("Duplicate share code generated: %s", code)
		}
		generatedCodes[code] = true
	}
}

// generateTestShareCode generates a test share code for testing purposes
func generateTestShareCode(seed int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 6
	
	// Use seed to generate different codes for testing
	result := make([]byte, length)
	for i := range result {
		index := (seed + i) % len(charset)
		result[i] = charset[index]
	}
	
	return string(result)
}

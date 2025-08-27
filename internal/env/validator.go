package env

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/devlink/internal/util"
)

// Validator handles .env file security validation
type Validator struct {
	// Patterns for detecting various security issues
	weakPasswordPatterns []*regexp.Regexp
	urlPatterns          []*regexp.Regexp
	emailPattern         *regexp.Regexp
	ipPattern            *regexp.Regexp
	// Minimum requirements
	minPasswordLength int
	maxValueLength    int
}

// NewValidator creates a new .env validator
func NewValidator() *Validator {
	return &Validator{
		weakPasswordPatterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)^(password|123456|qwerty|admin|root|test|demo)$`),
			regexp.MustCompile(`(?i)^(password|pass|pwd)\d*$`),
			regexp.MustCompile(`^[a-zA-Z0-9]{1,6}$`), // Very short passwords
		},
		urlPatterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)^https?://`),
			regexp.MustCompile(`(?i)^(postgresql|mysql|mongodb|redis)://`),
			regexp.MustCompile(`(?i)^(ftp|sftp|ssh)://`),
		},
		emailPattern:      regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
		ipPattern:         regexp.MustCompile(`^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`),
		minPasswordLength: 8,
		maxValueLength:    1000,
	}
}

// Validate validates an .env file for security issues
func (v *Validator) Validate(envFile *util.EnvFile) *util.ValidationResult {
	result := &util.ValidationResult{
		IsValid:       true,
		Errors:        []util.ValidationError{},
		Warnings:      []util.ValidationWarning{},
		SensitiveVars: []string{},
		RiskLevel:     util.RiskLevelLow,
	}

	// Check for parse errors first
	if len(envFile.ParseErrors) > 0 {
		result.IsValid = false
		result.RiskLevel = util.RiskLevelHigh
		for _, parseError := range envFile.ParseErrors {
			result.Errors = append(result.Errors, util.ValidationError{
				Variable:   "parse_error",
				Message:    parseError.Message,
				Severity:   "high",
				LineNumber: parseError.LineNumber,
			})
		}
	}

	// Validate each variable
	for _, variable := range envFile.Variables {
		v.validateVariable(variable, result)
	}

	// Determine overall risk level
	v.determineRiskLevel(result)

	return result
}

// validateVariable validates a single environment variable
func (v *Validator) validateVariable(variable util.EnvVariable, result *util.ValidationResult) {
	// Track sensitive variables
	if variable.IsSensitive {
		result.SensitiveVars = append(result.SensitiveVars, variable.Key)
	}

	// Check for empty values in sensitive variables
	if variable.IsSensitive && variable.Value == "" {
		result.Warnings = append(result.Warnings, util.ValidationWarning{
			Variable:   variable.Key,
			Message:    "sensitive variable has empty value",
			LineNumber: variable.LineNumber,
		})
	}

	// Check value length
	if len(variable.Value) > v.maxValueLength {
		result.Errors = append(result.Errors, util.ValidationError{
			Variable:   variable.Key,
			Message:    fmt.Sprintf("value too long (%d characters, max %d)", len(variable.Value), v.maxValueLength),
			Severity:   "medium",
			LineNumber: variable.LineNumber,
		})
		result.IsValid = false
	}

	// Validate sensitive variables more strictly
	if variable.IsSensitive {
		v.validateSensitiveVariable(variable, result)
	}

	// Check for common security issues
	v.checkCommonSecurityIssues(variable, result)
}

// validateSensitiveVariable performs additional validation for sensitive variables
func (v *Validator) validateSensitiveVariable(variable util.EnvVariable, result *util.ValidationResult) {
	// Check for weak passwords
	if v.isWeakPassword(variable.Key, variable.Value) {
		result.Errors = append(result.Errors, util.ValidationError{
			Variable:   variable.Key,
			Message:    "weak password detected",
			Severity:   "high",
			LineNumber: variable.LineNumber,
		})
		result.IsValid = false
	}

	// Check for hardcoded credentials
	if v.isHardcodedCredential(variable.Value) {
		result.Errors = append(result.Errors, util.ValidationError{
			Variable:   variable.Key,
			Message:    "hardcoded credential detected",
			Severity:   "critical",
			LineNumber: variable.LineNumber,
		})
		result.IsValid = false
	}

	// Check for exposed secrets in URLs
	if v.containsExposedSecret(variable.Value) {
		result.Errors = append(result.Errors, util.ValidationError{
			Variable:   variable.Key,
			Message:    "secret exposed in URL",
			Severity:   "critical",
			LineNumber: variable.LineNumber,
		})
		result.IsValid = false
	}

	// Check password complexity for password variables
	if v.isPasswordVariable(variable.Key) {
		v.validatePasswordComplexity(variable, result)
	}
}

// checkCommonSecurityIssues checks for common security problems
func (v *Validator) checkCommonSecurityIssues(variable util.EnvVariable, result *util.ValidationResult) {
	// Check for development values in production-like variables
	if v.isDevelopmentValue(variable.Value) && v.isProductionVariable(variable.Key) {
		result.Warnings = append(result.Warnings, util.ValidationWarning{
			Variable:   variable.Key,
			Message:    "development value in production variable",
			LineNumber: variable.LineNumber,
		})
	}

	// Check for localhost in database URLs
	if v.isDatabaseVariable(variable.Key) && strings.Contains(strings.ToLower(variable.Value), "localhost") {
		result.Warnings = append(result.Warnings, util.ValidationWarning{
			Variable:   variable.Key,
			Message:    "localhost detected in database URL",
			LineNumber: variable.LineNumber,
		})
	}

	// Check for HTTP URLs in production
	if v.isProductionVariable(variable.Key) && strings.HasPrefix(strings.ToLower(variable.Value), "http://") {
		result.Errors = append(result.Errors, util.ValidationError{
			Variable:   variable.Key,
			Message:    "HTTP URL in production variable (use HTTPS)",
			Severity:   "high",
			LineNumber: variable.LineNumber,
		})
		result.IsValid = false
	}
}

// isWeakPassword checks if a password is weak
func (v *Validator) isWeakPassword(key, value string) bool {
	// Skip if not a password variable
	if !v.isPasswordVariable(key) {
		return false
	}

	// Check against weak password patterns
	for _, pattern := range v.weakPasswordPatterns {
		if pattern.MatchString(value) {
			return true
		}
	}

	// Check length
	if len(value) < v.minPasswordLength {
		return true
	}

	// Check complexity
	return !v.hasGoodComplexity(value)
}

// isPasswordVariable checks if a variable is a password
func (v *Validator) isPasswordVariable(key string) bool {
	upperKey := strings.ToUpper(key)
	return strings.Contains(upperKey, "PASSWORD") ||
		strings.Contains(upperKey, "PASSWD") ||
		strings.Contains(upperKey, "PWD") ||
		strings.Contains(upperKey, "SECRET")
}

// isHardcodedCredential checks if a value looks like a hardcoded credential
func (v *Validator) isHardcodedCredential(value string) bool {
	hardcodedPatterns := []string{
		"admin", "root", "password", "123456", "qwerty",
		"test", "demo", "guest", "user", "default",
	}

	lowerValue := strings.ToLower(value)
	for _, pattern := range hardcodedPatterns {
		if lowerValue == pattern {
			return true
		}
	}

	return false
}

// containsExposedSecret checks if a URL contains exposed secrets
func (v *Validator) containsExposedSecret(value string) bool {
	// Check for passwords in URLs
	if strings.Contains(value, "://") && strings.Contains(value, "@") {
		// Extract the part between :// and @
		parts := strings.Split(value, "@")
		if len(parts) > 1 {
			authPart := strings.Split(parts[0], "://")
			if len(authPart) > 1 {
				credentials := authPart[1]
				if strings.Contains(credentials, ":") {
					// This looks like username:password in URL
					return true
				}
			}
		}
	}

	return false
}

// isDevelopmentValue checks if a value looks like a development value
func (v *Validator) isDevelopmentValue(value string) bool {
	devPatterns := []string{
		"localhost", "127.0.0.1", "dev", "development", "test",
		"demo", "example", "temp", "tmp", "fake", "mock",
	}

	lowerValue := strings.ToLower(value)
	for _, pattern := range devPatterns {
		if strings.Contains(lowerValue, pattern) {
			return true
		}
	}

	return false
}

// isProductionVariable checks if a variable looks like a production variable
func (v *Validator) isProductionVariable(key string) bool {
	prodPatterns := []string{
		"PROD", "PRODUCTION", "LIVE", "REAL", "ACTUAL",
	}

	upperKey := strings.ToUpper(key)
	for _, pattern := range prodPatterns {
		if strings.Contains(upperKey, pattern) {
			return true
		}
	}

	return false
}

// isDatabaseVariable checks if a variable is a database-related variable
func (v *Validator) isDatabaseVariable(key string) bool {
	dbPatterns := []string{
		"DATABASE", "DB_", "CONNECTION", "DSN", "URI", "URL",
	}

	upperKey := strings.ToUpper(key)
	for _, pattern := range dbPatterns {
		if strings.Contains(upperKey, pattern) {
			return true
		}
	}

	return false
}

// validatePasswordComplexity validates password complexity
func (v *Validator) validatePasswordComplexity(variable util.EnvVariable, result *util.ValidationResult) {
	if !v.hasGoodComplexity(variable.Value) {
		result.Warnings = append(result.Warnings, util.ValidationWarning{
			Variable:   variable.Key,
			Message:    "password lacks complexity (should include uppercase, lowercase, numbers, and special characters)",
			LineNumber: variable.LineNumber,
		})
	}
}

// hasGoodComplexity checks if a password has good complexity
func (v *Validator) hasGoodComplexity(password string) bool {
	if len(password) < v.minPasswordLength {
		return false
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// determineRiskLevel determines the overall risk level based on validation results
func (v *Validator) determineRiskLevel(result *util.ValidationResult) {
	criticalCount := 0
	highCount := 0
	mediumCount := 0

	for _, err := range result.Errors {
		switch err.Severity {
		case "critical":
			criticalCount++
		case "high":
			highCount++
		case "medium":
			mediumCount++
		}
	}

	switch {
	case criticalCount > 0:
		result.RiskLevel = util.RiskLevelCritical
	case highCount > 0:
		result.RiskLevel = util.RiskLevelHigh
	case mediumCount > 0:
		result.RiskLevel = util.RiskLevelMedium
	default:
		result.RiskLevel = util.RiskLevelLow
	}
}

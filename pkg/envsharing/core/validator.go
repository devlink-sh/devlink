package core

import (
	"regexp"
	"strings"
)

type Validator struct {
	weakPasswordPatterns []*regexp.Regexp
	urlPatterns          []*regexp.Regexp
	emailPattern         *regexp.Regexp
	ipPattern            *regexp.Regexp
	minPasswordLength    int
	maxValueLength       int
}

func NewValidator() *Validator {
	return &Validator{
		weakPasswordPatterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)^(password|123456|admin|test)$`),
			regexp.MustCompile(`^.{1,5}$`),
		},
		urlPatterns: []*regexp.Regexp{
			regexp.MustCompile(`https?://[^\s]+`),
		},
		emailPattern:      regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
		ipPattern:         regexp.MustCompile(`\b(?:[0-9]{1,3}\.){3}[0-9]{1,3}\b`),
		minPasswordLength: 8,
		maxValueLength:    1000,
	}
}

func (v *Validator) Validate(envFile *EnvFile) *ValidationResult {
	result := &ValidationResult{
		IsValid:       true,
		Errors:        []ValidationError{},
		SensitiveVars: []string{},
		RiskLevel:     RiskLow,
	}

	for _, variable := range envFile.Variables {
		v.validateVariable(variable, result)
		if variable.IsSensitive {
			result.SensitiveVars = append(result.SensitiveVars, variable.Key)
		}
	}

	v.determineRiskLevel(result)
	return result
}

func (v *Validator) validateVariable(variable EnvVariable, result *ValidationResult) {
	if len(variable.Value) > v.maxValueLength {
		result.Errors = append(result.Errors, ValidationError{
			Variable:   variable.Key,
			Message:    "value too long",
			Severity:   "warning",
			LineNumber: variable.LineNumber,
		})
	}

	if variable.IsSensitive {
		v.validateSensitiveVariable(variable, result)
	}

	v.checkCommonSecurityIssues(variable, result)
}

func (v *Validator) validateSensitiveVariable(variable EnvVariable, result *ValidationResult) {
	if v.isWeakPassword(variable.Key, variable.Value) {
		result.Errors = append(result.Errors, ValidationError{
			Variable:   variable.Key,
			Message:    "weak password detected",
			Severity:   "error",
			LineNumber: variable.LineNumber,
		})
	}

	if v.isHardcodedCredential(variable.Value) {
		result.Errors = append(result.Errors, ValidationError{
			Variable:   variable.Key,
			Message:    "hardcoded credential detected",
			Severity:   "error",
			LineNumber: variable.LineNumber,
		})
	}

	if v.containsExposedSecret(variable.Value) {
		result.Errors = append(result.Errors, ValidationError{
			Variable:   variable.Key,
			Message:    "exposed secret in URL",
			Severity:   "error",
			LineNumber: variable.LineNumber,
		})
	}
}

func (v *Validator) checkCommonSecurityIssues(variable EnvVariable, result *ValidationResult) {
	if v.isDevelopmentValue(variable.Value) && v.isProductionVariable(variable.Key) {
		result.Errors = append(result.Errors, ValidationError{
			Variable:   variable.Key,
			Message:    "development value in production variable",
			Severity:   "warning",
			LineNumber: variable.LineNumber,
		})
	}
}

func (v *Validator) isWeakPassword(key, value string) bool {
	if !v.isPasswordVariable(key) {
		return false
	}

	for _, pattern := range v.weakPasswordPatterns {
		if pattern.MatchString(value) {
			return true
		}
	}

	return len(value) < v.minPasswordLength
}

func (v *Validator) isPasswordVariable(key string) bool {
	passwordKeys := []string{"password", "pass", "pwd", "secret", "key"}
	keyLower := strings.ToLower(key)
	for _, pwdKey := range passwordKeys {
		if strings.Contains(keyLower, pwdKey) {
			return true
		}
	}
	return false
}

func (v *Validator) isHardcodedCredential(value string) bool {
	hardcodedValues := []string{
		"admin", "password", "123456", "test", "demo",
		"changeme", "default", "root", "guest",
	}
	valueLower := strings.ToLower(value)
	for _, hardcoded := range hardcodedValues {
		if valueLower == hardcoded {
			return true
		}
	}
	return false
}

func (v *Validator) containsExposedSecret(value string) bool {
	for _, pattern := range v.urlPatterns {
		if pattern.MatchString(value) {
			return strings.Contains(value, "://") &&
				(strings.Contains(value, "password=") ||
					strings.Contains(value, "token=") ||
					strings.Contains(value, "key="))
		}
	}
	return false
}

func (v *Validator) isDevelopmentValue(value string) bool {
	devValues := []string{"localhost", "127.0.0.1", "dev", "development", "test", "demo"}
	valueLower := strings.ToLower(value)
	for _, devValue := range devValues {
		if strings.Contains(valueLower, devValue) {
			return true
		}
	}
	return false
}

func (v *Validator) isProductionVariable(key string) bool {
	prodKeys := []string{"prod", "production", "live", "api", "database"}
	keyLower := strings.ToLower(key)
	for _, prodKey := range prodKeys {
		if strings.Contains(keyLower, prodKey) {
			return true
		}
	}
	return false
}

func (v *Validator) determineRiskLevel(result *ValidationResult) {
	errorCount := 0
	warningCount := 0

	for _, err := range result.Errors {
		if err.Severity == "error" {
			errorCount++
		} else if err.Severity == "warning" {
			warningCount++
		}
	}

	if errorCount > 0 {
		result.RiskLevel = RiskHigh
		result.IsValid = false
	} else if warningCount > 2 {
		result.RiskLevel = RiskMedium
	} else {
		result.RiskLevel = RiskLow
	}
}

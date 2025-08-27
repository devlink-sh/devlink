package util

import "time"

// EnvVariable represents a single environment variable
type EnvVariable struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	IsSensitive bool   `json:"is_sensitive"`
	LineNumber  int    `json:"line_number"`
	Comment     string `json:"comment,omitempty"`
}

// EnvFile represents a parsed .env file
type EnvFile struct {
	Variables    []EnvVariable `json:"variables"`
	RawContent   string        `json:"raw_content"`
	FilePath     string        `json:"file_path"`
	ParseErrors  []ParseError  `json:"parse_errors,omitempty"`
	Warnings     []string      `json:"warnings,omitempty"`
	TotalLines   int           `json:"total_lines"`
	ValidLines   int           `json:"valid_lines"`
	CommentLines int           `json:"comment_lines"`
	EmptyLines   int           `json:"empty_lines"`
}

// ParseError represents an error during .env file parsing
type ParseError struct {
	LineNumber int    `json:"line_number"`
	Message    string `json:"message"`
	Line       string `json:"line"`
}

// ValidationResult represents the result of security validation
type ValidationResult struct {
	IsValid       bool                `json:"is_valid"`
	Errors        []ValidationError   `json:"errors,omitempty"`
	Warnings      []ValidationWarning `json:"warnings,omitempty"`
	SensitiveVars []string            `json:"sensitive_vars,omitempty"`
	RiskLevel     RiskLevel           `json:"risk_level"`
}

// ValidationError represents a validation error
type ValidationError struct {
	Variable   string `json:"variable"`
	Message    string `json:"message"`
	Severity   string `json:"severity"`
	LineNumber int    `json:"line_number"`
}

// ValidationWarning represents a validation warning
type ValidationWarning struct {
	Variable   string `json:"variable"`
	Message    string `json:"message"`
	LineNumber int    `json:"line_number"`
}

// RiskLevel represents the security risk level
type RiskLevel string

const (
	RiskLevelLow      RiskLevel = "low"
	RiskLevelMedium   RiskLevel = "medium"
	RiskLevelHigh     RiskLevel = "high"
	RiskLevelCritical RiskLevel = "critical"
)

// FormatOptions represents options for formatting output
type FormatOptions struct {
	MaskSensitive   bool   `json:"mask_sensitive"`
	ShowComments    bool   `json:"show_comments"`
	ShowLineNumbers bool   `json:"show_line_numbers"`
	OutputFormat    string `json:"output_format"` // "text", "json", "yaml"
	IndentSize      int    `json:"indent_size"`
}

// ShareInfo represents information about a shared environment file
type ShareInfo struct {
	Code           string    `json:"code"`
	ExpiresAt      time.Time `json:"expires_at"`
	IsReadOnly     bool      `json:"is_read_only"`
	CreatedAt      time.Time `json:"created_at"`
	FileSize       int64     `json:"file_size"`
	VariableCount  int       `json:"variable_count"`
	SensitiveCount int       `json:"sensitive_count"`
}

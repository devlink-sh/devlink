package core

import "time"

type EnvVariable struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	IsSensitive bool   `json:"is_sensitive"`
	LineNumber  int    `json:"line_number"`
	Comment     string `json:"comment,omitempty"`
}

type EnvFile struct {
	Variables    []EnvVariable `json:"variables"`
	RawContent   string        `json:"raw_content"`
	FilePath     string        `json:"file_path"`
	ParseErrors  []ParseError  `json:"parse_errors,omitempty"`
	TotalLines   int           `json:"total_lines"`
	ValidLines   int           `json:"valid_lines"`
	CommentLines int           `json:"comment_lines"`
	EmptyLines   int           `json:"empty_lines"`
}

type ParseError struct {
	LineNumber int    `json:"line_number"`
	Message    string `json:"message"`
	Line       string `json:"line"`
}

type ValidationResult struct {
	IsValid       bool              `json:"is_valid"`
	Errors        []ValidationError `json:"errors,omitempty"`
	SensitiveVars []string          `json:"sensitive_vars,omitempty"`
	RiskLevel     RiskLevel         `json:"risk_level"`
}

type ValidationError struct {
	Variable   string `json:"variable"`
	Message    string `json:"message"`
	Severity   string `json:"severity"`
	LineNumber int    `json:"line_number"`
}

type RiskLevel string

const (
	RiskLow    RiskLevel = "low"
	RiskMedium RiskLevel = "medium"
	RiskHigh   RiskLevel = "high"
)

type FormatOptions struct {
	MaskSensitive   bool   `json:"mask_sensitive"`
	ShowComments    bool   `json:"show_comments"`
	ShowLineNumbers bool   `json:"show_line_numbers"`
	OutputFormat    string `json:"output_format"`
	IndentSize      int    `json:"indent_size"`
}

type Share struct {
	ID            string                 `json:"id"`
	ShareCode     string                 `json:"share_code"`
	EncryptedData *EncryptedData         `json:"encrypted_data"`
	CreatedAt     time.Time              `json:"created_at"`
	ExpiresAt     time.Time              `json:"expires_at"`
	IsReadOnly    bool                   `json:"is_read_only"`
	AccessCount   int                    `json:"access_count"`
	MaxAccesses   int                    `json:"max_accesses"`
	Metadata      map[string]interface{} `json:"metadata"`
}

type EncryptedData struct {
	Data      string `json:"data"`
	Nonce     string `json:"nonce"`
	Salt      string `json:"salt"`
	Version   string `json:"version"`
	Algorithm string `json:"algorithm"`
}

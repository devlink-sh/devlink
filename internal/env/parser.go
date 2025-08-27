package env

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/devlink/internal/util"
)

// Parser handles .env file parsing
type Parser struct {
	// Sensitive patterns for detecting sensitive variables
	sensitivePatterns []*regexp.Regexp
	// Comment patterns
	commentPattern *regexp.Regexp
	// Variable assignment pattern
	varPattern *regexp.Regexp
}

// NewParser creates a new .env parser
func NewParser() *Parser {
	return &Parser{
		sensitivePatterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)(password|passwd|pwd)`),
			regexp.MustCompile(`(?i)(secret|key|token|auth)`),
			regexp.MustCompile(`(?i)(api_key|apikey|access_key)`),
			regexp.MustCompile(`(?i)(private_key|privatekey|privkey)`),
			regexp.MustCompile(`(?i)(database_url|db_url|connection_string)`),
			regexp.MustCompile(`(?i)(redis_url|redis_password)`),
			regexp.MustCompile(`(?i)(jwt_secret|jwt_key)`),
			regexp.MustCompile(`(?i)(encryption_key|encrypt_key)`),
			regexp.MustCompile(`(?i)(aws_secret|aws_key|aws_access)`),
			regexp.MustCompile(`(?i)(google_api|github_token|gitlab_token)`),
		},
		commentPattern: regexp.MustCompile(`^\s*#.*$`),
		varPattern:     regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_]*)\s*=\s*(.*)$`),
	}
}

// ParseFile parses a .env file and returns an EnvFile struct
func (p *Parser) ParseFile(filePath string) (*util.EnvFile, error) {
	// Validate file path
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, fmt.Errorf("invalid file path: %w", err)
	}

	// Check if file exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s", filePath)
	}

	// Read file content
	content, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Parse the content
	return p.ParseContent(string(content), absPath)
}

// ParseContent parses .env content from a string
func (p *Parser) ParseContent(content, filePath string) (*util.EnvFile, error) {
	envFile := &util.EnvFile{
		RawContent: content,
		FilePath:   filePath,
		Variables:  []util.EnvVariable{},
		ParseErrors: []util.ParseError{},
		Warnings:   []string{},
	}

	scanner := bufio.NewScanner(strings.NewReader(content))
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimRight(scanner.Text(), "\r\n")
		envFile.TotalLines++

		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			envFile.EmptyLines++
			continue
		}

		// Handle comments
		if p.commentPattern.MatchString(line) {
			envFile.CommentLines++
			continue
		}

		// Parse variable assignment
		if variable, err := p.parseVariableLine(line, lineNumber); err != nil {
			envFile.ParseErrors = append(envFile.ParseErrors, util.ParseError{
				LineNumber: lineNumber,
				Message:    err.Error(),
				Line:       line,
			})
		} else if variable != nil {
			envFile.Variables = append(envFile.Variables, *variable)
			envFile.ValidLines++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading content: %w", err)
	}

	// Add warnings for common issues
	p.addWarnings(envFile)

	return envFile, nil
}

// parseVariableLine parses a single line for variable assignment
func (p *Parser) parseVariableLine(line string, lineNumber int) (*util.EnvVariable, error) {
	// Remove leading/trailing whitespace
	line = strings.TrimSpace(line)

	// Check if it's a comment
	if strings.HasPrefix(line, "#") {
		return nil, nil
	}

	// Find the first equals sign
	equalsIndex := strings.Index(line, "=")
	if equalsIndex == -1 {
		return nil, fmt.Errorf("no equals sign found")
	}

	// Extract key and value
	key := strings.TrimSpace(line[:equalsIndex])
	value := strings.TrimSpace(line[equalsIndex+1:])

	// Validate key
	if err := p.validateKey(key); err != nil {
		return nil, fmt.Errorf("invalid key: %w", err)
	}

	// Parse value (handle quotes)
	parsedValue, comment, err := p.parseValue(value)
	if err != nil {
		return nil, fmt.Errorf("invalid value: %w", err)
	}

	// Check if variable is sensitive
	isSensitive := p.isSensitiveVariable(key)

	return &util.EnvVariable{
		Key:         key,
		Value:       parsedValue,
		IsSensitive: isSensitive,
		LineNumber:  lineNumber,
		Comment:     comment,
	}, nil
}

// validateKey validates a variable key
func (p *Parser) validateKey(key string) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	// Check first character
	if !unicode.IsLetter(rune(key[0])) && key[0] != '_' {
		return fmt.Errorf("key must start with a letter or underscore")
	}

	// Check remaining characters
	for _, char := range key {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) && char != '_' {
			return fmt.Errorf("key contains invalid character: %c", char)
		}
	}

	return nil
}

// parseValue parses a value, handling quotes and comments
func (p *Parser) parseValue(value string) (string, string, error) {
	if value == "" {
		return "", "", nil
	}
	
	// Add bounds checking for large values
	const maxValueLength = 10000 // 10KB limit
	if len(value) > maxValueLength {
		return "", "", fmt.Errorf("value too long (max %d characters)", maxValueLength)
	}

	var result strings.Builder
	var comment strings.Builder
	inQuotes := false
	quoteChar := rune(0)
	escapeNext := false

	for i, char := range value {
		if escapeNext {
			switch char {
			case 'n':
				result.WriteRune('\n')
			case 't':
				result.WriteRune('\t')
			case 'r':
				result.WriteRune('\r')
			case '"', '\'', '\\':
				result.WriteRune(char)
			default:
				result.WriteRune('\\')
				result.WriteRune(char)
			}
			escapeNext = false
			continue
		}

		if char == '\\' {
			escapeNext = true
			continue
		}

		// Handle quotes
		if (char == '"' || char == '\'') && !inQuotes {
			inQuotes = true
			quoteChar = char
			continue
		}

		if char == quoteChar && inQuotes {
			inQuotes = false
			quoteChar = 0
			continue
		}

		// Handle comments (only outside quotes)
		if char == '#' && !inQuotes {
			comment.WriteString(value[i:])
			break
		}

		// Add character to result
		if inQuotes || char != '#' {
			result.WriteRune(char)
		}
	}

	if inQuotes {
		return "", "", fmt.Errorf("unclosed quotes")
	}

	if escapeNext {
		return "", "", fmt.Errorf("incomplete escape sequence")
	}

	return strings.TrimSpace(result.String()), strings.TrimSpace(comment.String()), nil
}

// isSensitiveVariable checks if a variable name indicates sensitive data
func (p *Parser) isSensitiveVariable(key string) bool {
	upperKey := strings.ToUpper(key)
	
	for _, pattern := range p.sensitivePatterns {
		if pattern.MatchString(upperKey) {
			return true
		}
	}

	return false
}

// addWarnings adds common warnings to the parsed file
func (p *Parser) addWarnings(envFile *util.EnvFile) {
	// Check for duplicate keys
	seenKeys := make(map[string]int)
	for _, variable := range envFile.Variables {
		if lineNum, exists := seenKeys[variable.Key]; exists {
			envFile.Warnings = append(envFile.Warnings, 
				fmt.Sprintf("duplicate key '%s' found at lines %d and %d", 
					variable.Key, lineNum, variable.LineNumber))
		} else {
			seenKeys[variable.Key] = variable.LineNumber
		}
	}

	// Check for empty values
	for _, variable := range envFile.Variables {
		if variable.Value == "" {
			envFile.Warnings = append(envFile.Warnings, 
				fmt.Sprintf("empty value for key '%s' at line %d", 
					variable.Key, variable.LineNumber))
		}
	}

	// Check for sensitive variables
	sensitiveCount := 0
	for _, variable := range envFile.Variables {
		if variable.IsSensitive {
			sensitiveCount++
		}
	}

	if sensitiveCount > 0 {
		envFile.Warnings = append(envFile.Warnings, 
			fmt.Sprintf("found %d potentially sensitive variables", sensitiveCount))
	}
}

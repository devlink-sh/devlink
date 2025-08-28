package core

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Parser struct {
	sensitivePatterns []*regexp.Regexp
	commentPattern    *regexp.Regexp
	varPattern        *regexp.Regexp
}

func NewParser() *Parser {
	return &Parser{
		sensitivePatterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)(password|pass|pwd)`),
			regexp.MustCompile(`(?i)(secret|key|token)`),
			regexp.MustCompile(`(?i)(auth|credential)`),
			regexp.MustCompile(`(?i)(api[_-]?key|apikey)`),
			regexp.MustCompile(`(?i)(private[_-]?key|privkey)`),
		},
		commentPattern: regexp.MustCompile(`^\s*#`),
		varPattern:     regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_]*)\s*=\s*(.*)$`),
	}
}

func (p *Parser) ParseFile(filePath string) (*EnvFile, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	envFile := p.ParseContent(string(content), filePath)
	return envFile, nil
}

func (p *Parser) ParseContent(content, filePath string) *EnvFile {
	lines := strings.Split(content, "\n")
	envFile := &EnvFile{
		FilePath:   filePath,
		RawContent: content,
		Variables:  []EnvVariable{},
	}

	for i, line := range lines {
		lineNumber := i + 1
		line = strings.TrimSpace(line)

		if line == "" {
			envFile.EmptyLines++
			continue
		}

		if p.commentPattern.MatchString(line) {
			envFile.CommentLines++
			continue
		}

		if variable, err := p.parseVariableLine(line, lineNumber); err == nil {
			variable.IsSensitive = p.isSensitiveVariable(variable.Key)
			envFile.Variables = append(envFile.Variables, *variable)
			envFile.ValidLines++
		} else {
			envFile.ParseErrors = append(envFile.ParseErrors, ParseError{
				LineNumber: lineNumber,
				Message:    err.Error(),
				Line:       line,
			})
		}
	}

	envFile.TotalLines = len(lines)
	return envFile
}

func (p *Parser) parseVariableLine(line string, lineNumber int) (*EnvVariable, error) {
	matches := p.varPattern.FindStringSubmatch(line)
	if len(matches) != 3 {
		return nil, fmt.Errorf("invalid variable format")
	}

	key := strings.TrimSpace(matches[1])
	value := strings.TrimSpace(matches[2])

	if err := p.validateKey(key); err != nil {
		return nil, err
	}

	parsedValue, comment, err := p.parseValue(value)
	if err != nil {
		return nil, err
	}

	return &EnvVariable{
		Key:        key,
		Value:      parsedValue,
		LineNumber: lineNumber,
		Comment:    comment,
	}, nil
}

func (p *Parser) validateKey(key string) error {
	if len(key) == 0 {
		return fmt.Errorf("key cannot be empty")
	}
	if len(key) > 100 {
		return fmt.Errorf("key too long")
	}
	return nil
}

func (p *Parser) parseValue(value string) (string, string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", "", nil
	}

	if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
		unquoted, err := strconv.Unquote(value)
		if err != nil {
			return value, "", nil
		}
		return unquoted, "", nil
	}

	if strings.HasPrefix(value, `'`) && strings.HasSuffix(value, `'`) {
		return value[1 : len(value)-1], "", nil
	}

	return value, "", nil
}

func (p *Parser) isSensitiveVariable(key string) bool {
	for _, pattern := range p.sensitivePatterns {
		if pattern.MatchString(key) {
			return true
		}
	}
	return false
}

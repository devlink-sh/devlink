package env

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/devlink/internal/util"
)

type BulkShareRequest struct {
	Files    []string          `json:"files"`
	Expiry   time.Duration     `json:"expiry"`
	ReadOnly bool              `json:"read_only"`
	Prefix   string            `json:"prefix"`
	GroupBy  string            `json:"group_by"`
	Compress bool              `json:"compress"`
	Metadata map[string]string `json:"metadata"`
}

type BulkShareResult struct {
	File      string    `json:"file"`
	ShareCode string    `json:"share_code"`
	Success   bool      `json:"success"`
	Error     string    `json:"error,omitempty"`
	ExpiresAt time.Time `json:"expires_at"`
	FileSize  int64     `json:"file_size"`
	Variables int       `json:"variables"`
	Sensitive int       `json:"sensitive"`
}

type BulkManager struct {
	parser   *Parser
	tokenGen *util.TokenGenerator
	config   *util.Config
}

func NewBulkManager() *BulkManager {
	config, err := util.LoadConfig()
	if err != nil {
		config = util.DefaultConfig()
	}

	return &BulkManager{
		parser:   NewParser(),
		tokenGen: util.NewTokenGenerator(),
		config:   config,
	}
}

func (bm *BulkManager) ShareFiles(request BulkShareRequest) []BulkShareResult {
	results := make([]BulkShareResult, 0, len(request.Files))

	for _, filePath := range request.Files {
		result := bm.shareSingleFile(filePath, request)
		results = append(results, result)
	}

	return results
}

func (bm *BulkManager) shareSingleFile(filePath string, request BulkShareRequest) BulkShareResult {
	result := BulkShareResult{
		File: filePath,
	}

	envFile, err := bm.parser.ParseFile(filePath)
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to parse file: %v", err)
		return result
	}

	shareCode, err := bm.generateShareCode(filePath, request)
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to generate share code: %v", err)
		return result
	}

	result.ShareCode = shareCode
	result.Success = true
	result.ExpiresAt = time.Now().Add(request.Expiry)
	result.FileSize = int64(len(envFile.RawContent))
	result.Variables = len(envFile.Variables)
	result.Sensitive = bm.countSensitiveVariables(envFile)

	return result
}

func (bm *BulkManager) generateShareCode(filePath string, request BulkShareRequest) (string, error) {
	baseName := filepath.Base(filePath)
	ext := filepath.Ext(baseName)
	name := strings.TrimSuffix(baseName, ext)

	if request.Prefix != "" {
		name = request.Prefix + "-" + name
	}

	shareCode, err := bm.tokenGen.GenerateShareCode()
	if err != nil {
		return "", err
	}

	return shareCode, nil
}

func (bm *BulkManager) countSensitiveVariables(envFile *util.EnvFile) int {
	count := 0
	for _, variable := range envFile.Variables {
		if variable.IsSensitive {
			count++
		}
	}
	return count
}

func (bm *BulkManager) ValidateBulkRequest(request BulkShareRequest) error {
	if len(request.Files) == 0 {
		return fmt.Errorf("no files specified")
	}

	if len(request.Files) > 50 {
		return fmt.Errorf("too many files (max 50)")
	}

	if request.Expiry <= 0 {
		request.Expiry = bm.config.GetDefaultExpiry()
	}

	if request.Expiry > 168*time.Hour {
		return fmt.Errorf("expiry too long (max 7 days)")
	}

	return nil
}

func (bm *BulkManager) GetBulkStatistics(results []BulkShareResult) map[string]interface{} {
	totalFiles := len(results)
	successfulShares := 0
	failedShares := 0
	totalVariables := 0
	totalSensitive := 0
	totalSize := int64(0)

	for _, result := range results {
		if result.Success {
			successfulShares++
			totalVariables += result.Variables
			totalSensitive += result.Sensitive
			totalSize += result.FileSize
		} else {
			failedShares++
		}
	}

	return map[string]interface{}{
		"total_files":       totalFiles,
		"successful_shares": successfulShares,
		"failed_shares":     failedShares,
		"total_variables":   totalVariables,
		"total_sensitive":   totalSensitive,
		"total_size_bytes":  totalSize,
		"success_rate":      float64(successfulShares) / float64(totalFiles),
	}
}

func (bm *BulkManager) GroupResultsByCategory(results []BulkShareResult) map[string][]BulkShareResult {
	groups := make(map[string][]BulkShareResult)

	for _, result := range results {
		category := bm.categorizeFile(result.File)
		groups[category] = append(groups[category], result)
	}

	return groups
}

func (bm *BulkManager) categorizeFile(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	baseName := strings.ToLower(filepath.Base(filePath))

	switch {
	case strings.Contains(baseName, "prod") || strings.Contains(baseName, "production"):
		return "production"
	case strings.Contains(baseName, "dev") || strings.Contains(baseName, "development"):
		return "development"
	case strings.Contains(baseName, "test") || strings.Contains(baseName, "testing"):
		return "testing"
	case strings.Contains(baseName, "staging"):
		return "staging"
	case ext == ".env":
		return "environment"
	case ext == ".config":
		return "configuration"
	default:
		return "other"
	}
}

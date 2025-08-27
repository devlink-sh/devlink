package util

import (
	"regexp"
	"strings"
)

type SearchFilter struct {
	Query       string   `json:"query"`
	Categories  []string `json:"categories"`
	Sensitive   *bool    `json:"sensitive"`
	ExactMatch  bool     `json:"exact_match"`
	CaseSensitive bool   `json:"case_sensitive"`
	Regex       bool     `json:"regex"`
}

type SearchResult struct {
	Variable    EnvVariable `json:"variable"`
	File        string      `json:"file"`
	LineNumber  int         `json:"line_number"`
	MatchType   string      `json:"match_type"`
	Score       float64     `json:"score"`
}

type SearchManager struct {
	envFiles []*EnvFile
}

func NewSearchManager() *SearchManager {
	return &SearchManager{
		envFiles: make([]*EnvFile, 0),
	}
}

func (sm *SearchManager) AddEnvFile(envFile *EnvFile) {
	sm.envFiles = append(sm.envFiles, envFile)
}

func (sm *SearchManager) Search(filter SearchFilter) []SearchResult {
	var results []SearchResult

	for _, envFile := range sm.envFiles {
		fileResults := sm.searchInFile(envFile, filter)
		results = append(results, fileResults...)
	}

	return sm.rankResults(results, filter)
}

func (sm *SearchManager) searchInFile(envFile *EnvFile, filter SearchFilter) []SearchResult {
	var results []SearchResult

	for _, variable := range envFile.Variables {
		if sm.matchesFilter(variable, filter) {
			result := SearchResult{
				Variable:   variable,
				File:       envFile.FilePath,
				LineNumber: variable.LineNumber,
				MatchType:  sm.determineMatchType(variable, filter),
				Score:      sm.calculateScore(variable, filter),
			}
			results = append(results, result)
		}
	}

	return results
}

func (sm *SearchManager) matchesFilter(variable EnvVariable, filter SearchFilter) bool {
	if filter.Sensitive != nil && variable.IsSensitive != *filter.Sensitive {
		return false
	}

	if len(filter.Categories) > 0 && !sm.matchesCategory(variable, filter.Categories) {
		return false
	}

	if filter.Query == "" {
		return true
	}

	return sm.matchesQuery(variable, filter)
}

func (sm *SearchManager) matchesCategory(variable EnvVariable, categories []string) bool {
	variableCategory := sm.categorizeVariable(variable.Key)
	for _, category := range categories {
		if strings.EqualFold(variableCategory, category) {
			return true
		}
	}
	return false
}

func (sm *SearchManager) categorizeVariable(key string) string {
	key = strings.ToUpper(key)
	
	switch {
	case strings.Contains(key, "DATABASE") || strings.Contains(key, "DB_"):
		return "database"
	case strings.Contains(key, "API_") || strings.Contains(key, "_KEY"):
		return "api"
	case strings.Contains(key, "AWS_") || strings.Contains(key, "GCP_") || strings.Contains(key, "AZURE_"):
		return "cloud"
	case strings.Contains(key, "DOCKER_") || strings.Contains(key, "KUBERNETES_"):
		return "container"
	case strings.Contains(key, "REDIS_") || strings.Contains(key, "CACHE_"):
		return "cache"
	case strings.Contains(key, "EMAIL_") || strings.Contains(key, "SMTP_"):
		return "email"
	case strings.Contains(key, "LOG_") || strings.Contains(key, "DEBUG"):
		return "logging"
	default:
		return "general"
	}
}

func (sm *SearchManager) matchesQuery(variable EnvVariable, filter SearchFilter) bool {
	query := filter.Query
	key := variable.Key
	value := variable.Value

	if !filter.CaseSensitive {
		query = strings.ToLower(query)
		key = strings.ToLower(key)
		value = strings.ToLower(value)
	}

	if filter.Regex {
		return sm.matchesRegex(key, value, query)
	}

	if filter.ExactMatch {
		return key == query || value == query
	}

	return strings.Contains(key, query) || strings.Contains(value, query)
}

func (sm *SearchManager) matchesRegex(key, value, pattern string) bool {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}

	return regex.MatchString(key) || regex.MatchString(value)
}

func (sm *SearchManager) determineMatchType(variable EnvVariable, filter SearchFilter) string {
	if filter.Query == "" {
		return "all"
	}

	query := filter.Query
	key := variable.Key
	value := variable.Value

	if !filter.CaseSensitive {
		query = strings.ToLower(query)
		key = strings.ToLower(key)
		value = strings.ToLower(value)
	}

	if filter.ExactMatch {
		if key == query {
			return "exact_key"
		}
		if value == query {
			return "exact_value"
		}
	}

	if strings.Contains(key, query) {
		return "key_match"
	}

	if strings.Contains(value, query) {
		return "value_match"
	}

	return "partial"
}

func (sm *SearchManager) calculateScore(variable EnvVariable, filter SearchFilter) float64 {
	score := 0.0

	if filter.Query == "" {
		return 1.0
	}

	query := filter.Query
	key := variable.Key
	value := variable.Value

	if !filter.CaseSensitive {
		query = strings.ToLower(query)
		key = strings.ToLower(key)
		value = strings.ToLower(value)
	}

	if filter.ExactMatch {
		if key == query {
			score += 10.0
		}
		if value == query {
			score += 8.0
		}
	} else {
		if strings.Contains(key, query) {
			score += 5.0
			if strings.HasPrefix(key, query) {
				score += 2.0
			}
		}
		if strings.Contains(value, query) {
			score += 3.0
		}
	}

	if variable.IsSensitive {
		score += 1.0
	}

	return score
}

func (sm *SearchManager) rankResults(results []SearchResult, filter SearchFilter) []SearchResult {
	for i := range results {
		for j := i + 1; j < len(results); j++ {
			if results[i].Score < results[j].Score {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	return results
}

func (sm *SearchManager) GetVariableSuggestions(partial string) []string {
	suggestions := make(map[string]bool)
	
	for _, envFile := range sm.envFiles {
		for _, variable := range envFile.Variables {
			if strings.HasPrefix(strings.ToLower(variable.Key), strings.ToLower(partial)) {
				suggestions[variable.Key] = true
			}
		}
	}

	result := make([]string, 0, len(suggestions))
	for suggestion := range suggestions {
		result = append(result, suggestion)
	}

	return result
}

func (sm *SearchManager) GetStatistics() map[string]interface{} {
	totalVariables := 0
	sensitiveVariables := 0
	categories := make(map[string]int)

	for _, envFile := range sm.envFiles {
		for _, variable := range envFile.Variables {
			totalVariables++
			if variable.IsSensitive {
				sensitiveVariables++
			}
			
			category := sm.categorizeVariable(variable.Key)
			categories[category]++
		}
	}

	return map[string]interface{}{
		"total_variables":     totalVariables,
		"sensitive_variables": sensitiveVariables,
		"files_count":         len(sm.envFiles),
		"categories":          categories,
	}
}

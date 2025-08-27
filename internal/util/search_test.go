package util

import (
	"testing"
)

func TestNewSearchManager(t *testing.T) {
	sm := NewSearchManager()
	if sm == nil {
		t.Fatal("NewSearchManager() returned nil")
	}
}

func TestSearchManager_AddEnvFile(t *testing.T) {
	sm := NewSearchManager()
	
	envFile := &EnvFile{
		Variables: []EnvVariable{
			{Key: "DATABASE_URL", Value: "postgresql://localhost:5432/mydb", IsSensitive: true},
			{Key: "API_KEY", Value: "sk-1234567890abcdef", IsSensitive: true},
			{Key: "NODE_ENV", Value: "development", IsSensitive: false},
		},
		FilePath: "test.env",
	}

	sm.AddEnvFile(envFile)
	
	stats := sm.GetStatistics()
	if stats["files_count"].(int) != 1 {
		t.Error("Should have 1 file")
	}
}

func TestSearchManager_Search(t *testing.T) {
	sm := NewSearchManager()
	
	envFile := &EnvFile{
		Variables: []EnvVariable{
			{Key: "DATABASE_URL", Value: "postgresql://localhost:5432/mydb", IsSensitive: true},
			{Key: "API_KEY", Value: "sk-1234567890abcdef", IsSensitive: true},
			{Key: "NODE_ENV", Value: "development", IsSensitive: false},
		},
		FilePath: "test.env",
	}

	sm.AddEnvFile(envFile)

	filter := SearchFilter{
		Query: "DATABASE",
	}

	results := sm.Search(filter)
	if len(results) == 0 {
		t.Error("Should find DATABASE_URL")
	}

	if results[0].Variable.Key != "DATABASE_URL" {
		t.Error("Should return DATABASE_URL")
	}
}

func TestSearchManager_SearchSensitive(t *testing.T) {
	sm := NewSearchManager()
	
	envFile := &EnvFile{
		Variables: []EnvVariable{
			{Key: "DATABASE_URL", Value: "postgresql://localhost:5432/mydb", IsSensitive: true},
			{Key: "API_KEY", Value: "sk-1234567890abcdef", IsSensitive: true},
			{Key: "NODE_ENV", Value: "development", IsSensitive: false},
		},
		FilePath: "test.env",
	}

	sm.AddEnvFile(envFile)

	sensitive := true
	filter := SearchFilter{
		Sensitive: &sensitive,
	}

	results := sm.Search(filter)
	if len(results) != 2 {
		t.Errorf("Should find 2 sensitive variables, got %d", len(results))
	}
}

func TestSearchManager_GetVariableSuggestions(t *testing.T) {
	sm := NewSearchManager()
	
	envFile := &EnvFile{
		Variables: []EnvVariable{
			{Key: "DATABASE_URL", Value: "postgresql://localhost:5432/mydb", IsSensitive: true},
			{Key: "API_KEY", Value: "sk-1234567890abcdef", IsSensitive: true},
			{Key: "NODE_ENV", Value: "development", IsSensitive: false},
		},
		FilePath: "test.env",
	}

	sm.AddEnvFile(envFile)

	suggestions := sm.GetVariableSuggestions("DAT")
	if len(suggestions) == 0 {
		t.Error("Should suggest DATABASE_URL")
	}

	if suggestions[0] != "DATABASE_URL" {
		t.Error("Should suggest DATABASE_URL")
	}
}

func TestSearchManager_GetStatistics(t *testing.T) {
	sm := NewSearchManager()
	
	envFile := &EnvFile{
		Variables: []EnvVariable{
			{Key: "DATABASE_URL", Value: "postgresql://localhost:5432/mydb", IsSensitive: true},
			{Key: "API_KEY", Value: "sk-1234567890abcdef", IsSensitive: true},
			{Key: "NODE_ENV", Value: "development", IsSensitive: false},
		},
		FilePath: "test.env",
	}

	sm.AddEnvFile(envFile)

	stats := sm.GetStatistics()
	if stats["total_variables"].(int) != 3 {
		t.Error("Should have 3 total variables")
	}

	if stats["sensitive_variables"].(int) != 2 {
		t.Error("Should have 2 sensitive variables")
	}

	if stats["files_count"].(int) != 1 {
		t.Error("Should have 1 file")
	}
}

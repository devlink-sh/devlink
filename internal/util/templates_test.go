package util

import (
	"testing"
)

func TestNewTemplateManager(t *testing.T) {
	tm := NewTemplateManager()
	if tm == nil {
		t.Fatal("NewTemplateManager() returned nil")
	}
}

func TestListTemplates(t *testing.T) {
	tm := NewTemplateManager()
	templates := tm.ListTemplates()

	if len(templates) == 0 {
		t.Error("Should have default templates")
	}

	foundNodeJS := false
	for _, template := range templates {
		if template.Name == "nodejs" {
			foundNodeJS = true
			break
		}
	}

	if !foundNodeJS {
		t.Error("Should have nodejs template")
	}
}

func TestGetTemplate(t *testing.T) {
	tm := NewTemplateManager()

	template, err := tm.GetTemplate("nodejs")
	if err != nil {
		t.Fatalf("GetTemplate failed: %v", err)
	}

	if template.Name != "nodejs" {
		t.Errorf("Expected template name 'nodejs', got '%s'", template.Name)
	}

	if len(template.Variables) == 0 {
		t.Error("Template should have variables")
	}
}

func TestGetTemplateNotFound(t *testing.T) {
	tm := NewTemplateManager()

	_, err := tm.GetTemplate("nonexistent")
	if err == nil {
		t.Error("Should return error for nonexistent template")
	}
}

func TestSearchTemplates(t *testing.T) {
	tm := NewTemplateManager()

	results := tm.SearchTemplates("node")
	if len(results) == 0 {
		t.Error("Should find nodejs template")
	}

	results = tm.SearchTemplates("backend")
	if len(results) == 0 {
		t.Error("Should find backend templates")
	}
}

func TestCreateEnvFileFromTemplate(t *testing.T) {
	tm := NewTemplateManager()

	envFile, err := tm.CreateEnvFileFromTemplate("nodejs")
	if err != nil {
		t.Fatalf("CreateEnvFileFromTemplate failed: %v", err)
	}

	if envFile == nil {
		t.Fatal("Should return env file")
	}

	if len(envFile.Variables) == 0 {
		t.Error("Env file should have variables")
	}

	if envFile.FilePath != "nodejs.env" {
		t.Errorf("Expected file path 'nodejs.env', got '%s'", envFile.FilePath)
	}
}

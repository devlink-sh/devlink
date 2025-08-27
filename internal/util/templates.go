package util

import (
	"fmt"
	"strings"
)

type Template struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Variables   []EnvVariable `json:"variables"`
	Tags        []string      `json:"tags"`
	Category    string        `json:"category"`
}

type TemplateManager struct {
	templates map[string]*Template
}

func NewTemplateManager() *TemplateManager {
	tm := &TemplateManager{
		templates: make(map[string]*Template),
	}
	tm.loadDefaultTemplates()
	return tm
}

func (tm *TemplateManager) loadDefaultTemplates() {
	templates := []*Template{
		{
			Name:        "nodejs",
			Description: "Node.js development environment",
			Category:    "backend",
			Tags:        []string{"node", "javascript", "backend"},
			Variables: []EnvVariable{
				{Key: "NODE_ENV", Value: "development", IsSensitive: false},
				{Key: "PORT", Value: "3000", IsSensitive: false},
				{Key: "DATABASE_URL", Value: "postgresql://localhost:5432/mydb", IsSensitive: true},
				{Key: "JWT_SECRET", Value: "your-jwt-secret", IsSensitive: true},
				{Key: "REDIS_URL", Value: "redis://localhost:6379", IsSensitive: true},
				{Key: "LOG_LEVEL", Value: "debug", IsSensitive: false},
			},
		},
		{
			Name:        "react",
			Description: "React frontend environment",
			Category:    "frontend",
			Tags:        []string{"react", "javascript", "frontend"},
			Variables: []EnvVariable{
				{Key: "REACT_APP_API_URL", Value: "http://localhost:3000/api", IsSensitive: false},
				{Key: "REACT_APP_ENV", Value: "development", IsSensitive: false},
				{Key: "REACT_APP_DEBUG", Value: "true", IsSensitive: false},
			},
		},
		{
			Name:        "python",
			Description: "Python development environment",
			Category:    "backend",
			Tags:        []string{"python", "backend"},
			Variables: []EnvVariable{
				{Key: "PYTHON_ENV", Value: "development", IsSensitive: false},
				{Key: "DATABASE_URL", Value: "postgresql://localhost:5432/mydb", IsSensitive: true},
				{Key: "SECRET_KEY", Value: "your-secret-key", IsSensitive: true},
				{Key: "DEBUG", Value: "True", IsSensitive: false},
				{Key: "ALLOWED_HOSTS", Value: "localhost,127.0.0.1", IsSensitive: false},
			},
		},
		{
			Name:        "docker",
			Description: "Docker development environment",
			Category:    "devops",
			Tags:        []string{"docker", "container", "devops"},
			Variables: []EnvVariable{
				{Key: "DOCKER_HOST", Value: "unix:///var/run/docker.sock", IsSensitive: false},
				{Key: "COMPOSE_PROJECT_NAME", Value: "myproject", IsSensitive: false},
				{Key: "DOCKER_REGISTRY", Value: "docker.io", IsSensitive: false},
			},
		},
		{
			Name:        "aws",
			Description: "AWS development environment",
			Category:    "cloud",
			Tags:        []string{"aws", "cloud", "infrastructure"},
			Variables: []EnvVariable{
				{Key: "AWS_ACCESS_KEY_ID", Value: "your-access-key", IsSensitive: true},
				{Key: "AWS_SECRET_ACCESS_KEY", Value: "your-secret-key", IsSensitive: true},
				{Key: "AWS_REGION", Value: "us-east-1", IsSensitive: false},
				{Key: "AWS_S3_BUCKET", Value: "my-bucket", IsSensitive: false},
			},
		},
	}

	for _, template := range templates {
		tm.templates[template.Name] = template
	}
}

func (tm *TemplateManager) ListTemplates() []*Template {
	templates := make([]*Template, 0, len(tm.templates))
	for _, template := range tm.templates {
		templates = append(templates, template)
	}
	return templates
}

func (tm *TemplateManager) GetTemplate(name string) (*Template, error) {
	template, exists := tm.templates[name]
	if !exists {
		return nil, fmt.Errorf("template not found: %s", name)
	}
	return template, nil
}

func (tm *TemplateManager) SearchTemplates(query string) []*Template {
	query = strings.ToLower(query)
	var results []*Template

	for _, template := range tm.templates {
		if strings.Contains(strings.ToLower(template.Name), query) ||
			strings.Contains(strings.ToLower(template.Description), query) ||
			strings.Contains(strings.ToLower(template.Category), query) ||
			tm.containsTag(template.Tags, query) {
			results = append(results, template)
		}
	}

	return results
}

func (tm *TemplateManager) containsTag(tags []string, query string) bool {
	for _, tag := range tags {
		if strings.Contains(strings.ToLower(tag), query) {
			return true
		}
	}
	return false
}

func (tm *TemplateManager) GetTemplatesByCategory(category string) []*Template {
	var results []*Template
	for _, template := range tm.templates {
		if strings.EqualFold(template.Category, category) {
			results = append(results, template)
		}
	}
	return results
}

func (tm *TemplateManager) CreateEnvFileFromTemplate(templateName string) (*EnvFile, error) {
	template, err := tm.GetTemplate(templateName)
	if err != nil {
		return nil, err
	}

	var content strings.Builder
	for _, variable := range template.Variables {
		content.WriteString(fmt.Sprintf("%s=%s\n", variable.Key, variable.Value))
	}

	return &EnvFile{
		Variables:    template.Variables,
		RawContent:   content.String(),
		FilePath:     fmt.Sprintf("%s.env", templateName),
		TotalLines:   len(template.Variables),
		ValidLines:   len(template.Variables),
		CommentLines: 0,
		EmptyLines:   0,
	}, nil
}

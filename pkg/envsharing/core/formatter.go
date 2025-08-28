package core

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
)

type Formatter struct{}

func NewFormatter() *Formatter {
	return &Formatter{}
}

func (f *Formatter) Format(envFile *EnvFile, options *FormatOptions) (string, error) {
	switch options.OutputFormat {
	case "json":
		return f.formatJSON(envFile, options)
	case "yaml":
		return f.formatYAML(envFile, options)
	default:
		return f.formatText(envFile, options)
	}
}

func (f *Formatter) formatText(envFile *EnvFile, options *FormatOptions) (string, error) {
	var output strings.Builder

	output.WriteString("📄 Environment File Analysis\n")
	output.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")
	output.WriteString(fmt.Sprintf("📁 File: %s\n", envFile.FilePath))
	output.WriteString("📊 Statistics:\n")
	output.WriteString(fmt.Sprintf("   • Total lines: %d\n", envFile.TotalLines))
	output.WriteString(fmt.Sprintf("   • Valid variables: %d\n", envFile.ValidLines))
	output.WriteString(fmt.Sprintf("   • Comment lines: %d\n", envFile.CommentLines))
	output.WriteString(fmt.Sprintf("   • Empty lines: %d\n", envFile.EmptyLines))

	sensitiveCount := 0
	for _, variable := range envFile.Variables {
		if variable.IsSensitive {
			sensitiveCount++
		}
	}
	output.WriteString(fmt.Sprintf("   • Sensitive variables: %d\n", sensitiveCount))

	if len(envFile.ParseErrors) > 0 {
		output.WriteString("\n⚠️  Warnings:\n")
		for _, err := range envFile.ParseErrors {
			output.WriteString(fmt.Sprintf("   • Line %d: %s\n", err.LineNumber, err.Message))
		}
	}

	output.WriteString("\n🔧 Variables:\n")
	output.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	for _, variable := range envFile.Variables {
		line := f.formatVariableLine(variable, options)
		output.WriteString(line + "\n")
	}

	return output.String(), nil
}

func (f *Formatter) formatVariableLine(variable EnvVariable, options *FormatOptions) string {
	value := variable.Value
	if options.MaskSensitive && variable.IsSensitive {
		value = f.maskValue(value)
	}

	sensitive := ""
	if variable.IsSensitive {
		sensitive = " 🔒"
	}

	lineNumber := ""
	if options.ShowLineNumbers {
		lineNumber = fmt.Sprintf("%3d: ", variable.LineNumber)
	}

	return fmt.Sprintf("%s%-15s = %-30s%s", lineNumber, variable.Key, value, sensitive)
}

func (f *Formatter) maskValue(value string) string {
	if len(value) <= 4 {
		return strings.Repeat("*", len(value))
	}
	return value[:1] + strings.Repeat("*", len(value)-2) + value[len(value)-1:]
}

func (f *Formatter) formatJSON(envFile *EnvFile, options *FormatOptions) (string, error) {
	safeData := f.createSafeCopy(envFile, options)
	data, err := json.MarshalIndent(safeData, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (f *Formatter) formatYAML(envFile *EnvFile, options *FormatOptions) (string, error) {
	safeData := f.createSafeCopy(envFile, options)
	data, err := yaml.Marshal(safeData)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (f *Formatter) createSafeCopy(envFile *EnvFile, options *FormatOptions) map[string]interface{} {
	variables := make(map[string]interface{})

	for _, variable := range envFile.Variables {
		value := variable.Value
		if options.MaskSensitive && variable.IsSensitive {
			value = f.maskValue(value)
		}
		variables[variable.Key] = value
	}

	return map[string]interface{}{
		"file_path":     envFile.FilePath,
		"total_lines":   envFile.TotalLines,
		"valid_lines":   envFile.ValidLines,
		"comment_lines": envFile.CommentLines,
		"empty_lines":   envFile.EmptyLines,
		"variables":     variables,
		"parse_errors":  envFile.ParseErrors,
	}
}

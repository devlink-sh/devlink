package env

import (
	"encoding/json"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/devlink/internal/util"
)

// Formatter handles .env file output formatting
type Formatter struct {
	// Default formatting options
	defaultOptions *util.FormatOptions
}

// NewFormatter creates a new .env formatter
func NewFormatter() *Formatter {
	return &Formatter{
		defaultOptions: &util.FormatOptions{
			MaskSensitive:  true,
			ShowComments:   true,
			ShowLineNumbers: false,
			OutputFormat:   "text",
			IndentSize:     2,
		},
	}
}

// Format formats an .env file according to the specified options
func (f *Formatter) Format(envFile *util.EnvFile, options *util.FormatOptions) (string, error) {
	if options == nil {
		options = f.defaultOptions
	}

	switch options.OutputFormat {
	case "text":
		return f.formatText(envFile, options)
	case "json":
		return f.formatJSON(envFile, options)
	case "yaml":
		return f.formatYAML(envFile, options)
	default:
		return f.formatText(envFile, options)
	}
}

// formatText formats the .env file as text
func (f *Formatter) formatText(envFile *util.EnvFile, options *util.FormatOptions) (string, error) {
	var result strings.Builder

	// Add header
	result.WriteString("📄 Environment File Analysis\n")
	result.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	// Add file info
	result.WriteString(fmt.Sprintf("📁 File: %s\n", envFile.FilePath))
	result.WriteString(fmt.Sprintf("📊 Statistics:\n"))
	result.WriteString(fmt.Sprintf("   • Total lines: %d\n", envFile.TotalLines))
	result.WriteString(fmt.Sprintf("   • Valid variables: %d\n", envFile.ValidLines))
	result.WriteString(fmt.Sprintf("   • Comment lines: %d\n", envFile.CommentLines))
	result.WriteString(fmt.Sprintf("   • Empty lines: %d\n", envFile.EmptyLines))
	result.WriteString(fmt.Sprintf("   • Sensitive variables: %d\n", f.countSensitiveVariables(envFile)))

	// Add warnings if any
	if len(envFile.Warnings) > 0 {
		result.WriteString("\n⚠️  Warnings:\n")
		for _, warning := range envFile.Warnings {
			result.WriteString(fmt.Sprintf("   • %s\n", warning))
		}
	}

	// Add parse errors if any
	if len(envFile.ParseErrors) > 0 {
		result.WriteString("\n❌ Parse Errors:\n")
		for _, err := range envFile.ParseErrors {
			result.WriteString(fmt.Sprintf("   • Line %d: %s\n", err.LineNumber, err.Message))
			if options.ShowLineNumbers {
				result.WriteString(fmt.Sprintf("     Content: %s\n", err.Line))
			}
		}
	}

	// Add variables
	if len(envFile.Variables) > 0 {
		result.WriteString("\n🔧 Variables:\n")
		result.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

		// Use tabwriter for aligned output
		var buf strings.Builder
		w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)

		for _, variable := range envFile.Variables {
			line := f.formatVariableLine(variable, options)
			fmt.Fprintln(w, line)
		}
		w.Flush()
		result.WriteString(buf.String())
	}

	return result.String(), nil
}

// formatVariableLine formats a single variable line
func (f *Formatter) formatVariableLine(variable util.EnvVariable, options *util.FormatOptions) string {
	var parts []string

	// Add line number if requested
	if options.ShowLineNumbers {
		parts = append(parts, fmt.Sprintf("%d", variable.LineNumber))
	}

	// Add variable key
	parts = append(parts, variable.Key)

	// Add equals sign
	parts = append(parts, "=")

	// Add value (masked if sensitive)
	value := variable.Value
	if variable.IsSensitive && options.MaskSensitive {
		value = f.maskValue(value)
	}
	parts = append(parts, value)

	// Add comment if present and requested
	if variable.Comment != "" && options.ShowComments {
		parts = append(parts, fmt.Sprintf("# %s", variable.Comment))
	}

	// Add sensitive indicator
	if variable.IsSensitive {
		parts = append(parts, "🔒")
	}

	return strings.Join(parts, "\t")
}

// maskValue masks a sensitive value
func (f *Formatter) maskValue(value string) string {
	if value == "" {
		return ""
	}

	if len(value) <= 4 {
		return strings.Repeat("*", len(value))
	}

	// Show first and last character, mask the rest
	return string(value[0]) + strings.Repeat("*", len(value)-2) + string(value[len(value)-1])
}

// formatJSON formats the .env file as JSON
func (f *Formatter) formatJSON(envFile *util.EnvFile, options *util.FormatOptions) (string, error) {
	// Create a safe copy for JSON output
	safeEnvFile := f.createSafeCopy(envFile, options)

	// Convert to JSON
	jsonData, err := json.MarshalIndent(safeEnvFile, "", strings.Repeat(" ", options.IndentSize))
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return string(jsonData), nil
}

// formatYAML formats the .env file as YAML
func (f *Formatter) formatYAML(envFile *util.EnvFile, options *util.FormatOptions) (string, error) {
	var result strings.Builder

	// Add file info
	result.WriteString(fmt.Sprintf("file_path: %s\n", envFile.FilePath))
	result.WriteString(fmt.Sprintf("statistics:\n"))
	result.WriteString(fmt.Sprintf("  total_lines: %d\n", envFile.TotalLines))
	result.WriteString(fmt.Sprintf("  valid_lines: %d\n", envFile.ValidLines))
	result.WriteString(fmt.Sprintf("  comment_lines: %d\n", envFile.CommentLines))
	result.WriteString(fmt.Sprintf("  empty_lines: %d\n", envFile.EmptyLines))
	result.WriteString(fmt.Sprintf("  sensitive_variables: %d\n", f.countSensitiveVariables(envFile)))

	// Add warnings
	if len(envFile.Warnings) > 0 {
		result.WriteString("warnings:\n")
		for _, warning := range envFile.Warnings {
			result.WriteString(fmt.Sprintf("  - %s\n", warning))
		}
	}

	// Add parse errors
	if len(envFile.ParseErrors) > 0 {
		result.WriteString("parse_errors:\n")
		for _, err := range envFile.ParseErrors {
			result.WriteString(fmt.Sprintf("  - line: %d\n", err.LineNumber))
			result.WriteString(fmt.Sprintf("    message: %s\n", err.Message))
			if options.ShowLineNumbers {
				result.WriteString(fmt.Sprintf("    content: %s\n", err.Line))
			}
		}
	}

	// Add variables
	if len(envFile.Variables) > 0 {
		result.WriteString("variables:\n")
		for _, variable := range envFile.Variables {
			result.WriteString(fmt.Sprintf("  %s:\n", variable.Key))
			
			value := variable.Value
			if variable.IsSensitive && options.MaskSensitive {
				value = f.maskValue(value)
			}
			result.WriteString(fmt.Sprintf("    value: %s\n", value))
			result.WriteString(fmt.Sprintf("    sensitive: %t\n", variable.IsSensitive))
			
			if options.ShowLineNumbers {
				result.WriteString(fmt.Sprintf("    line: %d\n", variable.LineNumber))
			}
			
			if variable.Comment != "" && options.ShowComments {
				result.WriteString(fmt.Sprintf("    comment: %s\n", variable.Comment))
			}
		}
	}

	return result.String(), nil
}

// createSafeCopy creates a safe copy of the env file for JSON output
func (f *Formatter) createSafeCopy(envFile *util.EnvFile, options *util.FormatOptions) map[string]interface{} {
	safeCopy := map[string]interface{}{
		"file_path": envFile.FilePath,
		"statistics": map[string]interface{}{
			"total_lines":         envFile.TotalLines,
			"valid_lines":         envFile.ValidLines,
			"comment_lines":       envFile.CommentLines,
			"empty_lines":         envFile.EmptyLines,
			"sensitive_variables": f.countSensitiveVariables(envFile),
		},
	}

	// Add warnings
	if len(envFile.Warnings) > 0 {
		safeCopy["warnings"] = envFile.Warnings
	}

	// Add parse errors
	if len(envFile.ParseErrors) > 0 {
		errors := make([]map[string]interface{}, len(envFile.ParseErrors))
		for i, err := range envFile.ParseErrors {
			errorMap := map[string]interface{}{
				"line_number": err.LineNumber,
				"message":     err.Message,
			}
			if options.ShowLineNumbers {
				errorMap["line"] = err.Line
			}
			errors[i] = errorMap
		}
		safeCopy["parse_errors"] = errors
	}

	// Add variables
	if len(envFile.Variables) > 0 {
		variables := make([]map[string]interface{}, len(envFile.Variables))
		for i, variable := range envFile.Variables {
			value := variable.Value
			if variable.IsSensitive && options.MaskSensitive {
				value = f.maskValue(value)
			}

			varMap := map[string]interface{}{
				"key":        variable.Key,
				"value":      value,
				"sensitive":  variable.IsSensitive,
			}

			if options.ShowLineNumbers {
				varMap["line_number"] = variable.LineNumber
			}

			if variable.Comment != "" && options.ShowComments {
				varMap["comment"] = variable.Comment
			}

			variables[i] = varMap
		}
		safeCopy["variables"] = variables
	}

	return safeCopy
}

// countSensitiveVariables counts the number of sensitive variables
func (f *Formatter) countSensitiveVariables(envFile *util.EnvFile) int {
	count := 0
	for _, variable := range envFile.Variables {
		if variable.IsSensitive {
			count++
		}
	}
	return count
}

// FormatValidationResult formats a validation result
func (f *Formatter) FormatValidationResult(result *util.ValidationResult, options *util.FormatOptions) (string, error) {
	var output strings.Builder

	// Add header
	output.WriteString("🔍 Security Validation Results\n")
	output.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	// Add overall status
	status := "✅ Valid"
	if !result.IsValid {
		status = "❌ Invalid"
	}
	output.WriteString(fmt.Sprintf("Status: %s\n", status))
	output.WriteString(fmt.Sprintf("Risk Level: %s\n", result.RiskLevel))

	// Add sensitive variables count
	if len(result.SensitiveVars) > 0 {
		output.WriteString(fmt.Sprintf("Sensitive Variables: %d\n", len(result.SensitiveVars)))
	}

	// Add errors
	if len(result.Errors) > 0 {
		output.WriteString("\n❌ Errors:\n")
		for _, err := range result.Errors {
			output.WriteString(fmt.Sprintf("   • %s (Line %d): %s [%s]\n", 
				err.Variable, err.LineNumber, err.Message, err.Severity))
		}
	}

	// Add warnings
	if len(result.Warnings) > 0 {
		output.WriteString("\n⚠️  Warnings:\n")
		for _, warning := range result.Warnings {
			output.WriteString(fmt.Sprintf("   • %s (Line %d): %s\n", 
				warning.Variable, warning.LineNumber, warning.Message))
		}
	}

	// Add sensitive variables list
	if len(result.SensitiveVars) > 0 {
		output.WriteString("\n🔒 Sensitive Variables Detected:\n")
		for _, varName := range result.SensitiveVars {
			output.WriteString(fmt.Sprintf("   • %s\n", varName))
		}
	}

	return output.String(), nil
}

package env

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/devlink/internal/env"
	"github.com/devlink/internal/util"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get [code]",
	Short: "📥 Get a shared environment file",
	Long: `📥 Get a shared environment file

Download and display a shared environment file using a share code.

Examples:
  devlink env get ABC123              # Get and display the file
  devlink env get ABC123 --output .env # Save to a file
  devlink env get ABC123 -o config.env # Save with short flag`,
	Args:          cobra.ExactArgs(1),
	RunE:          runGet,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func runGet(cmd *cobra.Command, args []string) error {
	shareCode := args[0]
	outputFile, _ := cmd.Flags().GetString("output")

	validationManager := util.NewValidationManager()
	if err := validationManager.ValidateShareCode(shareCode, "cli-client"); err != nil {
		return fmt.Errorf("❌ %w", err)
	}

	if outputFile != "" {
		if err := validateOutputFile(outputFile); err != nil {
			return fmt.Errorf("❌ %w", err)
		}
	}

	fmt.Printf("🔍 Retrieving: %s\n", shareCode)

	content := os.Getenv("DEVLINK_SAMPLE_CONTENT")
	if content == "" {
		content = `# Sample environment file
NODE_ENV=development
DEBUG=false
LOG_LEVEL=info
PORT=3000`
	}

	parser := env.NewParser()
	envFile, err := parser.ParseContent(content, "retrieved")
	if err != nil {
		return fmt.Errorf("❌ failed to parse content: %w", err)
	}

	formatter := env.NewFormatter()
	formatOptions := &util.FormatOptions{
		MaskSensitive:   true,
		ShowComments:    true,
		ShowLineNumbers: false,
		OutputFormat:    "text",
	}

	formattedOutput, err := formatter.Format(envFile, formatOptions)
	if err != nil {
		return fmt.Errorf("❌ failed to format output: %w", err)
	}

	if outputFile != "" {
		fmt.Printf("💾 Saving to: %s\n", outputFile)
		if err := writeToFile(outputFile, content); err != nil {
			return fmt.Errorf("❌ failed to save file: %w", err)
		}
		fmt.Printf("✅ Saved successfully to: %s\n", outputFile)
	} else {
		fmt.Println("\n📄 Environment file content:")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println(formattedOutput)
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	}

	return nil
}

func validateOutputFile(filePath string) error {
	if err := validateOutputPath(filePath); err != nil {
		return fmt.Errorf("invalid output path: %w", err)
	}

	dir := filepath.Dir(filePath)
	if dir != "." {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return fmt.Errorf("directory does not exist: %s", dir)
		}
	}

	if _, err := os.Stat(filePath); err == nil {
		if _, err := os.OpenFile(filePath, os.O_WRONLY, 0644); err != nil {
			return fmt.Errorf("file is not writable: %w", err)
		}
	}

	return nil
}

// validateOutputPath validates the output file path for security
func validateOutputPath(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	// Get absolute path
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("invalid file path: %w", err)
	}

	// Prevent path traversal attacks
	if strings.Contains(absPath, "..") {
		return fmt.Errorf("path traversal not allowed")
	}

	// Restrict to current directory and subdirectories
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot determine current directory: %w", err)
	}

	if !strings.HasPrefix(absPath, currentDir) {
		return fmt.Errorf("output path must be within current directory")
	}

	// Check for dangerous file extensions
	dangerousExts := []string{".exe", ".bat", ".sh", ".py", ".js", ".php"}
	ext := strings.ToLower(filepath.Ext(filePath))
	for _, dangerousExt := range dangerousExts {
		if ext == dangerousExt {
			return fmt.Errorf("dangerous file extension not allowed: %s", ext)
		}
	}

	return nil
}

func writeToFile(filePath, content string) error {
	dir := filepath.Dir(filePath)
	if dir != "." {
		if err := os.MkdirAll(dir, 0750); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// Use secure permissions (owner read/write only)
	if err := os.WriteFile(filePath, []byte(strings.TrimSpace(content)), 0600); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func init() {
	getCmd.Flags().StringP("output", "o", "", "Save to file (default: display only)")
}

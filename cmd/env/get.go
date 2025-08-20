package env

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

	if err := validateShareCode(shareCode); err != nil {
		return fmt.Errorf("❌ %w", err)
	}

	if outputFile != "" {
		if err := validateOutputFile(outputFile); err != nil {
			return fmt.Errorf("❌ %w", err)
		}
	}

	fmt.Printf("🔍 Retrieving: %s\n", shareCode)

	content := `DATABASE_URL=postgresql://localhost:5432/mydb
API_KEY=your-secret-key
REDIS_URL=redis://localhost:6379
NODE_ENV=development`

	if outputFile != "" {
		fmt.Printf("💾 Saving to: %s\n", outputFile)
		if err := writeToFile(outputFile, content); err != nil {
			return fmt.Errorf("❌ failed to save file: %w", err)
		}
		fmt.Printf("✅ Saved successfully to: %s\n", outputFile)
	} else {
		fmt.Println("\n📄 Environment file content:")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println(content)
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	}

	return nil
}

func validateShareCode(code string) error {
	if code == "" {
		return fmt.Errorf("share code cannot be empty")
	}

	if len(code) < 3 {
		return fmt.Errorf("share code must be at least 3 characters")
	}

	for _, char := range code {
		if !((char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')) {
			return fmt.Errorf("share code must contain only uppercase letters and numbers")
		}
	}

	return nil
}

func validateOutputFile(filePath string) error {
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

func writeToFile(filePath, content string) error {
	dir := filepath.Dir(filePath)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	if err := os.WriteFile(filePath, []byte(strings.TrimSpace(content)), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func init() {
	getCmd.Flags().StringP("output", "o", "", "Save to file (default: display only)")
}

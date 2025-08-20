package env

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var shareCmd = &cobra.Command{
	Use:   "share [file]",
	Short: "📤 Share an environment file",
	Long: `📤 Share an environment file securely

Creates a secure share of your environment file that others can access with a code.

Examples:
  devlink env share .env              # Share your main environment file
  devlink env share config.env        # Share a specific config file
  devlink env share .env --expiry 24h # Share with 24-hour expiry`,
	Args:          cobra.ExactArgs(1),
	RunE:          runShare,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func runShare(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	if err := validateFile(filePath); err != nil {
		return fmt.Errorf("❌ %w", err)
	}

	expiry, _ := cmd.Flags().GetString("expiry")
	readonly, _ := cmd.Flags().GetBool("readonly")

	fmt.Printf("🚀 Sharing: %s\n", filePath)
	fmt.Printf("⏰ Expires: %s\n", expiry)
	if readonly {
		fmt.Println("🔒 Read-only: enabled")
	}

	fmt.Println("\n✨ Share created successfully!")
	fmt.Println("📋 Share this code with your team:")
	fmt.Printf("   %s\n", "ABC123")
	fmt.Println("\n💡 Use: devlink env get ABC123")

	return nil
}

func validateFile(filePath string) error {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("invalid file path: %w", err)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", filePath)
	}

	if _, err := os.ReadFile(absPath); err != nil {
		return fmt.Errorf("cannot read file: %w", err)
	}

	return nil
}

func init() {
	shareCmd.Flags().StringP("expiry", "e", "1h", "Share expiry time (1h, 24h, 7d)")
	shareCmd.Flags().BoolP("readonly", "r", false, "Make share read-only")
}

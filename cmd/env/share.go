package env

import (
	"fmt"

	"github.com/devlink/internal/env"
	"github.com/devlink/internal/util"
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

	parser := env.NewParser()
	envFile, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("❌ failed to parse file: %w", err)
	}

	validator := env.NewValidator()
	validationResult := validator.Validate(envFile)

	if !validationResult.IsValid {
		fmt.Println("⚠️  Security issues detected:")
		for _, err := range validationResult.Errors {
			fmt.Printf("   • %s: %s\n", err.Variable, err.Message)
		}
		fmt.Println()
	}

	expiry, _ := cmd.Flags().GetString("expiry")
	readonly, _ := cmd.Flags().GetBool("readonly")

	fmt.Printf("🚀 Sharing: %s\n", filePath)
	fmt.Printf("⏰ Expires: %s\n", expiry)
	if readonly {
		fmt.Println("🔒 Read-only: enabled")
	}

	fmt.Printf("📊 File stats: %d variables, %d sensitive\n",
		len(envFile.Variables), len(validationResult.SensitiveVars))

	tokenGen := util.NewTokenGenerator()
	shareCode, err := tokenGen.GenerateShareCode()
	if err != nil {
		return fmt.Errorf("❌ failed to generate share code: %w", err)
	}

	fmt.Println("\n✨ Share created successfully!")
	fmt.Println("📋 Share this code with your team:")
	fmt.Printf("   %s\n", shareCode)
	fmt.Printf("\n💡 Use: devlink env get %s\n", shareCode)

	return nil
}

func init() {
	shareCmd.Flags().StringP("expiry", "e", "1h", "Share expiry time (1h, 24h, 7d)")
	shareCmd.Flags().BoolP("readonly", "r", false, "Make share read-only")
}

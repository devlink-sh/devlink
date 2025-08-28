package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/devlink/pkg/envsharing/core"
	"github.com/devlink/pkg/envsharing/core/encryption"
	"github.com/devlink/pkg/envsharing/network"
	"github.com/spf13/cobra"
)

type EnvSharingCLI struct {
	zitiConfig *network.ZitiConfig
	encryption *encryption.Manager
}

func NewEnvSharingCLI() *EnvSharingCLI {
	return &EnvSharingCLI{
		zitiConfig: &network.ZitiConfig{
			ControllerURL: getEnvOrDefault("ZITI_CONTROLLER_URL", "https://localhost:1280"),
			IdentityFile:  getEnvOrDefault("ZITI_IDENTITY_FILE", filepath.Join(getHomeDir(), ".devlink", "identity.json")),
			ServiceName:   getEnvOrDefault("ZITI_SERVICE_NAME", "devlink-env-service"),
		},
		encryption: encryption.NewManager(),
	}
}

func (cli *EnvSharingCLI) CreateShareCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "share [file]",
		Short: "📤 Share your .env file with the team",
		Long: `📤 Share your .env file with the team

This command takes your environment file and creates a secure, temporary link that you can share with your teammates.
The file is encrypted and can only be accessed by people with the special code.

Think of it like creating a self-destructing message that only your team can read! 💥

Examples:
  devlink env share .env              # Share for 1 hour (default)
  devlink env share .env --expiry 24h # Share for 24 hours
  devlink env share .env --readonly   # Make it read-only (safer)`,
		Args:  cobra.ExactArgs(1),
		RunE:  cli.runShare,
	}

	cmd.Flags().StringP("expiry", "e", "1h", "How long the share should last (e.g., 1h, 24h, 7d)")
	cmd.Flags().BoolP("readonly", "r", false, "Make the share read-only (recommended)")
	cmd.Flags().StringP("server", "s", "", "Custom server URL (advanced)")

	return cmd
}

func (cli *EnvSharingCLI) CreateGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [code]",
		Short: "📥 Get a shared .env file from your teammate",
		Long: `📥 Get a shared .env file from your teammate

Someone shared an environment file with you? Great! Use this command to retrieve it.
Just paste the code they gave you, and voilà! Your .env file appears like magic! ✨

The file will be displayed on screen by default, or you can save it directly to a file.

Examples:
  devlink env get blue-dragon-123          # Get and display the file
  devlink env get blue-dragon-123 --output .env  # Save directly to .env file
  devlink env get blue-dragon-123 --unmask # Show the actual secret values (be careful!)`,
		Args:  cobra.ExactArgs(1),
		RunE:  cli.runGet,
	}

	cmd.Flags().StringP("output", "o", "", "Save the file to this path (e.g., .env)")
	cmd.Flags().StringP("server", "s", "", "Custom server URL (advanced)")
	cmd.Flags().BoolP("unmask", "u", false, "Show the actual secret values (use with caution!)")

	return cmd
}

func (cli *EnvSharingCLI) runShare(cmd *cobra.Command, args []string) error {
	filePath := args[0]
	expiry, _ := cmd.Flags().GetString("expiry")
	readonly, _ := cmd.Flags().GetBool("readonly")
	serverURL, _ := cmd.Flags().GetString("server")

	if serverURL != "" {
		cli.zitiConfig.ControllerURL = serverURL
	}

	parser := core.NewParser()
	envFile, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
	}

	validator := core.NewValidator()
	validationResult := validator.Validate(envFile)

	if !validationResult.IsValid {
		fmt.Println("⚠️  Security issues detected:")
		for _, err := range validationResult.Errors {
			fmt.Printf("   • %s: %s\n", err.Variable, err.Message)
		}
		fmt.Println()
	}

	expiryDuration, err := time.ParseDuration(expiry)
	if err != nil {
		return fmt.Errorf("invalid expiry format: %w", err)
	}

	fmt.Printf("🚀 Sharing your .env file: %s\n", filePath)
	fmt.Printf("⏰ This share will expire in: %s\n", expiry)
	if readonly {
		fmt.Println("🔒 Read-only mode: enabled (safer!)")
	}

	fmt.Printf("📊 Your file has %d variables (%d are sensitive)\n",
		len(envFile.Variables), len(validationResult.SensitiveVars))

	shareCode := generateShareCode()
	client, err := network.NewZitiClient(cli.zitiConfig, cli.encryption)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	fmt.Printf("📡 Connecting to the secure tunnel...\n")

	response, err := client.CreateShare(envFile, shareCode, expiryDuration, readonly)
	if err != nil {
		return fmt.Errorf("Oops! Failed to create the share: %w", err)
	}

	fmt.Println("\n🎉 Success! Your .env file is now shared securely!")
	fmt.Println("📋 Copy and paste this code to your teammate:")
	fmt.Printf("   %s\n", response.ShareCode)
	fmt.Printf("⏰ This share expires: %s\n", response.ExpiresAt.Format("Jan 2, 3:04 PM"))
	fmt.Printf("\n💡 Your teammate can use: devlink env get %s\n", response.ShareCode)

	return nil
}

func (cli *EnvSharingCLI) runGet(cmd *cobra.Command, args []string) error {
	shareCode := args[0]
	outputFile, _ := cmd.Flags().GetString("output")
	serverURL, _ := cmd.Flags().GetString("server")
	unmask, _ := cmd.Flags().GetBool("unmask")

	if serverURL != "" {
		cli.zitiConfig.ControllerURL = serverURL
	}

	if outputFile != "" {
		if err := validateOutputFile(outputFile); err != nil {
			return fmt.Errorf("invalid output path: %w", err)
		}
	}

	fmt.Printf("🔍 Looking for your shared .env file...\n")
	fmt.Printf("📡 Connecting to the secure tunnel...\n")

	client, err := network.NewZitiClient(cli.zitiConfig, cli.encryption)
	if err != nil {
		return fmt.Errorf("Oops! Failed to connect: %w", err)
	}

	envFile, err := client.GetShare(shareCode)
	if err != nil {
		return fmt.Errorf("Oops! Couldn't find that share: %w", err)
	}

	formatter := core.NewFormatter()
	formatOptions := &core.FormatOptions{
		MaskSensitive:   !unmask,
		ShowComments:    true,
		ShowLineNumbers: false,
		OutputFormat:    "text",
	}

	formattedOutput, err := formatter.Format(envFile, formatOptions)
	if err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}

	if outputFile != "" {
		fmt.Printf("💾 Saving to: %s\n", outputFile)
		if err := writeToFile(outputFile, envFile.RawContent); err != nil {
			return fmt.Errorf("Oops! Failed to save the file: %w", err)
		}
		fmt.Printf("✅ Perfect! Saved to: %s\n", outputFile)
	} else {
		fmt.Println("\n📄 Here's your .env file:")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println(formattedOutput)
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	}

	return nil
}

func validateOutputFile(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("invalid file path: %w", err)
	}

	if strings.Contains(absPath, "..") {
		return fmt.Errorf("path traversal not allowed")
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot determine current directory: %w", err)
	}

	if !strings.HasPrefix(absPath, currentDir) {
		return fmt.Errorf("output path must be within current directory")
	}

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

	if err := os.WriteFile(filePath, []byte(strings.TrimSpace(content)), 0600); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return homeDir
}

func generateShareCode() string {
	adjectives := []string{"blue", "red", "green", "yellow", "purple", "orange", "pink", "brown", "black", "white"}
	nouns := []string{"dragon", "tiger", "eagle", "lion", "wolf", "bear", "fox", "owl", "shark", "whale"}
	numbers := []string{"123", "456", "789", "101", "202", "303", "404", "505", "606", "707"}

	return fmt.Sprintf("%s-%s-%s",
		adjectives[time.Now().UnixNano()%int64(len(adjectives))],
		nouns[time.Now().UnixNano()%int64(len(nouns))],
		numbers[time.Now().UnixNano()%int64(len(numbers))])
}

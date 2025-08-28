package cmd

import (
	"github.com/devlink/pkg/envsharing/cli"
	"github.com/spf13/cobra"
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "🔐 Share .env files securely with your team",
	Long: `🔐 Share .env files securely with your team

Ever wanted to share your environment variables without exposing them to the internet? 
Well, you're in luck! DevLink uses zero-trust networking to keep your secrets safe.

Think of it like passing a secret note through an invisible, encrypted tunnel. 
Only people with the right "key" can read it, and it disappears after use! ✨

Examples:
  devlink env share .env                    # Share your environment file
  devlink env get blue-dragon-123          # Get a shared environment file
  devlink env get blue-dragon-123 --output .env  # Save to a file`,
}

func init() {
	cli := cli.NewEnvSharingCLI()

	envCmd.AddCommand(cli.CreateShareCommand())
	envCmd.AddCommand(cli.CreateGetCommand())

	rootCmd.AddCommand(envCmd)
}

package env

import (
	"github.com/spf13/cobra"
)

var EnvCmd = &cobra.Command{
	Use:   "env",
	Short: "🔐 Share environment files securely",
	Long: `🔐 Environment file sharing for DevLink

Share your .env files securely with team members using simple codes.
Perfect for sharing database credentials, API keys, and configuration files.

Examples:
  devlink env share .env                    # Share your environment file
  devlink env get ABC123                    # Get a shared environment file
  devlink env get ABC123 --output .env      # Save to a file`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	EnvCmd.AddCommand(shareCmd)
	EnvCmd.AddCommand(getCmd)
	EnvCmd.AddCommand(templateCmd)
	EnvCmd.AddCommand(bulkCmd)
	EnvCmd.AddCommand(searchCmd)
	EnvCmd.AddCommand(completionCmd)
}

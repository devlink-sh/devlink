package cmd

import (
	"fmt"
	"os"

	"github.com/devlink/cmd/db"
	"github.com/devlink/cmd/directory"
	"github.com/devlink/cmd/env"
	"github.com/devlink/cmd/git"
	"github.com/devlink/cmd/hive"
	"github.com/devlink/cmd/pair"
	"github.com/devlink/cmd/registry"
	"github.com/spf13/cobra"
)

const banner = `
╔══════════════════════════════════════════════════════════════╗
║                                                              ║
║    ██████╗ ███████╗██╗   ██╗██╗     ██╗███╗   ██╗██╗  ██╗    ║
║   ██╔══██╗██╔════╝██║   ██║██║     ██║████╗  ██║██║ ██╔╝     ║
║   ██║  ██║█████╗  ██║   ██║██║     ██║██╔██╗ ██║█████╔╝      ║
║   ██║  ██║██╔══╝  ██║   ██║██║     ██║██║╚██╗██║██╔═██╗      ║
║   ██████╔╝███████╗╚██████╔╝███████╗██║██║ ╚████║██║  ██╗     ║
║   ╚═════╝ ╚══════╝ ╚═════╝ ╚══════╝╚═╝╚═╝  ╚═══╝╚═╝  ╚═╝     ║
║                                                              ║
║                    Development Workflow CLI                  ║
║                                                              ║
╚══════════════════════════════════════════════════════════════╝

`

var rootCmd = &cobra.Command{
	Use:   "devlink",
	Short: "🚀 DevLink - Streamline your development workflow",
	Long: fmt.Sprintf(`%s
DevLink is a powerful CLI tool that streamlines development workflows
by providing efficient link management, project organization, and 
developer utilities.

📚 Available Commands:
  • git       - Share and connect to Git repositories
  • pair      - Share your localhost with teammates
  • env       - Share development environments
  • db        - Database management and sharing
  • registry  - Docker registry management

💡 Quick Start:
  devlink git serve <repo-path>     # Share a repository
  devlink pair                      # Share localhost
  devlink help <command>            # Get detailed help

Use 'devlink help' to explore all available commands.`, banner),
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Global flags can be added here
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
	rootCmd.AddCommand(env.EnvCmd)
	rootCmd.AddCommand(db.DBCmd)
	rootCmd.AddCommand(pair.PairCmd)
	rootCmd.AddCommand(registry.RegistryCmd)
	rootCmd.AddCommand(git.GitCmd)
	rootCmd.AddCommand(directory.DirectoryCmd)
	rootCmd.AddCommand(hive.HiveCmd)
}

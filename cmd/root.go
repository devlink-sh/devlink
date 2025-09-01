package cmd

import (
	"fmt"
	"os"

	"github.com/devlink/cmd/db"
	"github.com/devlink/cmd/env"
	"github.com/devlink/cmd/pair"
	"github.com/spf13/cobra"
)

const banner = `
 ____             _     _       _    
|  _ \  _____   _| |   (_)_ __ | | __
| | | |/ _ \ \ / / |   | | '_ \| |/ /
| |_| |  __/\ V /| |___| | | | |   < 
|____/ \___| \_/ |_____|_|_| |_|_|\_\

🚀 DevLink - Development Workflow Management CLI
Version: 1.0.0
`

var rootCmd = &cobra.Command{
	Use:   "devlink",
	Short: "DevLink - A powerful CLI for development workflow management",
	Long: fmt.Sprintf(`%s
DevLink is a comprehensive CLI tool designed to streamline development workflows
by providing efficient link management, project organization, and developer utilities.

Use 'devlink help' to see available commands.`, banner),
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
}

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	// Version is the current version of DevLink CLI
	Version = "1.0.0"
	// AppName is the name of the application
	AppName = "DevLink"
)

const banner = `
 ____             _     _       _    
|  _ \  _____   _| |   (_)_ __ | | __
| | | |/ _ \ \ / / |   | | '_ \| |/ /
| |_| |  __/\ V /| |___| | | | |   < 
|____/ \___| \_/ |_____|_|_| |_|_|\_\

🚀 DevLink - Development Workflow Management CLI
Version: ` + Version + `
`

var rootCmd = &cobra.Command{
	Use:   "devlink",
	Short: "DevLink - A powerful CLI for development workflow management",
	Long: fmt.Sprintf(`%s
DevLink is a comprehensive CLI tool designed to streamline development workflows
by providing efficient link management, project organization, and developer utilities.

Use 'devlink help' to see available commands.`, banner),
	Run: func(cmd *cobra.Command, args []string) {
		showVersion, _ := cmd.Flags().GetBool("version")
		if showVersion {
			fmt.Printf("%s CLI v%s\n", AppName, Version)
			return
		}
		fmt.Print(banner)
		fmt.Println("\nWelcome to DevLink! Use 'devlink help' to get started.")
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolP("version", "V", false, "Show version information")
}

package env

import (
	"github.com/spf13/cobra"
)

// envCmd represents the env command
var EnvCmd = &cobra.Command{
	Use:   "env",
	Short: "Commands for sharing development environments",
}

func init() {
	EnvCmd.AddCommand(envShareCmd)
	EnvCmd.AddCommand(envGetCmd)
}

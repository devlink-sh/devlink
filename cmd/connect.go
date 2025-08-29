package cmd

import (
	"fmt"

	"github.com/gupta-nu/devlink/internal/gitutils"
	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:   "connect [url]",
	Short: "Connect to a remote git repo shared via devlink-git serve",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		url := args[0]
		fmt.Println("Connecting to remote repo:", url)

		if err := gitutils.CloneOrFetch(url); err != nil {
			return err
		}
		fmt.Println("Successfully connected to remote repo")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}

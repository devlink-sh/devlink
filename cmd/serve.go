package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gupta-nu/devlink/internal/gitutils"
	"github.com/gupta-nu/devlink/internal/signal"
	"github.com/gupta-nu/devlink/internal/zrok"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Expose the current git repo over a zrok tunnel",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 1. detect repo
		repoPath, err := gitutils.DetectRepo()
		if err != nil {
			return err
		}
		fmt.Println("Serving repo from:", repoPath)

		// 2. start zrok tunnel
		url, err := zrok.StartTunnel(repoPath)
		if err != nil {
			return err
		}
		fmt.Println("Tunnel available at:", url)

		// 3. wait for ctrl+c
		signal.WaitForInterrupt()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

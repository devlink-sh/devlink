package cmd

import (
	"fmt"
	"os"

	"github.com/gupta-nu/devlink/internal/git"
	"github.com/gupta-nu/devlink/internal/p2p"

	"github.com/spf13/cobra"
)

var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "Ephemeral git collaboration",
}

var shareCmd = &cobra.Command{
	Use:   "share",
	Short: "Share this repo temporarily",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, _ := os.Getwd()
		port := 9418 // git default

		go git.StartDaemon(cwd, port)

		addr, err := p2p.StartTunnel(port)
		if err != nil {
			return err
		}
		fmt.Printf("Share link:\n  git clone git://%s/\n", addr)
		return nil
	},
}

func init() {
	gitCmd.AddCommand(shareCmd)
	rootCmd.AddCommand(gitCmd)
}

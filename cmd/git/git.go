package git

import "github.com/spf13/cobra"

var GitCmd = &cobra.Command{
	Use:   "git",
	Short: "Git repository sharing commands",
	Long:  `Commands to share and connect to Git repositories in real time.`,
}

func init() {
	GitCmd.AddCommand(gitServeCmd)
	GitCmd.AddCommand(gitConnectCmd)
}

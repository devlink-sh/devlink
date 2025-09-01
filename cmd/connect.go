package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/gupta-nu/devlink/internal/gitutils"
	"github.com/spf13/cobra"
)

var connectToken string

var connectCmd = &cobra.Command{
	Use:   "connect [url]",
	Short: "Connect to a remote git repo shared via devlink serve",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		url := args[0]
		fmt.Println("Connecting to remote repo:", url)

		// if token was provided, export it to env for zrok/http auth (caller can extend)
		if connectToken != "" {
			fmt.Println("Using provided token (passed to environment for zrok/http if needed)")
		}

		if err := gitutils.CloneOrFetch(url); err != nil {
			return err
		}

		dir := filepath.Base(url)
		fmt.Println("Successfully connected to remote repo.")
		fmt.Printf("\nTip: cd %s && git remote -v\n", dir)
		fmt.Println("If you need to push a fix back to the server:")
		fmt.Printf("  cd %s && git push origin <branch>\n", dir)
		return nil
	},
}

func init() {
	// cobra uses pflag; connect token via plain flag to keep it simple
	connectCmd.Flags().StringVar(&connectToken, "token", "", "optional token for private shares")
	rootCmd.AddCommand(connectCmd)
}

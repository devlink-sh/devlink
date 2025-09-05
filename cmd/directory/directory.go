package directory

import "github.com/spf13/cobra"

var DirectoryCmd = &cobra.Command{
	Use:   "dir",
	Short: "Manage directories",
}

func init() {
	DirectoryCmd.AddCommand(directoryShareCmd)
	DirectoryCmd.AddCommand(directoryGetCmd)
}

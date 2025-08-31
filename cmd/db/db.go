package db

import "github.com/spf13/cobra"

var DBCmd = &cobra.Command{
	Use:   "db",
	Short: "Database commands",
	Long:  `Commands for managing the database`,
}

func init() {
	DBCmd.AddCommand(dbShareCmd)
	DBCmd.AddCommand(dbGetCmd)
}
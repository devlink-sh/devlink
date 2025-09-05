package hive

import "github.com/spf13/cobra"

// Centralized BaseURL
var BaseURL = "https://tazidgt171sl.share.zrok.io"

// Root "hive" command
var HiveCmd = &cobra.Command{
	Use:   "hive",
	Short: "Ephemeral staging environments for your team",
}

func init() {
	HiveCmd.AddCommand(hiveCreateCmd)
	HiveCmd.AddCommand(hiveContributeCmd)
	HiveCmd.AddCommand(hiveConnectCmd)
}

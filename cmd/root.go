package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "devlink",
	Short: "Devlink Git helps you share and connect repos instantly",
	Long:  `devlink — minimal MVP that wraps git daemon behind a zrok tunnel.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

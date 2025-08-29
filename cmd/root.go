package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "devlink-git",
	Short: "Devlink Git helps you share and connect repos instantly",
	Long:  `Devlink Git is a CLI tool that helps you expose local git repos via zrok tunnels and connect to others seamlessly.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

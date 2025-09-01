package registry

import (
	"github.com/spf13/cobra"
)

var RegistryCmd = &cobra.Command{
	Use:   "registry",
	Short: "Manage Docker registry shares",
}

func init() {
	RegistryCmd.AddCommand(registryShareCmd)
	RegistryCmd.AddCommand(registryGetCmd)
}

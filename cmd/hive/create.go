package hive

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

var hiveCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new Hive",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		// Use BaseURL for consistency
		resp, err := http.Get(fmt.Sprintf("%s/hives/create?name=%s", BaseURL, name))
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		tokenBytes, _ := io.ReadAll(resp.Body)
		token := string(tokenBytes)

		log.Printf("Hive '%s' created! Invite token: %s", name, token)
	},
}

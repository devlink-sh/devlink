package hive

import (
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

		resp, err := http.Get("http://localhost:8081/hives/create?name=" + name)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		tokenBytes, _ := io.ReadAll(resp.Body)
		token := string(tokenBytes)

		log.Printf("Hive '%s' created! Invite token: %s", name, token)

	},
}

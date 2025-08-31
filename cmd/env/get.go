package env

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/openziti/zrok/environment"
	"github.com/openziti/zrok/sdk/golang/sdk"
	"github.com/spf13/cobra"
)

var envGetCmd = &cobra.Command{
	Use:   "get <token>",
	Short: "Retrieve shared environment variables",
	Long:  `Connect to a shared environment and save the received .env file in the current directory.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		token := args[0]

		root, err := environment.LoadRoot()
		if err != nil {
			log.Fatal(err)
		}

		// create access so the service has a terminator to dial
		acc, err := sdk.CreateAccess(root, &sdk.AccessRequest{ShareToken: token})
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := sdk.DeleteAccess(root, acc); err != nil {
				log.Printf("error deleting access: %v", err)
			}
		}()

		// this returns a connected net.Conn (no Dial() needed)
		conn, err := sdk.NewDialer(token, root)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		destPath := filepath.Join(".", ".env.received")
		out, err := os.Create(destPath)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()

		n, err := io.Copy(out, conn) // stream until EOF
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Received %d bytes -> %s", n, destPath)
	},
}

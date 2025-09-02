// connect.go
package git

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/devlink/internal"
	"github.com/openziti/zrok/environment"
	"github.com/openziti/zrok/sdk/golang/sdk"
	"github.com/spf13/cobra"
)

var gitConnectCmd = &cobra.Command{
	Use:   "connect <share-token> <repo-name>",
	Short: "Connect to a shared Git repository",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		token := args[0]
		repoName := args[1]

		root, err := environment.LoadRoot()
		if err != nil {
			log.Fatal(err)
		}

		acc, err := sdk.CreateAccess(root, &sdk.AccessRequest{ShareToken: token})
		if err != nil {
			log.Fatal(err)
		}
		defer sdk.DeleteAccess(root, acc)

		listener, err := net.Listen("tcp", "127.0.0.1:9418")
		if err != nil {
			log.Fatal(err)
		}
		defer listener.Close()

		log.Printf("Git tunnel ready at git://127.0.0.1:9418/%s", repoName)

		// Handle Ctrl+C
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			log.Println("Shutting down git connect...")
			_ = listener.Close()
			os.Exit(0)
		}()

		for {
			client, err := listener.Accept()
			if err != nil {
				log.Printf("error accepting client: %v", err)
				continue
			}

			go func(c net.Conn) {
				remote, err := sdk.NewDialer(token, root)
				if err != nil {
					log.Printf("error dialing zrok: %v", err)
					c.Close()
					return
				}
				log.Printf("Forwarding local Git client...")
				internal.Pipe(c, remote)
			}(client)
		}
	},
}

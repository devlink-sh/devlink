package db

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/openziti/zrok/environment"
	"github.com/openziti/zrok/sdk/golang/sdk"
	"github.com/spf13/cobra"
)

var dbShareCmd = &cobra.Command{
	Use:   "share <port>",
	Short: "Share a local database",
	Long:  `Securely share a local database over zrok. Example: devlink db share 5432`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		port := args[0]

		root, err := environment.LoadRoot()
		if err != nil {
			log.Fatal(err)
		}

		// Create the share
		share, err := sdk.CreateShare(root, &sdk.ShareRequest{
			BackendMode: sdk.TcpTunnelBackendMode,
			ShareMode:   sdk.PrivateShareMode,
			Target:      "db",
		})
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Database share ready! Let others connect using:\n  devlink db get %s <local-port>", share.Token)

		// Start listener for incoming zrok connections
		listener, err := sdk.NewListener(share.Token, root)
		if err != nil {
			log.Fatal(err)
		}
		defer listener.Close()

		// Handle SIGINT/SIGTERM cleanly
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			log.Println("Shutting down db share...")
			if err := sdk.DeleteShare(root, share); err != nil {
				log.Printf("error deleting share: %v", err)
			}
			_ = listener.Close()
			os.Exit(0)
		}()

		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("error accepting zrok connection: %v", err)
				continue
			}

			go func(remote net.Conn) {
				local, err := net.Dial("tcp", "127.0.0.1:"+port)
				if err != nil {
					log.Printf("error dialing local DB: %v", err)
					remote.Close()
					return
				}
				log.Printf("Forwarding DB connection -> localhost:%s", port)
				Pipe(remote, local)
			}(conn)
		}

	},
}

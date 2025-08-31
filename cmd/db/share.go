package db

import (
	"io"
	"log"
	"net"

	"github.com/openziti/zrok/environment"
	"github.com/openziti/zrok/sdk/golang/sdk"
	"github.com/spf13/cobra"
)

var dbShareCmd = &cobra.Command{
	Use:   "share <port>",
	Short: "Share database",
	Long:  `Share database with other commands`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		port := args[0]

		root, err := environment.LoadRoot()
		if err != nil {
			log.Fatal(err)
		}
		share, err := sdk.CreateShare(root, &sdk.ShareRequest{
			BackendMode: sdk.TcpTunnelBackendMode,
			ShareMode:   sdk.PrivateShareMode,
			Target:      "db",
		})

		if err != nil {
			log.Fatal(err)
		}

		log.Printf("access db using 'devlink db get %s'", share.Token)

		listener, err := sdk.NewListener(share.Token, root)
		if err != nil {
			log.Fatal(err)
		}
		defer listener.Close()

		        go func() {
            for {
                conn, err := listener.Accept()
                if err != nil {
                    log.Printf("error accepting: %v", err)
                    return
                }

                // Connect to actual DB port locally
                local, err := net.Dial("tcp", "127.0.0.1:"+port)
                if err != nil {
                    log.Printf("error dialing local db: %v", err)
                    conn.Close()
                    continue
                }

                // Pipe connections (DB <-> zrok)
                go io.Copy(local, conn)
                go io.Copy(conn, local)
            }
        }()

        select {}

	},
}

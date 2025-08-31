package db

import (
	"io"
	"log"
	"net"

	"github.com/openziti/zrok/environment"
	"github.com/openziti/zrok/sdk/golang/sdk"
	"github.com/spf13/cobra"
)


var dbGetCmd = &cobra.Command{
    Use:   "get <token> <port>",
    Short: "Connect to a shared database",
    Args:  cobra.ExactArgs(2),
    Run: func(cmd *cobra.Command, args []string) {
        token := args[0]
        port := args[1]

        root, err := environment.LoadRoot()
        if err != nil {
            log.Fatal(err)
        }

        acc, err := sdk.CreateAccess(root, &sdk.AccessRequest{ShareToken: token})
        if err != nil {
            log.Fatal(err)
        }
        defer sdk.DeleteAccess(root, acc)

        conn, err := sdk.NewDialer(token, root)
        if err != nil {
            log.Fatal(err)
        }
        defer conn.Close()

        // Expose locally on given port
        listener, err := net.Listen("tcp", "127.0.0.1:"+port)
        if err != nil {
            log.Fatal(err)
        }
        log.Printf("DB tunnel ready at 127.0.0.1:%s", port)

        for {
            client, err := listener.Accept()
            if err != nil {
                log.Printf("error accepting local conn: %v", err)
                continue
            }

            go func() {
                defer client.Close()
                // pipe local client <-> zrok conn
                go io.Copy(client, conn)
                go io.Copy(conn, client)
            }()
        }
    },
}

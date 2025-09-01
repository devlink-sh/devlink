package pair

import (
	// "io"
	"log"
	"net"

	"github.com/openziti/zrok/environment"
	"github.com/openziti/zrok/sdk/golang/sdk"
	"github.com/spf13/cobra"
)

var pairShareCmd = &cobra.Command{
	Use:  "share <port>",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		port := args[0]

		root, err := environment.LoadRoot()
		if err != nil {
			log.Fatal(err)
		}

		share, err := sdk.CreateShare(root, &sdk.ShareRequest{
			BackendMode: sdk.TcpTunnelBackendMode,
			ShareMode:   sdk.PrivateShareMode,
			Target:      "pair",
		})
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Frontend share ready! Let others connect with:\n  devlink pair get %s <local-port>", share.Token)

		listener, err := sdk.NewListener(share.Token, root)
		if err != nil {
			log.Fatal(err)
		}
		defer listener.Close()

		for{
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("error accepting incoming connection: %v", err)
				continue
			}

			go func(remote net.Conn){
				local, err := net.Dial("tcp", "127.0.0.1:"+port)
				if err != nil {
					log.Printf("error connecting to local service: %v", err)
					_ = remote.Close()
					return
				}
				log.Printf("Forwarding client -> localhost:%s", port)
				Pipe(remote, local)
			}(conn)
		}
	},
}

// func Pipe(a, b net.Conn) {
// 	defer a.Close()
// 	defer b.Close()

// 	done := make(chan struct{}, 2)

// 	go func() {
// 		io.Copy(a, b)
// 		done <- struct{}{}
// 	}()
// 	go func() {
// 		io.Copy(b, a)
// 		done <- struct{}{}
// 	}()

// 	<-done
// }

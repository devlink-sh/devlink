package pair

import (
	"io"
	"log"
	"net"

	"github.com/openziti/zrok/environment"
	"github.com/openziti/zrok/sdk/golang/sdk"
	"github.com/spf13/cobra"
)

var pairGetCmd = &cobra.Command{
	Use:   "get <token> <port>",
	Short: "Connect to a shared frontend",
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

		listener, err := net.Listen("tcp", "127.0.0.1:"+port)
		if err != nil {
			log.Fatal(err)
		}
		defer listener.Close()

		log.Printf("Frontend available locally at http://127.0.0.1:%s", port)

		for {
			client, err := listener.Accept()
			if err != nil {
				log.Printf("error accepting local connection: %v", err)
				continue
			}

			go func(c net.Conn) {
				remote, err := sdk.NewDialer(token, root)
				if err != nil {
					log.Printf("error dialing zrok: %v", err)
					c.Close()
					return
				}
				log.Printf("Browser connected, tunneling traffic...")
				Pipe(c, remote)
			}(client)
		}
	},
}

func Pipe(a, b net.Conn) {
	defer a.Close()
	defer b.Close()

	done := make(chan struct{}, 2)

	go func() {
		io.Copy(a, b)
		done <- struct{}{}
	}()
	go func() {
		io.Copy(b, a)
		done <- struct{}{}
	}()

	<-done
}

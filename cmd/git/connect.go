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
		shareToken := args[0]
		repoName := args[1] // 👈 we require repo name explicitly

		root, err := environment.LoadRoot()
		if err != nil {
			log.Fatal(err)
		}

		conn, err := sdk.NewConn(shareToken, root)
		if err != nil {
			log.Fatal(err)
		}

		// Dial local TCP port 9418
		local, err := net.Listen("tcp", "127.0.0.1:9418")
		if err != nil {
			log.Fatal(err)
		}
		defer local.Close()

		log.Printf("Connected to shared repo! You can now run:\n\n   git clone git://127.0.0.1:9418/%s\n", repoName)

		// Handle Ctrl+C
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			log.Println("Shutting down git connect...")
			_ = conn.Close()
			_ = local.Close()
			os.Exit(0)
		}()

		for {
			remote, err := conn.Accept()
			if err != nil {
				log.Printf("error accepting remote: %v", err)
				continue
			}

			go func(r net.Conn) {
				l, err := net.Dial("tcp", "127.0.0.1:9418")
				if err != nil {
					log.Printf("error dialing local git: %v", err)
					r.Close()
					return
				}
				internal.Pipe(r, l)
			}(remote)
		}
	},
}

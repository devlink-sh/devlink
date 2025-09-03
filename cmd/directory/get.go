package directory

import (
	"io"
	"log"
	"net"
	"os/exec"
	"runtime"

	"github.com/openziti/zrok/environment"
	"github.com/openziti/zrok/sdk/golang/sdk"
	"github.com/spf13/cobra"
)

var directoryGetCmd = &cobra.Command{
	Use:   "get <token> <port>",
	Short: "Access a shared directory in your browser",
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

		// Local listener for browser
		listener, err := net.Listen("tcp", "127.0.0.1:"+port)
		if err != nil {
			log.Fatal(err)
		}

		url := "http://127.0.0.1:" + port
		log.Printf("Open %s in your browser", url)

		// Auto-open browser tab
		openBrowser(url)

		for {
			client, err := listener.Accept()
			if err != nil {
				log.Printf("accept error: %v", err)
				continue
			}

			go func(c net.Conn) {
				defer c.Close()

				// Each browser request gets its own tunnel conn
				tunnelConn, err := sdk.NewDialer(token, root)
				if err != nil {
					log.Printf("error creating tunnel conn: %v", err)
					return
				}
				defer tunnelConn.Close()

				// Bi-directional copy
				go io.Copy(tunnelConn, c)
				io.Copy(c, tunnelConn)
			}(client)
		}
	},
}

// Open default browser cross-platform
func openBrowser(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	default: // Linux, BSD, etc.
		cmd = "xdg-open"
		args = []string{url}
	}

	if err := exec.Command(cmd, args...).Start(); err != nil {
		log.Printf("Failed to open browser: %v", err)
	}
}

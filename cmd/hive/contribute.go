package hive

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/openziti/zrok/environment"
	"github.com/openziti/zrok/sdk/golang/sdk"
	"github.com/spf13/cobra"
)

var hiveContributeCmd = &cobra.Command{
	Use:   "contribute --service <name> --port <num> --hive <token>",
	Short: "Contribute a local service to the Hive",
	Run: func(cmd *cobra.Command, args []string) {
		service, _ := cmd.Flags().GetString("service")
		port, _ := cmd.Flags().GetString("port")
		hiveToken, _ := cmd.Flags().GetString("hive")

		if service == "" || port == "" || hiveToken == "" {
			log.Fatal("must provide --service, --port, and --hive")
		}

		root, err := environment.LoadRoot()
		if err != nil {
			log.Fatal(err)
		}

		share, err := sdk.CreateShare(root, &sdk.ShareRequest{
			BackendMode: sdk.TcpTunnelBackendMode,
			ShareMode:   sdk.PrivateShareMode,
			Target:      "hive-" + service,
		})
		if err != nil {
			log.Fatal(err)
		}
		url := fmt.Sprintf("http://localhost:8081/hives/contribute?hive=%s&service=%s&port=%s&token=%s",
			hiveToken, service, port, share.Token)
		_, err = http.Post(url, "text/plain", nil)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Contributed service '%s' on port %s\nShare token: %s", service, port, share.Token)
		log.Printf("Tell teammates: devlink hive connect --hive %s", hiveToken)

		listener, err := sdk.NewListener(share.Token, root)
		if err != nil {
			log.Fatal(err)
		}
		defer listener.Close()

		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("error accepting: %v", err)
				continue
			}

			go func(remote net.Conn) {
				local, err := net.Dial("tcp", "127.0.0.1:"+port)
				if err != nil {
					log.Printf("error connecting local service: %v", err)
					_ = remote.Close()
					return
				}
				Pipe(remote, local)
			}(conn)
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

func init() {
	hiveContributeCmd.Flags().String("service", "", "service name (e.g. api, frontend)")
	hiveContributeCmd.Flags().String("port", "", "local port to share")
	hiveContributeCmd.Flags().String("hive", "", "hive invite token")
}

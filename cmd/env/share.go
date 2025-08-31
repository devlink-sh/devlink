package env

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/openziti/zrok/environment"
	"github.com/openziti/zrok/sdk/golang/sdk"
)

var envShareCmd = &cobra.Command{
	Use:   "share",
	Short: "Share environment variables",
	Long:  `Share environment variables with other commands`,
	Run: func(cmd *cobra.Command, args []string) {
		root, err := environment.LoadRoot()
		if err != nil {
			log.Fatal(err)
		}

		envFile, err := os.ReadFile(".env")
		if err != nil {
			log.Fatal(err)
		}

		share, err :=  sdk.CreateShare(root, &sdk.ShareRequest{
			BackendMode: sdk.TcpTunnelBackendMode,
			ShareMode:  sdk.PrivateShareMode,
			Target:	"env",
		})

		if err != nil {
			log.Fatal(err)
		}

		log.Printf("access you env using 'devlink env get %s'", share.Token)

		listener, err := sdk.NewListener(share.Token, root)
		if err != nil {
			log.Fatal(err)
		}

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func(){
			<- c
			if err := sdk.DeleteShare(root, share); err != nil {
				log.Printf("error deleting share: %v", err)
			}
			_ = listener.Close()
			os.Exit(0)
		}()

		go func(){
			for{
				conn, err := listener.Accept()
				if err != nil {
					log.Printf("error accepting connection: %v", err)
					return
				}
				go func(c net.Conn){
					defer c.Close()
					_, _ = c.Write(envFile)
				}(conn)
			}
		}()

		select {}


	},
}

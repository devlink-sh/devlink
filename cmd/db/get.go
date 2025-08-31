package db

import (
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

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
		defer func() {
			if err := sdk.DeleteAccess(root, acc); err != nil {
				log.Printf("error deleting access: %v", err)
			}
		}()

		// Listen locally
		listener, err := net.Listen("tcp", "127.0.0.1:"+port)
		if err != nil {
			log.Fatal(err)
		}
		defer listener.Close()
		log.Printf("DB tunnel ready at 127.0.0.1:%s", port)

		// Handle SIGINT/SIGTERM cleanly
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			log.Println("Shutting down db get...")
			_ = listener.Close()
			os.Exit(0)
		}()

		for {
			client, err := listener.Accept()
			if err != nil {
				log.Printf("error accepting local client: %v", err)
				continue
			}

			go func(c net.Conn) {
				remote, err := sdk.NewDialer(token, root)
				if err != nil {
					log.Printf("error dialing zrok: %v", err)
					c.Close()
					return
				}
				log.Printf("Client connected, tunneling traffic...")
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
		_, _ = io.Copy(a, b)
		done <- struct{}{}
	}()
	go func() {
		_, _ = io.Copy(b, a)
		done <- struct{}{}
	}()

	<-done // wait for one side to finish
}

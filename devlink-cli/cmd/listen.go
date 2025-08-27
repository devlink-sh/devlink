package cmd

import (
	"fmt"
	"net"

	"github.com/devlink/internal/proxy"
	"github.com/devlink/internal/ziti"
	"github.com/spf13/cobra"
)

var listenCmd = &cobra.Command{
	Use:   "listen [<service-name>]",
	Short: "Listens for incoming database connections",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Implementation for listening for incoming database connections
		serviceName := args[0]

		localPort,_ := cmd.Flags().GetInt("local-port")

		appCtx, ok := ziti.AppContextFrom(cmd)
		if !ok {
			return fmt.Errorf("error getting ziti context")
		}
		
		localAddress := fmt.Sprintf("localhost:%d", localPort)
		localListener, err := net.Listen("tcp", localAddress)
		if err != nil {
			return fmt.Errorf("error starting local listener: %w", err)
		}
		defer localListener.Close()

		fmt.Printf("Listening on %s. Connect your client here.\n", localAddress)
		fmt.Println("Waiting for local connection...")


		localConn, err := localListener.Accept()
		if err != nil {
			return fmt.Errorf("error accepting local connection: %w", err)
		}
		defer localConn.Close()

		fmt.Println("Local Client connected...")

		zitiConn, err := appCtx.ZitiContext.Dial(serviceName)
		if err != nil {
			return fmt.Errorf("error dialing ziti connection: %w", err)
		}
		defer zitiConn.Close()

		proxy.Pipe(localConn, zitiConn)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listenCmd)
	listenCmd.Flags().Int("local-port", 15432, "The local port to expose the service on")
}
package directory

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/openziti/zrok/environment"
	"github.com/openziti/zrok/sdk/golang/sdk"
	"github.com/spf13/cobra"
)

var directoryShareCmd = &cobra.Command{
	Use:   "share <directory>",
	Short: "Share a directory with teammates (browse in browser)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dir := args[0]

		if _, err := os.Stat(dir); os.IsNotExist(err) {
			log.Fatalf("Directory not found: %s", dir)
		}

		root, err := environment.LoadRoot()
		if err != nil {
			log.Fatal(err)
		}

		share, err := sdk.CreateShare(root, &sdk.ShareRequest{
			BackendMode: sdk.TcpTunnelBackendMode,
			ShareMode:   sdk.PrivateShareMode,
			Target:      "dir",
		})
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Directory share ready! Teammates can run:\n  devlink dir get %s 8080\n", share.Token)

		listener, err := sdk.NewListener(share.Token, root)
		if err != nil {
			log.Fatal(err)
		}
		defer listener.Close()

		server := &http.Server{
			Handler: http.FileServer(http.Dir(dir)),
		}

		// Cleanup on exit
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			log.Println("Shutting down file share...")
			_ = sdk.DeleteShare(root, share)
			_ = server.Close()
			_ = listener.Close()
			os.Exit(0)
		}()

		log.Printf("Serving '%s' over Devlink...", dir)
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error serving: %v", err)
		}
	},
}

package registry

import (
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/openziti/zrok/environment"
	"github.com/openziti/zrok/sdk/golang/sdk"
	"github.com/spf13/cobra"
)

var registryShareCmd = &cobra.Command{
	Use:   "share <image>",
	Short: "Share a local docker image",
	Long:  `Stream a local docker image (docker save) to peers via zrok. Example: devlink registry share myapp:latest`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		image := args[0]

		root, err := environment.LoadRoot()
		if err != nil {
			log.Fatal(err)
		}

		share, err := sdk.CreateShare(root, &sdk.ShareRequest{
			BackendMode: sdk.TcpTunnelBackendMode,
			ShareMode:   sdk.PrivateShareMode,
			Target:      "registry",
		})
		if err != nil {
			log.Fatalf("unable to create share: %v", err)
		}

		log.Printf("Registry share ready! Let others pull using:\n  devlink registry get %s\nSharing image: %s", share.Token, image)

		listener, err := sdk.NewListener(share.Token, root)
		if err != nil {
			log.Fatalf("unable to create listener: %v", err)
		}
		defer listener.Close()

		// handle signals to cleanup share
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-sig
			log.Println("shutting down registry share...")
			if err := sdk.DeleteShare(root, share); err != nil {
				log.Printf("error deleting share: %v", err)
			}
			_ = listener.Close()
			os.Exit(0)
		}()

		// Accept loop: for every incoming connection, run `docker save` and stream to conn
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("error accepting zrok connection: %v", err)
				continue
			}

			go func(c io.ReadWriteCloser) {
				defer c.Close()

				log.Printf("client connected, streaming image %s ...", image)

				// docker save <image> -> stdout
				saveCmd := exec.Command("docker", "save", image)
				stdout, err := saveCmd.StdoutPipe()
				if err != nil {
					log.Printf("error getting docker save stdout: %v", err)
					return
				}
				if err := saveCmd.Start(); err != nil {
					log.Printf("error starting docker save: %v", err)
					return
				}

				// stream docker save stdout -> connection
				_, err = io.Copy(c, stdout)
				if err != nil {
					log.Printf("streaming error: %v", err)
				}

				// wait for docker save to finish
				if err := saveCmd.Wait(); err != nil {
					log.Printf("docker save exited with error: %v", err)
				}

				log.Printf("finished streaming image %s to client", image)
			}(conn)
		}
	},
}

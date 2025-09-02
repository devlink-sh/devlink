// serve.go
package git

import (
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/devlink/internal"
	"github.com/openziti/zrok/environment"
	"github.com/openziti/zrok/sdk/golang/sdk"
	"github.com/spf13/cobra"
)

var gitServeCmd = &cobra.Command{
	Use:   "serve <repo-path>",
	Short: "Share a local Git repository",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repoPath, err := filepath.Abs(args[0])
		if err != nil {
			log.Fatalf("failed to resolve repo path: %v", err)
		}

		parentDir := filepath.Dir(repoPath)

		// Start git daemon locally
		gitDaemon := exec.Command("git", "daemon",
			"--reuseaddr",
			"--base-path="+parentDir, // parent dir as base path
			"--export-all",
			"--verbose",
			"--enable=receive-pack", // allow push

		)
		gitDaemon.Stdout = os.Stdout
		gitDaemon.Stderr = os.Stderr

		if err := gitDaemon.Start(); err != nil {
			log.Fatalf("failed to start git daemon: %v", err)
		}
		log.Printf("Git daemon started for %s (listening on 127.0.0.1:9418)", repoPath)

		// Setup zrok share
		root, err := environment.LoadRoot()
		if err != nil {
			log.Fatal(err)
		}

		share, err := sdk.CreateShare(root, &sdk.ShareRequest{
			BackendMode: sdk.TcpTunnelBackendMode,
			ShareMode:   sdk.PrivateShareMode,
			Target:      "git",
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Git share ready! Teammates can connect via:\n  devlink git connect %s 9418", share.Token)

		// Accept zrok connections
		listener, err := sdk.NewListener(share.Token, root)
		if err != nil {
			log.Fatal(err)
		}
		defer listener.Close()

		// Handle Ctrl+C
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			log.Println("Shutting down git serve...")
			_ = sdk.DeleteShare(root, share)
			_ = listener.Close()
			_ = gitDaemon.Process.Kill()
			os.Exit(0)
		}()

		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("error accepting zrok connection: %v", err)
				continue
			}

			go func(remote net.Conn) {
				local, err := net.Dial("tcp", "127.0.0.1:9418") // git:// default port
				if err != nil {
					log.Printf("error dialing local git daemon: %v", err)
					remote.Close()
					return
				}
				log.Println("Forwarding git traffic...")
				internal.Pipe(remote, local)
			}(conn)
		}
	},
}

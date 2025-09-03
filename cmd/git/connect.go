package git

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/devlink/internal"
	"github.com/openziti/zrok/environment"
	"github.com/openziti/zrok/sdk/golang/sdk"
	"github.com/spf13/cobra"
)

func freeLocalListener() (net.Listener, int, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, 0, err
	}
	port := l.Addr().(*net.TCPAddr).Port
	return l, port, nil
}

var gitConnectCmd = &cobra.Command{
	Use:   "connect <token> <repo-name.git>",
	Short: "Connect to a shared Git repository",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		token := args[0]
		repoName := args[1]

		// Canonicalize repo name exactly once
		if !strings.HasSuffix(strings.ToLower(repoName), ".git") {
			repoName += ".git"
		}

		// Load zrok environment
		root, err := environment.LoadRoot()
		if err != nil {
			log.Fatal(err)
		}

		// Access session (directional: this side is the client)
		session, err := sdk.CreateAccess(root, &sdk.AccessRequest{ShareToken: token})
		if err != nil {
			log.Fatal(err)
		}
		defer sdk.DeleteAccess(root, session)

		// Bind a free local port for the git client to talk to
		listener, localPort, err := freeLocalListener()
		if err != nil {
			log.Fatalf("failed to acquire local port: %v", err)
		}
		defer listener.Close()

		// Print clone instructions with exact canonical values
		cloneURL := fmt.Sprintf("git://127.0.0.1:%d/%s", localPort, repoName)
		workDir := strings.TrimSuffix(repoName, ".git")
		clonePath := filepath.Join(".", workDir)

		log.Printf("Git tunnel ready!")
		log.Printf("Clone using:\n\n  git clone %s %s\n", cloneURL, clonePath)
		log.Printf("Keep this process running to use git fetch/pull/push.")

		// Ctrl+C shutdown
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			log.Println("Shutting down git connect...")
			_ = listener.Close()
			_ = sdk.DeleteAccess(root, session)
			os.Exit(0)
		}()

		// Accept loop with backoff on temporary errors
		var tempDelay time.Duration
		for {
			client, err := listener.Accept()
			if err != nil {
				if ne, ok := err.(interface{ Temporary() bool }); ok && ne.Temporary() {
					if tempDelay == 0 {
						tempDelay = 50 * time.Millisecond
					} else {
						tempDelay *= 2
						if tempDelay > time.Second {
							tempDelay = time.Second
						}
					}
					log.Printf("temporary accept error: %v; retrying in %v", err, tempDelay)
					time.Sleep(tempDelay)
					continue
				}
				log.Printf("fatal accept error: %v", err)
				break
			}
			tempDelay = 0

			// Dial remote through zrok
			go func(c net.Conn) {
				remote, err := sdk.NewDialer(token, root)
				if err != nil {
					log.Printf("error creating zrok dialer: %v", err)
					_ = c.Close()
					return
				}
				log.Printf("Forwarding git traffic (client<->remote)...")
				internal.Pipe(c, remote)
			}(client)
		}
	},
}

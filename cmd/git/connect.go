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

	"github.com/devlink/internal"
	"github.com/openziti/zrok/environment"
	"github.com/openziti/zrok/sdk/golang/sdk"
	"github.com/spf13/cobra"
)

var gitConnectCmd = &cobra.Command{
	Use:   "connect <token> <repo-name.git>",
	Short: "Connect to a shared Git repository",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		token := args[0]
		repoName := args[1]

		// Ensure .git suffix
		if !strings.HasSuffix(repoName, ".git") {
			repoName += ".git"
		}

		// Load zrok environment
		root, err := environment.LoadRoot()
		if err != nil {
			log.Fatal(err)
		}

		// Create an access session
		session, err := sdk.CreateAccess(root, &sdk.AccessRequest{ShareToken: token})
		if err != nil {
			log.Fatal(err)
		}
		defer sdk.DeleteAccess(root, session)

		// Start a local listener on 9418
		listener, err := net.Listen("tcp", "127.0.0.1:9418")
		if err != nil {
			log.Fatal(err)
		}
		defer listener.Close()

		// Print clone instructions
		cloneURL := fmt.Sprintf("git://127.0.0.1:9418/%s", repoName)
		workDir := strings.TrimSuffix(repoName, ".git")
		clonePath := filepath.Join(".", workDir)

		log.Printf("Git tunnel ready!")
		log.Printf("You can now clone the repo with:\n\n  git clone %s %s\n", cloneURL, clonePath)
		log.Printf("Keep this process running to use git pull/push.")

		// Handle Ctrl+C shutdown
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			log.Println("Shutting down git connect...")
			_ = listener.Close()
			_ = sdk.DeleteAccess(root, session)
			os.Exit(0)
		}()

		// Accept connections from local git client
		for {
			client, err := listener.Accept()
			if err != nil {
				log.Printf("error accepting client: %v", err)
				continue
			}

			// Dial remote through zrok
			go func(c net.Conn) {
				remote, err := sdk.NewDialer(token, root)
				if err != nil {
					log.Printf("error dialing zrok: %v", err)
					c.Close()
					return
				}
				log.Printf("Forwarding git traffic...")
				internal.Pipe(c, remote)
			}(client)
		}
	},
}

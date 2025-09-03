package git

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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

		// Ensure repoName ends with .git
		if !strings.HasSuffix(repoName, ".git") {
			repoName += ".git"
		}

		// Load zrok environment
		root, err := environment.LoadRoot()
		if err != nil {
			log.Fatal(err)
		}

		// Create an access session
		session, err := sdk.CreateAccess(root, &sdk.AccessRequest{
			Token: token,
		})
		if err != nil {
			log.Fatal(err)
		}
		defer sdk.DeleteAccess(root, session)

		// Repo will be cloned into current directory
		clonePath := filepath.Join(".", repoName)

		// Git clone command via tunnel
		cloneURL := fmt.Sprintf("git://127.0.0.1:9418/%s", repoName)
		cloneCmd := exec.Command("git", "clone", cloneURL, clonePath)
		cloneCmd.Stdout = os.Stdout
		cloneCmd.Stderr = os.Stderr

		log.Printf("Cloning %s into %s ...", cloneURL, clonePath)

		if err := cloneCmd.Run(); err != nil {
			log.Fatalf("git clone failed: %v", err)
		}

		log.Printf("Repository cloned successfully into %s", clonePath)
	},
}

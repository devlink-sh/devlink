package registry

import (
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/openziti/zrok/environment"
	"github.com/openziti/zrok/sdk/golang/sdk"
	"github.com/spf13/cobra"
)

var registryGetCmd = &cobra.Command{
	Use:   "get <token>",
	Short: "Pull a docker image from a remote registry share",
	Long:  `Connect to a registry share and load the streamed image into local Docker (runs docker load). Example: devlink registry get <token>`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		token := args[0]

		root, err := environment.LoadRoot()
		if err != nil {
			log.Fatal(err)
		}

		// create access (optional but good practice)
		acc, err := sdk.CreateAccess(root, &sdk.AccessRequest{ShareToken: token})
		if err != nil {
			log.Fatalf("unable to create access: %v", err)
		}
		// make sure to delete access when done
		defer func() {
			if err := sdk.DeleteAccess(root, acc); err != nil {
				log.Printf("error deleting access: %v", err)
			}
		}()

		// Connect to the share: returns a connected net.Conn (stream)
		conn, err := sdk.NewDialer(token, root)
		if err != nil {
			log.Fatalf("unable to dial share: %v", err)
		}
		defer conn.Close()

		log.Println("Connected — streaming image into local docker (docker load)...")

		// Start `docker load` and pipe the incoming stream into its stdin
		loadCmd := exec.Command("docker", "load")
		stdin, err := loadCmd.StdinPipe()
		if err != nil {
			log.Fatalf("error creating docker load stdin pipe: %v", err)
		}
		loadCmd.Stdout = os.Stdout
		loadCmd.Stderr = os.Stderr

		if err := loadCmd.Start(); err != nil {
			log.Fatalf("error starting docker load: %v", err)
		}

		// conn -> docker load stdin
		_, err = io.Copy(stdin, conn)
		if err != nil {
			log.Fatalf("error streaming image into docker load: %v", err)
		}

		// close stdin to signal docker load EOF
		_ = stdin.Close()

		if err := loadCmd.Wait(); err != nil {
			log.Fatalf("docker load failed: %v", err)
		}

		log.Println("Image loaded successfully into local Docker.")
	},
}

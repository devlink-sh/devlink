package hive

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/openziti/zrok/environment"
	"github.com/openziti/zrok/environment/env_core"
	"github.com/openziti/zrok/sdk/golang/sdk"
	"github.com/spf13/cobra"
)

// Service struct must match Hive Controller response
type Service struct {
	Name  string `json:"name"`
	Port  string `json:"port"`
	Token string `json:"token"`
}

var hiveConnectCmd = &cobra.Command{
	Use:   "connect --hive <token>",
	Short: "Connect to all services in a Hive",
	Run: func(cmd *cobra.Command, args []string) {
		hiveToken, _ := cmd.Flags().GetString("hive")
		if hiveToken == "" {
			log.Fatal("must provide --hive token")
		}

		// Load zrok root identity
		root, err := environment.LoadRoot()
		if err != nil {
			log.Fatal(err)
		}

		// Fetch services from Hive Controller
		url := fmt.Sprintf("http://localhost:8081/hives/services?hive=%s", hiveToken)
		resp, err := http.Get(url)
		if err != nil {
			log.Fatalf("error querying hive controller: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			log.Fatalf("controller error: %s", string(body))
		}

		var services map[string]Service
		if err := json.NewDecoder(resp.Body).Decode(&services); err != nil {
			log.Fatalf("decode error: %v", err)
		}

		if len(services) == 0 {
			log.Println("No services available in this Hive yet.")
			return
		}

		// Start listeners for each service
		for _, svc := range services {
			go startLocalListener(svc, root)
		}

		// Keep running until Ctrl+C
		log.Println("Connected to Hive! Press Ctrl+C to exit.")
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
	},
}

func startLocalListener(svc Service, root env_core.Root) {
	// Create Access for this service
	acc, err := sdk.CreateAccess(root, &sdk.AccessRequest{ShareToken: svc.Token})
	if err != nil {
		log.Printf("[%s] create access error: %v", svc.Name, err)
		return
	}
	defer sdk.DeleteAccess(root, acc)

	listener, err := net.Listen("tcp", "127.0.0.1:"+svc.Port)
	if err != nil {
		log.Printf("[%s] listener error: %v", svc.Name, err)
		return
	}
	defer listener.Close()

	log.Printf("Service '%s' ready at http://127.0.0.1:%s", svc.Name, svc.Port)

	for {
		client, err := listener.Accept()
		if err != nil {
			log.Printf("[%s] accept error: %v", svc.Name, err)
			continue
		}

		go func(c net.Conn) {
			defer c.Close()

			remote, err := sdk.NewDialer(svc.Token, root)
			if err != nil {
				log.Printf("[%s] dial error: %v", svc.Name, err)
				return
			}
			defer remote.Close()

			log.Printf("[%s] client connected", svc.Name)
			go io.Copy(remote, c)
			io.Copy(c, remote)
		}(client)
	}
}

func init() {
	hiveConnectCmd.Flags().String("hive", "", "hive invite token")
}

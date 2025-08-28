package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	// "time"

	"github.com/devlink/internal/proxy"
	"github.com/devlink/internal/ziti"
	"github.com/openziti/sdk-golang/ziti/edge"

	// "github.com/devlink/util"
	"github.com/spf13/cobra"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Share and manage database connections",
}

var dbShareCmd = &cobra.Command{
	Use:   "share",
	Short: "Share a database connection",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Implementation for sharing a database connection
		appCtx, ok := ziti.AppContextFrom(cmd)
		if !ok {
			return fmt.Errorf("error getting ziti context")
		}

		if appCtx.ZitiContext != nil {
			id := appCtx.ZitiContext.GetCredentials()
			if id != nil {
				log.Printf("🔑 Using Ziti identity: %s (fingerprint: %s)", id.Payload().Username, id.Payload().Password)
			} else {
				log.Printf("⚠️ No identity loaded in ZitiContext")
			}
		}

		dbType, _ := cmd.Flags().GetString("type")
		dbPort, _ := cmd.Flags().GetInt("port")

		serviceName := "devlink-service"
		if err := createService(serviceName); err != nil {
			return fmt.Errorf("failed to create service: %w", err)
		}
		log.Println("Service created successfully:", serviceName)

		// The controller may take a short moment to propagate the new service to the data-plane.
		// Retry Listen with backoff for a few seconds before giving up.
		var listener edge.Listener
		var err error
		for i := 0; i < 5; i++ {
			listener, err = appCtx.ZitiContext.Listen("devlink-service")
			if err == nil {
				break
			}
			log.Printf("⏳ Retrying listen... attempt %d, err: %v", i+1, err)
			time.Sleep(2 * time.Second)
		}
		if err != nil {
			return fmt.Errorf("error creating listener after retries: %w", err)
		}
		defer listener.Close()

		printConnectionString(serviceName, dbType, dbPort)
		fmt.Println("Waiting for receiver to connect...")

		zitiConn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("error accepting connection: %w", err)
		}
		fmt.Println("Receiver connected!")
		defer zitiConn.Close()

		localTarget := fmt.Sprintf("localhost:%d", dbPort)
		localDbConn, err := net.Dial("tcp", localTarget)
		if err != nil {
			return fmt.Errorf("failed to dial local database at '%s': %w", localTarget, err)
		}

		proxy.Pipe(zitiConn, localDbConn)
		fmt.Println("Connection closed.")

		return nil
	},
}

func createService(serviceName string) error {
	backendUrl := "http://localhost:8080/api/v1/create-service"

	request := map[string]string{
		"name": serviceName,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(backendUrl, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create service, status: %s", resp.Status)
	}

	return nil

}

func printConnectionString(serviceName, dbType string, remotePort int) {
	receiverLocalPort := remotePort + 10000

	fmt.Printf("\nShare code: %s\n\n", serviceName)
	fmt.Println("On the receiver's machine, run the following command...")
	fmt.Printf("devlink listen %s --local-port %d\n\n", serviceName, receiverLocalPort)

	fmt.Println("Then, use the following connection string in the database client:")

	var connString string
	switch dbType {
	case "postgres":
		connString = fmt.Sprintf("postgres://user:password@localhost:%d/mydatabase?sslmode=disable", receiverLocalPort)
	case "mysql":
		connString = fmt.Sprintf("mysql://user:password@tcp(localhost:%d)/mydatabase", receiverLocalPort)
	default:
		connString = fmt.Sprintf("Connect to localhost:%d", receiverLocalPort)
	}
	fmt.Printf("%s\n\n", connString)
}

func init() {
	rootCmd.AddCommand(dbCmd)
	dbCmd.AddCommand(dbShareCmd)
	dbShareCmd.Flags().String("type", "postgres", "Database type (postgres, mysql)")
	dbShareCmd.Flags().IntP("port", "p", 5432, "Database port")
	dbShareCmd.MarkFlagRequired("port")
}

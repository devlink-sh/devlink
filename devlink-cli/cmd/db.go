package cmd

import (
	"fmt"
	"net"

	"github.com/devlink/internal/proxy"
	"github.com/devlink/internal/ziti"
	"github.com/devlink/util"
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

		dbType, _ := cmd.Flags().GetString("type")
		dbPort, _ := cmd.Flags().GetInt("port")

		serviceName := util.GenerateHumanCode()
		listener, err := appCtx.ZitiContext.Listen(serviceName)
		if err != nil {
			return fmt.Errorf("error creating listener: %w", err)
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


func init(){
	rootCmd.AddCommand(dbCmd)
	dbCmd.AddCommand(dbShareCmd)
	dbShareCmd.Flags().String("type", "postgres", "Database type (postgres, mysql)")
	dbShareCmd.Flags().IntP("port", "p", 5432, "Database port")
	dbShareCmd.MarkFlagRequired("port")
}
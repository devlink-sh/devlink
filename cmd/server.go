package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devlink/pkg/envsharing/core/encryption"
	"github.com/devlink/pkg/envsharing/network"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "🚀 Start the secure sharing service",
	Long: `🚀 Start the secure sharing service

This starts the magic tunnel that lets you share .env files securely!
Think of it as opening a secret door that only your team can access.

The service runs in the background and waits for sharing requests.
Keep it running while you want to share files with your team! 🔐

Examples:
  devlink server                    # Start the service
  devlink server --service my-team  # Start with custom service name
  devlink server --verbose          # See what's happening behind the scenes`,
	RunE: runServer,
}

func runServer(cmd *cobra.Command, args []string) error {
	serviceName, _ := cmd.Flags().GetString("service")
	verbose, _ := cmd.Flags().GetBool("verbose")

	zitiConfig := &network.ZitiConfig{
		ControllerURL: getEnvOrDefault("ZITI_CONTROLLER_URL", "https://localhost:1280"),
		IdentityFile:  getEnvOrDefault("ZITI_IDENTITY_FILE", getHomeDir()+"/.devlink/identity.json"),
		ServiceName:   getEnvOrDefault("ZITI_SERVICE_NAME", "devlink-env-service"),
	}

	if serviceName != "" {
		zitiConfig.ServiceName = serviceName
	}

	if verbose {
		fmt.Printf("🔧 Server Configuration:\n")
		fmt.Printf("   • Controller: %s\n", zitiConfig.ControllerURL)
		fmt.Printf("   • Service: %s\n", zitiConfig.ServiceName)
		fmt.Printf("   • Identity: %s\n", zitiConfig.IdentityFile)
		fmt.Println()
	}

	encryptionManager := encryption.NewManager()
	service, err := network.NewZitiService(zitiConfig, encryptionManager)
	if err != nil {
		return fmt.Errorf("failed to create ziti service: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\n🛑 Shutting down the service...")
		cancel()
	}()

	serviceErr := make(chan error, 1)
	go func() {
		serviceErr <- service.Start()
	}()

	select {
	case err := <-serviceErr:
		if err != nil {
			return fmt.Errorf("ziti service failed to start: %w", err)
		}
	case <-time.After(2 * time.Second):
		fmt.Printf("✅ DevLink service is running!\n")
		fmt.Printf("🔐 Zero-trust tunnel is active\n")
		fmt.Printf("🌐 Service: %s\n", zitiConfig.ServiceName)
		fmt.Printf("📡 Controller: %s\n", zitiConfig.ControllerURL)
		fmt.Println()
		fmt.Println("💡 Keep this running and use Ctrl+C to stop when done")
	}

	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := service.Stop(shutdownCtx); err != nil {
		return fmt.Errorf("failed to stop ziti service gracefully: %w", err)
	}

	fmt.Println("✅ Service stopped gracefully")
	return nil
}

func init() {
	serverCmd.Flags().StringP("service", "s", "", "Ziti service name")
	serverCmd.Flags().BoolP("verbose", "v", false, "Enable verbose output")

	rootCmd.AddCommand(serverCmd)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return homeDir
}

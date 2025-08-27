package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/openziti/sdk-golang/ziti/enroll"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [invitation-token]",
	Short: "Initialize your devlink environment for first-time use",
	Long:  `Initializes the devlink environment by securely enrolling a new identity. You must provide a one-time invitation token from your administrator.`,
	Args:  cobra.ExactArgs(1), // Ensures the user provides exactly one argument
	RunE: func(cmd *cobra.Command, args []string) error {
		invitationToken := args[0]

		fmt.Println("🔄 Contacting provisioning service...")
		enrollmentJwt, err := requestEnrollmentToken(invitationToken)
		if err != nil {
			return fmt.Errorf("failed to get enrollment token from backend: %w", err)
		}
		fmt.Println("✅ Received enrollment token.")
		fmt.Printf("🔍 Token received: %s\n", string(enrollmentJwt))

		// This generates the private key locally and gets the certificate from the controller.
		token, _, err := enroll.ParseToken(string(enrollmentJwt))
		if err != nil {
			return fmt.Errorf("failed to parse enrollment token: %w", err)
		}

		flags := enroll.EnrollmentFlags{
			Token:  token,
			KeyAlg: "RSA", // Using RSA algorithm for the key
		}
		config, err := enroll.Enroll(flags)
		if err != nil {
			return fmt.Errorf("enrollment failed: %w", err)
		}

		// Convert the config to JSON bytes
		configBytes, err := json.Marshal(config)
		if err != nil {
			return fmt.Errorf("failed to marshal config to JSON: %w", err)
		}

		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		identityPath := filepath.Join(home, ".devlink", "identity.json")

		// Create the.devlink directory if it doesn't exist
		if err := os.MkdirAll(filepath.Dir(identityPath), 0700); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}

		// Write the identity file
		if err := os.WriteFile(identityPath, configBytes, 0600); err != nil {
			return fmt.Errorf("failed to write identity file: %w", err)
		}

		fmt.Printf("✅ DevLink environment initialized successfully!\nIdentity saved to: %s\n", identityPath)
		return nil
	},
}

// requestEnrollmentToken is a helper function to call your backend API.
func requestEnrollmentToken(invitationToken string) ([]byte, error) {
	backendURL := "http://localhost:8080/api/v1/provision"

	reqBody, _ := json.Marshal(map[string]string{"token": invitationToken})
	resp, err := http.Post(backendURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("backend returned a non-200 status: %s - %s", resp.Status, string(body))
	}

	return io.ReadAll(resp.Body)
}

func init() {
	rootCmd.AddCommand(initCmd)
}

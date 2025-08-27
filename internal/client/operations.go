package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/devlink/internal/env"
	"github.com/devlink/internal/util"
)

// CreateShare creates a new share on the server
func (c *Client) CreateShare(envFile *util.EnvFile, shareCode string, expiry time.Duration, readOnly bool) (*ShareResponse, error) {
	requestData := map[string]interface{}{
		"share_code": shareCode,
		"data":       envFile.RawContent,
		"expiry":     expiry.String(),
		"read_only":  readOnly,
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize request: %w", err)
	}

	req, err := http.NewRequest("POST", c.serverURL+"/share", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server error: %s - %s", resp.Status, string(responseBody))
	}

	var shareResponse ShareResponse
	if err := json.Unmarshal(responseBody, &shareResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &shareResponse, nil
}

// GetShare retrieves a share from the server
func (c *Client) GetShare(shareCode string) (*util.EnvFile, error) {
	req, err := http.NewRequest("GET", c.serverURL+"/get/"+shareCode, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server error: %s - %s", resp.Status, string(responseBody))
	}

	var getResponse GetResponse
	if err := json.Unmarshal(responseBody, &getResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	envFile, err := c.encryption.DecryptEnvFile(getResponse.EncryptedData, shareCode)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt env file: %w", err)
	}

	return envFile, nil
}

// GetShareWithMetadata retrieves a share with metadata from the server
func (c *Client) GetShareWithMetadata(shareCode string) (*GetResponse, error) {
	req, err := http.NewRequest("GET", c.serverURL+"/get/"+shareCode, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server error: %s - %s", resp.Status, string(responseBody))
	}

	var getResponse GetResponse
	if err := json.Unmarshal(responseBody, &getResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &getResponse, nil
}

// ShareFile shares an environment file from disk
func (c *Client) ShareFile(filePath string, shareCode string, expiry time.Duration, readOnly bool) (*ShareResponse, error) {
	parser := env.NewParser()
	envFile, err := parser.ParseFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse env file: %w", err)
	}

	return c.CreateShare(envFile, shareCode, expiry, readOnly)
}

// GetShareToFile retrieves a share and saves it to a file
func (c *Client) GetShareToFile(shareCode string, outputPath string) error {
	envFile, err := c.GetShare(shareCode)
	if err != nil {
		return fmt.Errorf("failed to get share: %w", err)
	}

	formatter := env.NewFormatter()
	formatOptions := &util.FormatOptions{
		MaskSensitive:   false,
		ShowComments:    true,
		ShowLineNumbers: false,
		OutputFormat:    "text",
	}

	formattedContent, err := formatter.Format(envFile, formatOptions)
	if err != nil {
		return fmt.Errorf("failed to format env file: %w", err)
	}

	if err := os.WriteFile(outputPath, []byte(formattedContent), 0600); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

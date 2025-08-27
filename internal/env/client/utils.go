package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// HealthCheck checks if the server is healthy
func (c *Client) HealthCheck() error {
	req, err := http.NewRequest("GET", c.serverURL+"/health", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server unhealthy: %s", resp.Status)
	}

	return nil
}

// GetServerStats retrieves server statistics
func (c *Client) GetServerStats() (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", c.serverURL+"/stats", nil)
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

	var stats map[string]interface{}
	if err := json.Unmarshal(responseBody, &stats); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return stats, nil
}

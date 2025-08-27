package client

import (
	"net/http"
	"time"

	"github.com/devlink/internal/util"
)

// NewClient creates a new HTTP client
func NewClient(config *util.Config, serverURL string) (*Client, error) {
	encryption := util.NewEncryptionManager(config)

	client := &Client{
		config:     config,
		serverURL:  serverURL,
		encryption: encryption,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	return client, nil
}

// ValidateShareCode validates a share code format
func (c *Client) ValidateShareCode(shareCode string) error {
	validationManager := util.NewValidationManager()
	return validationManager.ValidateShareCode(shareCode, "http-client")
}

// GetShareInfo retrieves share information without decrypting
func (c *Client) GetShareInfo(shareCode string) (*GetResponse, error) {
	response, err := c.GetShareWithMetadata(shareCode)
	if err != nil {
		return nil, err
	}
	return response, nil
}

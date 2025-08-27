package client

import (
	"net/http"
	"time"

	"github.com/devlink/internal/util"
)

// Client represents the HTTP client for environment variable sharing
type Client struct {
	config     *util.Config
	httpClient *http.Client
	encryption *util.EncryptionManager
	serverURL  string
}

// ShareResponse represents the response from the server when creating a share
type ShareResponse struct {
	Success   bool      `json:"success"`
	ShareID   string    `json:"share_id"`
	ShareCode string    `json:"share_code"`
	ExpiresAt time.Time `json:"expires_at"`
}

// GetResponse represents the response from the server when retrieving a share
type GetResponse struct {
	Success       bool                   `json:"success"`
	ShareCode     string                 `json:"share_code"`
	EncryptedData *util.EncryptedData    `json:"encrypted_data"`
	Metadata      map[string]interface{} `json:"metadata"`
}

package network

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/devlink/pkg/envsharing/core"
	"github.com/devlink/pkg/envsharing/core/encryption"
	"github.com/openziti/sdk-golang/ziti"
)

type ZitiClient struct {
	config     *ZitiConfig
	zitiCtx    ziti.Context
	encryption *encryption.Manager
}

type ShareResponse struct {
	Success   bool      `json:"success"`
	ShareID   string    `json:"share_id"`
	ShareCode string    `json:"share_code"`
	ExpiresAt time.Time `json:"expires_at"`
}

type GetResponse struct {
	Success       bool                   `json:"success"`
	ShareCode     string                 `json:"share_code"`
	EncryptedData *core.EncryptedData    `json:"encrypted_data"`
	Metadata      map[string]interface{} `json:"metadata"`
}

func NewZitiClient(zitiConfig *ZitiConfig, encryptionManager *encryption.Manager) (*ZitiClient, error) {
	zitiCtx, err := ziti.NewContextFromFile(zitiConfig.IdentityFile)
	if err != nil {
		return nil, fmt.Errorf("failed to create ziti context: %w", err)
	}

	return &ZitiClient{
		config:     zitiConfig,
		zitiCtx:    zitiCtx,
		encryption: encryptionManager,
	}, nil
}

func (c *ZitiClient) CreateShare(envFile *core.EnvFile, shareCode string, expiry time.Duration, readOnly bool) (*ShareResponse, error) {
	conn, err := c.zitiCtx.Dial(c.config.ServiceName)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ziti service: %w", err)
	}
	defer conn.Close()

	request := map[string]interface{}{
		"action":     "share",
		"share_code": shareCode,
		"data":       envFile.RawContent,
		"expiry":     expiry.String(),
		"read_only":  readOnly,
	}

	if err := c.sendRequest(conn, request); err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	response, err := c.readResponse(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if !response["success"].(bool) {
		return nil, fmt.Errorf("server error: %s", response["error"])
	}

	shareResponse := &ShareResponse{
		Success:   response["success"].(bool),
		ShareID:   response["share_id"].(string),
		ShareCode: response["share_code"].(string),
		ExpiresAt: time.Now().Add(expiry),
	}

	return shareResponse, nil
}

func (c *ZitiClient) GetShare(shareCode string) (*core.EnvFile, error) {
	conn, err := c.zitiCtx.Dial(c.config.ServiceName)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ziti service: %w", err)
	}
	defer conn.Close()

	request := map[string]interface{}{
		"action":     "get",
		"share_code": shareCode,
	}

	if err := c.sendRequest(conn, request); err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	response, err := c.readResponse(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if !response["success"].(bool) {
		return nil, fmt.Errorf("server error: %s", response["error"])
	}

	encryptedDataJSON, err := json.Marshal(response["encrypted_data"])
	if err != nil {
		return nil, fmt.Errorf("failed to marshal encrypted data: %w", err)
	}

	var encryptedData core.EncryptedData
	if err := json.Unmarshal(encryptedDataJSON, &encryptedData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal encrypted data: %w", err)
	}

	envFile, err := c.encryption.DecryptEnvFile(&encryptedData, shareCode)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt env file: %w", err)
	}

	return envFile, nil
}

func (c *ZitiClient) HealthCheck() error {
	conn, err := c.zitiCtx.Dial(c.config.ServiceName)
	if err != nil {
		return fmt.Errorf("failed to connect to ziti service: %w", err)
	}
	defer conn.Close()

	request := map[string]interface{}{
		"action": "health",
	}

	if err := c.sendRequest(conn, request); err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	response, err := c.readResponse(conn)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if !response["success"].(bool) {
		return fmt.Errorf("service unhealthy: %s", response["error"])
	}

	return nil
}

func (c *ZitiClient) sendRequest(conn net.Conn, request map[string]interface{}) error {
	requestData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	_, err = conn.Write(requestData)
	if err != nil {
		return fmt.Errorf("failed to write request: %w", err)
	}

	return nil
}

func (c *ZitiClient) readResponse(conn net.Conn) (map[string]interface{}, error) {
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))

	responseData, err := io.ReadAll(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(responseData, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response, nil
}

package network

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/devlink/pkg/envsharing/core"
	"github.com/devlink/pkg/envsharing/core/encryption"
	"github.com/google/uuid"
	"github.com/openziti/sdk-golang/ziti"
)

type ZitiService struct {
	config     *ZitiConfig
	zitiCtx    ziti.Context
	shares     map[string]*core.Share
	sharesMu   sync.RWMutex
	encryption *encryption.Manager
	listener   net.Listener
}

type ZitiConfig struct {
	ControllerURL string
	IdentityFile  string
	ServiceName   string
}

func NewZitiService(zitiConfig *ZitiConfig, encryptionManager *encryption.Manager) (*ZitiService, error) {
	zitiCtx, err := ziti.NewContextFromFile(zitiConfig.IdentityFile)
	if err != nil {
		return nil, fmt.Errorf("failed to create ziti context: %w", err)
	}

	return &ZitiService{
		config:     zitiConfig,
		zitiCtx:    zitiCtx,
		shares:     make(map[string]*core.Share),
		encryption: encryptionManager,
	}, nil
}

func (s *ZitiService) Start() error {
	listener, err := s.zitiCtx.Listen(s.config.ServiceName)
	if err != nil {
		return fmt.Errorf("failed to listen on ziti service: %w", err)
	}

	s.listener = listener
	go s.startCleanupRoutine()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("failed to accept connection: %w", err)
		}
		go s.handleConnection(conn)
	}
}

func (s *ZitiService) Stop(ctx context.Context) error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *ZitiService) handleConnection(conn net.Conn) {
	defer conn.Close()

	requestData, err := io.ReadAll(conn)
	if err != nil {
		s.sendError(conn, "failed to read request", err)
		return
	}

	var request map[string]interface{}
	if err := json.Unmarshal(requestData, &request); err != nil {
		s.sendError(conn, "invalid request format", err)
		return
	}

	action, ok := request["action"].(string)
	if !ok {
		s.sendError(conn, "missing action", fmt.Errorf("action field required"))
		return
	}

	switch action {
	case "share":
		s.handleShareRequest(conn, request)
	case "get":
		s.handleGetRequest(conn, request)
	case "health":
		s.handleHealthRequest(conn)
	default:
		s.sendError(conn, "unknown action", fmt.Errorf("unknown action: %s", action))
	}
}

func (s *ZitiService) handleShareRequest(conn net.Conn, request map[string]interface{}) {
	shareCode, _ := request["share_code"].(string)
	data, _ := request["data"].(string)
	expiryStr, _ := request["expiry"].(string)
	readOnly, _ := request["read_only"].(bool)

	expiry, err := time.ParseDuration(expiryStr)
	if err != nil {
		s.sendError(conn, "invalid expiry format", err)
		return
	}

	parser := core.NewParser()
	envFile := parser.ParseContent(data, "shared")

	share, err := s.CreateShare(envFile, shareCode, expiry, readOnly)
	if err != nil {
		s.sendError(conn, "failed to create share", err)
		return
	}

	response := map[string]interface{}{
		"success":    true,
		"share_id":   share.ID,
		"share_code": share.ShareCode,
		"expires_at": share.ExpiresAt,
	}

	s.sendResponse(conn, response)
}

func (s *ZitiService) handleGetRequest(conn net.Conn, request map[string]interface{}) {
	shareCode, _ := request["share_code"].(string)

	share, err := s.GetShare(shareCode)
	if err != nil {
		s.sendError(conn, "failed to get share", err)
		return
	}

	response := map[string]interface{}{
		"success":        true,
		"share_code":     share.ShareCode,
		"encrypted_data": share.EncryptedData,
		"metadata":       share.Metadata,
	}

	s.sendResponse(conn, response)
}

func (s *ZitiService) handleHealthRequest(conn net.Conn) {
	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now(),
		"version":   "1.0.0",
		"service":   s.config.ServiceName,
	}
	s.sendResponse(conn, response)
}

func (s *ZitiService) CreateShare(envFile *core.EnvFile, shareCode string, expiry time.Duration, readOnly bool) (*core.Share, error) {
	encryptedData, err := s.encryption.EncryptEnvFile(envFile, shareCode)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt env file: %w", err)
	}

	share := &core.Share{
		ID:            uuid.New().String(),
		ShareCode:     shareCode,
		EncryptedData: encryptedData,
		CreatedAt:     time.Now(),
		ExpiresAt:     time.Now().Add(expiry),
		IsReadOnly:    readOnly,
		AccessCount:   0,
		MaxAccesses:   1,
		Metadata: map[string]interface{}{
			"variable_count":  len(envFile.Variables),
			"sensitive_count": s.countSensitiveVariables(envFile),
			"file_size":       len(envFile.RawContent),
		},
	}

	s.sharesMu.Lock()
	s.shares[shareCode] = share
	s.sharesMu.Unlock()

	return share, nil
}

func (s *ZitiService) GetShare(shareCode string) (*core.Share, error) {
	s.sharesMu.RLock()
	share, exists := s.shares[shareCode]
	s.sharesMu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("share not found: %s", shareCode)
	}

	if time.Now().After(share.ExpiresAt) {
		s.removeShare(shareCode)
		return nil, fmt.Errorf("share expired: %s", shareCode)
	}

	if share.AccessCount >= share.MaxAccesses {
		s.removeShare(shareCode)
		return nil, fmt.Errorf("share access limit exceeded: %s", shareCode)
	}

	s.sharesMu.Lock()
	share.AccessCount++
	s.sharesMu.Unlock()

	if share.AccessCount >= share.MaxAccesses {
		s.removeShare(shareCode)
	}

	return share, nil
}

func (s *ZitiService) removeShare(shareCode string) {
	s.sharesMu.Lock()
	delete(s.shares, shareCode)
	s.sharesMu.Unlock()
}

func (s *ZitiService) countSensitiveVariables(envFile *core.EnvFile) int {
	count := 0
	for _, variable := range envFile.Variables {
		if variable.IsSensitive {
			count++
		}
	}
	return count
}

func (s *ZitiService) startCleanupRoutine() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.cleanupExpiredShares()
	}
}

func (s *ZitiService) cleanupExpiredShares() {
	now := time.Now()
	expiredShares := []string{}

	s.sharesMu.RLock()
	for shareCode, share := range s.shares {
		if now.After(share.ExpiresAt) {
			expiredShares = append(expiredShares, shareCode)
		}
	}
	s.sharesMu.RUnlock()

	for _, shareCode := range expiredShares {
		s.removeShare(shareCode)
	}
}

func (s *ZitiService) sendResponse(conn net.Conn, data interface{}) {
	responseData, err := json.Marshal(data)
	if err != nil {
		s.sendError(conn, "failed to marshal response", err)
		return
	}
	conn.Write(responseData)
}

func (s *ZitiService) sendError(conn net.Conn, message string, err error) {
	errorResponse := map[string]interface{}{
		"success": false,
		"error":   message,
		"details": err.Error(),
	}
	s.sendResponse(conn, errorResponse)
}

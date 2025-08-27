package server

import (
	"context"
	"fmt"
	"time"

	"github.com/devlink/internal/util"
)

// NewServer creates a new HTTP server
func NewServer(config *util.Config) (*Server, error) {
	encryption := util.NewEncryptionManager(config)

	server := &Server{
		config:     config,
		shares:     make(map[string]*Share),
		encryption: encryption,
		network:    &NetworkManager{},
	}

	setupRoutes(server)
	return server, nil
}

// Start starts the HTTP server
func (s *Server) Start() error {
	fmt.Printf("🚀 Starting DevLink server on port %d\n", s.config.ServerPort)
	go s.startCleanupRoutine()
	return s.server.ListenAndServe()
}

// Stop gracefully stops the server
func (s *Server) Stop(ctx context.Context) error {
	fmt.Println("🛑 Stopping DevLink server...")
	return s.server.Shutdown(ctx)
}

// CreateShare creates a new share with encrypted data
func (s *Server) CreateShare(envFile *util.EnvFile, shareCode string, expiry time.Duration, readOnly bool) (*Share, error) {
	encryptedData, err := s.encryption.EncryptEnvFile(envFile, shareCode)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt env file: %w", err)
	}

	share := &Share{
		ID:            util.GenerateUUID(),
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

// GetShare retrieves a share by share code
func (s *Server) GetShare(shareCode string) (*Share, error) {
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

func (s *Server) removeShare(shareCode string) {
	s.sharesMu.Lock()
	delete(s.shares, shareCode)
	s.sharesMu.Unlock()
}

func (s *Server) countSensitiveVariables(envFile *util.EnvFile) int {
	count := 0
	for _, variable := range envFile.Variables {
		if variable.IsSensitive {
			count++
		}
	}
	return count
}

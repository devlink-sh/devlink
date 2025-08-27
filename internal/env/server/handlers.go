package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/devlink/internal/env"
)

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now(),
		"version":   "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleShare(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		ShareCode string `json:"share_code"`
		Data      string `json:"data"`
		Expiry    string `json:"expiry"`
		ReadOnly  bool   `json:"read_only"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	expiry, err := time.ParseDuration(request.Expiry)
	if err != nil {
		http.Error(w, "Invalid expiry format", http.StatusBadRequest)
		return
	}

	parser := env.NewParser()
	envFile, err := parser.ParseContent(request.Data, "shared")
	if err != nil {
		http.Error(w, "Invalid environment data", http.StatusBadRequest)
		return
	}

	share, err := s.CreateShare(envFile, request.ShareCode, expiry, request.ReadOnly)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create share: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success":    true,
		"share_id":   share.ID,
		"share_code": share.ShareCode,
		"expires_at": share.ExpiresAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	shareCode := r.URL.Path[len("/get/"):]
	if shareCode == "" {
		http.Error(w, "Share code required", http.StatusBadRequest)
		return
	}

	share, err := s.GetShare(shareCode)
	if err != nil {
		http.Error(w, fmt.Sprintf("Share not found: %v", err), http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"success":        true,
		"share_code":     share.ShareCode,
		"encrypted_data": share.EncryptedData,
		"metadata":       share.Metadata,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.sharesMu.RLock()
	activeShares := 0
	totalAccesses := 0
	for _, share := range s.shares {
		if time.Now().Before(share.ExpiresAt) {
			activeShares++
			totalAccesses += share.AccessCount
		}
	}
	s.sharesMu.RUnlock()

	response := map[string]interface{}{
		"active_shares":  activeShares,
		"total_accesses": totalAccesses,
		"server_uptime":  time.Duration(0),
		"version":        "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

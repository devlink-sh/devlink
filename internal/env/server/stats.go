package server

import "time"

// GetStats returns server statistics
func (s *Server) GetStats() map[string]interface{} {
	s.sharesMu.RLock()
	defer s.sharesMu.RUnlock()

	activeShares := 0
	expiredShares := 0
	totalAccesses := 0

	for _, share := range s.shares {
		if time.Now().Before(share.ExpiresAt) {
			activeShares++
			totalAccesses += share.AccessCount
		} else {
			expiredShares++
		}
	}

	return map[string]interface{}{
		"total_shares":     len(s.shares),
		"active_shares":    activeShares,
		"expired_shares":   expiredShares,
		"total_accesses":   totalAccesses,
		"server_port":      s.config.ServerPort,
		"encryption_ready": s.encryption != nil,
	}
}

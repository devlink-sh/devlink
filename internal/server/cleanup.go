package server

import (
	"fmt"
	"time"
)

func (s *Server) startCleanupRoutine() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.cleanupExpiredShares()
	}
}

func (s *Server) cleanupExpiredShares() {
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

	if len(expiredShares) > 0 {
		fmt.Printf("🧹 Cleaned up %d expired shares\n", len(expiredShares))
	}
}

package service

import (
	"log"
	"time"

	"cctv-monitoring-backend/internal/repository"
)

// CleanupService handles periodic cleanup tasks
type CleanupService struct {
	tokenRepo repository.TokenRepository
}

// NewCleanupService creates a new cleanup service
func NewCleanupService(tokenRepo repository.TokenRepository) *CleanupService {
	return &CleanupService{
		tokenRepo: tokenRepo,
	}
}

// StartCleanupJob runs periodic cleanup of expired tokens
func (s *CleanupService) StartCleanupJob(interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		for range ticker.C {
			log.Println("Running token blacklist cleanup...")
			if err := s.tokenRepo.CleanupExpiredTokens(); err != nil {
				log.Printf("Error cleaning up expired tokens: %v", err)
			}
		}
	}()

	log.Printf("✓ Token cleanup job started (interval: %v)", interval)
}

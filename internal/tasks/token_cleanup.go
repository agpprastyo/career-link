package tasks

import (
	"context"
	"github.com/agpprastyo/career-link/pkg/logger"
	"time"
)

type TokenRepository interface {
	DeleteExpiredTokens(ctx context.Context, before time.Time) error
}

type TokenCleanupService struct {
	repo TokenRepository
	log  *logger.Logger
	quit chan struct{}
}

func NewTokenCleanupService(repo TokenRepository, log *logger.Logger) *TokenCleanupService {
	return &TokenCleanupService{
		repo: repo,
		log:  log,
		quit: make(chan struct{}),
	}
}

func (s *TokenCleanupService) Start() {
	go s.scheduleDailyCleanup()
}

func (s *TokenCleanupService) Stop() {
	close(s.quit)
}

func (s *TokenCleanupService) scheduleDailyCleanup() {
	// Calculate time until next midnight UTC+7
	now := time.Now()
	loc, err := time.LoadLocation("Asia/Jakarta") // UTC+7
	if err != nil {
		s.log.WithError(err).Error("Failed to load timezone, using local timezone")
		loc = time.Local
	}

	now = now.In(loc)
	nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, loc)
	duration := nextMidnight.Sub(now)

	// First run at next midnight
	timer := time.NewTimer(duration)

	for {
		select {
		case <-timer.C:
			// Execute cleanup
			s.cleanupExpiredTokens()

			// Schedule next run for tomorrow midnight
			timer.Reset(24 * time.Hour)
		case <-s.quit:
			timer.Stop()
			return
		}
	}
}

func (s *TokenCleanupService) cleanupExpiredTokens() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	now := time.Now()
	s.log.Info("Starting expired token cleanup")

	if err := s.repo.DeleteExpiredTokens(ctx, now); err != nil {
		s.log.WithError(err).Error("Failed to clean up expired tokens")
		return
	}

	s.log.Info("Successfully cleaned up expired tokens")
}

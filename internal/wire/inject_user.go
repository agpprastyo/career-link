package wire

import (
	"github.com/agpprastyo/career-link/config"
	"github.com/agpprastyo/career-link/internal/user/delivery"
	"github.com/agpprastyo/career-link/internal/user/repository"
	"github.com/agpprastyo/career-link/internal/user/usecase"
	"github.com/agpprastyo/career-link/pkg/database"
	"github.com/agpprastyo/career-link/pkg/logger"
	"github.com/agpprastyo/career-link/pkg/mail"
	"github.com/agpprastyo/career-link/pkg/minio"
	"github.com/agpprastyo/career-link/pkg/redis"
	"github.com/agpprastyo/career-link/pkg/token"
)

func provideUserRepository(db *database.PostgresDB, log *logger.Logger, mail *mail.Client, minio *minio.Client, verifyBaseURL string) *repository.UserRepository {
	return repository.NewUserRepository(db, log, mail, minio, verifyBaseURL)
}

func provideUserUseCase(repo *repository.UserRepository, log *logger.Logger, token token.Maker) *usecase.UserUseCase {
	return usecase.NewUserUseCase(repo, log, token)
}

func provideUserHandler(uc *usecase.UserUseCase, log *logger.Logger, cfg *config.AppConfig, tokenMaker token.Maker, redisClient *redis.Client, repo *repository.UserRepository) *delivery.Handler {
	return delivery.NewUserHandler(uc, log, cfg, tokenMaker, redisClient, repo)
}

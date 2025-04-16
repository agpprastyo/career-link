package usecase

import (
	"github.com/agpprastyo/career-link/internal/user/repository"
	"github.com/agpprastyo/career-link/pkg/logger"
	token2 "github.com/agpprastyo/career-link/pkg/token"
)

// UserUseCase implements users business logic
type UserUseCase struct {
	repo       repository.Repository
	log        *logger.Logger
	tokenMaker token2.Maker
}

// NewUserUseCase creates a new usersUseCase instance
func NewUserUseCase(repo repository.Repository, log *logger.Logger, tokenMaker token2.Maker) *UserUseCase {
	return &UserUseCase{
		repo:       repo,
		log:        log,
		tokenMaker: tokenMaker,
	}
}

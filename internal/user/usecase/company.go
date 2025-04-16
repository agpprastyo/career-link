package usecase

import (
	"context"
	"github.com/agpprastyo/career-link/internal/user/entity"
	"github.com/google/uuid"
)

// GetCompany returns a company by ID
func (uc *UserUseCase) GetCompany(ctx context.Context, id uuid.UUID) (*entity.Company, error) {
	company, err := uc.repo.GetCompanyByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return company, nil
}

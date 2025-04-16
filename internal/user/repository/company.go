package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/agpprastyo/career-link/internal/user/entity"
	"github.com/google/uuid"
)

// GetCompanyByID retrieves a company by its ID
func (r *UserRepository) GetCompanyByID(ctx context.Context, id uuid.UUID) (*entity.Company, error) {
	const query = `
		SELECT id, user_id, name, description, industry, website, email, phone, logo_url, status, size, is_verified, address_id, created_at, updated_at
		FROM companies
		WHERE id = $1
	`
	company := &entity.Company{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&company.ID,
		&company.UserID,
		&company.Name,
		&company.Description,
		&company.Industry,
		&company.Website,
		&company.Email,
		&company.Phone,
		&company.LogoURL,
		&company.Status,
		&company.Size,
		&company.IsVerified,
		&company.AddressID,
		&company.CreatedAt,
		&company.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		r.log.WithError(err).Error("Failed to get company by ID")
		return nil, err
	}

	return company, nil
}

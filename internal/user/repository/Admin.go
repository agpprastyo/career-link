package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/agpprastyo/career-link/internal/user/entity"
	"github.com/google/uuid"
)

// CreateAdminUser creates a new admin user
func (r *UserRepository) CreateAdminUser(ctx context.Context, admin *entity.Admin, usr *entity.User) error {
	// Set user role to admin
	usr.Role = entity.AdminRole
	usr.IsActive = true

	// Begin transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.log.WithError(err).Error("Failed to begin transaction")
		return err
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			r.log.WithError(err).Error("Failed to rollback transaction")
		}
	}(tx)

	// Insert the user first
	userQuery := `
        INSERT INTO users (id, username, email, password, role, avatar, is_active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
    `
	_, err = tx.ExecContext(ctx, userQuery,
		usr.ID,
		usr.Username,
		usr.Email,
		usr.Password,
		usr.Role,
		usr.Avatar,
		usr.IsActive)

	if err != nil {
		r.log.WithError(err).Error("Failed to create user record for admin")
		return err
	}

	// Insert admin-specific data if needed
	adminQuery := `
        INSERT INTO admins (id, user_id, role, last_login)
        VALUES ($1, $2, $3, $4)
    `
	_, err = tx.ExecContext(ctx, adminQuery,
		usr.ID,
		admin.UserID)

	if err != nil {
		r.log.WithError(err).Error("Failed to create admin details")
		return err
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		r.log.WithError(err).Error("Failed to commit transaction")
		return err
	}

	return nil
}

// GetAdminByUserID retrieves admin details by user ID
func (r *UserRepository) GetAdminByUserID(ctx context.Context, userID uuid.UUID) (*entity.Admin, error) {
	admin := &entity.Admin{}
	query := `
		SELECT id, user_id, role, last_login, is_active, created_at, updated_at
		FROM admins
		WHERE user_id = $1
	`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&admin.ID,
		&admin.UserID,
		&admin.Role,
		&admin.LastLogin,
		&admin.IsActive,
		&admin.CreatedAt,
		&admin.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No admin found
		}
		r.log.WithError(err).Error("Failed to get admin by user ID")
		return nil, err
	}

	return admin, nil
}

// UpdateAdminStatus updates the status of an admin user
func (r *UserRepository) UpdateAdminStatus(ctx context.Context, userID uuid.UUID) (bool, error) {
	var isActive bool
	query := `
		UPDATE admins
		SET is_active = NOT is_active
		WHERE user_id = $1
		RETURNING is_active
	`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&isActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, ErrUserNotFound
		}

		r.log.WithError(err).Error("Failed to update admin status")
		return false, err
	}

	return isActive, nil

}

// SoftDeleteAdmin marks an admin user as deleted
func (r *UserRepository) SoftDeleteAdmin(ctx context.Context, userID uuid.UUID) error {
	query := `
        UPDATE admins 
        SET is_active = false, deleted_at = NOW(), updated_at = NOW()
        WHERE user_id = $1 AND deleted_at IS NULL
        RETURNING id
    `

	var id uuid.UUID
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}

	return nil
}

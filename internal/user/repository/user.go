package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/agpprastyo/career-link/internal/common/pagination"
	"github.com/agpprastyo/career-link/internal/user/entity"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"time"
)

// GetUsersPaginated retrieves a paginated list of users
func (r *UserRepository) GetUsersPaginated(ctx context.Context, paging pagination.Pagination) ([]entity.User, int, error) {
	// Query to get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM users`

	err := r.db.GetContext(ctx, &total, countQuery)
	if err != nil {
		return nil, 0, errors.New("failed to get total count")
	}

	// Query to get paginated users
	var users []entity.User
	query := `SELECT id, username, email, role, avatar, is_active, created_at, updated_at FROM users
              ORDER BY created_at DESC` + paging.GetSQLLimitOffset()

	err = r.db.SelectContext(ctx, &users, query)
	log.Info("users: ", users)
	if err != nil {
		log.Error(err)
		return nil, 0, errors.New("failed to get users 1")
	}

	fmt.Println("users: ", users)
	fmt.Println("total: ", total)

	return users, total, nil
}

func (r *UserRepository) TokenExists(ctx context.Context, token string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM verification_tokens WHERE token = $1)`

	err := r.db.QueryRowContext(ctx, query, token).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// GetUserByEmail gets a user by email for authentication
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	const query = `
  SELECT id, username, email, password, role, avatar, is_active
  FROM users
  WHERE email = $1
 `

	var usr entity.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&usr.ID,
		&usr.Username,
		&usr.Email,
		&usr.Password,
		&usr.Role,
		&usr.Avatar,
		&usr.IsActive,
	)
	if err != nil {
		r.log.WithError(err).WithField("email", email).Error("Failed to get user by email")
		return nil, ErrUserNotFound
	}

	if &usr == nil {
		r.log.WithError(err).WithField("email", email).Error("User not found")
		return nil, ErrUserNotFound
	}

	return &usr, nil
}

// GetUserByUsername gets a user by username for authentication
func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	const query = `
		SELECT id, username, email, password, role, avatar, is_active
		FROM users
		WHERE username = $1
	`

	var usr entity.User
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&usr.ID,
		&usr.Username,
		&usr.Email,
		&usr.Password,
		&usr.Role,
		&usr.Avatar,
		&usr.IsActive,
	)
	if err != nil {
		r.log.WithError(err).WithField("username", username).Error("Failed to get user by username")
		return nil, ErrUserNotFound
	}

	return &usr, nil
}

// GetUserByID gets a user by ID
func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	const query = `
		SELECT id, username, email, password, role, avatar, is_active
		FROM users
		WHERE id = $1
	`

	var usr entity.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&usr.ID,
		&usr.Username,
		&usr.Email,
		&usr.Password,
		&usr.Role,
		&usr.Avatar,
		&usr.IsActive)
	if err != nil {
		r.log.WithError(err).WithField("user_id", id).Error("Failed to get user by ID")
		return nil, ErrUserNotFound
	}

	return &usr, nil
}

// CreateUser creates a new user
func (r *UserRepository) CreateUser(ctx context.Context, usr *entity.User) error {
	const query = `
		INSERT INTO users (id, username, email, password, role, avatar)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, query, usr.ID, usr.Username, usr.Email, usr.Password, usr.Role, usr.Avatar)
	if err != nil {
		r.log.WithError(err).WithField("user", usr).Error("Failed to create user")
		return err
	}

	return nil
}

// UpdateUser updates a user
func (r *UserRepository) UpdateUser(ctx context.Context, usr *entity.User) error {
	const query = `
		UPDATE users
		SET username = $1, email = $2, password = $3, role = $4, avatar = $5, is_active = $6
		WHERE id = $7
	`

	_, err := r.db.ExecContext(ctx, query, usr.Username, usr.Email, usr.Password, usr.Role, usr.Avatar, usr.IsActive, usr.ID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrUserNotFound
		case errors.Is(err, sql.ErrConnDone):
			return ErrDatabaseConnection
		default:
			return err
		}
	}

	return nil
}

// ActivateUser sets a user's active status to true by ID
func (r *UserRepository) ActivateUser(ctx context.Context, userID uuid.UUID) error {
	const query = `
        UPDATE users
        SET is_active = true
        WHERE id = $1
    `

	result, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		r.log.WithError(err).WithField("user_id", userID).Error("Failed to activate user")
		switch {
		case errors.Is(err, sql.ErrConnDone):
			return ErrDatabaseConnection
		default:
			return err
		}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) CreateToken(ctx context.Context, userID uuid.UUID, tokenID uuid.UUID, token string, tokenType entity.TokenType, expiry time.Time) error {

	const query = `
		INSERT INTO verification_tokens (id, user_id, token, type, expired_at, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		`

	_, err := r.db.ExecContext(ctx, query, tokenID, userID, token, tokenType, expiry)
	if err != nil {
		r.log.WithError(err).WithField("token", token).Error("Failed to create token")
		return err
	}

	return nil
}

// GetToken gets a token by token
func (r *UserRepository) GetToken(ctx context.Context, token string) (*entity.VerificationToken, error) {
	const query = `
		SELECT id, user_id, token, type, expired_at, created_at, used_at
		FROM verification_tokens
		WHERE token = $1
	`

	var t entity.VerificationToken
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&t.ID,
		&t.UserID,
		&t.Token,
		&t.Type,
		&t.ExpiredAt,
		&t.CreatedAt,
		&t.UsedAt,
	)
	if err != nil {
		r.log.WithError(err).WithField("token", token).Error("Failed to get token")
		return nil, err
	}

	return &t, nil
}

// IsTokenCorrect checks if the token is correct
func (r *UserRepository) IsTokenCorrect(ctx context.Context, token string, tokenType string) (bool, error) {
	const query = `
		SELECT id
		FROM verification_tokens
		WHERE token = $1 AND type = $2 AND expired_at > NOW()
	`

	var id uuid.UUID
	err := r.db.QueryRowContext(ctx, query, token, tokenType).Scan(&id)
	if err != nil {
		r.log.WithError(err).WithField("token", token).Error("Failed to check token")
		return false, err
	}

	return true, nil
}

// UpdateToken updates a token
func (r *UserRepository) UpdateToken(ctx context.Context, id uuid.UUID, usedAt time.Time) error {
	const query = `
		UPDATE verification_tokens
		SET used_at = $1
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, usedAt, id)
	if err != nil {
		r.log.WithError(err).WithField("token_id", id).Error("Failed to update token")
		return err
	}

	return nil
}

// DeleteExpiredTokens DeleteExpiredTokens(ctx context.Context, before time.Time) error
func (r *UserRepository) DeleteExpiredTokens(ctx context.Context, before time.Time) error {
	const query = `
		DELETE FROM verification_tokens
		WHERE expired_at < $1
	`

	_, err := r.db.ExecContext(ctx, query, before)
	if err != nil {
		r.log.WithError(err).WithField("before", before).Error("Failed to delete expired tokens")
		return err
	}

	return nil
}

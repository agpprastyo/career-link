package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/agpprastyo/career-link/internal/user/entity"
	"github.com/google/uuid"
	"strings"
	"time"
)

// StoreUserSession stores user data in Redis for session management
func (r *UserRepository) StoreUserSession(ctx context.Context, userIDStr string, sessionID string, user *entity.User) error {
	if r.redis == nil {
		r.log.Error("Redis client is nil")
		return errors.New("redis client not initialized")
	}

	// Parse user ID to UUID for consistency
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		r.log.WithError(err).WithField("user_id", userIDStr).Error("Failed to parse user ID")
		return err
	}

	// Marshal user data to JSON
	userJSON, err := json.Marshal(user)
	if err != nil {
		r.log.WithError(err).Error("Failed to marshal user data")
		return err
	}

	// Store user data in Redis
	userSessionKey := fmt.Sprintf("session:%s", userID)
	if err := r.redis.Set(ctx, userSessionKey, string(userJSON), time.Hour*24); err != nil {
		r.log.WithError(err).WithField("session_key", userSessionKey).Error("Failed to store user session")
		return err
	}

	// If user is admin, fetch and store admin data
	if user.Role == entity.AdminRole {
		admin, err := r.GetAdminByUserID(ctx, userID)
		if err == nil && admin != nil {
			adminJSON, err := json.Marshal(admin)
			if err == nil {
				adminSessionKey := fmt.Sprintf("admin:%s", userID)
				if err := r.redis.Set(ctx, adminSessionKey, string(adminJSON), time.Hour*24); err != nil {
					r.log.WithError(err).WithField("admin_key", adminSessionKey).Error("Failed to store admin session")
				}
			}
		}
	}

	// Store session to user mapping (for lookup by session ID)
	return r.redis.Set(ctx, sessionID, userIDStr, time.Hour*24)
}

// GetUserSession retrieves a user session from Redis
func (r *UserRepository) GetUserSession(ctx context.Context, sessionID string) (*entity.User, error) {
	// Get user ID from session
	userIDStr, err := r.redis.Get(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// Get user data from Redis
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}

	userSessionKey := fmt.Sprintf("session:%s", userID)
	userJSON, err := r.redis.Get(ctx, userSessionKey)
	if err != nil {
		return nil, err
	}

	// Unmarshal user data
	var user entity.User
	if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// DeleteUserSession removes a user session from Redis
func (r *UserRepository) DeleteUserSession(ctx context.Context, userIDKey string) error {
	// Delete based on the user ID key directly
	if err := r.redis.Del(ctx, userIDKey); err != nil {
		r.log.Error("Failed to delete user session")
		return errors.New("failed to delete user session")
	}

	// Also delete admin session if it exists
	userID := strings.TrimPrefix(userIDKey, "session:")
	adminSessionKey := fmt.Sprintf("admin:%s", userID)
	if err := r.redis.Del(ctx, adminSessionKey); err != nil {
		r.log.Error("Failed to delete admin session")
		return errors.New("failed to delete admin session")
	}

	return nil
}

// GetUserSessionByID retrieves user data directly by userID
func (r *UserRepository) GetUserSessionByID(ctx context.Context, userID string) (*entity.User, error) {
	userSessionKey := fmt.Sprintf("session:%s", userID)
	userJSON, err := r.redis.Get(ctx, userSessionKey)
	if err != nil {
		return nil, err
	}

	var user entity.User
	if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// GetAdminSessionByID retrieves admin data by userID
func (r *UserRepository) GetAdminSessionByID(ctx context.Context, userID string) (*entity.Admin, error) {
	adminSessionKey := fmt.Sprintf("admin:%s", userID)
	adminJSON, err := r.redis.Get(ctx, adminSessionKey)
	if err != nil {
		return nil, err
	}

	var admin entity.Admin
	if err := json.Unmarshal([]byte(adminJSON), &admin); err != nil {
		return nil, err
	}

	return &admin, nil
}

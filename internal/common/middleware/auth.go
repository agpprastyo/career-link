package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	responseError "github.com/agpprastyo/career-link/internal/common/errors"
	"github.com/agpprastyo/career-link/internal/user/entity"
	"github.com/agpprastyo/career-link/internal/user/repository"
	"github.com/agpprastyo/career-link/pkg/logger"
	"github.com/agpprastyo/career-link/pkg/redis"
	"github.com/agpprastyo/career-link/pkg/token"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"strings"
	"time"
)

func RequireAdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		usr := c.Locals("user").(entity.User)
		if usr.Role != entity.AdminRole {
			return responseError.RespondWithError(c, fiber.StatusForbidden, "Admin role is required")
		}
		return c.Next()
	}
}

// RequireSuperAdminMiddleware Super Admin Middleware
func RequireSuperAdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		//usr := c.Locals("user").(entity.User)
		admin := c.Locals("admin").(entity.Admin)
		if admin.Role != entity.AdminRoleSuper {
			return responseError.RespondWithError(c, fiber.StatusForbidden, "Super Admin role is required")
		}
		return c.Next()
	}
}

// RequireCompanyMiddleware Company Middleware
func RequireCompanyMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		usr := c.Locals("user").(entity.User)
		if usr.Role != entity.CompanyRole {
			return responseError.RespondWithError(c, fiber.StatusForbidden, "Company role is required")
		}
		return c.Next()
	}
}

// RequireAuthMiddleware creates a middleware that validates JWT tokens and stores user data in Redis
func RequireAuthMiddleware(tokenMaker token.Maker, redisClient *redis.Client, userRepo *repository.UserRepository, log *logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract and validate token (existing code)
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return responseError.RespondWithError(c, fiber.StatusUnauthorized, "Authorization header is required")
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 || strings.ToLower(fields[0]) != "bearer" {
			return responseError.RespondWithError(c, fiber.StatusUnauthorized, "Invalid authorization format")
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			return responseError.RespondWithError(c, fiber.StatusUnauthorized, "Invalid or expired token")
		}

		// Set basic user data from token
		c.Locals("user_id", payload.UserID)
		c.Locals("email", payload.Email)

		// Try to get user from Redis first
		ctx := c.Context()
		sessionID := fmt.Sprintf("session:%s", payload.UserID)
		userJSON, err := redisClient.Get(ctx, sessionID)

		// User found in Redis
		if err == nil && userJSON != "" {
			var userData entity.User
			if err := json.Unmarshal([]byte(userJSON), &userData); err == nil {
				c.Locals("user", userData)

				// Check if admin and get admin data from Redis
				if userData.Role == entity.AdminRole {
					adminSessionID := fmt.Sprintf("admin:%s", payload.UserID)
					adminJSON, err := redisClient.Get(ctx, adminSessionID)
					if err == nil && adminJSON != "" {
						var adminData entity.Admin
						if err := json.Unmarshal([]byte(adminJSON), &adminData); err == nil {
							c.Locals("admin", adminData)
						}
					}
				}

				// Update expiration asynchronously
				go func() {
					redisClient.Expire(ctx, sessionID, 24*time.Hour)
				}()
				return c.Next()
			}
			// If unmarshal fails, continue to get from DB
		}

		// User not in Redis, fetch from database
		userUUID, err := uuid.Parse(payload.UserID)
		if err != nil {
			return responseError.RespondWithError(c, fiber.StatusUnauthorized, "Invalid user ID in token")
		}

		userData, err := userRepo.GetUserByID(ctx, userUUID)
		if err != nil {
			log.WithError(err).Error("Failed to get user by ID")
			return responseError.RespondWithError(c, fiber.StatusUnauthorized, "User not found")
		}

		// Store user in context
		c.Locals("user", userData)

		// If user is admin, also fetch and cache admin data
		if userData.Role == entity.AdminRole {

			adminData, err := userRepo.GetAdminByUserID(ctx, userUUID)
			if err == nil {
				c.Locals("admin", adminData)

				// Cache admin data asynchronously
				go func() {
					adminJSON, err := json.Marshal(adminData)
					if err != nil {
						log.WithError(err).Error("Failed to marshal admin data")
						return
					}
					adminSessionID := fmt.Sprintf("admin:%s", payload.UserID)
					err = redisClient.Set(context.Background(), adminSessionID, adminJSON, 24*time.Hour)
					if err != nil {
						log.WithError(err).Error("Failed to cache admin data")
						return
					}
				}()
			}
		}

		// Cache user data in Redis asynchronously
		go func() {
			userJSON, err := json.Marshal(userData)
			if err != nil {
				log.WithError(err).Error("Failed to marshal user data")
				return
			}
			err = redisClient.Set(context.Background(), sessionID, userJSON, 24*time.Hour)
			if err != nil {
				log.WithError(err).Error("Failed to cache user data")
				return
			}
		}()

		return c.Next()
	}
}

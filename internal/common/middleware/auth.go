package middleware

import (
	"context"

	responseError "github.com/agpprastyo/career-link/internal/common/errors"
	"github.com/agpprastyo/career-link/internal/user/entity"
	"github.com/agpprastyo/career-link/internal/user/repository"
	"github.com/agpprastyo/career-link/pkg/logger"
	"github.com/agpprastyo/career-link/pkg/token"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"strings"
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
func RequireAuthMiddleware(tokenMaker token.Maker, userRepo *repository.UserRepository, log *logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract and validate token
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

		// Extract user data from token
		userID := payload.UserID
		c.Locals("user_id", userID)
		c.Locals("email", payload.Email)

		// Try to get user from session using repository
		ctx := c.Context()
		user, err := userRepo.GetUserSessionByID(ctx, userID)

		if err == nil && user != nil {
			// User found in session
			c.Locals("user", *user)

			// If admin, get admin data
			if user.Role == entity.AdminRole {
				admin, err := userRepo.GetAdminSessionByID(ctx, userID)
				if err == nil && admin != nil {
					c.Locals("admin", *admin)
				}
			}

			return c.Next()
		}

		// Fetch from database if not in session
		userUUID, err := uuid.Parse(userID)
		if err != nil {
			return responseError.RespondWithError(c, fiber.StatusUnauthorized, "Invalid user ID format")
		}

		user, err = userRepo.GetUserByID(ctx, userUUID)
		if err != nil {
			log.WithError(err).Error("User not found")
			return responseError.RespondWithError(c, fiber.StatusUnauthorized, "User not found")
		}

		c.Locals("user", *user)

		// Store in session asynchronously
		go func() {
			bgCtx := context.Background()
			err := userRepo.StoreUserSession(bgCtx, userID, accessToken, user)
			if err != nil {
				log.WithError(err).Error("Failed to store user session")
				return
			}
		}()

		return c.Next()
	}
}

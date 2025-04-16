package health

import (
	"context"
	"time"

	"github.com/agpprastyo/career-link/pkg/database"
	"github.com/agpprastyo/career-link/pkg/redis"
	"github.com/gofiber/fiber/v2"
)

// Handler manages health check functionality
type Handler struct {
	db    *database.PostgresDB
	redis *redis.Client
}

// NewHandler creates a new health check handler
func NewHandler(db *database.PostgresDB, redis *redis.Client) *Handler {
	return &Handler{
		db:    db,
		redis: redis,
	}
}

// Check handles health check requests
func (h *Handler) Check(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	// Check database connection
	dbErr := h.db.Ping(ctx)

	// Check Redis connection
	redisErr := h.redis.GetClient().Ping(ctx).Err()

	if dbErr != nil || redisErr != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status": "unhealthy",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "healthy",
	})
}

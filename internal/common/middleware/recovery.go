package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/agpprastyo/career-link/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus" // Add this import
)

// RecoveryMiddleware returns a middleware that recovers from panics
func RecoveryMiddleware(log *logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				// Get stack trace
				stack := debug.Stack()

				// Prepare error message
				err := fmt.Errorf("panic recovered: %v", r)

				// Log the panic with stack trace - use logrus.Fields instead of logger.Fields
				log.WithError(err).WithFields(logrus.Fields{
					"stack":  string(stack),
					"method": c.Method(),
					"path":   c.Path(),
					"ip":     c.IP(),
				}).Error("Request handler panic")

				// Return error response to client
				err = c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Internal server error",
				})
				if err != nil {
					return
				}
			}
		}()

		// Process request
		return c.Next()
	}
}

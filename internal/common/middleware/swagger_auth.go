package middleware

import (
	"github.com/agpprastyo/career-link/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
)

// NewSwaggerAuth creates a new middleware for Swagger authentication
func NewSwaggerAuth(cfg *config.AppConfig) fiber.Handler {
	return basicauth.New(basicauth.Config{
		Users: map[string]string{
			cfg.SwaggerAuth.Username: cfg.SwaggerAuth.Password,
		},
		Realm: "Swagger Documentation",
		Unauthorized: func(c *fiber.Ctx) error {
			// Explicitly set WWW-Authenticate header
			c.Set("WWW-Authenticate", `Basic realm="Swagger Documentation"`)
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized access to Swagger documentation")
		},
		Next: func(c *fiber.Ctx) bool {
			// Allow access to the Swagger UI without authentication
			if c.Path() == "/swagger/index.html" || c.Path() == "/swagger/" {
				return true
			}
			return false
		},
	})
}

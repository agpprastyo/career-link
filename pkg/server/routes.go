package server

import (
	"github.com/agpprastyo/career-link/internal/common/health"
	"github.com/gofiber/fiber/v2"
)

// setupRoutes configures all API routes
func (s *Server) setupRoutes() {
	// Create handlers
	healthHandler := health.NewHandler(s.db, s.redis)

	// Health check endpoint
	s.app.Get("/health", healthHandler.Check)

	// API routes group
	api := s.app.Group("/api/v1")

	// User routes
	api.Get("/user", s.handleGetUsers)

}

// Example route handler
func (s *Server) handleGetUsers(c *fiber.Ctx) error {
	// Implementation to be added
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}
